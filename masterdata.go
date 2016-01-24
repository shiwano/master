package main

import (
	"github.com/jeffail/gabs"
	"path/filepath"
	"sort"
	"strings"
)

// MasterData represents structured master data which is converted from CSV file.
type MasterData struct {
	fileName  string
	indent    string
	container *gabs.Container
}

func newMasterData(path string, jsonText string, indent int) (*MasterData, error) {
	container, err := gabs.ParseJSON([]byte(jsonText))
	if err != nil {
		return nil, err
	}
	masterData := &MasterData{
		fileName:  filepath.Base(path),
		indent:    strings.Repeat(" ", indent),
		container: container,
	}
	return masterData, nil
}

func newMasterDataFromCSV(csvTable *CSVTable, indent int) (*MasterData, error) {
	data, err := csvTable.data()
	if err != nil {
		return nil, err
	}

	container, err := gabs.Consume(data)
	if err != nil {
		return nil, err
	}

	masterData := &MasterData{
		fileName:  strings.Replace(csvTable.fileName, ".csv", ".json", 1),
		indent:    strings.Repeat(" ", indent),
		container: container,
	}
	return masterData, nil
}

func (m *MasterData) json() string {
	if m.indent == "" {
		return m.container.String()
	}
	return m.container.StringIndent("", m.indent)
}

func (m *MasterData) jsonSchema() string {
	schema := getJSONSchemaRecursively(m.container.Data())
	schema.Set(m.fileName, "title")
	schema.Set("http://json-schema.org/draft-04/schema#", "$schema")
	if m.indent == "" {
		return schema.String()
	}
	return schema.StringIndent("", m.indent)
}

func getJSONSchemaRecursively(obj interface{}) *gabs.Container {
	schema := gabs.New()

	switch obj.(type) {
	case []map[string]interface{}:
		objAsArray := obj.([]map[string]interface{})
		if len(objAsArray) > 0 {
			schema.Set("array", "type")
			schema.Set(getJSONSchemaRecursively(objAsArray[0]).Data(), "items")
		}
	case []interface{}:
		objAsArray := obj.([]interface{})
		if len(objAsArray) > 0 {
			schema.Set("array", "type")
			schema.Set(getJSONSchemaRecursively(objAsArray[0]).Data(), "items")
		}
	case map[string]interface{}:
		objAsMap := obj.(map[string]interface{})
		schema.Set("object", "type")
		var keys []string
		for key, v := range objAsMap {
			schema.SetP(getJSONSchemaRecursively(v).Data(), "properties."+key)
			keys = append(keys, key)
		}
		sort.Strings(keys)
		schema.Set(false, "additionalProperties")
		schema.Set(keys, "required")
	case string:
		schema.Set("string", "type")
	case bool:
		schema.Set("boolean", "type")
	default:
		schema.Set("number", "type")
	}
	return schema
}
