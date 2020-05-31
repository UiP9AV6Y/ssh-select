package cli

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplication(t *testing.T) {
	unit := NewParser("/usr/bin/test")

	assert.Equal(t, "test", unit.Application, "sanitized application name")
}

func TestParseArgvMissingArg(t *testing.T) {
	unit := NewParser("test")
	input := []string{
		"--ssh",
		"--version",
	}
	err := unit.ParseArgv(input)

	assert.EqualError(t, err, "missing argument for SSH binary")
}

func TestParseArgvEmpty(t *testing.T) {
	unit := NewParser("test")
	input := []string{}
	err := unit.ParseArgv(input)

	assert.Nil(t, err, "parsing without errors")
	assert.Equal(t, 0, len(unit.SshArgv), "SSH passthough arguments")
}

func TestParseArgvSettings(t *testing.T) {
	unit := NewParser("test")
	input := []string{
		"--known-hosts", "/tmp/known_hosts",
		"--ssh", "/usr/local/bin/rsh",
		"--version",
	}
	err := unit.ParseArgv(input)

	assert.Nil(t, err, "parsing without errors")
	assert.True(t, unit.Version, "version switch triggered")
	assert.Equal(t, "/usr/local/bin/rsh", unit.SshBinary, "custom SSH path")
	assert.Subset(t, unit.KnownHostsFiles, []string{"/tmp/known_hosts"}, "custom known hosts file")
}

func TestParseArgvSsh(t *testing.T) {
	unit := NewParser("test")
	input := []string{
		"--ssh", "/usr/local/bin/rsh",
	}
	err := unit.ParseArgv(input)

	assert.Nil(t, err, "parsing without errors")
	assert.Equal(t, "/usr/local/bin/rsh", unit.SshBinary, "custom SSH path")
}

func TestParseArgvHost(t *testing.T) {
	unit := NewParser("test")
	input := []string{
		"remote.host.test",
	}
	err := unit.ParseArgv(input)

	assert.Nil(t, err, "parsing without errors")
	assert.Equal(t, 1, len(unit.SshArgv), "SSH passthough arguments")
}

func TestParseArgvHostArgv(t *testing.T) {
	unit := NewParser("test")
	input := []string{
		"remote.host.test",
		"uname", "-a",
	}
	err := unit.ParseArgv(input)

	assert.Nil(t, err, "parsing without errors")
	assert.Equal(t, 3, len(unit.SshArgv), "SSH passthough arguments")
}

func TestParseArgvAll(t *testing.T) {
	unit := NewParser("test")
	input := []string{
		"--known-hosts", "/tmp/known_hosts",
		"-i", "/tmp/id_rsa",
		"remote.host.test",
		"uname", "-a",
	}
	err := unit.ParseArgv(input)

	assert.Nil(t, err, "parsing without errors")
	assert.Subset(t, unit.KnownHostsFiles, []string{"/tmp/known_hosts"}, "custom known hosts file")
	assert.Equal(t, 5, len(unit.SshArgv), "SSH passthough arguments")
}

func TestParseEnvMissingValue(t *testing.T) {
	unit := NewParser("test")
	input := []string{
		"SSH_SELECT_SSH_BINARY=",
	}
	err := unit.ParseEnv(input)

	assert.EqualError(t, err, "env variable SSH_SELECT_SSH_BINARY must not be empty")
}

func TestParseEnv(t *testing.T) {
	unit := NewParser("test")
	input := []string{
		"SSH_SELECT_NO_SEARCH_KNOWN_HOSTS=yes",
		"SSH_SELECT_SSH_BINARY=/opt/bin/ssh",
		"SSH_SELECT_KNOWN_HOSTS_FILE_test=/tmp/known_hosts",
	}
	err := unit.ParseEnv(input)

	assert.Nil(t, err, "parsing without errors")
	assert.True(t, unit.NoSearchKnownHosts, "known_hosts search is disabled")
	assert.Equal(t, "/opt/bin/ssh", unit.SshBinary, "custom SSH path")
	assert.Subset(t, unit.KnownHostsFiles, []string{"/tmp/known_hosts"}, "custom known hosts file")
}
