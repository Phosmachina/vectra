package generator

type Core struct {
	*Generator
}

func NewCore(cfg *Vectra) *Generator {

	generator := NewAbstractGenerator(
		"core",
		[]string{
			"DefaultLang",
		},
		Report{
			Files: []SourceFile{
				NewSourceFile("app.go", CorePart),
				NewDynSourceFile("go.mod.embed", "go.mod", CorePart),
				NewDynSourceFile("go.sum.embed", "go.sum", CorePart),
				NewSourceFile("src/model/i18n/i18n.go", CorePart),
				NewSourceFile("src/model/service/service.go", CorePart),
				NewSourceFile("src/model/storage/storage.go", CorePart),
				NewSourceFile("src/model/helpers.go", CorePart),
				NewSourceFile("src/controller/controller.go", CorePart),
			},
			Version: 1,
		}, cfg)

	n := &Core{}
	n.Generator = generator
	n.IGenerator = n

	return generator
}

func (i *Core) Generate() {
	i.Generator.Generate(nil)
}
