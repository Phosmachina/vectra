package generator

type CoreConfig struct {
	DefaultLang string `yaml:"default_lang"`
}

type Core struct {
	Generator
}

func NewCore(projectPath string, cfg *CoreConfig) *Core {

	config := cfg
	generator := NewAbstractGenerator(
		"core",
		Report{
			Files: []SourceFile{
				NewSourceFile("app.go", Critical),
			},
			Config:  config,
			Version: 1,
		}, projectPath)
	n := Core{}
	n.Generator = *generator

	return &n
}

func (i *Core) Generate() {
	i.Generator.Generate(i.nextReport.Config.(CoreConfig))
}
