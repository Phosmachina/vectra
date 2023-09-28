package main

import (
	"os"
	"strings"
	"text/template"
)

type TemplateData struct {
	PackageName string
	Types       []I18nType
}

type I18nType struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name        string
	IsDirectory bool
	Key         string
}

type folder struct {
	name   string
	key    string
	parent *folder
	items  map[string]*folder
}

func newFolder(name string, parent *folder) *folder {
	f := folder{name: name, items: map[string]*folder{},
		parent: parent}
	if parent == nil || parent.key == "" {
		f.key = name
	} else {
		f.key = parent.key + "." + name
	}
	return &f
}

func (f *folder) add(keys []string) {
	if len(keys) == 0 {
		return
	} else if len(keys) == 1 {
		if f.items == nil {
			f.items = map[string]*folder{}
		}
		f.items[keys[0]] = newFolder(keys[0], f)
		return
	} else {
		if f.items[keys[0]] == nil {
			f.items[keys[0]] = newFolder(keys[0], f)
		}
		f.items[keys[0]].add(keys[1:])
	}
}

func generateGoCodeToFile(data map[string]string, outputFilePath string) error {

	// Prepare the data for the template

	var root = newFolder("", nil)
	for k := range data {
		root.add(strings.Split(k, "."))
	}

	types := buildDataTemplate(root)

	// Define the template for code generation
	file, _ := os.ReadFile("i18n_to_Go.gohtml")
	tmpl := string(file)

	// Create a template data struct
	templateData := TemplateData{
		PackageName: "i18n",
		Types:       types,
	}

	// Create the output file
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Parse and execute the template
	tmplParsed, err := template.New("useless_name").Funcs(template.FuncMap{
		"Upper": Upper,
	}).Parse(tmpl)
	if err != nil {
		return err
	}

	err = tmplParsed.Execute(outputFile, templateData)
	if err != nil {
		return err
	}

	return nil
}

func buildDataTemplate(root *folder) []I18nType {
	var types []I18nType

	for _, f := range root.items {
		if len(f.items) == 0 {
			continue
		}

		var fields []Field
		for _, f := range f.items {
			fields = append(fields, Field{
				Name:        f.name,
				IsDirectory: len(f.items) > 0,
				Key:         f.key,
			})
		}
		types = append(types, I18nType{
			Name:   f.name,
			Fields: fields,
		})

		var endFields []Field
		for _, field := range fields {
			if !field.IsDirectory {
				endFields = append(endFields, field)
			}
		}

		types = append(types, buildDataTemplate(f)...)
	}

	return types
}

func Upper(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(string(str[0])) + str[1:]
}
