package notify

import (
	"github.com/andybons/hipchat"
	"github.com/crocos/rds-testrunner/config"
)

type HipChatNotifier struct {
	Config config.NotifyConfig
	Client hipchat.Client
}

func NewHipChatNotifier(notifyConfig config.NotifyConfig) (notifier NotifierInterface) {
	client := hipchat.Client{
		AuthToken: notifyConfig.Token,
	}

	notifier = &HipChatNotifier{
		Config: notifyConfig,
		Client: client,
	}

	return
}

func (n *HipChatNotifier) Message(message string, color string) (err error) {
	req := hipchat.MessageRequest{
		RoomId:        n.Config.Room,
		From:          n.Config.Name,
		Message:       message,
		Color:         color,
		MessageFormat: hipchat.FormatText,
		Notify:        true,
	}

	err = n.Client.PostMessage(req)

	return
}

func (n *HipChatNotifier) Info(message string) (err error) {
	err = n.Message(message, hipchat.ColorGreen)
	return
}

func (n *HipChatNotifier) Warning(message string) (err error) {
	err = n.Message(message, hipchat.ColorYellow)
	return
}

func (n *HipChatNotifier) Error(message string) (err error) {
	err = n.Message(message, hipchat.ColorRed)
	return
}
