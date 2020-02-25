package remote

import (
	"fmt"
	"strconv"

	prompt "github.com/c-bata/go-prompt"
)

type Host struct {
	User string
	Host string
	Port int
}

func (h *Host) String() string {
	if "" != h.User {
		return fmt.Sprintf("%s@%s", h.User, h.Host)
	}

	return h.Host
}

func (h *Host) Suggest() prompt.Suggest {
	return prompt.Suggest{
		Text: h.Host,
	}
}

func ParseHost(text string) (*Host, error) {
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
