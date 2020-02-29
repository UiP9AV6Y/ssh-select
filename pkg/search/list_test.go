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
  assert.ElementsMatch(t, unit.Select("2"), want, "searching for '2'")

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
  assert.ElementsMatch(t, unit.Select("0"), want, "searching for '0'")

  assert.Empty(t, unit.Select("xxx"), "searching for 'xxx'")
}

func TestGet(t *testing.T) {
  unit := NewList(nil)
  elements := newElements(3, unit)

  for _, want := range elements {
    got, ok := unit.Get(want.String())

    assert.Truef(t, ok, "expected to find %v", want)
    assert.Equal(t, got, want)
  }

  got, ok := unit.Get("xxx")
  assert.Falsef(t, ok, "expected to find nothing, got %v", got)
}

func TestLen(t *testing.T) {
  unit := NewList(nil)
  assert.Equal(t, unit.Len(), 0)

  unit.Add(newElement(1))
  assert.Equal(t, unit.Len(), 1)

  unit.Add(newElement(2))
  assert.Equal(t, unit.Len(), 2)

  unit.Add(newElement(1))
  assert.Equal(t, unit.Len(), 3)
}
