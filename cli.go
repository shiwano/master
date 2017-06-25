package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ttacon/chalk"
)

// Cli represents the master command.
type Cli struct {
	dir            string
	file           string
	outputDir      string
	schemaDir      string
	encoding       string
	fixEncoding    bool
	noOutputFile   bool
	outputSchema   bool
	skipValidation bool
	noSchemaSuffix bool
	silent         bool
}

func (c *Cli) run() {
	if c.fixEncoding {
		c.fixCSVEncoding()
	}
	if !c.noOutputFile {
		c.makeOutputDirs()
	}

	for _, masterData := range c.masterDataList() {
		if c.outputSchema {
			jsonSchemaPath := filepath.Join(c.schemaDir,
				strings.Replace(masterData.fileName, ".json", ".schema.json", 1))
			c.writeFile("Generated", jsonSchemaPath, []byte(masterData.jsonSchema()))
		}

		jsonText := masterData.json()

		if !c.skipValidation {
			c.validateJSON(masterData.fileName, jsonText)
		}
		if !c.noOutputFile {
			jsonPath := filepath.Join(c.outputDir, masterData.fileName)
			c.writeFile("Generated", jsonPath, []byte(jsonText))
		} else if c.hasSingleCSVFile() {
			c.log(jsonText)
		}
	}
}

func (c *Cli) log(args ...interface{}) {
	if !c.silent {
		fmt.Println(args...)
	}
}

func (c *Cli) validateJSON(fileName string, jsonText string) {
	var schemaPath string
	if c.noSchemaSuffix {
		schemaPath = filepath.Join(c.schemaDir, fileName)
	} else {
		schemaPath = filepath.Join(c.schemaDir, strings.Replace(fileName, ".json", ".schema.json", 1))
	}

	schemaData, err := ioutil.ReadFile(schemaPath)
	if err == nil {
		schemaText := string(schemaData)
		if err := validateJSON(jsonText, schemaText); err != nil {
			fatalf("Failed to validate generated JSON: %v\n%v", fileName, err)
		}
	}
}

func (c *Cli) makeOutputDirs() {
	if err := os.MkdirAll(c.outputDir, 0777); err != nil {
		fatalf("Failed to make output directories\n%v", err)
	}
	if err := os.MkdirAll(c.schemaDir, 0777); err != nil {
		fatalf("Failed to make schema directories\n%v", err)
	}
}

func (c *Cli) writeFile(message string, path string, data []byte) {
	if err := ioutil.WriteFile(path, data, 0777); err != nil {
		fatalf("Failed to write a file\n%v", err)
	}
	c.log(message, chalk.Cyan.Color(path))
}

func (c *Cli) readFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fatalf("Failed to read a file: %v\n%v", path, err)
	}
	return data
}

func (c *Cli) detectEncoding(path string, data []byte) string {
	encoding, err := detectEncoding(data)
	if err != nil {
		fatalf("Failed to detect file encoding: %v\n%v", path, err)
	}
	return encoding
}

func (c *Cli) decode(path string, encoding string, data []byte) []byte {
	decoded, err := decode(data, encoding)
	if err != nil {
		fatalf("Failed to decode data: %v\n%v", path, err)
	}
	return decoded
}

func (c *Cli) fixCSVEncoding() {
	for _, filePath := range c.csvFilePaths() {
		data := c.readFile(filePath)
		detected := c.detectEncoding(filePath, data)

		encoding := c.encoding
		if encoding == "auto" {
			encoding = "UTF-8"
		}

		if detected != encoding || detected == "UTF-8" {
			decoded := c.decode(filePath, detected, data)
			encoded, err := encode(decoded, encoding)
			if err != nil {
				fatalf("Failed to encode CSV data: %v\n%v", filePath, err)
			}
			c.writeFile("Fixed file encoding of", filePath, encoded)
		}
	}
}

func (c *Cli) masterDataList() []*MasterData {
	filePaths := c.csvFilePaths()
	result := make([]*MasterData, len(filePaths))

	for i, filePath := range filePaths {
		data := c.readFile(filePath)
		encoding := c.encoding
		if encoding == "auto" {
			encoding = c.detectEncoding(filePath, data)
		}
		decoded := c.decode(filePath, encoding, data)

		csvTable, err := newCSVTable(filePath, encoding, decoded)
		if err != nil {
			fatalf("Failed to parse CSV data: %v\n%v", filePath, err)
		}

		masterData, err := newMasterDataFromCSV(csvTable, 2)
		if err != nil {
			fatalf("Failed to convert master data from CSV data: %v\n%v", csvTable.fileName, err)
		}
		result[i] = masterData
	}
	return result
}

func (c *Cli) csvFilePaths() []string {
	if c.hasSingleCSVFile() {
		return []string{c.file}
	}
	filePaths, err := filepath.Glob(filepath.Join(c.dir, "*.csv"))
	if err != nil {
		fatalf("Failed to find CSV paths: %v\n%v", c.dir, err)
	}
	return filePaths
}

func (c *Cli) hasSingleCSVFile() bool {
	return c.file != "" && strings.HasSuffix(c.file, ".csv")
}
