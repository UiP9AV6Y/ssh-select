package provider

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
	"github.com/UiP9AV6Y/ssh-select/pkg/util"
)

var (
	hashed_indicator  = []byte("|")
	comment_indicator = []byte("#")
	host_sanitizer    = regexp.MustCompile(`(>%[0-9]+|_[a-z0-9]+)$`)
)

type KnownHostsProvider struct {
	ignoreMalformed bool
	file            string
}

func (p *KnownHostsProvider) String() string {
	return p.file
}

func (p *KnownHostsProvider) Parse() ([]remote.Host, error) {
	var fd *os.File
	var err error

	fd, err = os.Open(p.file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	return p.parseFile(fd)
}

func (p *KnownHostsProvider) parseFile(file *os.File) ([]remote.Host, error) {
	scanner := bufio.NewScanner(file)
	result := []remote.Host{}
	lineNo := 0

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		lineNo++

		if len(line) == 0 ||
			bytes.HasPrefix(line, hashed_indicator) ||
			bytes.HasPrefix(line, comment_indicator) {
			continue
		}

		if hosts, err := p.parseLine(line); err == nil {
			result = append(result, hosts...)
		} else if !p.ignoreMalformed {
			return nil, fmt.Errorf("Parse error on line %d in %s: %w", lineNo, file.Name(), err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *KnownHostsProvider) parseLine(line []byte) (result []remote.Host, err error) {
	var hosts []string

	_, hosts, _, _, _, err = ssh.ParseKnownHosts(line)

	if err != nil && !p.ignoreMalformed {
		return nil, err
	}

	result = make([]remote.Host, len(hosts))

	for i, host := range hosts {
		if result[i], err = p.parseHost(host); err != nil {
			if !p.ignoreMalformed {
				return nil, err
			}
		}
	}

	return result, nil
}

func (p *KnownHostsProvider) parseHost(host string) (remote.Host, error) {
	var parsed remote.Host
	var colonPort string
	var port int
	var err error

	if strings.HasPrefix(host, "[") {
		// port notation [1.2.3.4]:22
		i := strings.LastIndex(host, "]")
		if i < 0 {
			return parsed, errors.New("missing ']' in host")
		} else if i < 2 {
			return parsed, errors.New("host address must not be empty")
		}

		colonPort = host[i+1:]
		host = host[1:i]
	}

	if host == "" {
		return parsed, errors.New("host must not be empty")
	} else if host, err = p.sanitizeHost(host); err != nil {
		return parsed, err
	}

	if colonPort != "" {
		if 2 > len(colonPort) {
			return parsed, errors.New("port must not be empty")
		} else if port, err = strconv.Atoi(colonPort[1:]); err != nil {
			return parsed, err
		}
	}

	parsed = remote.Host{
		Host: host,
		Port: port,
	}

	return parsed, nil
}

func (p *KnownHostsProvider) sanitizeHost(host string) (string, error) {
	host = strings.TrimPrefix(host, "<")
	host = host_sanitizer.ReplaceAllLiteralString(host, "")

	return host, nil
}

func NewKnownHostsProvider(file string, ignoreMalformed bool) *KnownHostsProvider {
	provider := &KnownHostsProvider{
		ignoreMalformed: ignoreMalformed,
		file:            file,
	}

	return provider
}

func UserKnownHostsProvider(ignoreMalformed bool) *KnownHostsProvider {
	var provider *KnownHostsProvider

	if file, err := util.UserFilePath(false, ".ssh", "known_hosts"); err == nil {
		provider = &KnownHostsProvider{
			ignoreMalformed: ignoreMalformed,
			file:            file,
		}
	}

	return provider
}

func ConfigKnownHostsProvider(ignoreMalformed bool) *KnownHostsProvider {
	var provider *KnownHostsProvider

	if file, err := util.UserFilePath(true, "ssh", "known_hosts"); err == nil {
		provider = &KnownHostsProvider{
			ignoreMalformed: ignoreMalformed,
			file:            file,
		}

	}

	return provider
}

func SystemKnownHostsProvider(ignoreMalformed bool) *KnownHostsProvider {
	var provider *KnownHostsProvider

	if file, err := util.ResolvePath("/", "etc", "ssh", "ssh_known_hosts"); err == nil {
		provider = &KnownHostsProvider{
			ignoreMalformed: ignoreMalformed,
			file:            file,
		}
	}

	return provider
}
