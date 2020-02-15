package provider

import (
	"fmt"
	"bufio"
	"bytes"
	"os"
	"strconv"

	"golang.org/x/crypto/ssh"

	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
	"github.com/UiP9AV6Y/ssh-select/pkg/util"
)

var (
	hashed_indicator  = []byte("|")
	comment_indicator = []byte("#")
)

type KnownHostsProvider struct {
	file string
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

	return p.parseFile(fd, true)
}

func (p *KnownHostsProvider) parseFile(file *os.File, ignoreMalformed bool) ([]remote.Host, error) {
	scanner := bufio.NewScanner(file)
	result := []remote.Host {}
	lineNo := 0

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		lineNo++

		if len(line) == 0 ||
			bytes.HasPrefix(line, hashed_indicator) ||
			bytes.HasPrefix(line, comment_indicator) {
			continue
		}

		if hosts, err := p.parseLine(line, ignoreMalformed); err == nil {
			result = append(result, hosts...)
		} else if !ignoreMalformed {
			return nil, fmt.Errorf("Parse error on line %d in %s: %w", lineNo, file.Name(), err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *KnownHostsProvider) parseLine(line []byte, ignoreMalformed bool) (result []remote.Host, err error) {
	var hosts []string

	_, hosts, _, _, _, err = ssh.ParseKnownHosts(line)

	if err != nil && !ignoreMalformed {
		return nil, err
	}

	result = make([]remote.Host, len(hosts))

	for i, host := range hosts {
		if result[i], err = p.parseHost(host); err != nil {
			if !ignoreMalformed {
				return nil, err
			}
		}
	}

	return result, nil
}

func (p *KnownHostsProvider) parseHost(host string) (remote.Host, error) {
	h := host
	// TODO: parse port
	n, _ := strconv.Atoi("0")
	parsed := remote.Host {
		Host: h,
		Port: n,
	}

	return parsed, nil
}

func NewKnownHostsProvider(file string) *KnownHostsProvider {
	provider := &KnownHostsProvider{
		file: file,
	}

	return provider
}

func UserKnownHostsProvider() *KnownHostsProvider {
	var provider *KnownHostsProvider

	if file, err := util.UserFilePath(false, ".ssh", "known_hosts"); err == nil {
		provider = &KnownHostsProvider{
			file: file,
		}
	}

	return provider
}

func ConfigKnownHostsProvider() *KnownHostsProvider {
	var provider *KnownHostsProvider

	if file, err := util.UserFilePath(true, "ssh", "known_hosts"); err == nil {
		provider = &KnownHostsProvider{
			file: file,
		}

	}

	return provider
}

func SystemKnownHostsProvider() *KnownHostsProvider {
	var provider *KnownHostsProvider

	if file, err := util.ResolvePath("/", "etc", "ssh", "ssh_known_hosts"); err == nil {
		provider = &KnownHostsProvider{
			file: file,
		}
	}

	return provider
}
