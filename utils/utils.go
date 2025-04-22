package utils

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

func IsString(in interface{}) bool {
	if in == nil {
		return false
	}

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false
		}
		// Dereference pointer
		v = v.Elem()
	}

	return v.Kind() == reflect.String
}

// IsPointerToStruct returns true if in is a pointer to a struct.
func IsPointerToStruct(in interface{}) bool {
	if in == nil {
		return false
	}

	v := reflect.ValueOf(in)
	if v.Kind() != reflect.Ptr {
		return false
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false
		}
		// Dereference pointer
		v = v.Elem()
	}
	return v.IsValid() && v.Kind() == reflect.Struct
}

// getColumnName extracts the column name from the `gorm` tag of a struct field.
// If no column name is specified, it returns an empty string.
func GetColumnName(field reflect.StructField) string {
	if gormTag := field.Tag.Get("gorm"); gormTag != "" {
		// Extract the column name from the gorm tag
		for _, part := range strings.Split(gormTag, ";") {
			if strings.HasPrefix(part, "column:") {
				return strings.TrimPrefix(part, "column:")
			}
		}
	}
	return toSnakeCase(field.Name)
}

func GetStructName(i interface{}) string {
	t := reflect.TypeOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return strings.ToLower(t.Name()) + "s"
}

// identifierRegex defines an allowlist: only letters, digits, and underscores.
var identifierRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// sanitizeIdentifier ensures the given identifier (e.g., table/column name) is safe.
func SanitizeIdentifier(input string) (string, error) {
	if identifierRegex.MatchString(input) {
		return input, nil
	}
	return "", errors.New("invalid identifier: " + input)
}

// escapeString escapes single quotes in string values by doubling them.
func EscapeString(val string) string {
	return strings.ReplaceAll(val, "'", "''")
}

// toSnakeCase converts a string from CamelCase to snake_case.
func toSnakeCase(str string) string {
	var sb strings.Builder
	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 {
				sb.WriteRune('_')
			}
			sb.WriteRune(unicode.ToLower(r))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
