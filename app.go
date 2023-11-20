package main

import (
	"fmt"
	"github.com/Phosmachina/vectra/generator"
	"github.com/urfave/cli"
	"os"
	"strings"
)

func main() {

	var vectra *generator.Vectra

	app := cli.NewApp()
	app.Name = "vectra"
	app.Usage = "Manage Vectra projects: initialize projet, report current state, " +
		"and generate files following configuration"
	app.Version = "1.0.1"

	app.EnableBashCompletion = true

	app.Before = func(c *cli.Context) error {
		path := c.String("path")
		vectra = generator.NewVectra(path)
		c.App.Metadata = map[string]any{
			"select": c.String("select"),
		}
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "Initialize a folder with the default Vectra project file",
			Action: func(c *cli.Context) error {
				fmt.Println("Initializing Vectra project at", vectra.ProjectPath)
				vectra.Init()
				return nil
			},
		},
		{
			Name:  "gen",
			Usage: "Generate Vectra project templates",
			Flags: []cli.Flag{},
			Action: func(c *cli.Context) error {
				generators := strings.Split(c.App.Metadata["select"].(string), ",")
				if len(generators) == 1 && generators[0] == "" {
					fmt.Println("🔧 Generating all templates available.")
					vectra.FullGenerate()
				} else {
					for _, s := range generators {
						fmt.Println("🔧 Generating", s, "template.")
						vectra.Generate(s)
					}
				}
				return nil
			},
		},
		{
			Name: "watch",
			Usage: "Survey Sass, Pug files of Vectra project and execute pipeline" +
				"Check Docker, build Docker images and containers if not exist",
			Action: func(c *cli.Context) error {
				vectra.Watch()
				return nil
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path, p",
			Usage: "Path to the Vectra project directory.",
		},
		cli.StringFlag{
			Name: "select, s",
			Usage: "List of generator name separated by comma. " +
				"Empty value run all generators. (e.g.: services,controllers). " +
				"Available generators: base, types, services, controllers, " +
				"i18n (managed by watcher)",
		},
	}

	app.Action = func(c *cli.Context) error {
		generators := strings.Split(c.String("select"), ",")
		if len(generators) == 1 && generators[0] == "" {
			fmt.Println("Summarizing the state of deployment for Vectra project at",
				vectra.ProjectPath)
			vectra.FullReport()
		} else {
			for _, s := range generators {
				vectra.Report(s)
			}
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
