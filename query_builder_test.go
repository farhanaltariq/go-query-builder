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
		user := &User{}
		So(IsStruct(user), ShouldBeTrue)

		type Example struct {
		}
		ex := Example{}
		So(IsStruct(ex), ShouldBeTrue)
	})
	Convey("Given a non struct, should return false", t, func() {
		So(IsStruct(nil), ShouldBeFalse)
		So(IsStruct("something"), ShouldBeFalse)

		var user interface{}
		So(IsStruct(user), ShouldBeFalse)
	})
}

func TestQueryBuilderUpdate(t *testing.T) {
	Convey("Given a User struct", t, func() {
		Convey("when filled only username, should only update username and leave the rest unchanged", func() {
			PK := uint32(1)
			user := &User{ID: PK, Username: "john_doe"}
			tableName := "users"

			expectedQuery := `UPDATE ` + tableName + ` SET username = ? WHERE id = ?`
			expectedArgs := []interface{}{
				user.Username,
				user.ID,
			}
			query, args, err := GenerateUpdateQuery(&user, tableName, "id")
			So(query, ShouldEqual, expectedQuery)
			So(args, ShouldEqual, expectedArgs...)
			So(err, ShouldBeNil)
		})
	})
}
