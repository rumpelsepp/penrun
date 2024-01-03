// SPDX-FileCopyrightText: Stefan Tatschner
//
// SPDX-License-Identifier: MIT

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/alecthomas/kong-toml"
	"github.com/rumpelsepp/penrun/internal"
)

func main() {
	var (
		config penrun.Config
		ctx    *kong.Context
	)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		// TODO: Use exit codes from FreeBSD
		os.Exit(1)
	}

	configPath, err := penrun.FindConfigFile(cwd)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintln(os.Stderr, err)
			// TODO: Use exit codes from FreeBSD
			os.Exit(1)
		}
	}

	if configPath != "" {
		ctx = kong.Parse(&config, kong.Configuration(kongtoml.Loader, configPath))
	} else {
		ctx = kong.Parse(&config)
	}

	fr, err := penrun.RunCommand(&config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		// TODO: Use exit codes from FreeBSD
		ctx.Exit(1)
	}

	os.Exit(fr.ExitCode)
}
