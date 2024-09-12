/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package app

import (
	"github.com/opensourceways/message-manager/message/domain"
)

type MessageListDTO = domain.MessageListDO
type MessagePushDTO = domain.MessagePushDO
type MessageRecipientDTO = domain.MessageRecipientDO
type MessageSubscribeDTO = domain.MessageSubscribeDO
type CountDTO = domain.CountDO

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
type CmdToUpdateSubscribe = domain.CmdToUpdateSubscribe
type CmdToDeleteSubscribe = domain.CmdToDeleteSubscribe
