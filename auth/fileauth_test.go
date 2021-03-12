package auth

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/auxten/letmein/utils"
)

func TestAuth(t *testing.T) {
	up := &Auth{}
	Convey("load yaml", t, func() {
		err := up.LoadUserPass(utils.FJ(utils.GetProjectSrcDir(), "./test/userpass.yaml"))
		So(err, ShouldBeNil)
	})
	Convey("auth", t, func() {
		So(up.IsPermitted("auxten", "123456"), ShouldBeTrue)
		So(up.IsPermitted("auxten", "1236"), ShouldBeFalse)
	})
}
