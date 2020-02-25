package completer

import (
	prompt "github.com/c-bata/go-prompt"

	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
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

		needles := c.lookup.Select(doc.GetWordBeforeCursor())
		result := make([]prompt.Suggest, len(needles))

		for i, needle := range needles {
			result[i] = *needle.(*remote.Data).Suggestion
		}

		return result
	}

	return completer
}

func NewCompleter(lookup search.Search) *Completer {
	completer := &Completer{
		lookup: lookup,
	}

	return completer
}
