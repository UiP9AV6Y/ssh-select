package search

import (
	"fmt"
	"container/list"
	"strings"
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

func (l *List) Add(element fmt.Stringer) {
	l.container.PushBack(element)
	l.length++
}

func (l *List) Select(query string) []fmt.Stringer {
	result := []fmt.Stringer{}

	for e := l.container.Front(); e != nil; e = e.Next() {
		v := e.Value.(fmt.Stringer)

		if l.filter(v.String(), query) {
			result = append(result, v)
		}
	}

	return result
}

func (l *List) Get(query string) (fmt.Stringer, bool) {
	for e := l.container.Front(); e != nil; e = e.Next() {
		v := e.Value.(fmt.Stringer)

		if v.String() == query {
			return v, true
		}
	}

	return nil, false
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
