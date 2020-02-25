package search

import "fmt"

type Search interface {
	// Add an element to the search pool
	Add(element fmt.Stringer)
	// Find matches for the provided query
	Select(query string) []fmt.Stringer
	// Find the exact element
	Get(query string) (fmt.Stringer, bool)
	// Returns the number of elements
	Len() int
}
