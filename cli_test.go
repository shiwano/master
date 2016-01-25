package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"testing"
)

func TestCli(t *testing.T) {
	Convey("Cli", t, func() {
		cli := &Cli{
			outputDir: "./.tmp",
			schemaDir: "./.tmp",
			encoding:  "auto",
			silent:    true,
		}

		Convey("#run", func() {
			os.MkdirAll("./.tmp", 0777)
			data, _ := ioutil.ReadFile("./fixtures/masterdata.csv")
			ioutil.WriteFile("./.tmp/masterdata.csv", data, 0777)

			cli.file = "./.tmp/masterdata.csv"

			Convey("should output JSON files", func() {
				cli.run()
				actual, err := ioutil.ReadFile("./.tmp/masterdata.json")
				So(err, ShouldBeNil)
				expected, err := ioutil.ReadFile("./fixtures/masterdata.json")
				So(err, ShouldBeNil)
				So(actual, ShouldResemble, expected)
			})

			Convey("with fixEncoding option", func() {
				cli.fixEncoding = true

				Convey("should fix CSV file encoding", func() {
					cli.run()
					actual, err := ioutil.ReadFile("./.tmp/masterdata.csv")
					So(err, ShouldBeNil)
					expected, err := ioutil.ReadFile("./fixtures/masterdata-utf-8-bom.csv")
					So(err, ShouldBeNil)
					So(actual, ShouldResemble, expected)
				})

				Convey("should add BOM to UTF-8 CSV", func() {
					data, _ := ioutil.ReadFile("./fixtures/masterdata-utf-8.csv")
					ioutil.WriteFile("./.tmp/masterdata-utf-8.csv", data, 0777)
					cli.file = "./.tmp/masterdata-utf-8.csv"

					cli.run()
					actual, err := ioutil.ReadFile("./.tmp/masterdata-utf-8.csv")
					So(err, ShouldBeNil)
					expected, err := ioutil.ReadFile("./fixtures/masterdata-utf-8-bom.csv")
					So(err, ShouldBeNil)
					So(actual, ShouldResemble, expected)
				})
			})

			Convey("with outputSchema option", func() {
				cli.outputSchema = true

				Convey("should output JSON Schema files", func() {
					cli.run()
					actual, err := ioutil.ReadFile("./.tmp/masterdata.schema.json")
					So(err, ShouldBeNil)
					expected, err := ioutil.ReadFile("./fixtures/masterdata.schema.json")
					So(err, ShouldBeNil)
					So(actual, ShouldResemble, expected)
				})
			})

			Convey("with noOutputFile option", func() {
				cli.noOutputFile = true

				Convey("should not output JSON files", func() {
					cli.run()
					_, err := ioutil.ReadFile("./.tmp/masterdata.json")
					So(err, ShouldNotBeNil)
				})
			})

			Reset(func() {
				os.RemoveAll("./.tmp")
			})
		})

		Convey("#masterDataList", func() {
			Convey("should return master data list", func() {
				cli.file = "./fixtures/masterdata.csv"

				actual := cli.masterDataList()[0]
				So(actual.json(), ShouldContainSubstring, "ムーミン")
			})
		})

		Convey("#csvFilePaths", func() {
			Convey("should return target csv file paths", func() {
				cli.dir = "./fixtures"

				actual := cli.csvFilePaths()
				So(actual, ShouldResemble, []string{
					"fixtures/masterdata-utf-8-bom.csv",
					"fixtures/masterdata-utf-8.csv",
					"fixtures/masterdata.csv",
				})
			})
		})
	})
}
