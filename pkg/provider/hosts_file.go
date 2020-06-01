package provider

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
	"github.com/UiP9AV6Y/ssh-select/pkg/util"
)

var (
	hostsCommentIndicator = []byte("#")
)

type HostsFileProvider struct {
	file string
}

func (p *HostsFileProvider) String() string {
	return p.file
}

func (p *HostsFileProvider) Parse() ([]remote.Host, error) {
	var fd *os.File
	var err error

	fd, err = os.Open(p.file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	return p.parseFile(fd)
}

func (p *HostsFileProvider) parseFile(file *os.File) ([]remote.Host, error) {
	scanner := bufio.NewScanner(file)
	result := []remote.Host{}
	lineNo := 0

	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		lineNo++

		if len(line) == 0 ||
			bytes.HasPrefix(line, hostsCommentIndicator) {
			continue
		}

		if hosts, err := p.parseLine(line); err == nil {
			result = append(result, hosts...)
		} else {
			return nil, fmt.Errorf("Parse error on line %d in %s: %w", lineNo, file.Name(), err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *HostsFileProvider) parseLine(line []byte) (result []remote.Host, err error) {
	hosts := bytes.Fields(line)
	result = make([]remote.Host, 0, len(hosts))

	for _, host := range hosts {
		if bytes.HasPrefix(host, []byte(hostsCommentIndicator)) {
			break
		}

		result = append(result, remote.NewSimpleHost(string(host)))
	}

	return result, nil
}

func NewHostsFileProvider(file string) *HostsFileProvider {
	provider := &HostsFileProvider{
		file: file,
	}

	return provider
}

func SystemHostsFileProvider() *HostsFileProvider {
	var provider *HostsFileProvider

	if file, err := util.ResolvePath(systemHostsFile()); err == nil {
		provider = &HostsFileProvider{
			file: file,
		}
	}

	return provider
}

func systemHostsFile() string {
	switch runtime.GOOS {
	case "windows":
		path := filepath.Join("System32", "drivers", "etc", "hosts")
		dir := os.Getenv("SystemRoot")
		if dir == "" {
			dir = filepath.Join("C:", "Windows")
		}

		return filepath.Join(dir, path)
	case "beos":
		return "/boot/beos/etc/hosts"
	case "haiku":
		return "/system/settings/network/hosts"
	case "plan9":
		return "/lib/ndb/hosts"
	default: // Unix
		return "/etc/hosts"
	}
}
