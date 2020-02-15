package remote

import (
	"strconv"

	prompt "github.com/c-bata/go-prompt"
)

type Host struct {
	User string
	Host string
	Port int
}

func (h *Host) String() string {
	return h.Host
}

func (h *Host) Suggest() prompt.Suggest {
	return prompt.Suggest{
		Text: h.Host,
	}
}

func ParseSuggestText(text string) (*Host, error) {
	// TODO: parse user/port
	u := ""
	h := text
	p, _ := strconv.Atoi("0")
	host := &Host {
		User: u,
		Host: h,
		Port: p,
	}

	return host, nil
}
