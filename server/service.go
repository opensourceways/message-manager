package server

import (
	"github.com/opensourceways/message-manager/message/app"
)

type allServices struct {
	MessageListAppService      app.MessageListAppService
	MessagePushAppService      app.MessagePushAppService
	MessageRecipientAppService app.MessageRecipientAppService
	MessageSubscribeAppService app.MessageSubscribeAppService
}

// initServices init All service
func initServices() (services allServices, err error) {
	if err = initMessage(&services); err != nil {
		return
	}
	return
}
