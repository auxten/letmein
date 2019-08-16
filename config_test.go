package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfig_LoadConfig(t *testing.T) {
	conf := &Config{}
	Convey("load config", t, func() {
		err := conf.LoadConfig("./config.yaml")
		So(err, ShouldBeNil)
		So(conf.Auth, ShouldNotBeNil)
		So(conf.AwsSg, ShouldNotBeNil)
	})
}
