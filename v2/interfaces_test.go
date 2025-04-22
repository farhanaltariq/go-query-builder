package new

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

func TestSelect(t *testing.T) {
	queryBuilder := NewQueryBuilder()
	Convey("Given a User struct", t, func() {
		Convey("when given user struct, should get gorm tags in user struct", func() {
			query := queryBuilder.Select(&User{})
			So(query.query, ShouldEqual, "SELECT `id`, `username`, `email`, `password`")
			So(*query.table, ShouldEqual, "users")
		})
		Convey("when given no gorm-tag struct, should set error", func() {
			type CustomStruct struct {
				hello string
			}
			query := queryBuilder.Select(&CustomStruct{})
			So(query.error, ShouldBeNil)
		})
	})
}

func TestFrom(t *testing.T) {
	queryBuilder := NewQueryBuilder()
	Convey("Given a User struct", t, func() {
		q := queryBuilder.From("user")
		So(*q.table, ShouldEqual, "user")
	})
}

func TestCombination(t *testing.T) {
	queryBuilder := NewQueryBuilder()
	Convey("SELECT .. FROM ... WHERE ... DIR", t, func() {
		Convey("Using string", func() {
			q := queryBuilder.From("custom_tablename").Select("id").Desc("email").Where("`id` = 1")
			So(q.Raw(), ShouldEqual, "SELECT `id` FROM custom_tablename WHERE `id` = 1 ORDER BY `email` DESC;")
			So(q.Error(), ShouldBeNil)
		})
		Convey("Using struct", func() {
			q := queryBuilder.Select(&User{}).Desc("email").Where(&User{ID: 1})
			So(q.Raw(), ShouldEqual, "SELECT `id`, `username`, `email`, `password` FROM custom_tablename WHERE `id` = 1 ORDER BY `email` DESC;")
			So(q.Error(), ShouldBeNil)
		})
	})
}
