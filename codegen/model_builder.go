package codegen

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/stoewer/go-strcase"
	"text/template"
)

//go:embed model.tmpl
var templateModel string

type ModelBuilder struct {
	PackageName string

	Imports Imports

	ModelName      string
	CollectionName string

	Fields []Field
}

func ModelBuilderFromSchema(pkgName string, collectionName string, s *schema.Schema) (*ModelBuilder, error) {
	mb := &ModelBuilder{
		PackageName:    strcase.SnakeCase(pkgName),
		ModelName:      strcase.UpperCamelCase(collectionName),
		CollectionName: collectionName,
		Imports:        Imports{},
	}

	mb.Imports.addImport("github.com/pocketbase/pocketbase/models", "")

	for _, field := range s.Fields() {

		mbField := Field{
			FieldName:    field.Name,
			FunctionName: strcase.UpperCamelCase(field.Name),
		}

		if err := resolveSchemaField(mb, field, &mbField); err != nil {
			return nil, err
		}

		mb.Fields = append(mb.Fields, mbField)
	}

	return mb, nil
}

func (mb *ModelBuilder) Gen(t *template.Template) (string, error) {
	buffer := bytes.NewBuffer(nil)
	err := t.Execute(buffer, mb)

	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buffer.String(), nil
}

type Imports map[string]Import

func (i Imports) addImport(importPath string, alias string) {
	i[importPath] = Import{
		Alias: alias,
		Path:  importPath,
	}
}

func (i Imports) List() []Import {
	var imports []Import
	for _, v := range i {
		imports = append(imports, v)
	}

	return imports
}

type Import struct {
	Alias string
	Path  string
}

func (i Import) ToImportPath() string {
	if i.Alias != "" {
		return fmt.Sprintf(`%s "%s"`, i.Alias, i.Path)
	} else {
		return fmt.Sprintf(`"%s"`, i.Path)
	}
}

type Field struct {
	FieldName    string
	FunctionName string

	GoType string

	GetterComment string
	GetterCall    string
	HasGetterCast bool

	SetterComment string
}

func resolveSchemaField(builder *ModelBuilder, field *schema.SchemaField, f *Field) error {
	f.GetterComment = fmt.Sprintf(`// Get%s returns the value of the "%s" field`, f.FunctionName, field.Name)
	f.SetterComment = fmt.Sprintf(`// Set%s sets the value of the "%s" field`, f.FunctionName, field.Name)

	switch field.Type {
	case schema.FieldTypeText:
		f.GoType = "string"
		f.GetterCall = `GetString`

	case schema.FieldTypeEditor:
		f.GoType = "string"
		f.GetterCall = `GetString`
		f.GetterComment = fmt.Sprintf(`// Get%s returns the value of the "%s" field as HTML`, f.FunctionName,
			field.Name)
		f.SetterComment = fmt.Sprintf(`// Set%s sets the value of the "%s" field as HTML`, f.FunctionName,
			field.Name)

	case schema.FieldTypeUrl:
		f.GoType = "string"
		f.GetterCall = `GetString`
		f.GetterComment = fmt.Sprintf(`// Get%s returns the value of the "%s" field as URL`, f.FunctionName, field.Name)
		f.SetterComment = fmt.Sprintf(`// Set%s sets the value of the "%s" field as URL`, f.FunctionName, field.Name)

	case schema.FieldTypeEmail:
		f.GoType = "string"
		f.GetterCall = `GetString`
		f.GetterComment = fmt.Sprintf(`// Get%s returns the value of the "%s" field as email`, f.FunctionName,
			field.Name)
		f.SetterComment = fmt.Sprintf(`// Set%s sets the value of the "%s" field as email`, f.FunctionName, field.Name)
	case schema.FieldTypeFile:
		option := field.Options.(*schema.FileOptions)
		f.GetterComment = fmt.Sprintf(`// Get%s returns the value of the "%s" field as file`, f.FunctionName,
			field.Name)
		f.SetterComment = fmt.Sprintf(`// Set%s sets the value of the "%s" field as file`, f.FunctionName, field.Name)
		if option.IsMultiple() {
			f.GoType = "[]string"
			f.GetterCall = `GetStringSlice`
		} else {
			f.GoType = "string"
			f.GetterCall = `GetString`
		}

	case schema.FieldTypeNumber:
		f.GoType = "int"
		f.GetterCall = `GetInt`

	case schema.FieldTypeBool:
		f.GoType = "bool"
		f.GetterCall = `GetBool`

	case schema.FieldTypeDate:
		builder.Imports.addImport("github.com/pocketbase/pocketbase/tools/types", "")
		f.GoType = "types.DateTime"
		f.GetterCall = `GetDateTime`

	case schema.FieldTypeSelect:
		option := field.Options.(*schema.SelectOptions)
		f.GetterComment += fmt.Sprintf("\n// Possible values: %s", option.Values)
		f.SetterComment += fmt.Sprintf("\n// Possible values: %s", option.Values)
		if option.IsMultiple() {
			f.GoType = "[]string"
			f.GetterCall = `GetStringSlice`
		} else {
			f.GoType = "string"
			f.GetterCall = `GetString`
		}

	case schema.FieldTypeJson:
		builder.Imports.addImport("github.com/pocketbase/pocketbase/tools/types", "")
		f.GoType = "types.JsonRaw"
		f.GetterCall = `Get`
		f.HasGetterCast = true

	case schema.FieldTypeRelation:
		option := field.Options.(*schema.RelationOptions)
		f.GetterComment += fmt.Sprintf("\n// Relation collection related : %s",
			option.CollectionId /* TODO: resolve it name */)
		f.SetterComment += fmt.Sprintf("\n// Relation collection related : %s",
			option.CollectionId /* TODO: resolve it name */)
		if option.IsMultiple() {
			f.GoType = "[]string"
			f.GetterCall = `GetStringSlice`
		} else {
			f.GoType = "string"
			f.GetterCall = `GetString`
		}

	default:
		return fmt.Errorf("unknown field type %s", field.Type)
	}

	return nil
}
