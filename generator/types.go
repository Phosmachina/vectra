package generator

type ViewTypes struct {
	Types        []VectraType[SimpleAttribute] `yaml:"types"`
	Constructors []ViewTypeConstructor         `yaml:"constructors"`
	Bodies       map[string]string             `yaml:"-"`
}

type ViewTypeConstructor struct {
	Name       string            `yaml:"name"`
	IsPageCtx  bool              `yaml:"is_page_ctx"`
	Attributes []SimpleAttribute `yaml:"attributes"`
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

type ConfigurationAttribute struct {
	SimpleAttribute `yaml:"base"`
	IsTransient     bool `yaml:"is_transient"`
}

type Types struct {
	*Generator
}

func NewTypes(cfg *Vectra) *Generator {

	generator := NewAbstractGenerator(
		"types",
		[]string{
			"Configuration",
			"DefaultLang",
			"StorageTypes",
			"ViewTypes",
		},
		Report{
			Files: []SourceFile{
				NewSourceFile("src/model/storage/configuration.go.tmpl", FullGen),
				NewSourceFile("src/model/storage/types.go.tmpl", FullGen),
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

	i.vectra.ViewTypes.Bodies = extractFunctionBody(
		i.vectra.ProjectPath + "/src/view/go/view.go")

	i.Generator.Generate(map[string]any{
		"Configuration": map[string]any{
			"IsDev":         true,
			"DefaultLang":   i.vectra.DefaultLang,
			"Configuration": i.vectra.Configuration,
		},
		"StorageTypes": i.vectra.StorageTypes,
		"ViewTypes":    i.vectra.ViewTypes,
	})
}
