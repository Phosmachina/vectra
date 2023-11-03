package generator

type Controller struct {
	IsView bool    `yaml:"is_view"`
	Routes []Route `yaml:"routes"`
}

type Route struct {
	Kind   string `yaml:"kind"`
	Path   string `yaml:"path"`
	Target string `yaml:"target"`
}

type Controllers struct {
	*Generator
}

func NewControllers(cfg *Vectra) *Generator {

	generator := NewAbstractGenerator(
		"controllers",
		[]string{},
		Report{
			Files: []SourceFile{
				NewSourceFile("", FullGen),
			},
			Version: 1,
		}, cfg)

	n := &Controllers{}
	n.Generator = generator
	n.IGenerator = n

	return generator
}

func (i *Controllers) Generate() {
	i.Generator.Generate(map[string]any{})
}
