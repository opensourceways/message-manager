/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package controller

import "github.com/opensourceways/message-manager/message/app"

type deletePushConfigDTO struct {
	SubscribeId int   `gorm:"column:subscribe_id" json:"subscribe_id"`
	RecipientId int64 `gorm:"column:recipient_id" json:"recipient_id"`
}

func (req *deletePushConfigDTO) toCmd() (cmd app.CmdToDeletePushConfig, err error) {
	cmd.SubscribeId = req.SubscribeId
	cmd.RecipientId = req.RecipientId
	return cmd, nil
}

type newPushConfigDTO struct {
	SubscribeId      int   `json:"subscribe_id"`
	RecipientId      int64 `json:"recipient_id"`
	NeedMessage      bool  `json:"need_message"`
	NeedPhone        bool  `json:"need_phone"`
	NeedMail         bool  `json:"need_mail"`
	NeedInnerMessage bool  `json:"need_inner_message"`
}

func (req *newPushConfigDTO) toCmd() (cmd app.CmdToAddPushConfig, err error) {
	cmd.SubscribeId = req.SubscribeId
	cmd.RecipientId = req.RecipientId
	cmd.NeedMessage = req.NeedMessage
	cmd.NeedPhone = req.NeedPhone
	cmd.NeedMail = req.NeedMail
	cmd.NeedInnerMessage = req.NeedInnerMessage
	return cmd, nil
}

type updatePushConfigDTO struct {
	SubscribeId      []int  `json:"subscribe_id"`
	RecipientId      string `json:"recipient_id"`
	NeedMessage      bool   `json:"need_message"`
	NeedPhone        bool   `json:"need_phone"`
	NeedMail         bool   `json:"need_mail"`
	NeedInnerMessage bool   `json:"need_inner_message"`
}

func (req *updatePushConfigDTO) toCmd() (cmd app.CmdToUpdatePushConfig, err error) {
	cmd.SubscribeId = req.SubscribeId
	cmd.RecipientId = req.RecipientId
	cmd.NeedMessage = req.NeedMessage
	cmd.NeedPhone = req.NeedPhone
	cmd.NeedMail = req.NeedMail
	cmd.NeedInnerMessage = req.NeedInnerMessage
	return cmd, nil
}
