package command

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mitchellh/colorstring"
	"github.com/mitchellh/goamz/rds"
)

const WaitPeriodSecond = 10
const WaitTimeoutPeriod = 180 // 180 * 10 sec = 30 min

type OpenCommand struct {
	Base
	instanceType  string
	queryFilePath string
	noWarmUp      bool
}

func (c *OpenCommand) Run(args []string) int {
	exitCode, err := c.realRun(args)
	if err != nil {
		c.Ui.Error(err.Error())
		c.Notifier.Error(fmt.Sprintf("@%s: %s", os.Getenv("USER"), err.Error()))
	}

	return exitCode
}

func (c *OpenCommand) parseArgs(args []string) (err error) {
	cFlag := flag.NewFlagSet("open", flag.ContinueOnError)
	cFlag.BoolVar(&c.noWarmUp, "W", false, "Disable Warmup.")
	cFlag.StringVar(&c.queryFilePath, "q", "", "Execute and Measurement query file.")
	cFlag.StringVar(&c.instanceType, "t", "db.m3.medium", "DB Instance-type. (default: db.m3.medium)")
	cFlag.Usage = func() { fmt.Println(c.Help()) }

	err = cFlag.Parse(args)
	return
}

func (c *OpenCommand) realRun(args []string) (exitCode int, err error) {
	if err = c.parseArgs(args); err != nil {
		exitCode = 1
		return
	}

	latestDBSnapshot, err := c.GetLatestSnapshot()
	if err != nil {
		exitCode = 1
		return
	}

	instanceName, err := c.CreateInstanceFromSnapshot(&latestDBSnapshot)
	if err != nil {
		exitCode = 1
		return
	}

	res := c.WaitTestInstanceAvailable(instanceName)
	if res == false {
		exitCode = 1
		return
	}

	instances, err := c.GetInstances(instanceName)
	if err != nil {
		exitCode = 1
		return
	}
	endpoint := instances[0].Address

	if c.noWarmUp == false && c.Resource.Warmup != "" {
		warmQueries, err := c.LoadQueries(c.Resource.Warmup)
		if err != nil {
			exitCode = 1
			return exitCode, err
		}

		c.Ui.Output(fmt.Sprintf("Warmup DB %s from %s...", endpoint, c.Resource.Warmup))

		_, err = c.Measure(endpoint, warmQueries)
		if err != nil {
			exitCode = 1
			return exitCode, err
		}

		c.Ui.Output("End Warmup.")
	} else {
		c.Ui.Output("Skip Warmup.")
	}

	if c.queryFilePath != "" {
		c.Ui.Output("Start Query execute measurement.")
		queries, err := c.LoadQueries(c.queryFilePath)
		if err != nil {
			exitCode = 1
			return exitCode, err
		}

		intervals, err := c.Measure(endpoint, queries)
		if err != nil {
			exitCode = 1
			return exitCode, err
		}

		var total float64
		message := "Execute Result:\n"
		for k, interval := range intervals {
			total += interval.Seconds()
			message += fmt.Sprintf("  Query: %s\n   Time: %s\n", queries[k], interval.String())
		}

		hour := int(total) / 3600
		minute := (int(total) - hour*3600) / 60
		second := total - float64(hour*3600) - float64(minute*60)

		message += "  ------------------------------\n"
		message += fmt.Sprintf("  Total: %d h %d m %.3f sec\n", hour, minute, second)

		c.Ui.Output(message)

		c.Notifier.Info(fmt.Sprintf("@%s Finished Query measurement.\nEndPoint: %s\n\n%s", os.Getenv("USER"), endpoint, message))
	} else {
		c.Notifier.Info(fmt.Sprintf("@%s Finished mirror DB setup. EndPoint: %s", os.Getenv("USER"), endpoint))
	}

	c.Ui.Output("All test instance setup Finished.")
	c.Ui.Output(colorstring.Color(fmt.Sprintf("Endpoint: [bold][green]%s", endpoint)))

	exitCode = 0
	return
}

func (c *OpenCommand) Help() string {
	return `Usage: rds-testrunner open [<args>]

Create Test DB instance.

Options:
    -W          Disable Warmup.
    -q <file>   Execute and Measurement query file.
    -t <type>   DB Instannce-type. (default: db.m3.medium)
`
}

func (c *OpenCommand) Synopsis() string {
	return "Create test Database instance."
}

// RDS Utility methods.
func (c *OpenCommand) GetLatestSnapshot() (latestDBSnapshot rds.DBSnapshot, err error) {
	options := rds.DescribeDBSnapshots{
		DBInstanceIdentifier: c.Resource.InstanceIdentifier,
		SnapshotType:         "automated",
	}

	resp, err := c.Client.DescribeDBSnapshots(&options)
	if err != nil {
		resp = nil
		return
	}

	length := len(resp.DBSnapshots)
	if length < 1 {
		err = fmt.Errorf("Snapshot is empty")
		return
	}

	latestDBSnapshot = resp.DBSnapshots[len(resp.DBSnapshots)-1]
	return
}

func (c *OpenCommand) CreateInstanceFromSnapshot(snapshot *rds.DBSnapshot) (testInstanceName string, err error) {
	testInstanceName = c.makeNewTestInstanceName()

	c.Ui.Output(colorstring.Color(fmt.Sprintf("Create: [bold][green]%s", testInstanceName)))
	c.Ui.Output(colorstring.Color(fmt.Sprintf("instance type: [bold][green]%s", c.instanceType)))

	options := rds.RestoreDBInstanceFromDBSnapshot{
		DBInstanceIdentifier: testInstanceName,
		DBSnapshotIdentifier: snapshot.DBSnapshotIdentifier,
		DBInstanceClass:      c.instanceType,
		MultiAZ:              false,
	}

	_, err = c.Client.RestoreDBInstanceFromDBSnapshot(&options)

	if err != nil {
		return
	}

	return
}

func (c *OpenCommand) WaitTestInstanceAvailable(instanceName string) bool {
	c.Ui.Output("waiting few minutes running status.")
	loop := 0
	for {
		i, err := c.GetInstances(instanceName)
		if err != nil {
			return false
		}

		if len(i) > 0 {
			break
		}

		time.Sleep(WaitPeriodSecond * 1000 * 1000 * 1000)
		fmt.Printf(".")

		loop++
		if loop > WaitTimeoutPeriod {
			fmt.Println(colorstring.Color("[red]wait timeout"))
			return false
		}
	}

	fmt.Println("finish.")

	return true
}
