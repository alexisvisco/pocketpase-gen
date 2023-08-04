package codegen

import (
	"fmt"
	"text/template"
)

func CreateModelTemplateParser() (*template.Template, error) {
	t, err := template.New("model").Parse(templateModel)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return t, nil
}
