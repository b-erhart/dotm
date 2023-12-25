package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/u-root/u-root/pkg/cp"
	"github.com/urfave/cli/v2"
)

func distribute(ctx *cli.Context) error {
	return copyConfigEntries(ctx, false)
}

func fetch(ctx *cli.Context) error {
	return copyConfigEntries(ctx, true)
}

func copyConfigEntries(ctx *cli.Context, swapSrcDest bool) error {
	configPath := ctx.Args().Get(0)

	if configPath == "" {
		return cli.Exit("No config file was provided. Abort.", 10)
	}

	if !ctx.Bool("overwrite") {
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
		if swapSrcDest {
			src, dest = dest, src
		}

		if !copy(src, dest) {
			errCount += 1
		}
	}

	fmt.Printf("%d error(s) occured.", errCount)

	return nil
}

func copy(src, dest string) bool {
	success := true

	src = filepath.FromSlash(os.ExpandEnv(src))
	dest = filepath.FromSlash(os.ExpandEnv(dest))
	fmt.Printf("Copying \"%s\" to \"%s\"...\n", src, dest)

	opts := cp.Options{
		NoFollowSymlinks: true,
		PreCallback: func(_, dest string, srcfi os.FileInfo) error {
			parentDir := filepath.Dir(dest)

			if _, err := os.Stat(parentDir); os.IsNotExist(err) {
				fmt.Println("Parent directory of destination does not exist. Creating...")
				if err := os.MkdirAll(parentDir, os.ModePerm); err != nil {
					fmt.Printf("\nError - %v\n", err)
					success = false
					return cp.ErrSkip
				}
			}

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
