package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/urfave/cli/v2"
)

// Read configuration TOML file.
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
