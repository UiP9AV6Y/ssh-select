package main

import (
	"fmt"
	"os"
	"os/exec"

	prompt "github.com/c-bata/go-prompt"

	"github.com/UiP9AV6Y/ssh-select/pkg/cli"
	"github.com/UiP9AV6Y/ssh-select/pkg/completer"
	"github.com/UiP9AV6Y/ssh-select/pkg/provider"
	"github.com/UiP9AV6Y/ssh-select/pkg/remote"
	"github.com/UiP9AV6Y/ssh-select/pkg/search"
	"github.com/UiP9AV6Y/ssh-select/pkg/version"
)

func printHelp(cli *cli.Parser) {
	fmt.Printf("Usage: %s [SSH_ARG...]\n", cli.Application)
}

func newProviders(cli *cli.Parser) []provider.HostProvider {
	providers := []provider.HostProvider{}

	for _, file := range cli.KnownHostsFiles {
		provider := provider.NewKnownHostsProvider(file, true)
		providers = append(providers, provider)
	}

	if cli.NoSearchKnownHosts {
		return providers
	}

	if provider := provider.UserKnownHostsProvider(true); provider != nil {
		providers = append(providers, provider)
	}

	if provider := provider.ConfigKnownHostsProvider(true); provider != nil {
		providers = append(providers, provider)
	}

	if provider := provider.SystemKnownHostsProvider(true); provider != nil {
		providers = append(providers, provider)
	}

	return providers
}

func newParser() *cli.Parser {
	parser := cli.NewParser(os.Args[0])

	if err := parser.ParseEnv(os.Environ()); err != nil {
		fmt.Println("Invalid environment variable:", err)
		os.Exit(1)
	}

	if err := parser.ParseArgv(os.Args[1:]); err != nil {
		fmt.Println("Invalid argument:", err)
		os.Exit(1)
	}

	return parser
}

func newSearch(providers []provider.HostProvider) (search.Search, map[string]int, error) {
	lookup := search.NewList(nil)
	sources := make(map[string]int)

	for _, provider := range providers {
		hosts, err := provider.Parse()

		if err != nil {
			return nil, nil, err
		}

		sources[provider.String()] = len(hosts)

		for _, host := range hosts {
			lookup.Add(remote.NewData(host))
		}
	}

	return lookup, sources, nil
}

func main() {
	var parser *cli.Parser
	var providers []provider.HostProvider
	var client *remote.SshClient
	var complete *completer.Completer
	var suggestions prompt.Completer
	var executor prompt.Executor
	var sources map[string]int
	var lookup search.Search
	var choice *prompt.Prompt
	var cmd string
	var err error

	parser = newParser()

	if parser.Version {
		fmt.Println(version.Version(), version.Commit())
		os.Exit(0)
	} else if parser.Help {
		printHelp(parser)
		os.Exit(0)
	}

	cmd = parser.SshBinary
	if len(cmd) == 0 {
		if cmd, err = exec.LookPath("ssh"); err != nil {
			fmt.Println("Unable to locale ssh binary:", err)
			os.Exit(1)
		}
	}

	providers = newProviders(parser)

	lookup, sources, err = newSearch(providers)
	if err != nil {
		fmt.Println("Unable to prepare host completion:", err)
		os.Exit(1)
	}

	client = remote.NewSshClient(cmd, parser.SshArgv, parser.Environment)
	complete = completer.NewCompleter(lookup)
	executor = client.NewExecutor(lookup, false)
	suggestions = complete.NewSuggestions()
	choice = prompt.New(
		executor,
		suggestions,
		prompt.OptionPrefix(client.String()),
		prompt.OptionPrefixTextColor(prompt.DefaultColor),
	)

	fmt.Printf("%s\n", version.Application("SSH Select"))
	for source, count := range sources {
		fmt.Printf("Received %d hosts from %s\n", count, source)
	}
	fmt.Printf("Providing suggestions for %d hosts\n", complete.SuggestionCount())
	choice.Run()
}
