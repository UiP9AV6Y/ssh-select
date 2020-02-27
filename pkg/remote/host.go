package remote

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
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
	var user string
	var addr string
	var port int
	var err error

	at := strings.LastIndex(text, "@")
	if at < 0 {
		addr, port, err = parseHost(text)
	} else {
		user = text[:at]
		addr, port, err = parseHost(text[at+1:])
	}

	if err != nil {
		return nil, err
	}

	host := &Host {
		User: user,
		Host: addr,
		Port: port,
	}

	return host, nil
}

func parseHost(host string) (string, int, error) {
	var colonPort string
	var port int
	var err error

	if strings.HasPrefix(host, "[") {
		// IPv6 notation [::1%2]:22
		i := strings.LastIndex(host, "]")
		if i < 0 {
			return "", 0, errors.New("missing ']' in host")
		} else if i < 2 {
			return "", 0, errors.New("host address must not be empty")
		}

		colonPort = host[i+1:]
		host = host[1:i]
	} else if i := strings.LastIndex(host, ":"); i != -1 {
		if i < 1 {
			return "", 0, errors.New("host address must not be empty")
		}

		colonPort = host[i:]
		host = host[:i]
	}

	if host == "" {
		return "", 0, errors.New("host must not be empty")
	}

	if colonPort != "" {
		if 2 > len(colonPort) {
			return "", 0, errors.New("port must not be empty")
		} else if port, err = strconv.Atoi(colonPort[1:]); err != nil {
			return "", 0, err
		}
	}

	return host, port, nil
}
