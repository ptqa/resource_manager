package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestResources(t *testing.T) {
	c := Config{3000, 2, 4}
	a := Resources{}
	a.Init(c)
	Convey("Empty valid list", t, func() {
		list := a.List()
		So(list, ShouldContainSubstring, "\"allocated\":{}")
	})

	Convey("Allocates first resource", t, func() {
		answer, err := a.tryAllocate("alice", c.Workers)
		So(err, ShouldBeNil)
		So(answer, ShouldEqual, 1)
	})

	Convey("Allocates second resource", t, func() {
		answer, err := a.tryAllocate("bob", c.Workers)
		So(err, ShouldBeNil)
		So(answer, ShouldEqual, 2)
	})

	Convey("Limit exists", t, func() {
		_, err := a.tryAllocate("limit", c.Workers)
		So(err, ShouldNotBeNil)
	})

	Convey("List has valid users", t, func() {
		list := a.List()
		So(list, ShouldContainSubstring, "alice")
		So(list, ShouldContainSubstring, "bob")
	})

	Convey("Reseting gives empty list", t, func() {
		a.Reset(c.Workers)
		list := a.List()
		So(list, ShouldNotContainSubstring, "alice")
		So(list, ShouldNotContainSubstring, "bob")
	})

	Convey("Deallocating invalid id gives error", t, func() {
		err := a.tryDeallocate(1, c.Workers)
		So(err, ShouldNotBeNil)
	})
	Convey("Deallocating removes user", t, func() {
		john, _ := a.tryAllocate("john", c.Workers)
		a.tryAllocate("ted", c.Workers)
		err := a.tryDeallocate(john, c.Workers)
		list := a.List()
		So(err, ShouldBeNil)
		So(list, ShouldNotContainSubstring, "john")
		So(list, ShouldContainSubstring, "ted")
	})
	Convey("Search works properly", t, func() {
		a.Reset(c.Workers)
		a.tryAllocate("bob", c.Workers)
		a.tryAllocate("alice", c.Workers)
		empty_search := a.Search("empty")
		bob_search := a.Search("bob")
		alice_search := a.Search("alice")
		So(empty_search, ShouldContainSubstring, "[]")
		So(bob_search, ShouldContainSubstring, "[\"r1\"]")
		So(alice_search, ShouldContainSubstring, "[\"r2\"]")
	})
	Convey("Search works properly for several resources", t, func() {
		a.Reset(c.Workers)
		a.tryAllocate("bob", c.Workers)
		a.tryAllocate("bob", c.Workers)
		bob_search := a.Search("bob")
		So(bob_search, ShouldContainSubstring, "[\"r1\",\"r2\"]")
	})
}
