package domain

import (
	"github.com/opensourceways/message-manager/message/infrastructure"
)

type MessageListDTO = infrastructure.MessageListDTO
type MessagePushDTO = infrastructure.MessagePushDTO
type MessageRecipientDTO = infrastructure.MessageRecipientDTO
type MessageSubscribeDTO = infrastructure.MessageSubscribeDTO

type CmdToGetInnerMessageQuick = infrastructure.CmdToGetInnerMessageQuick
type CmdToGetInnerMessage = infrastructure.CmdToGetInnerMessage
type CmdToSetIsRead = infrastructure.CmdToSetIsRead
type CmdToAddPushConfig = infrastructure.CmdToAddPushConfig
type CmdToUpdatePushConfig = infrastructure.CmdToUpdatePushConfig
type CmdToDeletePushConfig = infrastructure.CmdToDeletePushConfig
type CmdToAddRecipient = infrastructure.CmdToAddRecipient
type CmdToUpdateRecipient = infrastructure.CmdToUpdateRecipient
type CmdToDeleteRecipient = infrastructure.CmdToDeleteRecipient
type CmdToSyncUserInfo = infrastructure.CmdToSyncUserInfo
type CmdToGetSubscribe = infrastructure.CmdToGetSubscribe
type CmdToAddSubscribe = infrastructure.CmdToAddSubscribe
type CmdToDeleteSubscribe = infrastructure.CmdToDeleteSubscribe
