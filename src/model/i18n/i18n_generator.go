package i18n

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

type Folder struct {
	Name    string
	FullKey string
	Parent  *Folder
	Items   map[string]*Folder
}

func newFolder(name string, parent *Folder) *Folder {
	f := Folder{Name: name, Items: map[string]*Folder{},
		Parent: parent}
	if parent == nil || parent.FullKey == "" {
		f.FullKey = name
	} else {
		f.FullKey = parent.FullKey + "." + name
	}
	return &f
}

func (f *Folder) add(keys []string) {
	if len(keys) == 0 {
		return
	} else if len(keys) == 1 {
		if f.Items == nil {
			f.Items = map[string]*Folder{}
		}
		f.Items[keys[0]] = newFolder(keys[0], f)
		return
	} else {
		if f.Items[keys[0]] == nil {
			f.Items[keys[0]] = newFolder(keys[0], f)
		}
		f.Items[keys[0]].add(keys[1:])
	}
}

func GenerateCode(data map[string]string) error {

	var root = newFolder("", nil)
	for k := range data {
		root.add(strings.Split(k, "."))
	}

	types := buildDataTemplate(root)
	err := generateGoCodeToFile(TemplateData{Types: types})
	if err != nil {
		return err
	}
	err = generatePugCodeToFile(root)
	if err != nil {
		return err
	}

	return nil
}

func generateGoCodeToFile(data TemplateData) error {

	file, _ := os.ReadFile("template/i18n/Go.tmpl")
	tmpl := string(file)

	outputFile, err := os.Create("gen/i18n_gen.go")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Parse and execute the template
	tmplParsed, err := template.New("i18n_tmpl").Funcs(
		template.FuncMap{"Upper": Upper},
	).Parse(tmpl)
	if err != nil {
		return err
	}

	err = tmplParsed.Execute(outputFile, data)
	if err != nil {
		return err
	}

	return nil
}

func generatePugCodeToFile(data *Folder) error {

	file, _ := os.ReadFile("template/i18n/Pug.tmpl")
	tmpl := string(file)

	outputFile, err := os.Create("gen/i18n_gen.pug")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	tmplParsed, err := template.New("i18n_tmpl").Funcs(
		template.FuncMap{"Upper": Upper},
	).Parse(tmpl)
	if err != nil {
		return err
	}

	err = tmplParsed.Execute(outputFile, data)
	if err != nil {
		return err
	}

	return nil
}

func buildDataTemplate(root *Folder) []I18nType {
	var types []I18nType

	for _, f := range root.Items {
		if len(f.Items) == 0 {
			continue
		}

		var fields []Field
		for _, f := range f.Items {
			fields = append(fields, Field{
				Name:        f.Name,
				IsDirectory: len(f.Items) > 0,
				Key:         f.FullKey,
			})
		}
		types = append(types, I18nType{
			Name:   f.Name,
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
