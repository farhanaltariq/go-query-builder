package new

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"query-builder/utils"
)


// --- QueryBuilder Implementation ---

type QueryBuilder struct {
	*modes
	query      string
	whereQuery *string
	whereValue map[string]interface{}
	table      *string
	*direction
	error
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

// Select accepts either a comma-separated string of columns or a pointer to a struct.
// When a struct is provided, it uses the gorm tag to determine the columns.
func (db *QueryBuilder) Select(model interface{}) *QueryBuilder {
	db.modes = &modes{mode: MODE_SELECT}

	// If model is a string, assume it's a comma-separated list of column names.
	if utils.IsString(model) {
		input := model.(string)
		columns := strings.Split(input, ",")
		var sanitizedCols []string
		for _, col := range columns {
			col = strings.TrimSpace(col)
			sanitized, err := utils.SanitizeIdentifier(col)
			if err != nil {
				db.error = err
				return db
			}
			sanitizedCols = append(sanitizedCols, "`"+sanitized+"`")
		}
		db.query = "SELECT " + strings.Join(sanitizedCols, ", ")
		return db
	}

	if !utils.IsStruct(model) {
		db.error = errors.New("model must be a pointer to struct or a string")
		return db
	}

	var q strings.Builder
	q.WriteString("SELECT ")

	modelValue := reflect.ValueOf(model).Elem()
	getCounter := 0

	for i := 0; i < modelValue.NumField(); i++ {
		fieldType := modelValue.Type().Field(i)

		// Get the column name from the gorm tag
		columnName := utils.GetColumnName(fieldType)
		if columnName == "" {
			continue
		}
		sanitized, err := utils.SanitizeIdentifier(columnName)
		if err != nil {
			db.error = err
			return db
		}

		// Append comma-separated list of columns.
		if getCounter > 0 {
			q.WriteString(", ")
		}
		q.WriteString("`" + sanitized + "`")
		getCounter++
	}

	if getCounter == 0 {
		db.error = errors.New("struct must have at least 1 field with gorm tag")
		return db
	}

	db.query = q.String()

	// If table not explicitly set, use the struct name as the table name.
	if db.table == nil {
		table := utils.GetStructName(model)
		sanitizedTable, err := utils.SanitizeIdentifier(table)
		if err != nil {
			db.error = err
			return db
		}
		db.table = &sanitizedTable
	}

	return db
}

// Where accepts either a raw string (trusted input) or a pointer to a struct.
// When a struct is provided, it builds conditions from non-zero fields.
func (db *QueryBuilder) Where(model interface{}) *QueryBuilder {
	// If model is a raw string, use it as-is (ensure you only use trusted input).
	if utils.IsString(model) {
		wq := model.(string)
		db.whereQuery = &wq
		return db
	}
	if !utils.IsStruct(model) {
		db.error = errors.New("model must be a pointer to struct or a string")
		return db
	}

	var q strings.Builder
	modelValue := reflect.ValueOf(model).Elem()
	getCounter := 0

	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Field(i)
		fieldType := modelValue.Type().Field(i)

		// Get the column name from the gorm tag
		columnName := utils.GetColumnName(fieldType)
		if columnName == "" || reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			continue // Skip fields without a gorm tag or with zero value.
		}
		sanitized, err := utils.SanitizeIdentifier(columnName)
		if err != nil {
			db.error = err
			return db
		}

		// Append conditions using AND.
		if getCounter > 0 {
			q.WriteString(" AND ")
		}

		// Handle string types by quoting and escaping.
		switch field.Kind() {
		case reflect.String:
			val := field.String()
			escapedVal := utils.EscapeString(val)
			q.WriteString("`" + sanitized + "` = '" + escapedVal + "'")
		default:
			q.WriteString("`" + sanitized + "` = " + fmt.Sprint(field.Interface()))
		}
		getCounter++
	}

	if getCounter == 0 {
		db.error = errors.New("struct must have at least 1 field with gorm tag and non-zero value")
		return db
	}
	qStr := q.String()
	db.whereQuery = &qStr
	return db
}

// From sets the table name; it sanitizes the table name to allow only safe characters.
func (db *QueryBuilder) From(table string) *QueryBuilder {
	if db.modes != nil && db.mode != MODE_SELECT {
		db.error = errors.New("can only be used in select mode")
		return db
	}
	sanitized, err := utils.SanitizeIdentifier(table)
	if err != nil {
		db.error = err
		return db
	}
	db.table = &sanitized
	return db
}

// Asc sets the ORDER BY clause in ascending order after sanitizing the column name.
func (db *QueryBuilder) Asc(column string) *QueryBuilder {
	sanitized, err := utils.SanitizeIdentifier(column)
	if err != nil {
		db.error = err
		return db
	}
	db.direction = &direction{columnSort: sanitized, dir: DIR_ASCENDING}
	return db
}

// Desc sets the ORDER BY clause in descending order after sanitizing the column name.
func (db *QueryBuilder) Desc(column string) *QueryBuilder {
	sanitized, err := utils.SanitizeIdentifier(column)
	if err != nil {
		db.error = err
		return db
	}
	db.direction = &direction{columnSort: sanitized, dir: DIR_DESCENDING}
	return db
}

// buildSelect builds the complete SELECT query string.
func buildSelect(db *QueryBuilder) string {
	var q strings.Builder
	q.WriteString(db.query + " FROM " + *db.table)
	if db.whereQuery != nil {
		q.WriteString(" WHERE " + *db.whereQuery)
	}
	if db.direction != nil {
		q.WriteString(" ORDER BY `" + db.direction.columnSort + "` " + db.direction.dir)
	}
	q.WriteString(";")
	return q.String()
}

func setError(db *QueryBuilder, err error) {
	db.error = err
}

// Raw returns the complete SQL query string for the current query builder state.
func (db *QueryBuilder) Raw() string {
	if db.modes == nil {
		setError(db, errors.New("cannot build queries"))
		db.query = ""
		return db.query
	}
	switch db.mode {
	case MODE_SELECT:
		return buildSelect(db)
	}
	return ""
}

// Error returns any error encountered during query building.
func (db *QueryBuilder) Error() error {
	return db.error
}
