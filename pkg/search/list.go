package search

import (
	"container/list"
	"strings"

	prompt "github.com/c-bata/go-prompt"

	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
)

type ListFilter func(s, substr string) bool

type List struct {
	length    int
	filter    ListFilter
	container *list.List
}

func (l *List) Len() int {
	return l.length
}

func (l *List) Add(host remote.Host) {
	l.container.PushBack(host.Suggest())
	l.length++
}

func (l *List) Select(query string) []prompt.Suggest {
	result := []prompt.Suggest{}

	for e := l.container.Front(); e != nil; e = e.Next() {
		v := e.Value.(prompt.Suggest)

		if l.filter(v.Text, query) {
			result = append(result, v)
		}
	}

	return result
}

func NewList(filter ListFilter) *List {
	if nil == filter {
		filter = strings.Contains
	}

	list := &List{
		length:    0,
		filter:    filter,
		container: list.New(),
	}

	return list
}
