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

/*
Generate update query only non-zero fields of a model if the identifier is given.
It generate raw SQL queries to perform the update.
ex.: UPDATE table SET model... WHERE identifier = ?
*/
func GenerateUpdateQuery(model interface{}, table, identifier string) (query string, args []interface{}, err error) {
	if model == nil {
		return "", nil, errors.New("model cannot be nil")
	}

	var q strings.Builder

	return q.String(), nil, errors.New("update model")
}
