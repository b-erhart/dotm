package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name:  "dotm",
		Usage: "an incredibly simple dotfiles manager",
		Authors: []*cli.Author{
			{
				Name:  "Benedikt Erhart",
				Email: "contact@b-erhart.de",
			},
		},
		Suggest:         true,
		HideHelpCommand: true,
		Commands: []*cli.Command{
			{
				Name:      "distribute",
				Usage:     "copy dotfiles to specified locations",
				ArgsUsage: "CONFIG_FILE",
				Action:    distribute,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "overwrite",
						Aliases: []string{"o"},
						Usage:   "skip confirmation for overwriting existing destinations",
					},
				},
				HideHelpCommand: true,
			},
			{
				Name:      "fetch",
				Usage:     "copy dotfiles from specified locations",
				ArgsUsage: "CONFIG_FILE",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "overwrite",
						Aliases: []string{"o"},
						Usage:   "skip confirmation for overwriting existing destinations",
					},
				},
				Action:          fetch,
				HideHelpCommand: true,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
