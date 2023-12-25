package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"github.com/u-root/u-root/pkg/cp"
	"github.com/urfave/cli/v2"
)

func readConfig(configPath string) (map[string]string, error) {
	configFile, err := os.Open(configPath)

	if err != nil {
		return nil, cli.Exit(fmt.Sprintf("Unable to read config file. Abort.\n  > %v", err), 20)
	}

	var config map[string]string
	decoder := toml.NewDecoder(configFile)
	err = decoder.Decode(&config)

	if err != nil {
		var decodeErr *toml.DecodeError
		if errors.As(err, &decodeErr) {
			return nil, cli.Exit(fmt.Sprintf("Unable to read toml data.\n\n%v", decodeErr.String()), 20)
		}

		return nil, cli.Exit(fmt.Sprintf("Unable to read toml data.\n  > %v", err), 20)
	}

	return config, nil
}

func copy(src, dest string) bool {
	success := true

	src = os.ExpandEnv(src)
	dest = os.ExpandEnv(dest)
	fmt.Printf("Copying \"%s\" to \"%s\"...\n", src, dest)

	opts := cp.Options{
		NoFollowSymlinks: true,
		PreCallback: func(_, dest string, srcfi os.FileInfo) error {
			if _, err := os.Stat(dest); os.IsNotExist(err) {
				return nil
			}

			fmt.Print("Destination already exists. Deleting...")
			err := os.RemoveAll(dest)

			if err != nil {
				fmt.Printf("\nError - %v\n", err)
				success = false
				return cp.ErrSkip
			}

			fmt.Println(" Done.")

			return nil
		},
	}

	err := opts.CopyTree(src, dest)

	if err != nil {
		success = false
		fmt.Printf("Error - %v\n", err)
	} else {
		fmt.Println("Done.")
	}

	fmt.Println()

	return success
}

func distribute(ctx *cli.Context) error {
	configPath := ctx.Args().Get(0)

	if configPath == "" {
		return cli.Exit("No config file was provided. Abort.", 10)
	}

	if (!ctx.Bool("overwrite")) {
		fmt.Print("Warning: Existing destinations will be deleted and replaced. Continue (y/n)? ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')

		if err != nil {
			return cli.Exit(fmt.Sprintf("\nUnable to process input - %v", err), 30)
		}

		input = strings.Replace(input, "\r\n", "", -1)
		input = strings.Replace(input, "\n", "", -1)

		if input != "y" && input != "yes" {
			return cli.Exit("Aborting.", 0)
		}
	}

	config, err := readConfig(configPath)

	if err != nil {
		return err
	}

	errCount := 0
	for src, dest := range config {
		if !copy(src, dest) {
			errCount += 1
		}
	}

	fmt.Printf("%d error(s) occured.", errCount)

	return nil
}

func fetch(ctx *cli.Context) error {
	configPath := ctx.Args().Get(0)

	if configPath == "" {
		return cli.Exit("No config file was provided. Abort.", 10)
	}

	if (!ctx.Bool("overwrite")) {
		fmt.Print("Warning: Existing destinations will be deleted and replaced. Continue (y/n)? ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')

		if err != nil {
			return cli.Exit(fmt.Sprintf("Unable to process input - %v", err), 30)
		}

		input = strings.Replace(input, "\r\n", "", -1)
		input = strings.Replace(input, "\n", "", -1)

		if input != "y" && input != "yes" {
			return cli.Exit("Aborting.", 0)
		}
	}

	config, err := readConfig(configPath)

	if err != nil {
		return err
	}

	errCount := 0
	for src, dest := range config {
		if !copy(dest, src) {
			errCount += 1
		}
	}

	fmt.Printf("%d error(s) occurred.", errCount)

	return nil
}

func main() {
	distr := cli.Command{
		Name:      "distribute",
		Usage:     "copy dotfiles to specified locations",
		ArgsUsage: "CONFIG_FILE",
		Action:    distribute,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "overwrite",
				Aliases: []string{"o"},
				Usage: "skip confirmation for overwriting existing destinations",
			},
		},
		HideHelpCommand: true,
	}

	ft := cli.Command{
		Name:      "fetch",
		Usage:     "copy dotfiles from specified locations",
		ArgsUsage: "CONFIG_FILE",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name: "overwrite",
				Aliases: []string{"o"},
				Usage: "skip confirmation for overwriting existing destinations",
			},
		},
		Action:    fetch,
		HideHelpCommand: true,
	}

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
			&distr,
			&ft,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
