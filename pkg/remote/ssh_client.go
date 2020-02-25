package remote

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	prompt "github.com/c-bata/go-prompt"

	"github.com/UiP9AV6Y/ssh-select/pkg/search"
	"github.com/UiP9AV6Y/ssh-select/pkg/util"
)

type SshClient struct {
	cmd    string
	config string
	argv   []string
	env    []string
}

func (c *SshClient) Connect(host *Host) error {
	return syscall.Exec(c.cmd, c.CmdArray(host, true), c.env)
}

func (c *SshClient) NewExecutor(lookup search.Search, noop bool) prompt.Executor {
	executor := func(host string) {
		var target *Host
		var err error

		host = strings.TrimSpace(host)

		if host == "" {
			fmt.Println("No host selected; exiting")
			os.Exit(0)
		}

		if value, ok := lookup.Get(host); ok {
			target = value.(*Data).Host
		} else if target, err = ParseHost(host); err != nil {
			fmt.Println("Malformed connection target:", err)
			os.Exit(2)
		}

		if noop {
			fmt.Println("Connecting to", target, "(NOOP)")
			os.Exit(0)
		}

		fmt.Println("Connecting to", target)

		if err = c.Connect(target); err != nil {
			fmt.Println("Unable to connect:", err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	return executor
}

func (c *SshClient) CmdArray(host *Host, full bool) []string {
	s := []string{c.cmd}
	s = append(s, c.argv...)

	if full {
		if c.config != "" {
			s = append(s, "-f", c.config)
		}
	}

	if host != nil {
		if 0 < host.Port {
			s = append(s, "-p", strconv.Itoa(host.Port))
		}

		s = append(s, host.String())
	} else {
		s = append(s, "")
	}

	return s
}

func (c *SshClient) CmdLine(host *Host) string {
	s := c.CmdArray(host, false)

	return strings.Join(s, " ")
}

func (c *SshClient) String() string {
	return c.CmdLine(nil)
}

func NewSshClient(cmd string, argv []string, env []string) *SshClient {
	// config is optional and does not
	// interfer with initialization
	config, _ := util.UserFilePath(true, "ssh", "config")

	client := &SshClient{
		config: config,
		cmd:    cmd,
		argv:   argv,
		env:    env,
	}

	return client
}
