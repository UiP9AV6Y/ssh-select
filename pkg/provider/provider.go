package provider

import (
	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
)

type HostProvider interface {
	// Returns a string representation of the provider
	String() string
	// Returns the parsed hosts from the underlying source
	Parse() ([]remote.Host, error)
}
