package util

import (
	"os"
	"path/filepath"
)

func ResolvePath(elem ...string) (path string, err error) {
	path, err = filepath.Abs(filepath.Join(elem...))

	if err != nil {
		return "", err
	}

	_, err = os.Stat(path)

	if err != nil {
		return "", err
	}

	return path, nil
}

func UserFilePath(configDir bool, elems ...string) (config string, err error) {
	if configDir {
		config, err = os.UserConfigDir()
	} else {
		config, err = os.UserHomeDir()
	}

	if err != nil {
		return "", err
	}

	elements := []string{config}
	elements = append(elements, elems...)

	return ResolvePath(elements...)
}
