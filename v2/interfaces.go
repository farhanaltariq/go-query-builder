package new

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"query-builder/utils"
)

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

func (db *QueryBuilder) Select(model interface{}) *QueryBuilder {
	db.modes = &modes{mode: MODE_SELECT}
	if utils.IsString(model) {
		db.query = "SELECT " + model.(string)
		return db
	}
	if !utils.IsStruct(model) {
		db.error = errors.New("model must be a pointer to struct or a string")
		return db
	}

	getCounter := 0
	var q strings.Builder
	q.WriteString("SELECT ")

	modelValue := reflect.ValueOf(model).Elem()

	for i := 0; i < modelValue.NumField(); i++ {
		fieldType := modelValue.Type().Field(i)

		// Get the column name from the gorm tag
		columnName := utils.GetColumnName(fieldType)
		if columnName == "" {
			continue
		}
		// Add non-zero fields to the query
		if getCounter > 0 {
			q.WriteString(", ")
		}
		q.WriteString("`" + columnName + "`")
		getCounter++
	}

	if getCounter == 0 {
		db.error = errors.New("struct must have at least 1 field with gorm tag")
		return db
	}
	db.query = q.String()
	if db.table == nil {
		table := utils.GetStructName(model)
		db.table = &table
	}
	return db
}

func (db *QueryBuilder) Where(model interface{}) *QueryBuilder {
	if utils.IsString(model) {
		wq := model.(string)
		db.whereQuery = &wq
		return db
	}
	if !utils.IsStruct(model) {
		db.error = errors.New("model must be a pointer to struct or a string")
		return db
	}

	getCounter := 0
	var q strings.Builder

	modelValue := reflect.ValueOf(model).Elem()

	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Field(i)
		fieldType := modelValue.Type().Field(i)

		// Get the column name from the gorm tag
		columnName := utils.GetColumnName(fieldType)
		if columnName == "" || reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
			continue // Skip fields without a gorm tag and empty value
		}
		// Add non-zero fields to the query
		if getCounter > 0 {
			q.WriteString("AND ")
		}
		q.WriteString("`" + columnName + "` = " + fmt.Sprint(field.Interface()))
		getCounter++
	}

	if getCounter == 0 {
		db.error = errors.New("struct must have at least 1 field with gorm tag")
		return db
	}
	qPtr := q.String()
	db.whereQuery = &qPtr
	return db
}

func (db *QueryBuilder) From(table string) *QueryBuilder {
	if db.modes != nil && db.mode != MODE_SELECT {
		db.error = errors.New("can only used in select mode")
		return db
	}
	db.table = &table
	return db
}

func (db *QueryBuilder) Asc(column string) *QueryBuilder {
	db.direction = &direction{columnSort: column, dir: DIR_ASCENDING}
	return db
}

func (db *QueryBuilder) Desc(column string) *QueryBuilder {
	db.direction = &direction{columnSort: column, dir: DIR_DESCENDING}
	return db
}

func buildSelect(db *QueryBuilder) string {
	var q strings.Builder
	q.WriteString(db.query + " FROM " + *db.table)
	if db.whereQuery != nil {
		q.WriteString(" WHERE " + *db.whereQuery)
	}
	if db.direction != nil {
		q.WriteString(" ORDER BY `" + db.columnSort + "` " + db.dir)
	}
	q.WriteString(";")
	return q.String()
}

func setError(db *QueryBuilder, err error) {
	db.error = err
}

func (db *QueryBuilder) Raw() string {
	var q strings.Builder
	if db.modes == nil {
		setError(db, errors.New("cannot build queries"))
		db.query = ""
		return db.query
	}
	switch db.mode {
	case MODE_SELECT:
		return buildSelect(db)
	}
	return q.String()
}

func (db *QueryBuilder) Error() error {
	return db.error
}
