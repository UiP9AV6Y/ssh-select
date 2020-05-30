package version

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	version = "1.2.3"
	commit = "mock"
	date = "1970-01-01T00:00:00Z00:00"
}

func TestApplication(t *testing.T) {
	actual := Application("test")
	expected := "test (1.2.3-mock) [1970-01-01T00:00:00Z00:00]"

	assert.Equal(t, expected, actual)
}
