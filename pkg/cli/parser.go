package cli

type Parser struct {
	version    bool
	help       bool
	knownHosts []string
	sshArgv    []string
	env        []string
}

func (p *Parser) Version() bool {
	return p.version
}

func (p *Parser) Help() bool {
	return p.help
}

func (p *Parser) KnownHostFiles() []string {
	return p.knownHosts
}

func (p *Parser) SshBinary() string {
	// TODO: add parameter for custom binary path
	return ""
}

func (p *Parser) SshArgv() []string {
	return p.sshArgv
}

func (p *Parser) Environment() []string {
	return p.env
}

func (p *Parser) ParseArgv(argv []string) error {
	// TODO: separate into passthrough and local arguments
	p.sshArgv = argv

	return nil
}

func (p *Parser) ParseEnv(env []string) error {
	// TODO: parse variables for local arguments
	p.env = env

	return nil
}

func NewParser() *Parser {
	parser := &Parser{
		version:    false,
		help:       false,
		knownHosts: []string{},
		sshArgv:    []string{},
		env:        []string{},
	}

	return parser
}
