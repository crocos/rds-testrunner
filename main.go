package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/crocos/rds-testrunner/command"
	"github.com/crocos/rds-testrunner/config"
	"github.com/crocos/rds-testrunner/notify"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/colorstring"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/rds"
)

const VERSION = "0.1.0"

func main() {
	exitCode, err := realMain()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

func realMain() (exitCode int, err error) {
	resourceName := "default"
	args := make([]string, 0, len(os.Args))
	for k, arg := range os.Args {
		if arg == "" {
			continue
		}
		if arg == "-r" || arg == "--r" {
			resourceName = os.Args[k+1]

			os.Args[k] = ""
			os.Args[k+1] = ""
			continue
		}
		args = append(args, arg)
	}

	credential, err := config.GetAwsCredential()

	auth, err := aws.GetAuth(credential.Key, credential.Secret)
	if err != nil {
		exitCode = 1
		return
	}

	// get resource
	resourceData, err := config.GetResource(resourceName)
	if err != nil {
		exitCode = 1
		return
	}

	// get notify config
	notifyConfig, err := config.GetNotifyConfig()
	if err != nil {
		exitCode = 1
		return
	}

	// get rds client
	client := rds.New(auth, aws.Regions[resourceData.InstanceRegion])

	// get cli
	c := cli.NewCLI("rds-testrunner", VERSION)
	c.Args = args[1:]

	// get output ui
	ui := &cli.PrefixedUi{
		AskPrefix:    "",
		OutputPrefix: colorstring.Color("[bold][green][out] "),
		InfoPrefix:   colorstring.Color("[bold][yellow][info] "),
		ErrorPrefix:  colorstring.Color("[bold][red][error] "),
		Ui:           &cli.BasicUi{Reader: os.Stdin, Writer: os.Stdout, ErrorWriter: os.Stderr},
	}

	// get notifier
	notifier := notify.NewNotifier(notifyConfig)

	// create command utils
	base := command.Base{
		Ui:       ui,
		Resource: resourceData,
		Client:   client,
		Notifier: notifier,
	}

	c.Commands = map[string]cli.CommandFactory{
		"open": func() (cli.Command, error) {
			return &command.OpenCommand{
				Base: base,
			}, nil
		},
		"close": func() (cli.Command, error) {
			return &command.CloseCommand{
				Base: base,
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &command.ListCommand{
				Base: base,
			}, nil
		},
	}

	c.HelpFunc = helpFunc(args[0])
	exitCode, err = c.Run()
	return
}

func helpFunc(app string) cli.HelpFunc {
	return func(commands map[string]cli.CommandFactory) string {
		help := fmt.Sprintf("Usage: %s [--version] [--help] [-r] <command> [<command args>]\n\n", app)
		help += "DB mirroring and Query execution measurement tool.\n\n"
		help += "Available commands are:\n"

		keys := make([]string, 0, len(commands))
		maxKeyLen := 0
		for key, _ := range commands {
			if len(key) > maxKeyLen {
				maxKeyLen = len(key)
			}

			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			commandFunc, ok := commands[key]
			if !ok {
				panic("command not found: " + key)
			}

			command, err := commandFunc()
			if err != nil {
				log.Printf("[ERR] cli: Command '%s' failed to load: %s", key, err)
				continue
			}

			key = fmt.Sprintf("%s%s", key, strings.Repeat(" ", maxKeyLen-len(key)))
			help += fmt.Sprintf("    %s    %s\n", key, command.Synopsis())
		}

		help += "\nOptions:\n"
		help += fmt.Sprintf("    -r <resource>    Choice 'resource' name on your config file. if not set, use 'default'\n")

		return help
	}
}
