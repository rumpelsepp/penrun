// SPDX-FileCopyrightText: Stefan Tatschner
//
// SPDX-License-Identifier: MIT

package penrun

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Extracted from: https://github.com/go-git/go-git/issues/74

func isGitDir(path string) (bool, error) {
	markers := []string{"HEAD", "objects", "refs"}

	for _, marker := range markers {
		_, err := os.Stat(filepath.Join(path, marker))
		if err == nil {
			continue
		}
		if !errors.Is(err, os.ErrNotExist) {
			// unknown error
			return false, err
		} else {
			return false, nil
		}
	}

	return true, nil
}

func detectGitPath(path string) (string, error) {
	// normalize the path
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	for {
		fi, err := os.Stat(filepath.Join(path, ".git"))
		if err == nil {
			if !fi.IsDir() {
				return "", fmt.Errorf(".git exist but is not a directory")
			}
			return filepath.Join(path, ".git"), nil
		}
		if !os.IsNotExist(err) {
			// unknown error
			return "", err
		}

		// detect bare repo
		ok, err := isGitDir(path)
		if err != nil {
			return "", err
		}
		if ok {
			return path, nil
		}

		if parent := filepath.Dir(path); parent == path {
			return "", fmt.Errorf(".git not found")
		} else {
			path = parent
		}
	}
}

func detectGitRoot(path string) (string, error) {
	p, err := detectGitPath(".")
	if err != nil {
		return "", err
	}
	return filepath.Dir(p), nil
}
