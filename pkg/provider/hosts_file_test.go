package provider

import (
	"path/filepath"
	"testing"

	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
)

func TestParseBasicHosts(t *testing.T) {
	unit := newHostsFileProviderUnit("basic")
	want := []remote.Host{
		{
			Host: "8.8.8.8",
		},
		{
			Host: "gogol.net",
		},
		{
			Host: "127.0.0.1",
		},
		{
			Host: "localhost",
		},
		{
			Host: "127.0.1.1",
		},
		{
			Host: "thishost.mydomain.org",
		},
		{
			Host: "thishost",
		},
		{
			Host: "192.168.1.10",
		},
		{
			Host: "foo.mydomain.org",
		},
		{
			Host: "foo",
		},
		{
			Host: "192.168.1.13",
		},
		{
			Host: "bar.mydomain.org",
		},
		{
			Host: "bar",
		},
		{
			Host: "146.82.138.7",
		},
		{
			Host: "master.debian.org",
		},
		{
			Host: "master",
		},
		{
			Host: "209.237.226.90",
		},
		{
			Host: "www.opensource.org",
		},
		{
			Host: "::1",
		},
		{
			Host: "localhost",
		},
		{
			Host: "ip6-localhost",
		},
		{
			Host: "ip6-loopback",
		},
		{
			Host: "ff02::1",
		},
		{
			Host: "ip6-allnodes",
		},
		{
			Host: "ff02::2",
		},
		{
			Host: "ip6-allrouters",
		},
	}

	testParse(t, unit, want)
}

func newHostsFileProviderUnit(fixture string) *HostsFileProvider {
	fixture = filepath.Join("testdata", fixture+".hosts")

	return NewHostsFileProvider(fixture)
}
