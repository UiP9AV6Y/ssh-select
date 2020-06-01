package provider

import (
	"testing"

	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
)

func testParse(t *testing.T, unit HostProvider, want []remote.Host) {
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
