package main

import (
	"github.com/Phosmachina/vectra/generator"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {

	var vectra *generator.Vectra

	app := cli.NewApp()
	app.Name = "vectra"
	app.Usage = "Manage Vectra projects"
	app.Version = "1.0.0"

	app.EnableBashCompletion = true

	app.Before = func(c *cli.Context) error {
		path := c.String("path")
		vectra = generator.NewVectra(path)
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "Initialize a folder with the default Vectra project file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "path, p",
					Value: "./",
					Usage: "Path to the directory where the Vectra project file will be initialized",
				},
			},
			Action: func(c *cli.Context) error {
				path := c.String("path")
				log.Println("Initializing Vectra project at", path)
				vectra.Init()
				return nil
			},
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				log.Println("Initializing Vectra project at", vectra.ProjectPath)
				return nil
			},
		},
		{
			Name:  "gen",
			Usage: "Generate Vectra project templates",
			Action: func(c *cli.Context) error {
				log.Println("Generate all parts of Vectra.")
				vectra.FullGenerate()
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:  "i18n",
					Usage: "Generate the i18n part of the Vectra project",
					Action: func(c *cli.Context) error {
						log.Println("Generating i18n template.")
						vectra.Generate("i18n")
						return nil
					},
				},
				{
					Name:  "controllers",
					Usage: "Generate controllers part of the Vectra project",
					Action: func(c *cli.Context) error {
						log.Println("Generating controllers template.")
						vectra.Generate("controllers")
						return nil
					},
				},
				{
					Name:  "base",
					Usage: "Generate base part of the Vectra project",
					Action: func(c *cli.Context) error {
						log.Println("Copy base parts following configuration.")
						vectra.Generate("base")
						return nil
					},
				},
				{
					Name:  "types",
					Usage: "Generate types part of the Vectra project",
					Action: func(c *cli.Context) error {
						log.Println("Generate types parts following configuration.")
						vectra.Generate("types")
						return nil
					},
				},
				{
					Name:  "services",
					Usage: "Generate services part of the Vectra project",
					Action: func(c *cli.Context) error {
						log.Println("Generate services parts following configuration.")
						vectra.Generate("services")
						return nil
					},
				},
				// Add more subcommands as needed
			},
		},
		{
			Name:  "watch",
			Usage: "Survey Sass, Pug files of Vectra project and execute pipeline.",
			Action: func(c *cli.Context) error {
				vectra.Watch()
				return nil
			},
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "path, p",
			Value: "./",
			Usage: "Path to the Vectra project file or directory",
		},
	}

	app.Action = func(c *cli.Context) error {
		path := c.String("path")
		log.Println("Summarizing the state of deployment for Vectra project at", path)
		vectra.FullReport()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Println(err)
	}
}
