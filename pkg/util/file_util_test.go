package util

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePathFile(t *testing.T) {
	cwd, _ := os.Getwd()
	actual, err := ResolvePath("testdata", "file.txt")
	expected := filepath.Join(cwd, "testdata", "file.txt")

	assert.Nil(t, err, "File found without error")
	assert.Equal(t, expected, actual)
}

func TestResolvePathError(t *testing.T) {
	_, err := ResolvePath("testdata", "not-found.json")

	assert.NotNil(t, err, "File not found")
}
