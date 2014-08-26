package notify

import "github.com/crocos/rds-testrunner/config"

func NewNotifier(notifyConfig config.NotifyConfig) NotifierInterface {
	if notifyConfig.Type == "hipchat" {
		return NewHipChatNotifier(notifyConfig)
	} else {
		return NewDummyNotifier(notifyConfig)
	}
}

type NotifierInterface interface {
	Info(string) error
	Warning(string) error
	Error(string) error
}

type DummyNotifier struct {
	Config config.NotifyConfig
}

func NewDummyNotifier(notifyConfig config.NotifyConfig) (notify NotifierInterface) {
	return &DummyNotifier{Config: notifyConfig}
}

func (n *DummyNotifier) Info(message string) (err error) {
	return
}

func (n *DummyNotifier) Warning(message string) (err error) {
	return
}

func (n *DummyNotifier) Error(message string) (err error) {
	return
}
