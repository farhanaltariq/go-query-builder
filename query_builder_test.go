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

func TestIsStructChecker(t *testing.T) {
	Convey("Given a struct, should return true", t, func() {
		PK := uint32(1)
		user := &User{}
		So(IsStruct(user), ShouldBeTrue)

		user = &User{ID: PK, Username: "john_doe"}
		So(IsStruct(user), ShouldBeTrue)
		type Example struct{}
		ex := Example{}
		So(IsStruct(ex), ShouldBeTrue)
	})
	Convey("Given a non struct, should return false", t, func() {
		So(IsStruct(nil), ShouldBeFalse)
		So(IsStruct("something"), ShouldBeFalse)

		var user interface{}
		So(IsStruct(user), ShouldBeFalse)
		So(IsStruct(42), ShouldBeFalse)
	})
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
	})
}

func TestQueryBuilderGet(t *testing.T) {
	Convey("Given a User struct", t, func() {
		PK := uint32(1)
		newVal := "newVal"
		Convey("when given user struct, should get gorm tags in user struct", func() {
			user := &User{ID: PK, Username: newVal}
			field := `id, username, email, password`
			tableName := "users"

			expectedQuery := `SELECT ` + field + ` FROM ` + tableName + `;`
			query,  err := GenerateGetQuery(user, tableName)
			So(err, ShouldBeNil)
			So(query, ShouldEqual, expectedQuery)
		})
	})
}
