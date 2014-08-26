package command

import (
	"flag"
	"fmt"
)

type ListCommand struct {
	Base
}

func (c *ListCommand) Run(args []string) int {
	exitCode, err := c.realRun(args)
	if err != nil {
		c.Ui.Error(err.Error())
	}

	return exitCode
}

func (c *ListCommand) parseArgs(args []string) (err error) {
	cFlag := flag.NewFlagSet("list", flag.ContinueOnError)
	cFlag.Usage = func() { fmt.Println(c.Help()) }

	err = cFlag.Parse(args)
	return
}

func (c *ListCommand) realRun(args []string) (exitCode int, err error) {
	if err = c.parseArgs(args); err != nil {
		exitCode = 1
		return
	}

	instances, err := c.GetInstances(c.GetTestInstancePrefix())
	if err != nil {
		exitCode = 1
		return
	}

	c.Ui.Output("Your TestDB instance list:")
	for _, instance := range instances {
		fmt.Println(instance.DBInstanceIdentifier)
	}

	exitCode = 0
	return
}

func (c *ListCommand) Help() string {
	return `Usage: rds-testrunner list

List available Test DB instances.
`
}

func (c *ListCommand) Synopsis() string {
	return "List test Database instance."
}
