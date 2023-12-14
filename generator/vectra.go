package generator

import (
	"bufio"
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
		DefaultLang: "en",
		SpriteConfig: SpriteConfig{
			SvgFolderPath:   "static/svg",
			OutputSpriteSvg: "static/svg/sprite",
		},
		NetConfDev: NetworkConfig{
			Domain: "localhost",
			Port:   8100,
			IsIPv6: false,
		},
		WithGitignore:        true,
		WithDockerDeployment: true,
		WithI18nExample:      true,
		WithSassExample:      true,
		WithPugExample:       true,
		StorageTypes: []VectraType[SimpleAttribute]{
			{
				"Role",
				[]SimpleAttribute{
					{Name: "Name", Type: "string"},
					{Name: "Level", Type: "int"}},
			},
			{
				"User",
				[]SimpleAttribute{
					{Name: "IsActivated", Type: "bool"},
					{Name: "Password", Type: "[]byte"},
					{Name: "Firstname", Type: "string"},
					{Name: "Lastname", Type: "string"},
					{Name: "Email", Type: "string"},
					{Name: "Sessions", Type: "map[string]SessionItem"}},
			},
		},
		ViewTypes: ViewTypes{
			Types: []VectraType[SimpleAttribute]{
				{
					"GlobalCtx",
					[]SimpleAttribute{
						{Name: "IsDev", Type: "bool"},
						{Name: "Domain", Type: "string"},
						{Name: "Port", Type: "int"},
						{Name: "TabTitle", Type: "string"},
						{Name: "Lang", Type: "string"},
						{Name: "Langs", Type: "[]string"},
						{Name: "User", Type: "UserCtx"},
					},
				},
				{
					"UserCtx",
					[]SimpleAttribute{
						{Name: "ID", Type: "string"},
						{Name: "Role", Type: "Role"},
						{Name: "IsActivated", Type: "bool"},
						{Name: "Firstname", Type: "string"},
						{Name: "Lastname", Type: "string"},
						{Name: "Email", Type: "string"},
					},
				},
			},
			Constructors: []ViewTypeConstructor{
				{
					"NewGlobalCtx",
					false,
					[]SimpleAttribute{
						{Name: "tabSuffix", Type: "string"},
						{Name: "userId", Type: "string"},
					},
				},
				{
					"newUserCtx",
					false,
					[]SimpleAttribute{
						{Name: "userId", Type: "string"},
					},
				},
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
					{Kind: "Post", Path: "/update/lang", Target: "updateLang"},
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
						{SimpleAttribute{Name: "Email", Type: "string"}, "trim,lcase", "required,email"},
						{SimpleAttribute{Name: "Password", Type: "string"}, "", "required"},
					}},
					{Name: "ConnectAdminExch", Attributes: []AttributeWithTag{
						{SimpleAttribute{Name: "Password", Type: "string"}, "", "required"},
						{SimpleAttribute{Name: "Email", Type: "string"}, "trim,lcase", "required,email"},
						{SimpleAttribute{Name: "Token", Type: "string"}, "", "required"},
					}},
					{Name: "UserExch", Attributes: []AttributeWithTag{
						{SimpleAttribute{Name: "ID", Type: "string"}, "", ""},
						{SimpleAttribute{Name: "Password", Type: "string"}, "trim,lcase", "required"},
						{SimpleAttribute{Name: "Firstname", Type: "string"}, "trim,lcase", "required"},
						{SimpleAttribute{Name: "Lastname", Type: "string"}, "trim,lcase", "required"},
						{SimpleAttribute{Name: "Email", Type: "string"}, "trim,lcase", "required,email"},
					}},
					{Name: "LangExch", Attributes: []AttributeWithTag{
						{SimpleAttribute{Name: "Lang", Type: "string"}, "trim,lcase", "required"},
					}},
					{Name: "ReasonExch", Attributes: []AttributeWithTag{
						{SimpleAttribute{Name: "Reason", Type: "string"}, "", ""},
					}},
				},
			},
		},
		Configuration: []ConfigurationAttribute{
			{SimpleAttribute{Name: "Domain", Type: "string"}, false},
			{SimpleAttribute{Name: "Port", Type: "int"}, false},
			{SimpleAttribute{Name: "IsIPv6", Type: "bool"}, false},
			{SimpleAttribute{Name: "TabPrefix", Type: "string"}, false},
			{SimpleAttribute{Name: "CurrentLang", Type: "string"}, true},
			{SimpleAttribute{Name: "Roles", Type: "map[string]int"}, false},
			{SimpleAttribute{Name: "AccessRules", Type: "[]AccessRule"}, false},
		},
	}
)

type NetworkConfig struct {
	Domain string `yaml:"domain"`
	Port   int    `yaml:"port"`
	IsIPv6 bool   `yaml:"isIPv6"`
}

type Vectra struct {
	generators  map[string]IGenerator `yaml:"-"`
	ProjectPath string                `yaml:"-"`
	isProdGen   bool                  `yaml:"-"`

	SpriteConfig         SpriteConfig                  `yaml:"sprite_config"`
	NetConfProd          NetworkConfig                 `yaml:"net_conf_prod"`
	NetConfDev           NetworkConfig                 `yaml:"net_conf_dev"`
	ProjectName          string                        `yaml:"project_name"`
	DefaultLang          string                        `yaml:"default_lang"`
	WithGitignore        bool                          `yaml:"with_gitignore"`
	WithDockerDeployment bool                          `yaml:"with_docker_deployment"`
	WithI18nExample      bool                          `yaml:"with_i18n_example"`
	WithSassExample      bool                          `yaml:"with_sass_example"`
	WithPugExample       bool                          `yaml:"with_pug_example"`
	StorageTypes         []VectraType[SimpleAttribute] `yaml:"storage_types"`
	ViewTypes            ViewTypes                     `yaml:"view_types"`
	Controllers          []Controller                  `yaml:"controllers"`
	Services             []Service                     `yaml:"services"`
	Configuration        []ConfigurationAttribute      `yaml:"configuration"`
}

func NewVectra(projectPath string) *Vectra {

	var vectra Vectra
	data, err := os.ReadFile(filepath.Join(
		projectPath,
		FolderProject,
		"project.yml",
	))
	if err != nil {
		fmt.Println("No configuration file found ; use the default one.")
		vectra = defaultVectra
	}
	err = yaml.Unmarshal(data, &vectra)
	if err != nil {
		log.Fatal("Failed to parse the project configuration file ; check syntax.")
	}
	if vectra.ProjectName == "" {
		vectra.ProjectName = filepath.Base(projectPath)
	}

	fullPath, err := filepath.Abs(projectPath)
	if err == nil {
		vectra.ProjectPath = fullPath
	} else {
		vectra.ProjectPath = projectPath
	}

	vectra.generators = generatorsToMap(
		NewI18n(&vectra),
		NewBase(&vectra),
		NewTypes(&vectra),
		NewServices(&vectra),
		NewControllers(&vectra),
	)

	return &vectra
}

func (v *Vectra) Watch() {

	fmt.Println("========= Check docker =========")

	if !IsDockerInstalled() {
		fmt.Println("Docker is not correctly installed.")
		return
	}

	images := []string{"Autoprefixer", "Pug", "Sass", "MinifyJS", "MinifyCSS"}
	for _, image := range images {
		imageName := "phosmachina/" + strings.ToLower(image)
		err := CreateDockerImage(image+".Dockerfile", imageName)
		if err != nil {
			fmt.Println("Failed to create image: " + image)
			return
		}
		containerName := v.ProjectName + "_" + image
		err = CreateDockerContainer(containerName, v.ProjectPath, imageName)
		if err != nil {
			fmt.Println("Failed to create container: " + containerName)
			return
		}
		err = StartDockerContainer(containerName)
		if err != nil {
			fmt.Println("Failed to start container: " + containerName)
			return
		}
	}

	fmt.Println("=========   Watching   =========")

	go watchPug(v)

	go watchJS(v)

	go watchI18n(v)

	go watchSass(v)

	<-make(chan struct{})
}

func (v *Vectra) Init() {

	err := os.MkdirAll(filepath.Join(v.ProjectPath, FolderReport), 0755)
	if err != nil {
		fmt.Println("Failed to create the project directory:", err)
		return
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
		log.Println("Failed to write the project config.")
	}
}

func (v *Vectra) FullReport() {
	for _, g := range v.generators {
		g.PrintReport()
	}
}

func (v *Vectra) FullGenerate() {

	fmt.Println("Warning: Full generation may override many files. Do you wish to continue? (yes/no)")
	text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	response := strings.ToLower(strings.TrimSpace(text))

	if response != "yes" {

		fmt.Println("Full generation aborted.")
		return
	}

	for _, g := range v.generators {
		g.Generate()
	}
	generateSpriteSvg(v)
}

func (v *Vectra) Generate(key string) {
	if key == "sprite" {
		generateSpriteSvg(v)
		return
	}
	generator, ok := v.generators[key]
	if !ok {
		fmt.Println("The generator", key, "does not exist.")
		return
	}
	generator.Generate()
}

func (v *Vectra) Report(key string) {
	generator, ok := v.generators[key]
	if !ok {
		fmt.Println("The generator", key, "does not exist.")
		return
	}
	generator.PrintReport()
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

func (v *Vectra) DeserializeFromMap(data map[string]any) error {
	for path, value := range data {
		if err := v.setField(path, value); err != nil {
			return err
		}
	}
	return nil
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

func generatorsToMap(g ...*Generator) map[string]IGenerator {

	m := map[string]IGenerator{}
	for _, generator := range g {
		m[generator.Name] = generator.IGenerator
	}

	return m
}
