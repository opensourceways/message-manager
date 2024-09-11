package postgresql

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

const (
	MaxConn = 500
	MaxIdle = 250
	Life    = 2
)

// TestSetDefault the unit test for the function
func TestSetDefault(t *testing.T) {
	convey.Convey("test SetDefault success", t, func() {
		testCfg := &Config{}
		testCfg.SetDefault()
		convey.So(testCfg.MaxConn, convey.ShouldEqual, MaxConn)
		convey.So(testCfg.MaxIdle, convey.ShouldEqual, MaxIdle)
		convey.So(testCfg.Life, convey.ShouldEqual, Life)
	})
}
