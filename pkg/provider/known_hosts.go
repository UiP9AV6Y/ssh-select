package provider

import (
	"bufio"
	"os"

	"golang.org/x/crypto/ssh"

	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
	"github.com/UiP9AV6Y/ssh-select/pkg/util"
)

type KnownHostsProvider struct {
	file string
}

func (p *KnownHostsProvider) String() string {
	return p.file
}

func (p *KnownHostsProvider) Parse() ([]remote.Host, error) {
	hosts, err := ParseKnownHosts(p.file, false)

	if err != nil {
		return nil, err
	}

	parsed := make([]remote.Host, len(hosts))

	for i := range hosts {
		parsed[i] = remote.Host(hosts[i])
	}

	return parsed, nil
}

func NewKnownHostsProvider(file string) *KnownHostsProvider {
	provider := &KnownHostsProvider{
		file: file,
	}

	return provider
}

func ParseKnownHosts(file string, ignoreMalformed bool) ([]string, error) {
	var fd *os.File
	var hosts []string
	var err error
	var result []string

	fd, err = os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		_, hosts, _, _, _, err = ssh.ParseKnownHosts(scanner.Bytes())

		if err == nil {
			result = append(result, hosts...)
		} else if !ignoreMalformed {
			return nil, err
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func UserKnownHostsProvider() *KnownHostsProvider {
	var provider *KnownHostsProvider

	if file, err := util.UserFilePath(false, ".ssh", "config"); err == nil {
		provider = &KnownHostsProvider{
			file: file,
		}
	}

	return provider
}

func ConfigKnownHostsProvider() *KnownHostsProvider {
	var provider *KnownHostsProvider

	if file, err := util.UserFilePath(true, "ssh", "config"); err == nil {
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
