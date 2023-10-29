package model

import (
	"Vectra/model/generator"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

var (
	defaultVectra = Vectra{
		generators:  nil,
		DefaultLang: "en",
	}
)

type Vectra struct {
	generators  map[string]generator.IGenerator `yaml:"-"`
	ProjectPath string                          `yaml:"-"`

	DefaultLang string `yaml:"default_lang"`
}

func NewVectra(projectPath string) *Vectra {

	var vectra Vectra
	data, err := os.ReadFile(filepath.Join(
		projectPath,
		generator.FolderProject,
		"project.yml",
	))
	if err != nil {
		vectra = defaultVectra
	}
	if yaml.Unmarshal(data, vectra) != nil {
		vectra = defaultVectra
	}
	vectra.ProjectPath = projectPath

	core := generator.NewCore(projectPath, vectra.toCoreConfig())
	i18n := generator.NewI18n(projectPath, vectra.toI18nConfig())
	generators := map[string]generator.IGenerator{
		core.Name: core,
		i18n.Name: i18n,
	}
	vectra.generators = generators

	return &vectra
}

func (v *Vectra) toCoreConfig() *generator.CoreConfig {
	return &generator.CoreConfig{
		DefaultLang: v.DefaultLang,
	}
}

func (v *Vectra) toI18nConfig() *generator.I18nConfig {
	return &generator.I18nConfig{
		DefaultLang: v.DefaultLang,
	}
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

	err := os.MkdirAll(filepath.Join(v.ProjectPath, generator.FolderProject), 0755)
	if err != nil {
		// TODO print error
	}

	data, err := yaml.Marshal(v)
	if err != nil {
		return
	}
	path := filepath.Join(
		v.ProjectPath,
		generator.FolderProject,
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
