/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package domain

import (
	"github.com/opensourceways/message-manager/message/infrastructure"
)

type MessageListDO = infrastructure.MessageListDAO
type MessagePushDO = infrastructure.MessagePushDAO
type MessageRecipientDO = infrastructure.MessageRecipientDAO
type MessageSubscribeDO = infrastructure.MessageSubscribeDAO
type MessageSubscribeDOWithPushConfig = infrastructure.MessageSubscribeDAOWithPushConfig
type CountDO = infrastructure.CountDAO

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
type CmdToUpdateSubscribe = infrastructure.CmdToUpdateSubscribe
type CmdToDeleteSubscribe = infrastructure.CmdToDeleteSubscribe
