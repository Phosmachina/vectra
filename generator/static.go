package generator

type Static struct {
	*Generator
}

func NewStatic(cfg *Vectra) *Generator {

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
	if cfg.WithPugExample {
		files = append(files, NewSourceFile("view/pug/shared/layout.pug", Copy))
		files = append(files, NewSourceFile("view/pug/shared/mixins.pug", Copy))
		files = append(files, NewSourceFile("view/pug/index.pug", Skeleton))
		files = append(files, NewSourceFile("view/pug/init.pug", Copy))
		files = append(files, NewSourceFile("view/pug/login.pug", Copy))
	}
	if cfg.WithGitignore {
		files = append(files, NewSourceFile(".gitignore", Copy))
	}

	generator := NewAbstractGenerator(
		"static",
		[]string{
			"WithSassExample",
			"WithI18nExample",
			"WithPugExample",
			"WithGitignore",
		},
		Report{
			Files:   files,
			Version: 1,
		}, cfg)

	n := &Static{}
	n.Generator = generator
	n.IGenerator = n

	return generator
}

func (i *Static) Generate() {
	i.Generator.Generate(nil)
}
