package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"testing"
)

func TestEncoding(t *testing.T) {
	Convey("encoding", t, func() {
		Convey(".encode", func() {
			Convey("should encode the given bytes with the specified encoding", func() {
				data, _ := ioutil.ReadFile("./fixtures/utf-8.txt")
				actual, err := encode(data, "shift-jis")
				expected, _ := ioutil.ReadFile("./fixtures/shift-jis.txt")

				So(err, ShouldBeNil)
				So(string(actual), ShouldEqual, string(expected))
			})

			Convey("with UTF-8 charset", func() {
				Convey("should add BOM", func() {
					data, _ := ioutil.ReadFile("./fixtures/utf-8.txt")
					actual, err := encode(data, "utf-8")
					expected, _ := ioutil.ReadFile("./fixtures/utf-8-bom.txt")

					So(err, ShouldBeNil)
					So(string(actual), ShouldEqual, string(expected))
				})
			})
		})

		Convey(".decode", func() {
			Convey("should decode the given bytes with the specified encoding", func() {
				data, _ := ioutil.ReadFile("./fixtures/shift-jis.txt")
				actual, err := decode(data, "shift-jis")
				expected, _ := ioutil.ReadFile("./fixtures/utf-8.txt")

				So(err, ShouldBeNil)
				So(string(actual), ShouldEqual, string(expected))
			})

			Convey("with UTF-8 bytes", func() {
				Convey("should remove BOM", func() {
					data, _ := ioutil.ReadFile("./fixtures/utf-8-bom.txt")
					actual, err := decode(data, "utf-8")
					expected, _ := ioutil.ReadFile("./fixtures/utf-8.txt")

					So(err, ShouldBeNil)
					So(string(actual), ShouldEqual, string(expected))
				})
			})
		})

		Convey(".detectEncoding", func() {
			Convey("should detect the encoding automatically", func() {
				data, _ := ioutil.ReadFile("./fixtures/shift-jis.txt")
				encoding, err := detectEncoding(data)

				So(err, ShouldBeNil)
				So(encoding, ShouldEqual, "Shift_JIS")
			})
		})
	})
}
