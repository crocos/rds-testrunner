package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/crocos/rds-testrunner/config"
	"github.com/crocos/rds-testrunner/notify"
	"github.com/mitchellh/cli"
	"github.com/mitchellh/goamz/rds"
)

type Base struct {
	Ui       cli.Ui
	Resource config.Resource
	Client   *rds.Rds
	Notifier notify.HipChatNotifier
}

func (c *Base) GetTestInstancePrefix() (prefix string) {
	return fmt.Sprintf("test-%s", os.Getenv("USER"))
}

func (c *Base) GetInstances(needle string) (instances []rds.DBInstance, err error) {
	options := rds.DescribeDBInstances{}
	resp, err := c.Client.DescribeDBInstances(&options)
	if err != nil {
		return
	}

	for _, instance := range resp.DBInstances {
		if instance.DBInstanceStatus != "available" {
			continue
		}
		if strings.Contains(instance.DBInstanceIdentifier, needle) == false {
			continue
		}

		instances = append(instances, instance)
	}

	return
}

func (c *Base) makeNewTestInstanceName() (testName string) {
	prefix := c.GetTestInstancePrefix()
	datetime := time.Now().Format("20060102150405")

	testName = fmt.Sprintf("%s-%s-%s", prefix, datetime, c.Resource.InstanceIdentifier)

	return
}

func (c *Base) Measure(endpoint string, queries []string) (intervals []time.Duration, err error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/", c.Resource.User, c.Resource.Password, endpoint))
	if err != nil {
		return
	}
	defer db.Close()

	intervals = make([]time.Duration, 0, len(queries))
	for _, query := range queries {
		start := time.Now()

		result, err := db.Query(query)
		if err != nil {
			return intervals, err // return in-block err
		}

		end := time.Now()
		intervals = append(intervals, end.Sub(start))

		result.Close()
	}

	return
}

func (c *Base) LoadQueries(file string) (queries []string, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")

	queries = make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		queries = append(queries, strings.TrimSpace(line))
	}

	return
}
