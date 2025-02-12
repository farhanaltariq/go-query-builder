package querybuilder

import (
	"errors"
	"reflect"
	"strings"
)

// IsStruct returns true if in is a struct or a pointer to a struct.
func IsStruct(in interface{}) bool {
	if in == nil {
		return false
	}

	v := reflect.ValueOf(in)
	// If it's a pointer, check for nil pointer first
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false
		}

		// Dereference pointer
		v = v.Elem()
	}
	return v.Kind() == reflect.Struct
}

// getColumnName extracts the column name from the `gorm` tag of a struct field.
// If no column name is specified, it returns an empty string.
func getColumnName(field reflect.StructField) string {
	if gormTag := field.Tag.Get("gorm"); gormTag != "" {
		// Extract the column name from the gorm tag
		for _, part := range strings.Split(gormTag, ";") {
			if strings.HasPrefix(part, "column:") {
				return strings.TrimPrefix(part, "column:")
			}
		}
	}
	return ""
}

// Generate update query only non-zero fields of a model if the identifier is given.
// It generate raw SQL queries to perform the update.
// ex.: UPDATE table SET model... WHERE identifier = ?
func GenerateUpdateQuery(model interface{}, table, identifier string) (query string, args []interface{}, err error) {
	if !IsStruct(model) {
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
		columnName := getColumnName(fieldType)
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

	if updateValueIdentifier == nil {
		return "", nil, errors.New("update identiifier must be set")
	}

	q.WriteString(" WHERE " + identifier + " = ?")
	args = append(args, updateValueIdentifier)

	return q.String(), args, nil
}
