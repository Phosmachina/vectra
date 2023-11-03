package generator

type Service struct {
	Name    string
	Errors  []string
	Methods []Method
}

type Method struct {
	Name    string
	Inputs  []SimpleAttribute
	Outputs []string
}

type Services struct {
	*Generator
}

func NewServices(cfg *Vectra) *Generator {

	generator := NewAbstractGenerator(
		"services",
		[]string{},
		Report{
			Files: []SourceFile{
				NewSourceFile("", FullGen),
			},
			Version: 1,
		}, cfg)

	n := &Services{}
	n.Generator = generator
	n.IGenerator = n

	return generator
}

func (i *Services) Generate() {
	i.Generator.Generate(map[string]any{})
}
