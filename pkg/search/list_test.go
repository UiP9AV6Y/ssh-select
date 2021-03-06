package search

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type element struct {
	Value string
}

func (e element) String() string {
	return e.Value
}

func newElement(id int) element {
	value := "element-" + strconv.Itoa(id)

	return element{
		Value: value,
	}
}

func newElements(count int, unit *List) []element {
	elements := make([]element, count)

	for i := range elements {
		elements[i] = newElement(i)
		unit.Add(elements[i])
	}

	return elements
}

func TestSelect(t *testing.T) {
	unit := NewList(nil)
	elements := newElements(21, unit)
	want := []element{
		elements[2],
		elements[12],
		elements[20],
	}
	assert.ElementsMatch(t, want, unit.Select("2"), "searching for '2'")

	assert.Empty(t, unit.Select("xxx"), "searching for 'xxx'")
}

func TestSelectCustom(t *testing.T) {
	unit := NewList(strings.HasSuffix)
	elements := newElements(21, unit)
	want := []element{
		elements[0],
		elements[10],
		elements[20],
	}
	assert.ElementsMatch(t, want, unit.Select("0"), "searching for '0'")

	assert.Empty(t, unit.Select("xxx"), "searching for 'xxx'")
}

func TestGet(t *testing.T) {
	unit := NewList(nil)
	elements := newElements(3, unit)

	for _, want := range elements {
		got, ok := unit.Get(want.String())

		assert.Truef(t, ok, "expected to find %v", want)
		assert.Equal(t, want, got)
	}

	got, ok := unit.Get("xxx")
	assert.Falsef(t, ok, "expected to find nothing, got %v", got)
}

func TestLen(t *testing.T) {
	unit := NewList(nil)
	assert.Equal(t, 0, unit.Len())

	unit.Add(newElement(1))
	assert.Equal(t, 1, unit.Len())

	unit.Add(newElement(2))
	assert.Equal(t, 2, unit.Len())

	unit.Add(newElement(1))
	assert.Equal(t, 2, unit.Len())

	unit.Add(newElement(3))
	assert.Equal(t, 3, unit.Len())

	unit.Add(newElement(2))
	assert.Equal(t, 3, unit.Len())

	unit.Add(newElement(5))
	assert.Equal(t, 4, unit.Len())

	unit.Add(newElement(4))
	assert.Equal(t, 5, unit.Len())

	unit.Add(newElement(3))
	assert.Equal(t, 5, unit.Len())
}

func TestAdd(t *testing.T) {
	unit := NewList(nil)
	assert.Equal(t, 0, unit.Len())

	unit.Add(element{
		Value: "aaaa",
	})
	unit.Add(element{
		Value: "bbbb",
	})
	unit.Add(element{
		Value: "aaa",
	})
	unit.Add(element{
		Value: "aa",
	})
	unit.Add(element{
		Value: "a",
	})
	unit.Add(element{
		Value: "aa",
	})
	unit.Add(element{
		Value: "aaa",
	})
	unit.Add(element{
		Value: "bbbb",
	})
	unit.Add(element{
		Value: "bbb",
	})
	unit.Add(element{
		Value: "bb",
	})
	unit.Add(element{
		Value: "b",
	})
	unit.Add(element{
		Value: "aa",
	})
	unit.Add(element{
		Value: "aaa",
	})
	unit.Add(element{
		Value: "aaaa",
	})
	unit.Add(element{
		Value: "cccc",
	})
	unit.Add(element{
		Value: "dddd",
	})
	unit.Add(element{
		Value: "ffff",
	})
	unit.Add(element{
		Value: "bbbb",
	})

	assert.Equal(t, 11, unit.Len())
}
