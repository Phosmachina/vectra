package generator

import (
	"github.com/go-ini/ini"
	"os"
	"path/filepath"
	"strings"
)

type I18n struct {
	*Generator
	dic map[string]string
}

func NewI18n(cfg *Vectra) *Generator {

	generator := NewAbstractGenerator(
		"i18n",
		[]string{
			"DefaultLang",
		},
		Report{
			Files: []SourceFile{
				NewSourceFile("src/model/i18n/i18n_gen.go.tmpl", FullGen),
				NewSourceFile("src/view/pug/shared/i18n_completion_variables.pug.tmpl",
					FullGen),
			},
			Version: 1,
		}, cfg)

	n := &I18n{}
	n.Generator = generator
	n.IGenerator = n

	return generator
}

func (i *I18n) Generate() {

	i.dic = make(map[string]string)

	path := filepath.Join(i.projectPath, "data", "i18n", i.vectra.DefaultLang)
	_ = i.loadData(path, "")

	var root = newFolder("", nil)

	for k := range i.dic {
		root.add(strings.Split(k, "."))
	}

	types := buildDataTemplate(root)

	i.Generator.Generate(map[string]any{
		"i18n_gen":                  TemplateData{Types: types},
		"i18n_completion_variables": root,
	})
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

			cfg, _ := ini.LoadSources(ini.LoadOptions{}, data)
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
