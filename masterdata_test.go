package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMasterData(t *testing.T) {
	Convey("MasterData", t, func() {
		Convey("#json", func() {
			Convey("should return the JSON string", func() {
				masterData, _ := newMasterData("foo.json", `[{"str":"foo"},{"str":"bar"}]`, 0)
				So(masterData.json(), ShouldEqual, `[{"str":"foo"},{"str":"bar"}]`)
			})
		})

		Convey("#jsonSchema", func() {
			Convey("should return the valid JSON Schema string", func() {
				jsonText := `[
				{ "str": "foo", "number": 1, "bool": true },
				{ "str": "bar", "number": 2, "bool": false }
				]`
				masterData, _ := newMasterData("foo.json", jsonText, 0)
				So(validateJSON(jsonText, masterData.jsonSchema()), ShouldBeNil)
			})
		})
	})
}
