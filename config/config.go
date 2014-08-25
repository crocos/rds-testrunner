package config

import (
	"fmt"
	"io/ioutil"
	"os"

	parser "github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl"
)

const RTRRC = "rtrrc"
const DEFALULT_CONFIG = "/etc/rds-testrunner.conf"

type Resource struct {
	InstanceIdentifier string
	InstanceRegion     string
	User               string
	Password           string
	Warmup             string
}

type NotifyConfig struct {
	Token string
	Name  string
	Room  string
}

func configFile() (file string, err error) {
	if home := os.Getenv("HOME"); home != "" {
		file = fmt.Sprintf("%s/.%s", home, RTRRC)

		if _, err = os.Stat(file); os.IsNotExist(err) {
			file = DEFALULT_CONFIG

			if _, err = os.Stat(file); os.IsNotExist(err) {
				err = fmt.Errorf("Config File does not exist.")
				return
			}
		}
	}

	err = nil
	return
}

func LoadConfig() (config *hcl.Object, err error) {
	file, err := configFile()
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	config, err = parser.Parse(string(data))

	return
}

func GetResource(name string) (resourceData Resource, err error) {
	config, err := LoadConfig()
	if err != nil {
		return
	}

	priv := config.Get("resource", false)
	if priv == nil {
		err = fmt.Errorf("resource not defined.")
		return
	}

	for {
		resource := priv.Get(name, false)

		if resource != nil {
			// require
			if resource.Get("instance", false) != nil {
				resourceData.InstanceIdentifier = resource.Get("instance", false).Value.(string)
			} else {
				err = fmt.Errorf("instance is required in resource '%s'", name)
				break
			}

			if resource.Get("region", false) != nil {
				resourceData.InstanceRegion = resource.Get("region", false).Value.(string)
			} else {
				err = fmt.Errorf("region is require in resource '%s'", name)
			}

			if resource.Get("user", false) != nil {
				resourceData.User = resource.Get("user", false).Value.(string)
			}
			if resource.Get("password", false) != nil {
				resourceData.Password = resource.Get("password", false).Value.(string)
			}

			if resource.Get("warmup", false) != nil {
				resourceData.Warmup = resource.Get("warmup", false).Value.(string)
			}

			break
		}

		if priv.Next == nil {
			err = fmt.Errorf("no resource found '%s'", name)
			break
		}
		priv = priv.Next
	}

	return
}

// Now, hipchat only. TODO: other platform notifier
func GetNotifyConfig() (notifyConfig NotifyConfig, err error) {
	config, err := LoadConfig()
	if err != nil {
		return
	}

	notify := config.Get("notify", false)
	if notify == nil {
		return // notify not used.
	}

	if notify.Get("token", false) != nil {
		notifyConfig.Token = notify.Get("token", false).Value.(string)
	} else {
		err = fmt.Errorf("notify 'token' is required.")
		return
	}

	if notify.Get("room", false) != nil {
		notifyConfig.Room = notify.Get("room", false).Value.(string)
	} else {
		err = fmt.Errorf("notify 'room' is required.")
		return
	}

	if notify.Get("name", false) != nil {
		notifyConfig.Name = notify.Get("name", false).Value.(string)
	} else {
		notifyConfig.Name = "RDS TestRunner"
	}

	return
}
