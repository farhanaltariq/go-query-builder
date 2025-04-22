package querybuilder

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type User struct {
	ID       uint32 `gorm:"column:id"`
	Username string `gorm:"column:username"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
}

func TestQueryBuilderUpdate(t *testing.T) {
	Convey("Given a User struct", t, func() {
		PK := uint32(1)
		newVal := "newVal"
		Convey("when filled only username, should only update username and leave the rest unchanged", func() {
			user := &User{ID: PK, Username: newVal}
			tableName := "users"

			expectedQuery := `UPDATE ` + tableName + ` SET username = ? WHERE id = ?`
			expectedArgs := []interface{}{
				newVal,
				user.ID,
			}
			query, args, err := GenerateUpdateQuery(user, tableName, "id")
			So(err, ShouldBeNil)
			So(query, ShouldEqual, expectedQuery)
			So(args, ShouldEqual, expectedArgs)
		})
		Convey("when filled only password, should only update password and leave the rest unchanged", func() {
			user := &User{ID: PK, Password: newVal}
			tableName := "users"

			expectedQuery := `UPDATE ` + tableName + ` SET password = ? WHERE id = ?`
			expectedArgs := []interface{}{
				newVal,
				user.ID,
			}
			query, args, err := GenerateUpdateQuery(user, tableName, "id")
			So(err, ShouldBeNil)
			So(query, ShouldEqual, expectedQuery)
			So(args, ShouldEqual, expectedArgs)
		})
		Convey("when filled username and password, should only update username and password and leave the rest unchanged", func() {
			user := &User{ID: PK, Password: newVal, Username: newVal}
			tableName := "users"

			expectedQuery := `UPDATE ` + tableName + ` SET username = ?, password = ? WHERE id = ?`
			expectedArgs := []interface{}{
				newVal,
				newVal,
				user.ID,
			}
			query, args, err := GenerateUpdateQuery(user, tableName, "id")
			So(err, ShouldBeNil)
			So(query, ShouldEqual, expectedQuery)
			So(args, ShouldEqual, expectedArgs)
		})
	})

	Convey("Given invalid struct", t, func() {
		user := &User{}
		tableName := "users"

		Convey("when given invalid input, should return error", func() {
			query, args, err := GenerateUpdateQuery("invalidtext", tableName, "id")
			So(err, ShouldNotBeNil)
			So(query, ShouldBeEmpty)
			So(args, ShouldBeEmpty)
		})
		Convey("when given struct with no update value, should return error", func() {
			query, args, err := GenerateUpdateQuery(user, tableName, "id")
			So(err, ShouldNotBeNil)
			So(query, ShouldBeEmpty)
			So(args, ShouldBeEmpty)
		})
	})
}

func TestQueryBuilderGet(t *testing.T) {
	Convey("Given a User struct", t, func() {
		PK := uint32(1)
		newVal := "newVal"
		Convey("when given user struct, should get gorm tags in user struct", func() {
			user := &User{}
			field := `id, username, email, password`
			tableName := "users"

			expectedQuery := `SELECT ` + field + ` FROM ` + tableName
			query, err := GenerateGetQuery(user, tableName)
			So(err, ShouldBeNil)
			So(query, ShouldEqual, expectedQuery)
		})
		Convey("when given user struct with id, should get gorm tags in user struct using where clauses", func() {
			user := &User{ID: PK}
			field := `id, username, email, password`
			tableName := "users"

			expectedQuery := `SELECT ` + field + ` FROM ` + tableName + ` WHERE id = 1`
			query, err := GenerateGetQuery(user, tableName)
			So(err, ShouldBeNil)
			So(query, ShouldEqual, expectedQuery)
		})
		Convey("when given user struct with multiple field filled, should get gorm tags in user struct using where clauses", func() {
			user := &User{ID: PK, Username: newVal}
			field := `id, username, email, password`
			tableName := "users"

			expectedQuery := `SELECT ` + field + ` FROM ` + tableName + ` WHERE id = 1 AND username = ` + newVal
			query, err := GenerateGetQuery(user, tableName)
			So(err, ShouldBeNil)
			So(query, ShouldEqual, expectedQuery)
		})
	})
	Convey("when given invalid input, should return error", t, func() {
		tableName := "users"
		query, err := GenerateGetQuery("invalidtext", tableName)
		So(err, ShouldNotBeNil)
		So(query, ShouldBeEmpty)
	})
}
