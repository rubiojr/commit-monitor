package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/muesli/beehive/bees"
	_ "github.com/muesli/beehive/bees/githubbee"
	_ "github.com/muesli/beehive/bees/hellobee"
	"github.com/muesli/beehive/cfg"
	"github.com/muesli/beehive/reactor"
	"github.com/muesli/termenv"
	_ "github.com/rubiojr/commit-monitor/bees/stdoutbee"
)

// Repositories to monitor (owner, repo name)
var repos = [][]string{
	{"rails", "rails"},
	{"torvalds", "linux"},
	{"kubernetes", "kubernetes"},
	{"prometheus", "prometheus"},
	{"golang", "go"},
	{"LineageOS", "android_frameworks_base"},
	{"rubiojr", "hello"},
}

func main() {
	// configuration boilerplate using the memory backend
	// so we don't touch the filesystem.
	config, err := cfg.New("mem://")
	if err != nil {
		panic(err)
	}
	config.Bees = []bees.BeeConfig{}
	config.Chains = []bees.Chain{}

	// the stdout bee prints to stdout via the
	stdoutBee := newStdoutBee()
	config.Bees = append(config.Bees, stdoutBee)

	// Create the Action and add it to the config
	// Every chain will re-use the same action to print commits to stdout
	// using the 'print' action
	action := bees.Action{}
	action.ID = "print-to-stdout"
	action.Bee = stdoutBee.Name
	action.Name = "print"
	action.Options = bees.Placeholders{
		{
			Name: "text",
			// prints something like: ** New commit in owner/repo ** for every commit
			Value: formatText(`{{ Bold (Foreground "#FF0000" "**") }}`) +
				" New commit in {{.repo}} " +
				formatText(`{{ Bold (Foreground "#FF0000" "**") }}`) +
				"\n" +
				"{{.message}}\n",
			Type: "string",
		},
	}
	config.Actions = []bees.Action{action}

	// Iterate over all the repositories we want to monitor
	// and create a new chain that will link the 'commit' event
	// to the 'print-to-stdout' action.
	for _, repo := range repos {
		nwo := strings.Join(repo, "/") // owner/repository
		// the GitHub bee is in charge of monitoring events
		// for the given repository
		bee := newGitHubBee(repo[0], repo[1])
		config.Bees = append(config.Bees, bee)

		// Create the event
		event := bees.Event{}
		event.Name = "commit"
		event.Bee = bee.Name

		// Create the chain and add it to the existing chains
		chain := bees.Chain{}
		chain.Name = "commits-" + nwo
		chain.Description = "Print commits for " + nwo
		chain.Actions = []string{action.ID} // Action to print the commit
		chain.Event = &event
		chain.Filters = []string{}
		config.Chains = append(config.Chains, chain)
	}

	// Debugging level, prints debug messages from bees
	// reactor.SetLogLevel(5)
	reactor.Run(config)
}

func newGitHubBee(owner, repo string) bees.BeeConfig {
	options := bees.BeeOptions{
		bees.BeeOption{Name: "accesstoken", Value: os.Getenv("GITHUB_TOKEN")},
		bees.BeeOption{Name: "owner", Value: owner},
		bees.BeeOption{Name: "repository", Value: repo},
	}
	bc, err := bees.NewBeeConfig(owner+"-"+repo, "githubbee", fmt.Sprintf("monitor %s/%s commits", owner, repo), options)
	if err != nil {
		panic(err)
	}

	return bc
}

func newStdoutBee() bees.BeeConfig {
	options := bees.BeeOptions{}
	bc, err := bees.NewBeeConfig("stdout", "stdoutbee", "test", options)
	if err != nil {
		panic(err)
	}

	return bc
}

func formatText(text string) string {
	// load template helpers
	f := termenv.TemplateFuncs(termenv.ColorProfile())
	tpl := template.New("tpl").Funcs(f)

	// parse and render
	tpl, err := tpl.Parse(text)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	tpl.Execute(&buf, nil)
	return buf.String()
}
