package remote

import (
	prompt "github.com/c-bata/go-prompt"
)

type Data struct {
	Host *Host
	Suggestion *prompt.Suggest
}

func (n *Data) String() string {
  return n.Suggestion.Text
}

func NewData(host *Host) *Data {
	suggestion := &prompt.Suggest{
		Text: host.Host,
	}
  needle := &Data{
    Host: host,
    Suggestion: suggestion,
  }

  return needle
}
