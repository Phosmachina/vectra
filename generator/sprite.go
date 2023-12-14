package generator

import (
	"fmt"
	"github.com/beevik/etree"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/xml"
	"os"
	"path/filepath"
	"strings"
)

type SpriteConfig struct {
	SvgFolderPath   string `yaml:"svg_path"`
	OutputSpriteSvg string `yaml:"output_sprite_svg"`
}

func generateSpriteSvg(cfg *Vectra) {

	var files []string
	root := filepath.Join(cfg.ProjectPath, cfg.SpriteConfig.SvgFolderPath)

	err := filepath.Walk(root, visit(&files))
	if err != nil {
		panic(err)
	}

	// The base of your sprite file
	sprite := etree.NewDocument()
	sprite.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	svg := sprite.CreateElement("svg")
	svg.CreateAttr("xmlns", "http://www.w3.org/2000/svg")
	defs := svg.CreateElement("defs")

	// Pug file base
	mixins := "mixin svg-%s()\n\tsvg(%s viewBox='%s')\n\t\tuse(href='#%s')\n\n"

	// It will store all the mixins
	var pugContent string

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("File reading error", err)
			return
		}

		doc := etree.NewDocument()
		if err := doc.ReadFromBytes(data); err != nil {
			panic(err)
		}

		svg := doc.SelectElement("svg")
		if svg == nil {
			fmt.Println("svg element is not found in file", file)
			continue
		}

		path := strings.TrimSuffix(strings.TrimPrefix(file, root+"/"), ".svg") // extracting svg element id from path
		path = strings.ReplaceAll(path, "/", "-")
		symbol := defs.CreateElement("symbol")
		symbol.CreateAttr("id", path)

		// To unwrap svg tag, copy the child elements directly under symbol
		for _, child := range svg.ChildElements() {
			symbol.AddChild(child.Copy())
		}

		viewBox := svg.SelectAttr("viewBox") // we need viewBox to create mixin
		if viewBox == nil {
			fmt.Println("viewBox attribute is not found in svg", file)
			continue
		}

		var dimensions string // stores width or/and height

		width := svg.SelectAttr("width") // we need width to create mixin
		if width != nil {
			dimensions += "width='" + width.Value + "' "
		}

		height := svg.SelectAttr("height") // we need height to create mixin
		if height != nil {
			dimensions += "height='" + height.Value + "' "
		}

		if dimensions == "" {
			fmt.Println("Both width and height attributes are found in svg", file)
			continue
		}

		pugContent += fmt.Sprintf(mixins, path, strings.TrimSpace(dimensions), viewBox.Value, path)
	}

	err = os.WriteFile(filepath.Join(
		cfg.ProjectPath, "src", "view", "pug", "component", "svg.pug"),
		[]byte(pugContent), 0644)
	if err != nil {
		fmt.Println("Pug file writing error", err)
	}

	// Create a minifier
	m := minify.New()
	m.AddFunc("text/xml", xml.Minify)

	// Minify the XML and handle errors
	b, _ := sprite.WriteToBytes()
	minified, err := m.Bytes("text/xml", b)
	if err != nil {
		fmt.Println("XML Minification error", err)
		return
	}

	// Write the minified output to file
	err = os.WriteFile(filepath.Join(
		cfg.ProjectPath, cfg.SpriteConfig.OutputSpriteSvg),
		minified, 0644)
	if err != nil {
		fmt.Println("Minified Sprite file writing error", err)
	}
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".svg") {
			*files = append(*files, path)
		}
		return nil
	}
}
