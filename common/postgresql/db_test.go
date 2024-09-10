package postgresql

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/smartystreets/goconvey/convey"
	"gorm.io/gorm"
)

func TestInit(t *testing.T) {
	mockMethodUseMock := gomonkey.ApplyMethodReturn(&gorm.DB{}, "Use", nil)
	defer mockMethodUseMock.Reset()

	convey.Convey("test Init success", t, func() {
		mockFunc := gomonkey.ApplyFuncReturn(gorm.Open, &gorm.DB{}, nil)
		defer mockFunc.Reset()

		mockFunc2 := gomonkey.ApplyFuncReturn(os.Remove, nil)
		defer mockFunc2.Reset()

		mockMethod := gomonkey.ApplyMethodReturn(&gorm.DB{}, "DB", &sql.DB{}, nil)
		defer mockMethod.Reset()

		convey.So(Init(&Config{Dbcert: "test"}), convey.ShouldBeNil)
		convey.So(DB(), convey.ShouldNotBeNil)
	})

	convey.Convey("test Init failed, open config failed", t, func() {
		testErr := errors.New("open config error")
		mockFunc := gomonkey.ApplyFuncReturn(gorm.Open, &gorm.DB{}, testErr)
		defer mockFunc.Reset()

		convey.So(Init(&Config{}), convey.ShouldResemble, testErr)
	})

	convey.Convey("test Init failed, db error", t, func() {
		mockFunc := gomonkey.ApplyFuncReturn(gorm.Open, &gorm.DB{}, nil)
		defer mockFunc.Reset()

		testErr := errors.New("db error")
		mockMethod := gomonkey.ApplyMethodReturn(&gorm.DB{}, "DB", &sql.DB{}, testErr)
		defer mockMethod.Reset()

		convey.So(Init(&Config{}), convey.ShouldResemble, testErr)
	})
}
