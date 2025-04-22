package utils

import (
	"reflect"
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
	Convey("Given a pointer to struct, should return true", t, func() {
		PK := uint32(1)
		user := &User{}
		So(IsPointerToStruct(user), ShouldBeTrue)

		user = &User{ID: PK, Username: "john_doe"}
		So(IsPointerToStruct(user), ShouldBeTrue)
	})
	Convey("Given a non struct, should return false", t, func() {
		So(IsPointerToStruct(nil), ShouldBeFalse)
		So(IsPointerToStruct("something"), ShouldBeFalse)

		var user interface{}
		So(IsPointerToStruct(user), ShouldBeFalse)
		So(IsPointerToStruct(42), ShouldBeFalse)
	})
	Convey("Given struct, should return false", t, func() {
		type Example struct{}
		ex := Example{}
		So(IsPointerToStruct(ex), ShouldBeFalse)
	})
	Convey("Given pointer to nil value, should return false", t, func() {
		var p *int = nil
		So(IsPointerToStruct(p), ShouldBeFalse)
	})
}

func TestIsString(t *testing.T) {
	Convey("Given string, should return true", t, func() {
		So(IsString("someString"), ShouldBeTrue)
	})
	Convey("Given a pointer toa  string, should return true", t, func() {
		s := "string"
		So(IsString(&s), ShouldBeTrue)
	})
	Convey("Given a non string, should return false", t, func() {
		So(IsString(nil), ShouldBeFalse)
		var user interface{}
		So(IsString(user), ShouldBeFalse)
		So(IsString(42), ShouldBeFalse)
	})
	Convey("Given pointer to nil value, should return false", t, func() {
		var p *string = nil
		So(IsString(p), ShouldBeFalse)
	})
}

type Sample struct {
	WithColumnTag  string `gorm:"column:custom_name"`
	WithOtherTags  string `gorm:"type:varchar(100);not null"`
	WithMixedTags  string `gorm:"type:varchar(100);column:mixed_name;not null"`
	WithoutGormTag string
	EmptyGormTag   string `gorm:""`
	CamelCaseField string
	AnotherExample string
}

func TestGetColumnName(t *testing.T) {
	Convey("Given a struct with various gorm tags", t, func() {
		sampleType := reflect.TypeOf(Sample{})

		tests := map[string]string{
			"WithColumnTag":  "custom_name",
			"WithOtherTags":  "with_other_tags",
			"WithMixedTags":  "mixed_name",
			"WithoutGormTag": "without_gorm_tag",
			"EmptyGormTag":   "empty_gorm_tag",
			"CamelCaseField": "camel_case_field",
			"AnotherExample": "another_example",
		}

		for fieldName, expected := range tests {
			Convey("When GetColumnName is called on "+fieldName, func() {
				field, found := sampleType.FieldByName(fieldName)
				So(found, ShouldBeTrue)

				actual := GetColumnName(field)
				So(actual, ShouldEqual, expected)
			})
		}
	})
}

func TestEscapeString(t *testing.T) {
	Convey("Given quoted string, should double the quote", t, func() {
		So(EscapeString("Ar'a"), ShouldEqual, "Ar''a")
	})
}

func TestStructName(t *testing.T) {
	Convey("Get struct name", t, func() {
		So(GetStructName(Sample{}), ShouldEqual, "samples")
		So(GetStructName(&Sample{}), ShouldEqual, "samples")
	})
}

func TestSanitizedIdentifier(t *testing.T) {
	Convey("sanitize", t, func() {
		res, err := SanitizeIdentifier("aA19_")
		So(res, ShouldNotBeEmpty)
		So(err, ShouldBeNil)
		res, err = SanitizeIdentifier("aA19_'")
		So(res, ShouldBeEmpty)
		So(err, ShouldNotBeNil)
		res, err = SanitizeIdentifier("aA19_#")
		So(res, ShouldBeEmpty)
		So(err, ShouldNotBeNil)
	})
}
