package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestValidator(t *testing.T) {
	const schema = `{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"type": "array",
		"items": {
			"type": "object",
			"properties": {
				"id": { "type": "integer" }
			},
			"additionalProperties": false,
			"required": [
				"id"
			]
		}
	}`

	Convey("validator", t, func() {
		Convey(".validate", func() {
			Convey("with valid JSON text", func() {
				Convey("should return no error", func() {
					err := validateJSON(`[{"id": 1}, {"id": 2}]`, schema)
					So(err, ShouldBeNil)
				})
			})

			Convey("with invalid JSON text", func() {
				Convey("should return error", func() {
					err := validateJSON(`{"id": "foo"}`, schema)
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}
