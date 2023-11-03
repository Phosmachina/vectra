package generator

import "path/filepath"

type ViewTypes struct {
	Types        []VectraType[SimpleAttribute] `yaml:"types"`
	Constructors []VectraType[SimpleAttribute] `yaml:"constructors"`
	Bodies       map[string]string             `yaml:"-"`
}

type VectraType[T any] struct {
	Name       string `yaml:"name"`
	Attributes []T    `yaml:"attributes"`
}

type SimpleAttribute struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type AttributeWithTag struct {
	SimpleAttribute `yaml:"base"`
	ModTag          string `yaml:"mod"`
	ValidatorTag    string `yaml:"validator"`
}

type Types struct {
	*Generator
}

func NewTypes(cfg *Vectra) *Generator {

	generator := NewAbstractGenerator(
		"types",
		[]string{
			"StorageTypes",
			"ExchangeTypes",
			"ViewTypes",
		},
		Report{
			Files: []SourceFile{
				//NewSourceFile("src/model/storage/storage_gen.go.tmpl", FullGen),
				NewSourceFile("src/view/go/view.go.tmpl", Skeleton),
			},
			Version: 1,
		}, cfg)

	n := &Types{}
	n.Generator = generator
	n.IGenerator = n

	return generator
}

func (i *Types) Generate() {

	i.vectra.ViewTypes.Bodies = extractFunctionBody(filepath.Join(i.vectra.ProjectPath,
		"src", "view", "go", "view.go"))
	i.Generator.Generate(map[string]any{
		"StorageTypes":  i.vectra.StorageTypes,
		"ExchangeTypes": i.vectra.ExchangeTypes,
		"ViewTypes":     i.vectra.ViewTypes,
	})
}
