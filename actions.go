package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/u-root/u-root/pkg/cp"
	"github.com/urfave/cli/v2"
)

// Copy dotfiles from repository to system locations.
func distribute(ctx *cli.Context) error {
	return copyConfigEntries(ctx, false)
}

// Copy dotfiles from system locations to repository.
func fetch(ctx *cli.Context) error {
	return copyConfigEntries(ctx, true)
}

// Copy dotfiles specified in config file. If swapSrcDest is false, files are copied from repository to system. Otherwise files are copied in the reverse direction.
func copyConfigEntries(ctx *cli.Context, swapSrcDest bool) error {
	configPath := ctx.Args().Get(0)

	if configPath == "" {
		return cli.Exit("No config file was provided. Abort.", 10)
	}

	err := checkOverwriteConfirmation(ctx)

	if err != nil {
		return err
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

	switch errCount {
	case 0:
		fmt.Println("Success.")
	case 1:
		fmt.Printf("%d error occured.\n", errCount)
	default:
		fmt.Printf("%d errors occured.\n", errCount)
	}

	return nil
}

// Check if user confirmed that overwriting existing data at copy destinations is fine.
func checkOverwriteConfirmation(ctx *cli.Context) error {
	if !ctx.Bool("overwrite") {
		fmt.Print("Warning: Existing destinations will be deleted and replaced. Continue (y/n)? ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')

		if err != nil {
			return cli.Exit(fmt.Sprintf("\nUnable to process input - %v", err), 30)
		}

		input = strings.ReplaceAll(input, "\r\n", "")
		input = strings.ReplaceAll(input, "\n", "")

		if input != "y" && input != "yes" {
			return cli.Exit("Aborting.", 0)
		}
	}

	return nil
}

// If path starts with "~", it is replaced with the path to the home directory.
func expandHomeDir(path string) (string, error) {
	path = filepath.FromSlash(path)
	prefix := "~" + string(os.PathSeparator)

	if !strings.HasPrefix(path, prefix) {
		return path, nil
	}

	homeVar := "HOME"
	if runtime.GOOS == "windows" {
		homeVar = "HOMEPATH"
	}

	envHome, envHomeSet := os.LookupEnv(homeVar)
	if envHomeSet && envHome != "" {
		path = strings.Replace(path, prefix, envHome+string(os.PathSeparator), 1)
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return path, err
	}

	path = strings.Replace(path, prefix, usr.HomeDir+string(os.PathSeparator), 1)

	return path, nil
}

// Copy a file or directory from src to dest.
func copy(src, dest string) bool {
	success := true

	src, err := expandHomeDir(os.ExpandEnv(src))
	if err != nil {
		fmt.Printf("Error - unable to expand \"~\" in \"%s\" (%v)\n", src, err)
		return false
	}

	dest, err = expandHomeDir(os.ExpandEnv(dest))
	if err != nil {
		fmt.Printf("Error - unable to expand \"~\" in \"%s\" (%v)\n", dest, err)
		return false
	}

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

	err = opts.CopyTree(src, dest)

	if err != nil {
		success = false
		fmt.Printf("Error - %v\n", err)
	} else {
		fmt.Println("Done.")
	}

	fmt.Println()

	return success
}
