/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package server

import (
	"github.com/gin-gonic/gin"

	"github.com/opensourceways/message-manager/message/app"
	messagectl "github.com/opensourceways/message-manager/message/controller"
	"github.com/opensourceways/message-manager/message/infrastructure"
)

func initMessage(services *allServices) error {
	services.MessageListAppService = app.NewMessageListAppService(
		infrastructure.MessageListAdapter(),
	)
	services.MessagePushAppService = app.NewMessagePushAppService(
		infrastructure.MessagePushAdapter(),
	)
	services.MessageRecipientAppService = app.NewMessageRecipientAppService(
		infrastructure.MessageRecipientAdapter(),
	)
	services.MessageSubscribeAppService = app.NewMessageSubscribeAppService(
		infrastructure.MessageSubscribeAdapter(),
	)

	return nil
}

// setRouteOfMessage is registering controller of moderation in api
func setRouteOfMessage(rg *gin.Engine, services *allServices) {
	messagectl.AddRouterForMessageListController(
		rg,
		services.MessageListAppService,
	)
	messagectl.AddRouterForMessagePushController(
		rg,
		services.MessagePushAppService,
	)
	messagectl.AddRouterForMessageRecipientController(
		rg,
		services.MessageRecipientAppService,
	)
	messagectl.AddRouterForMessageSubscribeController(
		rg,
		services.MessageSubscribeAppService,
	)
}
