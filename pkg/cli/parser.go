package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar"
)

type Parser struct {
	Version         bool
	Help            bool
	Quiet           bool
	NoSearch        bool
	NoOp            bool
	SshBinary       string
	Application     string
	ZoneFiles       []string
	HostsFiles      []string
	KnownHostsFiles []string
	SshArgv         []string
	Environment     []string
}

type optArg int

const (
	optArgNone optArg = iota
	optArgSshBinary
	optArgKnownHostsFile
	optArgHostsFile
	optArgZoneFile
)
const envVarPrefix = "SSH_SELECT_"

var optArgNames = []string{
	"",
	"SSH binary",
	"known hosts file",
	"hosts file",
	"zone file",
}

func (o optArg) String() string {
	return optArgNames[o]
}

func (p *Parser) ParseArgv(argv []string) error {
	var opt optArg
	var next bool
	var err error

	for _, arg := range argv {
		next, err = p.parseOptArg(opt, arg)
		if err != nil {
			return err
		} else if next {
			opt = optArgNone
			continue
		}

		opt, err = p.parseOpt(arg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) parseOpt(value string) (optArg, error) {
	switch value {
	case "--version":
		p.Version = true
	case "--help":
		p.Help = true
	case "--quiet":
		p.Quiet = true
	case "--no-search":
		p.NoSearch = true
	case "--no-op":
		p.NoOp = true
	case "--ssh":
		return optArgSshBinary, nil
	case "--known-hosts":
		return optArgKnownHostsFile, nil
	case "--hosts":
		return optArgHostsFile, nil
	case "--zone":
		return optArgZoneFile, nil
	default:
		p.SshArgv = append(p.SshArgv, value)
	}

	return optArgNone, nil
}

func (p *Parser) parseOptArg(optArg optArg, value string) (bool, error) {
	if optArg != optArgNone && strings.HasPrefix(value, "-") {
		return false, fmt.Errorf("missing argument for %s", optArg)
	}

	switch optArg {
	case optArgSshBinary:
		p.SshBinary = value
		return true, nil
	case optArgKnownHostsFile:
		return true, p.appendFileGlob(value, &p.KnownHostsFiles)
	case optArgHostsFile:
		return true, p.appendFileGlob(value, &p.HostsFiles)
	case optArgZoneFile:
		return true, p.appendFileGlob(value, &p.ZoneFiles)
	}

	return false, nil
}

func (p *Parser) ParseEnv(env []string) (err error) {
	for _, pair := range env {
		trimmed := strings.TrimPrefix(pair, envVarPrefix)

		if trimmed == pair {
			p.Environment = append(p.Environment, pair)
		} else {
			kv := strings.SplitN(trimmed, "=", 2)
			err = p.parseEnvArg(kv[0], kv[1])
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) parseEnvArg(key, value string) (err error) {
	if len(value) == 0 {
		return fmt.Errorf("env variable %s%s must not be empty", envVarPrefix, key)
	}

	switch {
	case key == "NO_OP":
		p.NoOp = true
	case key == "NO_SEARCH":
		p.NoSearch = true
	case key == "QUIET":
		p.Quiet = true
	case key == "SSH_BINARY":
		p.SshBinary = value
	case strings.HasPrefix(key, "KNOWN_HOSTS_FILE_"):
		return p.appendFileGlob(value, &p.KnownHostsFiles)
	case strings.HasPrefix(key, "HOSTS_FILE_"):
		return p.appendFileGlob(value, &p.HostsFiles)
	case strings.HasPrefix(key, "ZONE_FILE_"):
		return p.appendFileGlob(value, &p.ZoneFiles)
	}

	return nil
}

func (p *Parser) appendFileGlob(value string, target *[]string) error {
	if matches, err := doublestar.Glob(value); err != nil {
		return err
	} else if len(matches) > 0 {
		*target = append(*target, matches...)
	}

	return nil
}

func NewParser(application string) *Parser {
	parser := &Parser{
		Application: filepath.Base(application),
	}

	return parser
}
