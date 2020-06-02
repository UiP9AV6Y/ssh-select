package provider

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
)

var (
	dumpCommentIndicator = byte(';')
	dumpWildcardReplace = "wildcard"
)

type NamedDumpProvider struct {
	file             string
	processWildcards bool
}

func (p *NamedDumpProvider) String() string {
	return p.file
}

func (p *NamedDumpProvider) Parse() ([]remote.Host, error) {
	var fd *os.File
	var err error

	fd, err = os.Open(p.file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	return p.parseFile(fd)
}

func (p *NamedDumpProvider) parseFile(file *os.File) ([]remote.Host, error) {
	scanner := bufio.NewScanner(file)
	result := []remote.Host{}
	lineNo := 0
	sink := func(host string) {
		fqdn := strings.TrimSuffix(host, fullyQualifiedIndicator)

		if strings.Contains(fqdn, "*") {
			if !p.processWildcards {
				return
			}

			fqdn = strings.ReplaceAll(fqdn, "*", dumpWildcardReplace)
		}

		result = append(result, remote.NewSimpleHost(fqdn))
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNo++

		if len(line) == 0 ||
			line[0] == dumpCommentIndicator {
			continue
		}

		if err := p.parseLine(line, sink); err != nil {
			return nil, fmt.Errorf("Parse error on line %d in %s: %w", lineNo, file.Name(), err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *NamedDumpProvider) parseLine(line string, host func(string)) error {
	columns := strings.Fields(line)
	length := len(columns)

	if length < 5 {
		return fmt.Errorf("No resource type available")
	}

	switch columns[3] {
	case "PTR":
		// the zone $ORIGIN is usually not a plain address,
		// therefor it makes no sense to provide it as host
		// (e.g. 1.0.10.in-addr.arpa.)
		//host(columns[0])
		host(columns[length-1])
	case "A", "AAAA", "NS", "MX", "CNAME":
		host(columns[0])
		host(columns[length-1])
	}

	return nil
}

func NewNamedDumpProvider(file string, processWildcards bool) *NamedDumpProvider {
	provider := &NamedDumpProvider{
		file:             file,
		processWildcards: processWildcards,
	}

	return provider
}
