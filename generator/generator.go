package generator

import (
	"crypto/md5"
	"embed"
	"encoding/hex"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
	"time"
)

var (
	//go:embed template
	EmbedFS embed.FS

	FolderProject  = ".vectra"
	FolderTemplate = "template"
)

const (
	Same = 1 << iota
	Edited
	Deleted
	Info
)

const (
	Copy     = 1 << iota // For demo purpose or something like that
	Critical             // For core of Vectra edit with caution
	Static               // Don't be edited at all
	Skeleton             // Body functions to be filled later
)

type Report struct {
	Files   []SourceFile   `yaml:"files"`
	Config  map[string]any `yaml:"config"`
	Version int8           `yaml:"version"`
}

type SourceFile struct {
	TemplatePath string `yaml:"-"`
	RealPath     string `yaml:"path"`
	Hash         string `yaml:"hash"`
	Kind         int8   `yaml:"kind"`
	isTmpl       bool   `yaml:"-"`
}

func NewSourceFile(path string, kind int8) SourceFile {
	return NewDynSourceFile(path, strings.TrimSuffix(path, ".tmpl"), kind)
}

func NewDynSourceFile(path string, realpath string, kind int8) SourceFile {
	return SourceFile{
		TemplatePath: filepath.Join(FolderTemplate, path),
		RealPath:     realpath,
		Kind:         kind,
		isTmpl:       strings.HasSuffix(path, ".tmpl"),
	}
}

type IGenerator interface {
	Generate()
	PrintReport()
}

type Generator struct {
	Name        string
	projectPath string

	vectra          *Vectra
	configSelectors []string
	lastReport      Report
	nextReport      Report
}

func NewAbstractGenerator(
	name string,
	configSelectors []string,
	report Report,
	vectra *Vectra,
) *Generator {

	report.Config = vectra.GetFieldsAsMap(configSelectors)

	generator := Generator{
		Name:            name,
		vectra:          vectra,
		configSelectors: configSelectors,
		projectPath:     vectra.ProjectPath,
		nextReport:      report,
	}

	generator.init()

	return &generator
}

func (g *Generator) init() {
	data, err := os.ReadFile(filepath.Join(
		g.projectPath,
		FolderProject,
		g.Name+"_report.yml",
	))
	if err != nil {
		g.lastReport = Report{}
		return
	}
	if err := yaml.Unmarshal(data, &g.lastReport); err != nil {
		// TODO print error
	}
}

func (g *Generator) Generate(data any) {

	for _, file := range g.nextReport.Files {

		outputFile := filepath.Join(g.projectPath, file.RealPath)

		dir := filepath.Dir(outputFile)
		if os.MkdirAll(dir, 0744) != nil {
			// TODO print error
			continue
		}

		if !file.isTmpl {
			f, _ := EmbedFS.Open(file.TemplatePath)
			stat, _ := f.Stat()

			var err error
			if stat.IsDir() {
				err = copyDir(file.TemplatePath, outputFile)
			} else {
				err = copyFile(file.TemplatePath, outputFile)
			}
			if err != nil {
				// TODO print error
			}

			continue
		}
		in, err := EmbedFS.ReadFile(file.TemplatePath)
		if err != nil {
			// TODO print error
			continue
		}

		out, err := os.Create(outputFile)
		if err != nil {
			// TODO print error
			continue
		}
		defer out.Close()

		parsed, err := template.New("tmpl").Funcs(
			template.FuncMap{"Upper": Upper},
		).Funcs(
			template.FuncMap{"OnlyPrefix": OnlyPrefix},
		).Funcs(
			template.FuncMap{"IsNotPlural": IsNotPlural},
		).Parse(string(in))
		if err != nil {
			// TODO print error
		}

		if parsed.Execute(out, data) != nil {
			// TODO print error
		}
	}

	g.updateReport()
}

func (g *Generator) updateReport() {

	for i, file := range g.nextReport.Files {
		hash, err := calculateHash(filepath.Join(g.projectPath, file.RealPath))
		if err != nil {
			continue
		}
		g.nextReport.Files[i].Hash = hash
	}

	data, err := yaml.Marshal(g.nextReport)
	if err != nil {
		return
	}
	path := filepath.Join(
		g.projectPath,
		FolderProject,
		g.Name+"_report.yml",
	)

	if os.WriteFile(path, data, 0644) != nil {
		fmt.Printf(
			"Failed to write report for %s generator at: %s\n",
			g.Name,
			path,
		)
	}
}

func (g *Generator) PrintReport() {

	fmt.Printf("======= %s generator report (v%d) =======\n", g.Name, g.nextReport.Version)

	if len(g.lastReport.Files) == 0 {
		fmt.Println("No report for this generator found.")
		return
	}

	if g.isWaitingForGeneration() {
		printLogPrefix(Info, 0)
		fmt.Println(
			" The generator could be run to update files following the new configuration.")
	}

	if !g.isUpToDate() {
		printLogPrefix(Info, 0)
		fmt.Printf(
			" The generator could be run to update files following the new version ("+
				"v%d → v%d).\n",
			g.lastReport.Version,
			g.nextReport.Version,
		)
	}

	fmt.Println()

	for i, file := range g.lastReport.Files {
		hash, err := calculateHash(filepath.Join(g.projectPath, file.RealPath))
		if err != nil {
			printLogPrefix(Deleted, file.Kind)
		} else if hash != g.lastReport.Files[i].Hash {
			printLogPrefix(Edited, file.Kind)
		} else {
			printLogPrefix(Same, file.Kind)
		}
		fmt.Printf(" %s\n", file.RealPath)
	}

}

func (g *Generator) isUpToDate() bool {
	return g.lastReport.Version == g.nextReport.Version
}

func (g *Generator) isWaitingForGeneration() bool {
	return !reflect.DeepEqual(g.nextReport.Config, g.lastReport.Config)
}

func printLogPrefix(kind int8, fileKind int8) {

	var s string

	switch kind {
	case Same:
		switch fileKind {
		case Copy:
			s = "💡 [SAME]"
		case Critical:
			fallthrough
		case Static:
			s = "✅️ [SAME]"
		case Skeleton:
			s = "⚠️ [SAME]"
		}
	case Deleted:
		switch fileKind {
		case Copy:
			s = "💡 [DELETED]"
		case Critical:
			fallthrough
		case Static:
			fallthrough
		case Skeleton:
			s = "❌️ [DELETED]"
		}
	case Edited:
		switch fileKind {
		case Critical:
			s = "⚠️ [EDITED]"
		case Static:
			s = "❌️ [EDITED]"
		case Copy:
			fallthrough
		case Skeleton:
			s = "✅️ [EDITED]"
		}
	case Info:
		s = "💡 [INFO]"
	}

	fmt.Print(s)
}

//region Template helpers

func Upper(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(string(str[0])) + str[1:]
}

func IsNotPlural(str string) bool {
	return !strings.HasSuffix(str, "_plural")
}

func OnlyPrefix(str string) string {
	if len(str) == 0 {
		return str
	}

	if strings.HasSuffix(str, "_plural") {
		return strings.TrimSuffix(str, "_plural")
	}
	if strings.HasSuffix(str, "_singular") {
		return strings.TrimSuffix(str, "_singular")
	}

	return str
}

//endregion

//region File helpers

func calculateHash(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	m := md5.New()
	m.Write(data)
	hash := hex.EncodeToString(m.Sum(nil))

	return hash, nil
}

func copyFile(src, dst string) error {

	sourceFile, err := EmbedFS.Open(src)
	if err != nil {
		return err
	}
	defer func(sourceFile fs.File) {
		_ = sourceFile.Close()
	}(sourceFile)

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destinationFile *os.File) {
		_ = destinationFile.Close()
	}(destinationFile)

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

func copyDir(src, dst string) error {

	f, _ := EmbedFS.Open(src)
	info, err := f.Stat()
	if err != nil {
		return err
	}

	if err = os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}

	entries, err := EmbedFS.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		destinationPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(sourcePath, destinationPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(sourcePath, destinationPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func getLastModifiedTimes(filePaths []string) ([]time.Time, error) {
	var modifiedTimes []time.Time

	for _, filePath := range filePaths {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return nil, err
		}
		modifiedTime := fileInfo.ModTime()
		modifiedTimes = append(modifiedTimes, modifiedTime)
	}

	return modifiedTimes, nil
}

//endregion
