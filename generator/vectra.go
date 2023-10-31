package generator

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

var (
	defaultVectra = Vectra{
		generators:  nil,
		DefaultLang: "en",
	}
)

type Vectra struct {
	generators  map[string]IGenerator `yaml:"-"`
	ProjectPath string                `yaml:"-"`

	DefaultLang string `yaml:"default_lang"`

	DevPort          int    `yaml:"dev_port"`
	ProductionPort   int    `yaml:"production_port"`
	ProductionDomain string `yaml:"production_domain"`

	WithGitignore     bool `yaml:"with_gitignore"`
	WithDockerfile    bool `yaml:"with_dockerfile"`
	WithDockerCompose bool `yaml:"with_docker_compose"`

	WithI18nExample bool `yaml:"with_i18n_example"`
	WithSassExample bool `yaml:"with_sass_example"`
	WithPugExample  bool `yaml:"with_pug_example"`

	WithIdeaConfig bool `yaml:"with_idea_config"`
	WithDockerPipe bool `yaml:"with_docker_pipe"`
}

func NewVectra(projectPath string) *Vectra {

	var vectra Vectra
	data, err := os.ReadFile(filepath.Join(
		projectPath,
		FolderProject,
		"project.yml",
	))
	if err != nil {
		vectra = defaultVectra
	}
	if yaml.Unmarshal(data, &vectra) != nil {
		vectra = defaultVectra
	}
	vectra.ProjectPath = projectPath

	core := NewCore(&vectra)
	i18n := NewI18n(&vectra)
	generators := map[string]IGenerator{
		core.Name: core,
		i18n.Name: i18n,
	}
	vectra.generators = generators

	return &vectra
}

type Configuration struct {
	Storage StorageType
}

type StorageType struct {
	Types map[string][]Attribute
}

type Attribute struct {
	Name string
	Type string
	Tag  string
}

func (v *Vectra) Init() {

	err := os.MkdirAll(filepath.Join(v.ProjectPath, FolderProject), 0755)
	if err != nil {
		// TODO print error
	}

	data, err := yaml.Marshal(v)
	if err != nil {
		return
	}
	path := filepath.Join(
		v.ProjectPath,
		FolderProject,
		"project.yml",
	)
	if os.WriteFile(path, data, 0644) != nil {
		log.Printf("Failed to write project config.\n")
	}
}

func (v *Vectra) FullReport() {
	for _, g := range v.generators {
		g.PrintReport()
	}
}

func (v *Vectra) FullGenerate() {
	for _, g := range v.generators {
		g.Generate()
	}
}

func (v *Vectra) Generate(key string) {
	v.generators[key].Generate()
}

func (v *Vectra) Report(key string) {
	v.generators[key].PrintReport()
}

func (v *Vectra) GetFieldsAsMap(paths []string) map[string]any {
	result := make(map[string]any)

	configValue := reflect.ValueOf(*v)

	for _, path := range paths {
		fieldNames := strings.Split(path, ".")
		current := configValue

		for _, fieldName := range fieldNames {
			fieldValue := current.FieldByName(fieldName)
			if !fieldValue.IsValid() {
				break
			}

			current = fieldValue
		}

		if current.IsValid() {
			result[path] = current.Interface()
		}
	}

	return result
}

func (v *Vectra) setField(path string, value any) error {
	fields := strings.Split(path, ".")
	current := reflect.ValueOf(v).Elem()

	for i, field := range fields {
		fieldValue := current.FieldByName(field)
		if !fieldValue.IsValid() {
			return fmt.Errorf("invalid field: %s", field)
		}

		if i == len(fields)-1 {
			// This is the last part of the path, set the value
			fieldType := fieldValue.Type()
			val := reflect.ValueOf(value)
			if !val.Type().AssignableTo(fieldType) {
				return fmt.Errorf("invalid value type for field: %s", field)
			}
			fieldValue.Set(val)
		} else {
			// Not the last part, navigate to the nested struct
			if fieldValue.Kind() == reflect.Ptr {
				// If it's a pointer and nil, create a new instance
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				}
				fieldValue = fieldValue.Elem()
			} else if fieldValue.Kind() == reflect.Struct {
				// If it's a struct, continue with the next part of the path
				current = fieldValue
				continue
			} else {
				// Invalid path, can't navigate further
				return fmt.Errorf("invalid path: %s", path)
			}

			current = fieldValue
		}
	}

	return nil
}

func (v *Vectra) DeserializeFromMap(data map[string]any) error {
	for path, value := range data {
		if err := v.setField(path, value); err != nil {
			return err
		}
	}
	return nil
}
