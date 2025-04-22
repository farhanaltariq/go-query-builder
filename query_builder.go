package querybuilder

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"query-builder/utils"
)

// Generate update query only non-zero fields of a model if the identifier is given.
// It generate raw SQL queries to perform the update.
// ex.: UPDATE table SET model... WHERE identifier = ?
func GenerateUpdateQuery(model interface{}, table, identifier string) (query string, args []interface{}, err error) {
	if !utils.IsPointerToStruct(model) {
		return "", nil, errors.New("model must be a struct")
	}

	updateCounter := 0
	var updateValueIdentifier interface{}
	var q strings.Builder
	q.WriteString("UPDATE " + table + " SET ")

	modelValue := reflect.ValueOf(model).Elem()

	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Field(i)
		fieldType := modelValue.Type().Field(i)

		// Get the column name from the gorm tag
		columnName := utils.GetColumnName(fieldType)
		if columnName == "" || reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			continue // Skip fields without a gorm tag and empty value
		}

		if columnName == identifier {
			updateValueIdentifier = field.Interface()
			continue
		}

		// Add non-zero fields to the query
		if updateCounter > 0 {
			q.WriteString(", ")
		}
		q.WriteString(columnName)
		q.WriteString(" = ?")
		args = append(args, field.Interface())
		updateCounter++
	}

	// If no fields to update, return early
	if updateCounter == 0 {
		return "", nil, errors.New("no fields to update")
	}

	q.WriteString(" WHERE " + identifier + " = ?")
	args = append(args, updateValueIdentifier)

	return q.String(), args, nil
}

// Generate get query only get fields of a model if the gorm tag given.
// It generate raw SQL queries to perform the select.
// Put some data to model struct to add WHERE clauses
// ex.: SELCT column... FROM table;
// ex.: SELCT column... FROM table WHERE ...;
func GenerateGetQuery(model interface{}, table string) (query string, err error) {
	if !utils.IsPointerToStruct(model) {
		return "", errors.New("model must be a struct")
	}

	getCounter := 0
	whereCounter := 0
	var q, w strings.Builder
	q.WriteString("SELECT ")
	w.WriteString("WHERE ")

	modelValue := reflect.ValueOf(model).Elem()

	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Field(i)
		fieldType := modelValue.Type().Field(i)

		// Get the column name from the gorm tag
		columnName := utils.GetColumnName(fieldType)
		if !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			if whereCounter > 0 {
				w.WriteString(" AND ")
			}
			w.WriteString(columnName + " = ")
			w.WriteString(fmt.Sprint(field.Interface()))
			whereCounter++
		}
		// Add non-zero fields to the query
		if getCounter > 0 {
			q.WriteString(", ")
		}
		q.WriteString(columnName)
		getCounter++
	}

	q.WriteString(From(table))
	if whereCounter != 0 {
		q.WriteString(" " + w.String())
	}

	return q.String(), nil
}

// Generate additional WHERE AND statements
func And(model interface{}) (query string, err error) {
	return "", nil
}

// Set from the query strings
func From(table string) (query string) {
	return " FROM " + table
}
