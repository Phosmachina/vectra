package generator

import (
	"fmt"
	"strings"
)

type Controller struct {
	Name   string            `yaml:"name"`
	IsView bool              `yaml:"is_view"`
	Routes []Route           `yaml:"routes"`
	Bodies map[string]string `yaml:"-"`
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

	var files []SourceFile
	for _, controller := range cfg.Controllers {
		kindStr := "service"
		if controller.IsView {
			kindStr = "view"
		}
		files = append(files,
			NewDynSourceFile(
				fmt.Sprintf("src/controller/%s_controller.go.tmpl", kindStr),
				fmt.Sprintf("src/controller/%s_controller.go",
					strings.ToLower(controller.Name)),
				Skeleton),
		)
	}

	generator := NewAbstractGenerator(
		"controllers",
		[]string{
			"Controllers",
		},
		Report{
			Files:   files,
			Version: 1,
		}, cfg)

	n := &Controllers{}
	n.Generator = generator
	n.IGenerator = n

	return generator
}

func (i *Controllers) Generate() {

	for _, controller := range i.vectra.Controllers {
		controller.Bodies = extractFunctionBody(
			i.vectra.ProjectPath + "src/controller/" +
				fmt.Sprintf("%s_controller.go", strings.ToLower(controller.Name)),
		)
	}

	i.Generator.Generate(i.vectra.Controllers)
}
