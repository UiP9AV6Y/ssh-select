package completer

import (
	prompt "github.com/c-bata/go-prompt"

	"github.com/UiP9AV6Y/ssh-select/pkg/search"
)

type Completer struct {
	lookup search.Search
}

func (c *Completer) SuggestionCount() int {
	return c.lookup.Len()
}

func (c *Completer) NewSuggestions() prompt.Completer {
	completer := func(doc prompt.Document) []prompt.Suggest {
		if doc.TextBeforeCursor() == "" {
			return []prompt.Suggest{}
		}

		return c.lookup.Select(doc.GetWordBeforeCursor())
	}

	return completer
}

func NewCompleter(lookup search.Search) *Completer {
	completer := &Completer{
		lookup: lookup,
	}

	return completer
}
