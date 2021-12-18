package path

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	ErrTargetBlank = errors.New("target cannot be blank")
)

// ExpandAndValidate expands a directory, resolving ~ and relative references, and then validates that the directory
// exists. It returns the absolute path and any errors.
func ExpandAndValidate(path string) (string, error) {
	if path == "" {
		return "", ErrTargetBlank
	}

	usr, err := user.Current()
	home := usr.HomeDir

	// deal with ~
	if path == "~" {
		if err != nil {
			return "", err
		}
		// In case of "~", which won't be caught by the "else if"
		path = home
	} else if strings.HasPrefix(path, "~/") {
		if err != nil {
			return "", err
		}

		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		path = filepath.Join(home, path[2:])
	}

	// resolves '..', '.', and relative paths
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	file, err := os.Stat(abs)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("target directory does not exist")
	}

	if !file.IsDir() {
		return "", fmt.Errorf("target is not a directory")
	}

	return abs, nil
}
