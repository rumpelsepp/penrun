// SPDX-FileCopyrightText: Stefan Tatschner
//
// SPDX-License-Identifier: MIT

package penrun

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	CLI           []string `kong:"arg,passthrough,help='Command to invoke'"`
	ArtifactsDir  string   `kong:"short='d',help='Set artifacts_dir to this directory'"`
	ArtifactsBase string   `kong:"short='b',help='Set artifacts_base where the artifacts hierarchy (including LATEST) will be created'"`
	Batched       bool     `kong:"short='B',help='Run several commands in parallel'"`
	PreHook       string   `kong:"help='Shell script that runs before the command'"`
	PostHook      string   `kong:"help='Shell script that runs after the command'"`
}

const ConfigFileName = "penrun.toml"

func fileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

func FindConfigFile(path string) (string, error) {
	if p := os.Getenv("PENRUN_CONFIG"); p != "" {
		return p, nil
	}

	confPath := filepath.Join(path, ConfigFileName)
	exists, err := fileExists(confPath)
	if err != nil {
		return "", err
	}
	if exists {
		return confPath, nil
	}

	if gitRoot, err := detectGitRoot(path); err != nil {
		confPath := filepath.Join(gitRoot, ConfigFileName)
		exists, err := fileExists(confPath)
		if err != nil {
			return "", err
		}
		if exists {
			return confPath, nil
		}
	}

	confDir, err := os.UserConfigDir()
	if err == nil {
		confPath := filepath.Join(confDir, "penrun", ConfigFileName)
		exists, err := fileExists(confPath)
		if err != nil {
			return "", err
		}
		if exists {
			return confPath, nil
		}
	}

	confPath = filepath.Join("/etc/penrun", "penrun.toml")
	exists, err = fileExists(confPath)
	if err != nil {
		return "", err
	}
	if exists {
		return confPath, nil
	}

	return "", os.ErrNotExist
}

func ParseConfigFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
