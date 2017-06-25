package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCSV(t *testing.T) {
	Convey("csv", t, func() {
		Convey(".newCSVColumns", func() {
			Convey("should return new CSVColumns", func() {
				csvRecords := [][]string{
					[]string{"str", "num", "mixed", "bool"},
					[]string{"foo", "1", "1", "TRUE"},
					[]string{"bar", "2.01", "qux", "FALSE"},
					[]string{"baz", "3", "3", ""},
				}

				actual, err := newCSVColumns(csvRecords)
				So(err, ShouldBeNil)
				So(actual, ShouldResemble, []*CSVColumn{
					&CSVColumn{index: 0, name: "str", isString: true, isBool: false},
					&CSVColumn{index: 1, name: "num", isString: false, isBool: false},
					&CSVColumn{index: 2, name: "mixed", isString: true, isBool: false},
					&CSVColumn{index: 3, name: "bool", isString: false, isBool: true},
				})
			})
		})

		Convey(".newCSVTable", func() {
			Convey("with valid data", func() {
				csvData := []byte("str,num,bool\nfoo,1,TRUE\nbar,2,FALSE")

				Convey("should return a new CSVTable", func() {
					actual, err := newCSVTable("foo/test.csv", "utf-8", csvData)
					So(err, ShouldBeNil)
					So(actual, ShouldResemble, &CSVTable{
						fileName: "test.csv",
						encoding: "utf-8",
						columns: []*CSVColumn{
							&CSVColumn{index: 0, name: "str", isString: true, isBool: false},
							&CSVColumn{index: 1, name: "num", isString: false, isBool: false},
							&CSVColumn{index: 2, name: "bool", isString: false, isBool: true},
						},
						rows: [][]interface{}{
							[]interface{}{"foo", 1.0, true},
							[]interface{}{"bar", 2.0, false},
						},
					})
				})
			})

			Convey("with one record data", func() {
				csvData := []byte("str,num")

				Convey("should return a error", func() {
					actual, err := newCSVTable("test.csv", "utf-8", csvData)
					So(err, ShouldNotBeNil)
					So(actual, ShouldBeNil)
				})
			})

			Convey("with empty data", func() {
				csvData := []byte("")

				Convey("should return a error", func() {
					actual, err := newCSVTable("test.csv", "utf-8", csvData)
					So(err, ShouldNotBeNil)
					So(actual, ShouldBeNil)
				})
			})

			Convey("with bumpy data", func() {
				csvData := []byte("str,num\na")

				Convey("should return a error", func() {
					actual, err := newCSVTable("test.csv", "utf-8", csvData)
					So(err, ShouldNotBeNil)
					So(actual, ShouldBeNil)
				})
			})
		})
	})
}

func TestCSVColumn(t *testing.T) {
	Convey("CSVColumn", t, func() {
		Convey("#validate", func() {
			Convey("with valid column name", func() {
				column := &CSVColumn{name: "foo"}
				So(column.validate(), ShouldBeNil)
			})

			Convey("with invalid column name", func() {
				column := &CSVColumn{name: "0.foo"}
				So(column.validate(), ShouldNotBeNil)
			})
		})
	})
}

func TestCSVTable(t *testing.T) {
	Convey("CSVTable", t, func() {
		Convey("#data", func() {
			Convey("with normal key-value data", func() {
				csvData := []byte("str,num,bool\nfoo,1,TRUE\nbar,2,FALSE")
				csvTable, _ := newCSVTable("test.csv", "utf-8", csvData)

				Convey("should return map data", func() {
					actual, err := csvTable.data()
					So(err, ShouldBeNil)
					So(actual, ShouldResemble, []map[string]interface{}{
						map[string]interface{}{"str": "foo", "num": 1.0, "bool": true},
						map[string]interface{}{"str": "bar", "num": 2.0, "bool": false},
					})
				})
			})

			Convey("with structured data", func() {
				csvData := []byte("obj.foo,obj.bar\na,b\nc,d")
				csvTable, _ := newCSVTable("test.csv", "utf-8", csvData)

				Convey("should return map data", func() {
					actual, err := csvTable.data()
					So(err, ShouldBeNil)
					So(actual, ShouldResemble, []map[string]interface{}{
						map[string]interface{}{"obj": map[string]interface{}{"foo": "a", "bar": "b"}},
						map[string]interface{}{"obj": map[string]interface{}{"foo": "c", "bar": "d"}},
					})
				})
			})

			Convey("with array data which has some missings", func() {
				csvData := []byte("items.0,items.2\na,b\nc,d")
				csvTable, _ := newCSVTable("test.csv", "utf-8", csvData)

				Convey("should return map data", func() {
					actual, err := csvTable.data()
					So(err, ShouldBeNil)
					So(actual, ShouldResemble, []map[string]interface{}{
						map[string]interface{}{"items": []interface{}{"a", "b"}},
						map[string]interface{}{"items": []interface{}{"c", "d"}},
					})
				})
			})

			Convey("with array data", func() {
				csvData := []byte("items.0,items.1,items.2\na,b,c\nd,e,f")
				csvTable, _ := newCSVTable("test.csv", "utf-8", csvData)

				Convey("should return map data", func() {
					actual, err := csvTable.data()
					So(err, ShouldBeNil)
					So(actual, ShouldResemble, []map[string]interface{}{
						map[string]interface{}{"items": []interface{}{"a", "b", "c"}},
						map[string]interface{}{"items": []interface{}{"d", "e", "f"}},
					})
				})
			})

			Convey("with two-dimentional array data", func() {
				csvData := []byte("items.0.0,items.0.1\na,b\nc,d")
				csvTable, _ := newCSVTable("test.csv", "utf-8", csvData)

				Convey("should return map data", func() {
					actual, err := csvTable.data()
					So(err, ShouldBeNil)
					So(actual, ShouldResemble, []map[string]interface{}{
						map[string]interface{}{"items": []interface{}{[]interface{}{"a", "b"}}},
						map[string]interface{}{"items": []interface{}{[]interface{}{"c", "d"}}},
					})
				})
			})

			Convey("with structured array data", func() {
				csvData := []byte("items.0.foo,items.0.bar,items.1.foo,items.1.bar\na,b,c,d\ne,f,g,h")
				csvTable, _ := newCSVTable("test.csv", "utf-8", csvData)

				Convey("should return map data", func() {
					actual, err := csvTable.data()
					So(err, ShouldBeNil)
					So(actual, ShouldResemble, []map[string]interface{}{
						map[string]interface{}{"items": []interface{}{
							map[string]interface{}{"foo": "a", "bar": "b"},
							map[string]interface{}{"foo": "c", "bar": "d"},
						}},
						map[string]interface{}{"items": []interface{}{
							map[string]interface{}{"foo": "e", "bar": "f"},
							map[string]interface{}{"foo": "g", "bar": "h"},
						}},
					})
				})
			})

			Convey("with structured array data which has some missings", func() {
				csvData := []byte("items.0.foo,items.0.bar,items.1.foo,items.1.bar\na,b,c,d\n,,g,h")
				csvTable, _ := newCSVTable("test.csv", "utf-8", csvData)

				Convey("should return map data", func() {
					actual, err := csvTable.data()
					So(err, ShouldBeNil)
					So(actual, ShouldResemble, []map[string]interface{}{
						map[string]interface{}{"items": []interface{}{
							map[string]interface{}{"foo": "a", "bar": "b"},
							map[string]interface{}{"foo": "c", "bar": "d"},
						}},
						map[string]interface{}{"items": []interface{}{
							map[string]interface{}{"foo": "g", "bar": "h"},
						}},
					})
				})
			})

			Convey("with column name which is minus number", func() {
				csvData := []byte("items.-1,items.-2\na,b\nd,e")
				csvTable, _ := newCSVTable("test.csv", "utf-8", csvData)

				Convey("should return map data", func() {
					actual, err := csvTable.data()
					So(err, ShouldBeNil)
					So(actual, ShouldResemble, []map[string]interface{}{
						map[string]interface{}{"items": map[string]interface{}{"-1": "a", "-2": "b"}},
						map[string]interface{}{"items": map[string]interface{}{"-1": "d", "-2": "e"}},
					})
				})
			})
		})
	})
}
