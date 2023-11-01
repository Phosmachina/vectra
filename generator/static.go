package generator

type Static struct {
	Generator
}

func NewStatic(cfg *Vectra) *Static {

	files := []SourceFile{
		NewSourceFile("static/favicon.ico", Skeleton),
		NewSourceFile("static/js/main.js", Copy),
	}
	if cfg.WithSassExample {
		files = append(files, NewSourceFile("static/css/", Copy))
	}
	if cfg.WithI18nExample {
		files = append(files, NewSourceFile("data/i18n/", Copy))
	}
	if cfg.WithGitignore {
		files = append(files, NewSourceFile(".gitignore", Copy))
	}

	generator := NewAbstractGenerator(
		"static",
		[]string{
			"WithSassExample",
			"WithI18nExample",
			"WithGitignore",
		},
		Report{
			Files:   files,
			Version: 1,
		}, cfg,
	)
	n := Static{}
	n.Generator = *generator

	return &n
}

func (i *Static) Generate() {
	i.Generator.Generate(nil)
}
