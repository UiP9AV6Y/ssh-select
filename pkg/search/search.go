package search

import (
	prompt "github.com/c-bata/go-prompt"

	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
)

type Search interface {
	// Add a host object to the search pool
	Add(host remote.Host)
	// Find matches for the provided query
	Select(query string) []prompt.Suggest
	// Returns the number of elements
	Len() int
}
