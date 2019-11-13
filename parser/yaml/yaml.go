package yaml

import (
	"bytes"
	"fmt"
	"runtime"

	"github.com/ghodss/yaml"
)

type Parser struct{}

func (yp *Parser) Unmarshal(p []byte, v interface{}) error {
	subDocuments := separateSubDocuments(p)
	if len(subDocuments) > 1 {
		if err := unmarshalMultipleDocuments(subDocuments, v); err != nil {
			return fmt.Errorf("unmarshal multiple documents: %w", err)
		}

		return nil
	}

	if err := yaml.Unmarshal(p, v); err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}

	return nil
}

func separateSubDocuments(data []byte) [][]byte {
	linebreak := "\n"
	windowsLineEnding := bytes.Contains(data, []byte("\r\n"))
	if windowsLineEnding && runtime.GOOS == "windows" {
		linebreak = "\r\n"
	}

	return bytes.Split(data, []byte(linebreak+"---"+linebreak))
}

func unmarshalMultipleDocuments(subDocuments [][]byte, v interface{}) error {
	var documentStore []interface{}
	for _, subDocument := range subDocuments {
		var documentObject interface{}
		if err := yaml.Unmarshal(subDocument, &documentObject); err != nil {
			return fmt.Errorf("unmarshal subdocument yaml: %w", err)
		}

		documentStore = append(documentStore, documentObject)
	}

	yamlConfigBytes, err := yaml.Marshal(documentStore)
	if err != nil {
		return fmt.Errorf("marshal yaml document: %w", err)
	}

	if err := yaml.Unmarshal(yamlConfigBytes, v); err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}

	return nil
}