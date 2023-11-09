package generator

type Base struct {
	*Generator
}

func NewBase(cfg *Vectra) *Generator {

	files := []SourceFile{
		NewSourceFile("static/favicon.ico", Skeleton),
		NewSourceFile("static/js/main.js", Copy),
		NewSourceFile("app.go", CorePart),
		NewDynSourceFile("go.mod.embed", "go.mod", CorePart),
		NewDynSourceFile("go.sum.embed", "go.sum", CorePart),
		NewSourceFile("src/model/i18n/i18n.go", CorePart),
		NewSourceFile("src/model/service/service.go", CorePart),
		NewSourceFile("src/model/storage/storage.go", CorePart),
		NewSourceFile("src/model/helpers.go", CorePart),
		NewSourceFile("src/controller/controller.go", CorePart),
	}
	if cfg.WithSassExample {
		files = append(files, NewSourceFile("static/css/", Copy))
	}
	if cfg.WithI18nExample {
		files = append(files, NewSourceFile("data/i18n/", Copy))
	}
	if cfg.WithPugExample {
		files = append(files, NewSourceFile("src/view/pug/shared/layout.pug", Copy))
		files = append(files, NewSourceFile("src/view/pug/shared/mixins.pug", Copy))
		files = append(files, NewSourceFile("src/view/pug/index.pug", Skeleton))
		files = append(files, NewSourceFile("src/view/pug/init.pug", Copy))
		files = append(files, NewSourceFile("src/view/pug/login.pug", Skeleton))
	}
	if cfg.WithGitignore {
		files = append(files, NewSourceFile(".gitignore", Copy))
	}
	if cfg.WithDockerPipe {
		files = append(files, NewSourceFile(".pipe/docker-compose.yml", Copy))
		files = append(files, NewSourceFile(".pipe/Dockerfile", Copy))
	}
	if cfg.WithIdeaConfig {
		files = append(files, NewSourceFile(".idea/", Copy))
	}
	if cfg.WithDockerDeployment {
		files = append(files, NewDynSourceFile("docker-compose.yml.tmpl", "docker-compose.yml", Copy))
		files = append(files, NewDynSourceFile("Dockerfile.tmpl", "Dockerfile", Copy))
	}

	generator := NewAbstractGenerator(
		"base",
		[]string{
			"WithSassExample",
			"WithI18nExample",
			"WithPugExample",
			"WithGitignore",
			"WithDockerPipe",
			"WithIdeaConfig",
			"WithDockerDeployment",
		},
		Report{
			Files:   files,
			Version: 1,
		}, cfg)

	n := &Base{}
	n.Generator = generator
	n.IGenerator = n

	return generator
}

func (i *Base) Generate() {
	i.Generator.Generate(map[string]any{
		"ProductionPort": i.vectra.ProductionPort,
	})
}