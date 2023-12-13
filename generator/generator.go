package generator

import (
	"bytes"
	"crypto/md5"
	"embed"
	"encoding/hex"
	"fmt"
	"github.com/serenize/snaker"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
	"time"
)

var (
	//go:embed all:template
	EmbedFS embed.FS

	FolderProject  = ".vectra"
	FolderReport   = filepath.Join(FolderProject, "report")
	FolderTemplate = "template"
)

const (
	Same = 1 << iota
	Edited
	Deleted
	Info
)

const (
	Copy     = 1 << iota // Default content could work as is and could be edited.
	CorePart             // For the core of Vectra edit with caution.
	FullGen              // Don't be edited at all.
	Skeleton             // Stub content to be filled or edited later.
)

type Report struct {
	Files   []SourceFile   `yaml:"files"`
	Config  map[string]any `yaml:"config"`
	Version int8           `yaml:"version"`
}

type SourceFile struct {
	RealPath     string `yaml:"path"`
	Hash         string `yaml:"hash"`
	Kind         int8   `yaml:"kind"`
	templatePath string `yaml:"-"`
	isTmpl       bool   `yaml:"-"`
}

func NewSourceFile(path string, kind int8) SourceFile {
	return NewDynSourceFile(path, strings.TrimSuffix(path, ".tmpl"), kind)
}

func NewDynSourceFile(path string, realpath string, kind int8) SourceFile {
	return SourceFile{
		templatePath: filepath.Join(FolderTemplate, path),
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
	IGenerator
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
		FolderReport,
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

func (g *Generator) Generate(ctx any) {

	var ctxs []any
	if ctx != nil && reflect.TypeOf(ctx).Kind() == reflect.Slice {
		sliceValue := reflect.ValueOf(ctx)
		length := sliceValue.Len()
		ctxs = make([]any, length)
		for i := 0; i < length; i++ {
			ctxs[i] = sliceValue.Index(i).Interface()
		}
	}

	nbCtx := len(ctxs)
	if nbCtx > 1 && nbCtx != len(g.nextReport.Files) {
		log.Println("Incoherent number of context: generation cancelled.")
		return
	}

	for i, file := range g.nextReport.Files {

		outputFile := filepath.Join(g.projectPath, file.RealPath)

		dir := filepath.Dir(outputFile)
		if os.MkdirAll(dir, 0744) != nil {
			// TODO print error
			continue
		}

		if !file.isTmpl {
			f, err := EmbedFS.Open(file.templatePath)
			if err != nil {
				log.Println("Failed to handle", file.templatePath)
				continue
			}
			stat, _ := f.Stat()

			if stat.IsDir() {
				err = copyDir(file.templatePath, outputFile)
			} else {
				err = copyFile(file.templatePath, outputFile)
			}
			if err != nil {
				// TODO print error
			}

			continue
		}
		in, err := EmbedFS.ReadFile(file.templatePath)
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
			template.FuncMap{"CamelToSnake": snaker.CamelToSnake},
		).Funcs(
			template.FuncMap{"TrimPluralization": TrimPluralization},
		).Funcs(
			template.FuncMap{"KeyExist": KeyExist},
		).Funcs(
			template.FuncMap{"TrimNewPrefix": TrimNewPrefix},
		).Funcs(
			template.FuncMap{"IsNotPlural": IsNotPlural},
		).Parse(string(in))
		if err != nil {
			// TODO print error
		}

		buf := new(bytes.Buffer)
		if nbCtx == 0 {
			err = parsed.Execute(buf, ctx)
		} else {
			data := ctxs[i]
			err = parsed.Execute(buf, data)
		}
		if err != nil {
			// TODO print error

			continue
		}

		ar := buf.Bytes()
		if strings.HasSuffix(file.RealPath, ".go") && formatGoCode(buf) == nil {
			ar = buf.Bytes()
		}
		_, err = out.Write(ar)
		if err != nil {
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
		FolderReport,
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
	// BUG getAsMap is not recursive: all sub types are not map.
	//  Maybe marshal to yaml before...
	equal := reflect.DeepEqual(g.nextReport.Config, g.lastReport.Config)
	return !equal
}

func printLogPrefix(kind int8, fileKind int8) {

	var s string

	switch kind {
	case Same:
		switch fileKind {
		case Copy:
			s = "💡 [SAME]"
		case CorePart:
			fallthrough
		case FullGen:
			s = "✅️ [SAME]"
		case Skeleton:
			s = "⚠️ [SAME]"
		}
	case Deleted:
		switch fileKind {
		case Copy:
			s = "💡 [DELETED]"
		case CorePart:
			fallthrough
		case FullGen:
			fallthrough
		case Skeleton:
			s = "❌️ [DELETED]"
		}
	case Edited:
		switch fileKind {
		case CorePart:
			s = "⚠️ [EDITED]"
		case FullGen:
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

func TrimPluralization(str string) string {
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

func TrimNewPrefix(str string) string {
	if len(str) < 3 {
		return str
	}

	return strings.TrimPrefix(strings.TrimPrefix(str, "new"), "New")
}

func KeyExist(key string, m map[string]string) bool {
	_, ok := m[key]
	return ok
}

//endregion

//region File helpers

func calculateHash(path string) (string, error) {
	m := md5.New()

	err := filepath.WalkDir(path, func(filePath string, fileInfo os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if fileInfo.IsDir() {
			return nil
		}

		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		m.Write(fileData)
		return nil
	})

	if err != nil {
		return "", err
	}

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

func extractFunctionBody(path string) map[string]string {

	bodiesByName := map[string]string{}

	file, err := os.Open(path)
	if err != nil {
		return map[string]string{}
	}
	defer file.Close()

	// Parse the Go file
	node, err := parser.ParseFile(token.NewFileSet(), "", file, parser.AllErrors)
	if err != nil {
		return map[string]string{}
	}

	// Iterate through the functions and extract their bodies
	for _, function := range node.Decls {
		if function, ok := function.(*ast.FuncDecl); ok {
			bodiesByName[function.Name.Name] = extractPartOfFile(function.Body.Lbrace+1, function.Body.Rbrace-1, file)
		}
	}

	return bodiesByName
}

func extractPartOfFile(start, end token.Pos, file *os.File) string {
	if start >= end {
		return ""
	}
	_, err := file.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Read the file content into a buffer
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	// Extract the specified portion of the file content
	body := buffer[start:end]

	// Convert the formatted body to string and remove leading/trailing whitespace
	return string(body)
}

func formatGoCode(input *bytes.Buffer) error {
	// Parse the input Go code
	set := token.NewFileSet()
	node, err := parser.ParseFile(set, "", input, parser.ParseComments)
	if err != nil {
		log.Printf("Error parsing Go code: %v", err)
		return err
	}

	// Format the parsed node into the input buffer
	input.Reset() // Clear the buffer before writing formatted code
	if err := format.Node(input, set, node); err != nil {
		log.Printf("Error formatting Go code: %v", err)
		return err
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
