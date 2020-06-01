package provider

import (
	"path/filepath"
	"testing"

	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
)

func TestParseBasicZone(t *testing.T) {
	unit := newZoneFileProviderUnit("basic")
	want := []remote.Host{
		{
			Host: "dns1.example.com",
		},
		{
			Host: "example.com",
		},
		{
			Host: "dns2.example.com",
		},
		{
			Host: "example.com",
		},
		{
			Host: "10.0.1.1",
		},
		{
			Host: "dns1.example.com",
		},
		{
			Host: "aaaa:bbbb::1",
		},
		{
			Host: "dns1.example.com",
		},
		{
			Host: "10.0.1.2",
		},
		{
			Host: "dns2.example.com",
		},
		{
			Host: "aaaa:bbbb::2",
		},
		{
			Host: "dns2.example.com",
		},
		{
			Host: "mail.example.com",
		},
		{
			Host: "example.com",
		},
		{
			Host: "mail2.example.com",
		},
		{
			Host: "example.com",
		},
		{
			Host: "10.0.1.5",
		},
		{
			Host: "mail.example.com",
		},
		{
			Host: "aaaa:bbbb::5",
		},
		{
			Host: "mail.example.com",
		},
		{
			Host: "10.0.1.6",
		},
		{
			Host: "mail2.example.com",
		},
		{
			Host: "aaaa:bbbb::6",
		},
		{
			Host: "mail2.example.com",
		},
		{
			Host: "10.0.1.10",
		},
		{
			Host: "services.example.com",
		},
		{
			Host: "aaaa:bbbb::10",
		},
		{
			Host: "services.example.com",
		},
		{
			Host: "10.0.1.11",
		},
		{
			Host: "services.example.com",
		},
		{
			Host: "aaaa:bbbb::11",
		},
		{
			Host: "services.example.com",
		},
		{
			Host: "services.example.com",
		},
		{
			Host: "ftp.example.com",
		},
		{
			Host: "services.example.com",
		},
		{
			Host: "www.example.com",
		},
	}

	testParse(t, unit, want)
}

func TestParseReverseZone(t *testing.T) {
	unit := newZoneFileProviderUnit("reverse")
	want := []remote.Host{
		{
			Host: "dns1.example.com",
		},
		{
			Host: "1.0.10.in-addr.arpa",
		},
		{
			Host: "dns1.example.com",
		},
		{
			Host: "dns2.example.com",
		},
		{
			Host: "server1.example.com",
		},
		{
			Host: "server2.example.com",
		},
		{
			Host: "ftp.example.com",
		},
		{
			Host: "ftp.example.com",
		},
	}

	testParse(t, unit, want)
}

func newZoneFileProviderUnit(fixture string) *ZoneFileProvider {
	fixture = filepath.Join("testdata", fixture+".zone")

	return NewZoneFileProvider(fixture)
}
