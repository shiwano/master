package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	numberValuePattern = regexp.MustCompile("^[0-9]+\\.?[0-9]*$")
	csvColumnPattern   = regexp.MustCompile("(?i)^[^0-9.]+(\\.[0-9a-z]+)*$")
)

// CSVColumn represents a column of CSVTable.
type CSVColumn struct {
	index    int
	name     string
	isString bool
}

func newCSVColumns(records [][]string) ([]*CSVColumn, error) {
	columnLength := len(records[0])
	columns := make([]*CSVColumn, columnLength)
	for i, value := range records[0] {
		columns[i] = &CSVColumn{index: i, name: value}

		if err := columns[i].validate(); err != nil {
			return nil, err
		}
	}

	for _, record := range records[1:] {
		if len(record) != columnLength {
			return nil, fmt.Errorf("Record length is not enough: %v", record)
		}
		for i, value := range record {
			if value != "" && !numberValuePattern.MatchString(value) {
				columns[i].isString = true
			}
		}
	}
	return columns, nil
}

func (c *CSVColumn) validate() error {
	if !csvColumnPattern.MatchString(c.name) {
		return fmt.Errorf("Invalid column name: %v", c.name)
	}
	return nil
}

// CSVTable represents structured CSV data table.
type CSVTable struct {
	fileName string
	encoding string
	columns  []*CSVColumn
	rows     [][]interface{}
}

func newCSVTable(path string, encoding string, data []byte) (*CSVTable, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("CSV data should have 2 rows at a minimum: %v", path)
	}

	columns, err := newCSVColumns(records)
	if err != nil {
		return nil, err
	}

	rows := make([][]interface{}, len(records)-1)
	for recordIndex, record := range records[1:] {
		row := make([]interface{}, len(record))
		rows[recordIndex] = row

		for i, value := range record {
			strValue := fmt.Sprintf("%v", value)

			if columns[i].isString {
				row[i] = strValue
			} else if strValue == "" {
				row[i] = 0
			} else {
				floatValue, err := strconv.ParseFloat(strValue, 64)
				if err != nil {
					return nil, err
				}
				row[i] = floatValue
			}
		}
	}
	csvTable := &CSVTable{
		fileName: filepath.Base(path),
		encoding: encoding,
		columns:  columns,
		rows:     rows,
	}
	return csvTable, err
}

func (c *CSVTable) data() ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, len(c.rows))

	for rowIndex, row := range c.rows {
		root := make(map[string]interface{})
		result[rowIndex] = root

		for i, value := range row {
			column := c.columns[i]
			c.getMapData(root, strings.Split(column.name, "."), value)
		}
		c.removeEmptyArrayItemRecursively(root)
	}
	return result, nil
}

func (c *CSVTable) removeEmptyArrayItemRecursively(container map[string]interface{}) {
	for key, value := range container {
		if valueAsMap, ok := value.(map[string]interface{}); ok {
			c.removeEmptyArrayItemRecursively(valueAsMap)
		} else if valueAsArray, ok := value.([]interface{}); ok {
			array := []interface{}{}

			for _, arrayItem := range valueAsArray {
				if arrayItemAsMap, ok := arrayItem.(map[string]interface{}); ok {
					for _, v := range arrayItemAsMap {
						if v != 0 && v != "" {
							array = append(array, arrayItem)
							break
						}
					}
				} else if arrayItem != nil {
					array = append(array, arrayItem)
				}
			}
			container[key] = array
		}
	}
}

func (c *CSVTable) getMapData(container map[string]interface{},
	keys []string, value interface{}) interface{} {
	key := keys[0]
	var nextKey string
	if len(keys) >= 2 {
		nextKey = keys[1]
	}

	if len(keys) == 1 {
		container[key] = value
	} else if arrayIndex, err := strconv.Atoi(nextKey); err == nil {
		if container[key] == nil {
			container[key] = make([]interface{}, arrayIndex+1)
		}
		array := container[key].([]interface{})
		container[key] = c.getArrayData(array, arrayIndex, keys[2:], value)
	} else {
		if container[key] == nil {
			container[key] = make(map[string]interface{})
		}
		newContainer := container[key].(map[string]interface{})
		container[key] = newContainer
		c.getMapData(newContainer, keys[1:], value)
	}
	return container
}

func (c *CSVTable) getArrayData(array []interface{}, index int,
	keys []string, value interface{}) []interface{} {
	if len(array) <= index {
		array = append(array, make([]interface{}, index-len(array)+1)...)
	}

	if len(keys) == 0 {
		array[index] = value
		return array
	}

	if nextArrayIndex, err := strconv.Atoi(keys[0]); err == nil {
		if array[index] == nil {
			array[index] = make([]interface{}, nextArrayIndex+1)
		}
		nextArray := array[index].([]interface{})
		array[index] = c.getArrayData(nextArray, nextArrayIndex, keys[1:], value)
	} else if array[index] == nil {
		array[index] = c.getMapData(make(map[string]interface{}), keys, value)
	} else {
		newContainer := array[index].(map[string]interface{})
		array[index] = c.getMapData(newContainer, keys, value)
	}
	return array
}
