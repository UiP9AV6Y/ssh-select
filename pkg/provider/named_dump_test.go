package provider

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
)

func TestParseBasicNamed(t *testing.T) {
	unit := newNamedDumpProviderUnit("basic", false)
	want := []remote.Host{
		{
			Host: "example.com",
		},
		{
			Host: "ns01.example.com",
		},
		{
			Host: "kubernetes.example.com",
		},
		{
			Host: "www.example.com",
		},
		{
			Host: "example.com",
		},
		{
			Host: "ns01.example.com",
		},
		{
			Host: "10.1.1.1",
		},
		{
			Host: "mail.example.com",
		},
		{
			Host: "10.1.1.1",
		},
		{
			Host: "example.com",
		},
		{
			Host: "mail.example.com",
		},
		{
			Host: "node01.example.com",
		},
	}

	testParse(t, unit, want)
}

func TestParseNamedWildcard(t *testing.T) {
	unit := newNamedDumpProviderUnit("basic", true)
	want := remote.Host{
		Host: "wildcard.services.example.com",
	}
	got, err := unit.Parse()

	assert.NoError(t, err, "parsing succeeded")
	assert.Contains(t, got, want, "wildcard entry has been processed")
}

func newNamedDumpProviderUnit(fixture string, processWildcards bool) *NamedDumpProvider {
	fixture = filepath.Join("testdata", fixture+".named")

	return NewNamedDumpProvider(fixture, processWildcards)
}
