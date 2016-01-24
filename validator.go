package main

import (
	"errors"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
)

func validateJSON(jsonText string, schemaText string) error {
	schemaLoader := gojsonschema.NewStringLoader(schemaText)
	docLoader := gojsonschema.NewStringLoader(jsonText)

	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		errMessage := "The JSON data is not valid:\n"
		for _, desc := range result.Errors() {
			errMessage += fmt.Sprintf("  %v\n", desc)
		}
		return errors.New(errMessage)
	}
	return nil
}
