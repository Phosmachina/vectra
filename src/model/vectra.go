package model

const (
	TemplateFolderPath = "template"
)

type Vectra struct {
	projectPath string
}

func NewVectra(projectPath string) *Vectra {
	return &Vectra{projectPath: projectPath}
}

type Configuration struct {
}
