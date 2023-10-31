package generator

type Core struct {
	Generator
}

func NewCore(cfg *Vectra) *Core {

	generator := NewAbstractGenerator(
		"core",
		[]string{
			"DefaultLang",
		},
		Report{
			Files: []SourceFile{
				NewSourceFile("app.go", Critical),
				NewDynSourceFile("go.mod.embed", "go.mod", Critical),
				NewDynSourceFile("go.sum.embed", "go.sum", Critical),
				NewSourceFile("src/model/i18n/i18n.go", Critical),
				NewSourceFile("src/model/service/service.go", Critical),
				NewSourceFile("src/model/storage/storage.go", Critical),
				NewSourceFile("src/model/helpers.go", Critical),
				NewSourceFile("src/controller/controller.go", Critical),
				NewSourceFile("src/view/go/view.go", Critical),
			},
			Version: 1,
		}, cfg,
	)
	n := Core{}
	n.Generator = *generator

	return &n
}

func (i *Core) Generate() {
	i.Generator.Generate(i.nextReport.Config)
}
