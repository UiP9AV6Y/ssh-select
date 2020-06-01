package provider

import (
	"path/filepath"
	"testing"

	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
)

func TestParseBasicKnownHosts(t *testing.T) {
	unit := newKnownHostsProviderUnit("basic")
	want := []remote.Host{
		{
			Host: "ldn01.example.net",
		},
		{
			Host: "127.162.222.67",
		},
		{
			Host: "nyc01.example.org",
		},
		{
			Host: "127.230.83.95",
		},
		{
			Host: "syd01.example.com",
		},
		{
			Host: "127.124.68.128",
		},
	}

	testParse(t, unit, want)
}

func TestParseVariantsKnownHosts(t *testing.T) {
	unit := newKnownHostsProviderUnit("variants")
	want := []remote.Host{
		{
			Host: "regular.hostname",
		},
		{
			Host: "127.0.0.1",
		},
		{
			Host: "127.0.0.2",
		},
		{
			Host: "only.hostname",
		},
		{
			Host: "port.hostname",
			Port: 22,
		},
		{
			Host: "127.0.0.3",
			Port: 22,
		},
		{
			Host: "127.0.0.4",
			Port: 22,
		},
		{
			Host: "only.port.hostname",
			Port: 22,
		},
		{
			Host: "::1",
			Port: 22,
		},
		{
			Host: "fe80::200:5eff:fe00:5342",
		},
	}

	testParse(t, unit, want)
}

func TestParseDuplicatesKnownHosts(t *testing.T) {
	unit := newKnownHostsProviderUnit("duplicates")
	want := []remote.Host{
		{
			Host: "node01",
		},
		{
			Host: "node01.example.org",
		},
		{
			Host: "node01",
		},
		{
			Host: "node01.example.com",
		},
		{
			Host: "node02",
		},
		{
			Host: "node02",
		},
		{
			Host: "node02",
		},
		{
			Host: "node03",
		},
		{
			Host: "node03.example.com",
		},
	}

	testParse(t, unit, want)
}

func TestParsePuppetKnownHosts(t *testing.T) {
	unit := newKnownHostsProviderUnit("puppet")
	want := []remote.Host{
		{
			Host: "bastion-ed25519",
		},
		{
			Host: "bastion-ed25519.example.com",
		},
		{
			Host: "127.193.74.241",
		},
		{
			Host: "bastion-ecdsa",
		},
		{
			Host: "bastion-ecdsa.example.com",
		},
		{
			Host: "127.193.74.242",
		},
		{
			Host: "bastion-rsa",
		},
		{
			Host: "bastion-rsa.example.com",
		},
		{
			Host: "127.193.74.243",
		},
		{
			Host: "bastion-dsa",
		},
		{
			Host: "bastion-dsa.example.com",
		},
		{
			Host: "127.193.74.244",
		},
	}

	testParse(t, unit, want)
}

func newKnownHostsProviderUnit(fixture string) *KnownHostsProvider {
	fixture = filepath.Join("testdata", fixture+".known_hosts")

	return NewKnownHostsProvider(fixture, true)
}
