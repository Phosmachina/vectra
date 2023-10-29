package generator

import (
	"github.com/go-ini/ini"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type I18nConfig struct {
	DefaultLang string `yaml:"default_lang"`
}

type I18n struct {
	Generator
	dic map[string]string
}

func NewI18n(projectPath string, cfg *I18nConfig) *I18n {

	config := cfg
	generator := NewAbstractGenerator("i18n", Report{
		Files: []SourceFile{
			NewSourceFile("src/model/i18n/i18n_gen.go.tmpl", Static),
			NewSourceFile("src/view/go/view_gen.go.tmpl", Static),
		},
		Config:  config,
		Version: 1,
	}, projectPath)
	n := I18n{}
	n.Generator = *generator

	n.dic = make(map[string]string)

	return &n
}

func (i *I18n) Generate() {
	config := i.nextReport.Config.(I18nConfig)

	path := filepath.Join(i.projectPath, config.DefaultLang)
	_ = i.loadData(path, "")

	var root = newFolder("", nil)

	for k := range i.dic {
		root.add(strings.Split(k, "."))
	}

	//i.Generator.Generate()

	types := buildDataTemplate(root)
	err := generateGoCodeToFile(TemplateData{Types: types})
	if err != nil {
		// TODO print error
	}
	err = generatePugCodeToFile(root)
	if err != nil {
		// TODO print error
	}

}

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

func (i *I18n) loadData(path string, prefix string) error {

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		key := entry.Name()
		if entry.IsDir() {
			err := i.loadData(filepath.Join(path, key), prefix+key+".")
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(key, ".ini") {
			fullKey := prefix + strings.TrimSuffix(key, ".ini")

			data, err := os.ReadFile(filepath.Join(path, key))
			if err != nil {
				return err
			}

			cfg, err := ini.LoadSources(ini.LoadOptions{}, data)
			for _, k := range cfg.Section("").Keys() {
				i.dic[fullKey+"."+k.Name()] = k.Value()
			}
		}
	}

	return nil
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
	).Funcs(
		template.FuncMap{"OnlyPrefix": OnlyPrefix},
	).Funcs(
		template.FuncMap{"IsNotPlural": IsNotPlural},
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
	).Funcs(
		template.FuncMap{"OnlyPrefix": OnlyPrefix},
	).Funcs(
		template.FuncMap{"IsNotPlural": IsNotPlural},
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
