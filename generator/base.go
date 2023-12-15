package generator

type NetworkConfig struct {
	Domain string `yaml:"domain"`
	Port   int    `yaml:"port"`
	IsIPv6 bool   `yaml:"isIPv6"`
}

type Base struct {
	*Generator
}

func NewBase(cfg *Vectra) *Generator {

	files := []SourceFile{
		NewSourceFile("data/config/configuration.yml.tmpl", Copy),
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
		files = append(files, NewSourceFile("src/view/pug/component/", Copy))
		files = append(files, NewSourceFile("src/view/pug/shared/layout.pug", Skeleton))
		files = append(files, NewSourceFile("src/view/pug/shared/mixins.pug", Skeleton))
		files = append(files, NewSourceFile("src/view/pug/shared/sprite.pug", Skeleton))
		files = append(files, NewSourceFile("src/view/pug/index.pug", Skeleton))
		files = append(files, NewSourceFile("src/view/pug/init.pug", Copy))
		files = append(files, NewSourceFile("src/view/pug/login.pug", Copy))
	}
	if cfg.WithGitignore {
		files = append(files, NewSourceFile(".gitignore", Copy))
	}
	if cfg.WithDockerDeployment && cfg.isProdGen {
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
			"WithDockerDeployment",
			"NetConfDev",
			"NetConfProd",
			"DefaultLang",
		},
		Report{
			Files:   files,
			Version: 2,
		}, cfg)

	n := &Base{}
	n.Generator = generator
	n.IGenerator = n

	return generator
}

func (i *Base) Generate() {

	ctx := map[string]any{"DefaultLang": i.vectra.DefaultLang}

	if i.vectra.isProdGen {
		ctx["Domain"] = i.vectra.NetConfProd.Domain
		ctx["Port"] = i.vectra.NetConfProd.Port
		ctx["IsIPv6"] = i.vectra.NetConfProd.IsIPv6
	} else {
		ctx["Domain"] = i.vectra.NetConfDev.Domain
		ctx["Port"] = i.vectra.NetConfDev.Port
		ctx["IsIPv6"] = i.vectra.NetConfDev.IsIPv6
	}

	i.Generator.Generate(ctx)
}
