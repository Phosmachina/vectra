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
		DefaultLang:          "en",
		DevPort:              8080,
		ProductionPort:       8100,
		WithGitignore:        true,
		WithDockerDeployment: true,
		WithI18nExample:      true,
		WithSassExample:      true,
		WithPugExample:       true,
		WithIdeaConfig:       true,
		WithDockerPipe:       true,
		StorageTypes: []VectraType[SimpleAttribute]{
			{Name: "Role", Attributes: []SimpleAttribute{
				{Name: "Name", Type: "string"},
				{Name: "Level", Type: "int"}}},
			{Name: "User", Attributes: []SimpleAttribute{
				{Name: "IsActivated", Type: "bool"},
				{Name: "Password", Type: "[]byte"},
				{Name: "Firstname", Type: "string"},
				{Name: "Lastname", Type: "string"},
				{Name: "Email", Type: "string"},
				{Name: "Sessions", Type: "map[string]SessionItem"}}},
		},
		ViewTypes: ViewTypes{
			Types: []VectraType[SimpleAttribute]{
				{Name: "GlobalCtx", Attributes: []SimpleAttribute{
					{Name: "IsDev", Type: "bool"},
					{Name: "Domain", Type: "string"},
					{Name: "TabTitle", Type: "string"},
					{Name: "Lang", Type: "string"},
					{Name: "User", Type: "UserCtx"},
				}},
				{Name: "UserCtx", Attributes: []SimpleAttribute{
					{Name: "ID", Type: "string"},
					{Name: "Role", Type: "Role"},
					{Name: "IsActivated", Type: "bool"},
					{Name: "Firstname", Type: "string"},
					{Name: "Lastname", Type: "string"},
					{Name: "Email", Type: "string"},
				}},
			},
			Constructors: []VectraType[SimpleAttribute]{
				{Name: "NewGlobalCtx", Attributes: []SimpleAttribute{
					{Name: "tabSuffix", Type: "string"},
					{Name: "userId", Type: "string"},
				}},
				{Name: "newUserCtx", Attributes: []SimpleAttribute{
					{Name: "userId", Type: "string"},
				}},
			},
		},
		Controllers: []Controller{
			{Name: "View",
				IsView: true,
				Routes: []Route{
					{Kind: "Get", Path: "/", Target: "root"},
					{Kind: "Get", Path: "/init", Target: "init"},
					{Kind: "Get", Path: "/login", Target: "login"},
					{Kind: "Get", Path: "/sign", Target: "sign"},
				},
			},
			{Name: "ApiV1",
				IsView: false,
				Routes: []Route{
					{Kind: "Post", Path: "/activate/admin", Target: "activateAdmin"},
					{Kind: "Post", Path: "/login", Target: "login"},
				},
			},
		},
		Services: []Service{
			{
				Name: "ApiV1",
				Errors: []string{
					"ErrorNotFirstLaunch",
					"ErrorInvalidToken",
					"ErrorUnauthorised",
					"ErrorUserExist",
					"ErrorUserDisabled",
					"ErrorInvalidUserRef",
				},
				Methods: []Method{
					{Name: "IsFirstLaunch",
						Inputs:  []SimpleAttribute{},
						Outputs: []string{"bool"}},
					{Name: "IsConnected",
						Inputs: []SimpleAttribute{
							{Name: "session", Type: "*session.Session"},
						},
						Outputs: []string{"Role"}},
					{Name: "ActivateAdmin",
						Inputs: []SimpleAttribute{
							{Name: "info", Type: "ConnectAdminExch"},
						},
						Outputs: []string{"error"}},
					{Name: "CreateUser",
						Inputs: []SimpleAttribute{
							{Name: "info", Type: "UserExch"},
						},
						Outputs: []string{"error"}},
					{Name: "Connect",
						Inputs: []SimpleAttribute{
							{Name: "info", Type: "ConnectExch"},
							{Name: "cookie", Type: "string"},
							{Name: "ua", Type: "string"},
						},
						Outputs: []string{"error", "*ObjWrapper[User]"}},
				},
				ExchangeTypes: []VectraType[AttributeWithTag]{
					{Name: "ConnectExch", Attributes: []AttributeWithTag{
						{SimpleAttribute: SimpleAttribute{Name: "Email", Type: "string"},
							ModTag: "trim,lcase", ValidatorTag: "required,email"},
						{SimpleAttribute: SimpleAttribute{Name: "Password", Type: "string"},
							ModTag: "", ValidatorTag: "required"},
					}},
					{Name: "ConnectAdminExch", Attributes: []AttributeWithTag{
						{SimpleAttribute: SimpleAttribute{Name: "Password", Type: "string"},
							ModTag: "", ValidatorTag: "required"},
						{SimpleAttribute: SimpleAttribute{Name: "Email", Type: "string"},
							ModTag: "trim,lcase", ValidatorTag: "required,email"},
						{SimpleAttribute: SimpleAttribute{Name: "Token", Type: "string"},
							ModTag: "", ValidatorTag: "required"},
					}},
					{Name: "UserExch", Attributes: []AttributeWithTag{
						{SimpleAttribute: SimpleAttribute{Name: "ID", Type: "string"},
							ModTag: "", ValidatorTag: ""},
						{SimpleAttribute: SimpleAttribute{Name: "Password", Type: "string"},
							ModTag: "trim,lcase", ValidatorTag: "required"},
						{SimpleAttribute: SimpleAttribute{Name: "Firstname", Type: "string"},
							ModTag: "trim,lcase", ValidatorTag: "required"},
						{SimpleAttribute: SimpleAttribute{Name: "Lastname", Type: "string"},
							ModTag: "trim,lcase", ValidatorTag: "required"},
						{SimpleAttribute: SimpleAttribute{Name: "Email", Type: "string"},
							ModTag: "trim,lcase", ValidatorTag: "required,email"},
					}},
					{Name: "ReasonExch", Attributes: []AttributeWithTag{
						{SimpleAttribute: SimpleAttribute{Name: "Reason", Type: "string"},
							ModTag: "", ValidatorTag: ""},
					}},
				},
			},
		},
	}
)

type Vectra struct {
	generators  map[string]IGenerator `yaml:"-"`
	ProjectPath string                `yaml:"-"`

	DefaultLang          string                        `yaml:"default_lang"`
	DevPort              int                           `yaml:"dev_port"`
	ProductionPort       int                           `yaml:"production_port"`
	ProductionDomain     string                        `yaml:"production_domain"`
	WithGitignore        bool                          `yaml:"with_gitignore"`
	WithDockerDeployment bool                          `yaml:"with_docker_deployment"`
	WithI18nExample      bool                          `yaml:"with_i18n_example"`
	WithSassExample      bool                          `yaml:"with_sass_example"`
	WithPugExample       bool                          `yaml:"with_pug_example"`
	WithIdeaConfig       bool                          `yaml:"with_idea_config"`
	WithDockerPipe       bool                          `yaml:"with_docker_pipe"`
	StorageTypes         []VectraType[SimpleAttribute] `yaml:"storage_types"`
	ViewTypes            ViewTypes                     `yaml:"view_types"`
	Controllers          []Controller                  `yaml:"controllers"`
	Services             []Service                     `yaml:"services"`
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
	vectra.generators = generatorsToMap(
		NewI18n(&vectra),
		NewBase(&vectra),
		NewTypes(&vectra),
		NewServices(&vectra),
		NewControllers(&vectra),
	)

	return &vectra
}

func generatorsToMap(g ...*Generator) map[string]IGenerator {

	m := map[string]IGenerator{}
	for _, generator := range g {
		m[generator.Name] = generator.IGenerator
	}

	return m
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
		log.Println("Failed to write project config.")
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
