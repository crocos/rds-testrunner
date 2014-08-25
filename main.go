package main

import (
	"log"
	"os"

	"github.com/crocos/rds-testrunner/command"
	"github.com/crocos/rds-testrunner/config"
	"github.com/crocos/rds-testrunner/notify"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/colorstring"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/rds"
)

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
	c := cli.NewCLI("rds-testrunner", "0.1.0")
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
	notifier := notify.NewHipChatNotifier(notifyConfig)

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

	exitCode, err = c.Run()
	return
}
