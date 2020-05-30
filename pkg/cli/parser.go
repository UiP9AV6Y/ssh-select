package cli

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Parser struct {
	Version         bool
	Help            bool
	SshBinary       string
	Application     string
	KnownHostsFiles []string
	SshArgv         []string
	Environment     []string
}

type optArg int

const (
	optArgNone optArg = iota
	optArgSshBinary
	optArgKnownHostsFile
)
const envVarPrefix = "SSH_SELECT_"

var optArgNames = []string{
	"",
	"SSH binary",
	"known hosts file",
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
	case "--ssh":
		return optArgSshBinary, nil
	case "--known-hosts":
		return optArgKnownHostsFile, nil
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
		p.KnownHostsFiles = append(p.KnownHostsFiles, value)
		return true, nil
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
	case key == "SSH_BINARY":
		p.SshBinary = value
	case strings.HasPrefix(key, "KNOWN_HOSTS_FILE_"):
		p.KnownHostsFiles = append(p.KnownHostsFiles, value)
	}

	return nil
}

func NewParser(application string) *Parser {
	parser := &Parser{
		Version:         false,
		Help:            false,
		SshBinary:       "",
		Application:     filepath.Base(application),
		KnownHostsFiles: []string{},
		SshArgv:         []string{},
		Environment:     []string{},
	}

	return parser
}
