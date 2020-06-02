package remote

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	testCases := map[string]Host{
		"unit.test": {
			Host: "unit.test",
		},
		"port.test": {
			Host: "port.test",
			Port: 233,
		},
		"user@unit.test": {
			Host: "unit.test",
			User: "user",
		},
		"user@port.test": {
			Host: "port.test",
			User: "user",
			Port: 233,
		},
	}

	for expected, unit := range testCases {
		assert.Equal(t, expected, unit.String(), "string rendering")
	}
}

func TestCommandlineArgs(t *testing.T) {
	testCases := map[string]Host{
		"unit.test": {
			Host: "unit.test",
		},
		"-p 233 unit.test": {
			Host: "unit.test",
			Port: 233,
		},
		"-J proxy.test unit.test": {
			Host:     "unit.test",
			JumpHost: "proxy.test",
		},
		"-J user@proxy.test unit.test": {
			Host:     "unit.test",
			JumpHost: "proxy.test",
			JumpUser: "user",
		},
		"-J proxy.test:233 unit.test": {
			Host:     "unit.test",
			JumpHost: "proxy.test",
			JumpPort: 233,
		},
		"-J user@proxy.test:233 unit.test": {
			Host:     "unit.test",
			JumpUser: "user",
			JumpHost: "proxy.test",
			JumpPort: 233,
		},
		"-i /tmp/id_rsa unit.test": {
			Host:    "unit.test",
			KeyFile: "/tmp/id_rsa",
		},
		"-o UserKnownHostsFile=/dev/null -o GlobalKnownHostsFile=/dev/null unit.test": {
			Host: "unit.test",
			Options: []string{
				"UserKnownHostsFile=/dev/null",
				"GlobalKnownHostsFile=/dev/null",
			},
		},
	}

	for expected, unit := range testCases {
		argv := strings.Join(unit.CommandlineArgs(true), " ")

		assert.Equal(t, expected, argv, "command line argument rendering")
	}
}

func TestJumpProxy(t *testing.T) {
	testCases := map[string]Host{
		"": {
			Host: "unit.test",
		},
		"jump.host": {
			Host:     "unit.test",
			JumpHost: "jump.host",
		},
		"jump.host:2200": {
			Host:     "unit.test",
			JumpHost: "jump.host",
			JumpPort: 2200,
		},
		"user@jump.host": {
			Host:     "unit.test",
			JumpUser: "user",
			JumpHost: "jump.host",
		},
		"user@jump.host:2200": {
			Host:     "unit.test",
			JumpUser: "user",
			JumpHost: "jump.host",
			JumpPort: 2200,
		},
	}

	for expected, unit := range testCases {
		assert.Equal(t, expected, unit.JumpProxy(), "jump host rendering")
	}
}

func TestParseHostError(t *testing.T) {
	testCases := []string{
		"localhost:",
		"[::1%2]:",
		"[::1%2",
		"user@",
		"user@:23",
		"port:IO",
		"",
	}

	for _, input := range testCases {
		if got, err := ParseHost(input); err == nil {
			t.Fatalf("expected error, got '%v'", got)
		}
	}
}

func TestParseHost(t *testing.T) {
	testCases := map[string]Host{
		"localhost": {
			Host: "localhost",
		},
		"only.host": {
			Host: "only.host",
		},
		"host-with.port:23": {
			Host: "host-with.port",
			Port: 23,
		},
		"127.0.0.1": {
			Host: "127.0.0.1",
		},
		"127.1.2.3:23": {
			Host: "127.1.2.3",
			Port: 23,
		},
		"[::1%2]": {
			Host: "::1%2",
		},
		"[2001:0db8:85a3:0000:0000:8a2e:0370:7334%2]:23": {
			Host: "2001:0db8:85a3:0000:0000:8a2e:0370:7334%2",
			Port: 23,
		},
		"user@localhost": {
			User: "user",
			Host: "localhost",
		},
		"user@user.host": {
			User: "user",
			Host: "user.host",
		},
		"user@user-host.port:23": {
			User: "user",
			Host: "user-host.port",
			Port: 23,
		},
		"user@127.2.4.6": {
			User: "user",
			Host: "127.2.4.6",
		},
		"user@127.127.127.127:23": {
			User: "user",
			Host: "127.127.127.127",
			Port: 23,
		},
		"user@[::ffff%2]": {
			User: "user",
			Host: "::ffff%2",
		},
		"user@[::f00d%2]:23": {
			User: "user",
			Host: "::f00d%2",
			Port: 23,
		},
	}

	for input, want := range testCases {
		got, err := ParseHost(input)

		if err != nil {
			t.Fatalf("ParseHost: %v", err)
		}

		compare(t, &want, got)
	}
}

func compare(t *testing.T, want, got *Host) {
	if want.User != got.User {
		t.Errorf("got User=%q; want %q", got.User, want.User)
	}

	if want.Port != got.Port {
		t.Errorf("got Port=%q; want %q", got.Port, want.Port)
	}

	if want.Host != got.Host {
		t.Errorf("got Host=%q; want %q", got.Host, want.Host)
	}

	if want.JumpUser != got.JumpUser {
		t.Errorf("got JumpUser=%q; want %q", got.JumpUser, want.JumpUser)
	}

	if want.JumpHost != got.JumpHost {
		t.Errorf("got JumpHost=%q; want %q", got.JumpHost, want.JumpHost)
	}

	if want.JumpPort != got.JumpPort {
		t.Errorf("got JumpPort=%q; want %q", got.JumpPort, want.JumpPort)
	}

	if want.KeyFile != got.KeyFile {
		t.Errorf("got KeyFile=%q; want %q", got.KeyFile, want.KeyFile)
	}

	if len(want.Options) != len(got.Options) {
		t.Fatalf("got %d options; want %d options", len(got.Options), len(want.Options))
	}

	for i, actual := range got.Options {
		if want.Options[i] != actual {
			t.Errorf("got Option=%q; want %q", actual, want.Options[i])
		}
	}
}
