package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConf(t *testing.T) {
	c := Config{}
	Convey("Can load config_test.json without errors", t, func() {
		err := c.Load("test_configs/config_test.json")
		So(err, ShouldBeNil)
	})
	Convey("Gives error if config file does not exist", t, func() {
		err := c.Load("config_blabla.json")
		So(err, ShouldNotBeNil)
	})
	Convey("Gives error if config file is not valid json", t, func() {
		err := c.Load("test_configs/config_fail.json")
		So(err, ShouldNotBeNil)
	})
	Convey("Gives error if config file is not valid", t, func() {
		err := c.Load("test_configs/config_fail2.json")
		So(err, ShouldNotBeNil)
	})
}
