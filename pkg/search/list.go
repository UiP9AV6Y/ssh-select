package search

import (
	"container/list"
	"fmt"
	"strings"
)

type ListFilter func(s, substr string) bool

type List struct {
	filter    ListFilter
	container *list.List
}

func (l *List) Len() int {
	return l.container.Len()
}

func (l *List) Add(element fmt.Stringer) {
	needle := element.String()

	for e := l.container.Front(); e != nil; e = e.Next() {
		v := e.Value.(fmt.Stringer)
		c := strings.Compare(v.String(), needle)

		if c == 0 {
			return
		} else if c > 0 {
			l.container.InsertAfter(element, e)
			return
		}
	}

	l.container.PushBack(element)
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
		filter:    filter,
		container: list.New(),
	}

	return list
}
