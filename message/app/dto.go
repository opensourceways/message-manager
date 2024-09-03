package app

import (
	"github.com/opensourceways/message-manager/message/domain"
)

type MessageListDTO = domain.MessageListDTO
type MessagePushDTO = domain.MessagePushDTO
type MessageRecipientDTO = domain.MessageRecipientDTO
type MessageSubscribeDTO = domain.MessageSubscribeDTO

type CmdToGetInnerMessageQuick = domain.CmdToGetInnerMessageQuick
type CmdToGetInnerMessage = domain.CmdToGetInnerMessage
type CmdToSetIsRead = domain.CmdToSetIsRead
type CmdToAddPushConfig = domain.CmdToAddPushConfig
type CmdToUpdatePushConfig = domain.CmdToUpdatePushConfig
type CmdToDeletePushConfig = domain.CmdToDeletePushConfig
type CmdToAddRecipient = domain.CmdToAddRecipient
type CmdToUpdateRecipient = domain.CmdToUpdateRecipient
type CmdToDeleteRecipient = domain.CmdToDeleteRecipient
type CmdToSyncUserInfo = domain.CmdToSyncUserInfo
type CmdToGetSubscribe = domain.CmdToGetSubscribe
type CmdToAddSubscribe = domain.CmdToAddSubscribe
type CmdToDeleteSubscribe = domain.CmdToDeleteSubscribe
