package provider

import (
	"path/filepath"
	"testing"

	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
)

func TestParseBasic(t *testing.T) {
	unit := newUnit("basic")
	want := []remote.Host{
		remote.Host {
			Host: "ldn01.example.net",
		},
		remote.Host {
			Host: "127.162.222.67",
		},
		remote.Host {
			Host: "nyc01.example.org",
		},
		remote.Host {
			Host: "127.230.83.95",
		},
		remote.Host {
			Host: "syd01.example.com",
		},
		remote.Host {
			Host: "127.124.68.128",
		},
	}

	testParse(t, unit, want)
}

func TestParseVariants(t *testing.T) {
	unit := newUnit("variants")
	want := []remote.Host{
		remote.Host{
			Host: "regular.hostname",
		},
		remote.Host{
			Host: "127.0.0.1",
		},
		remote.Host{
			Host: "127.0.0.2",
		},
		remote.Host{
			Host: "only.hostname",
		},
		remote.Host{
			Host: "port.hostname",
			Port: 22,
		},
		remote.Host{
			Host: "127.0.0.3",
			Port: 22,
		},
		remote.Host{
			Host: "127.0.0.4",
			Port: 22,
		},
		remote.Host{
			Host: "only.port.hostname",
			Port: 22,
		},
		remote.Host{
			Host: "::1%2",
			Port: 22,
		},
		remote.Host{
			Host: "fe80::200:5eff:fe00:5342%3",
		},
	}

	testParse(t, unit, want)
}

func TestParseDuplicates(t *testing.T) {
	unit := newUnit("duplicates")
	want := []remote.Host{
		remote.Host{
			Host: "node01",
		},
		remote.Host{
			Host: "node01.example.org",
		},
		remote.Host{
			Host: "node01",
		},
		remote.Host{
			Host: "node01.example.com",
		},
		remote.Host{
			Host: "node02",
		},
		remote.Host{
			Host: "node02",
		},
		remote.Host{
			Host: "node02",
		},
		remote.Host{
			Host: "node03",
		},
		remote.Host{
			Host: "node03.example.com",
		},
	}

	testParse(t, unit, want)
}


func testParse(t *testing.T, unit *KnownHostsProvider, want []remote.Host) {
	got, err := unit.Parse()

	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("got %d elements; want %d elements", len(got), len(want))
	}

	for i, actual := range got {
		compare(t, &want[i], &actual)
	}
}

func newUnit(fixture string) *KnownHostsProvider {
	fixture = filepath.Join("testdata", fixture + ".known_hosts")

	return NewKnownHostsProvider(fixture, true)
}

func compare(t *testing.T, want, got *remote.Host) {
	if want.User != got.User {
		t.Errorf("got User=%q; want %q", got.User, want.User)
	}

	if want.Port != got.Port {
		t.Errorf("got Port=%q; want %q", got.Port, want.Port)
	}

	if want.Host != got.Host {
		t.Errorf("got Host=%q; want %q", got.Host, want.Host)
	}
}
