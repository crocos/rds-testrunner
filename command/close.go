package command

import (
	"flag"
	"fmt"

	"github.com/mitchellh/colorstring"
	"github.com/mitchellh/goamz/rds"
)

type CloseCommand struct {
	Base
}

func (c *CloseCommand) Run(args []string) int {
	exitCode, err := c.realRun(args)
	if err != nil {
		c.Ui.Error(err.Error())
	}

	return exitCode
}

func (c *CloseCommand) parseArgs(args []string) (err error) {
	cFlag := flag.NewFlagSet("close", flag.ContinueOnError)
	cFlag.Usage = func() { fmt.Println(c.Help()) }

	err = cFlag.Parse(args)
	return
}

func (c *CloseCommand) realRun(args []string) (exitCode int, err error) {
	if err = c.parseArgs(args); err != nil {
		exitCode = 1
		return
	}

	instances, err := c.Base.GetInstances(c.Base.GetTestInstancePrefix())
	if err != nil {
		exitCode = 1
		return
	}

	if len(instances) < 1 {
		c.Ui.Info("Your TestDB Instance not exist.")
		exitCode = 1
		return
	}

	c.Ui.Output("Delete instances list:")
	for _, instance := range instances {
		fmt.Println(instance.DBInstanceIdentifier)
	}

	fmt.Println("")
	resp, err := c.Ui.Ask("Are you OK? [yes/NO]:")
	if err != nil {
		exitCode = 1
		return
	}

	if resp == "yes" || resp == "YES" || resp == "Yes" || resp == "y" || resp == "Y" {
		err = c.DeleteInstances(instances)
		if err != nil {
			exitCode = 1
			return
		}
	}

	exitCode = 0
	return
}

func (c *CloseCommand) Help() string {
	return `Usage: rds-testrunner close

Delete all Test DB instances.
`
}

func (c *CloseCommand) Synopsis() string {
	return "Destroy test Database instance."
}

func (c *CloseCommand) DeleteInstances(instances []rds.DBInstance) (err error) {
	for _, instance := range instances {
		options := rds.DeleteDBInstance{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
			SkipFinalSnapshot:    true,
		}
		_, err = c.Client.DeleteDBInstance(&options)
		if err != nil {
			return
		}

		err = nil
		c.Ui.Output(colorstring.Color(fmt.Sprintf("Delete [bold][green]%s", instance.DBInstanceIdentifier)))
	}

	return
}
