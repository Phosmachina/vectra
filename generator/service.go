package generator

import (
	"fmt"
	"strings"
)

type Service struct {
	Name          string                         `yaml:"name"`
	Errors        []string                       `yaml:"errors"`
	Methods       []Method                       `yaml:"methods"`
	ExchangeTypes []VectraType[AttributeWithTag] `yaml:"exchange_types"`
	Bodies        map[string]string              `yaml:"-"`
}

type Method struct {
	Name    string            `yaml:"name"`
	Inputs  []SimpleAttribute `yaml:"inputs"`
	Outputs []string          `yaml:"outputs"`
}

type Services struct {
	*Generator
}

func NewServices(cfg *Vectra) *Generator {

	var files []SourceFile
	for _, service := range cfg.Services {
		files = append(files,
			NewDynSourceFile(
				"src/model/service/service.go.tmpl",
				fmt.Sprintf("src/model/service/%s_service.go",
					strings.ToLower(service.Name)),
				Skeleton),
		)
	}

	generator := NewAbstractGenerator(
		"services",
		[]string{
			"Services",
		},
		Report{
			Files:   files,
			Version: 1,
		}, cfg)

	n := &Services{}
	n.Generator = generator
	n.IGenerator = n

	return generator
}

func (i *Services) Generate() {

	for _, service := range i.vectra.Services {
		service.Bodies = extractFunctionBody(
			i.vectra.ProjectPath + "src/model/service/" +
				fmt.Sprintf("%s_service.go", strings.ToLower(service.Name)),
		)
	}

	i.Generator.Generate(i.vectra.Services)
}
