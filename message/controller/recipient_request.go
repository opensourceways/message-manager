/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package controller

import "github.com/opensourceways/message-manager/message/app"

type syncUserInfoDTO struct {
	Mail          string `json:"mail"`
	Phone         string `json:"phone"`
	CountryCode   string `json:"country_code"`
	UserName      string `json:"user_name"`
	GiteeUserName string `json:"gitee_user_name"`
}

func (req *syncUserInfoDTO) toCmd() (cmd app.CmdToSyncUserInfo, err error) {
	cmd.Mail = req.Mail
	cmd.Phone = req.Phone
	cmd.CountryCode = req.CountryCode
	cmd.UserName = req.UserName
	cmd.GiteeUserName = req.GiteeUserName
	return
}

type newRecipientDTO struct {
	Name    string `gorm:"column:recipient_name" json:"recipient_id"`
	Mail    string `gorm:"column:mail" json:"mail"`
	Message string `gorm:"column:message" json:"message"`
	Phone   string `gorm:"column:phone" json:"phone"`
	Remark  string `gorm:"column:remark" json:"remark"`
}

func (req *newRecipientDTO) toCmd() (cmd app.CmdToAddRecipient, err error) {
	cmd.Name = req.Name
	cmd.Mail = req.Mail
	cmd.Message = req.Message
	cmd.Phone = req.Phone
	cmd.Remark = req.Remark
	return
}

type updateRecipientDTO struct {
	Id      string `gorm:"column:id"   json:"id"`
	Name    string `gorm:"column:recipient_name" json:"recipient_id"`
	Mail    string `gorm:"column:mail" json:"mail"`
	Message string `gorm:"column:message" json:"message"`
	Phone   string `gorm:"column:phone" json:"phone"`
	Remark  string `gorm:"column:remark" json:"remark"`
}

func (req *updateRecipientDTO) toCmd() (cmd app.CmdToUpdateRecipient, err error) {
	cmd.Id = req.Id
	cmd.Name = req.Name
	cmd.Mail = req.Mail
	cmd.Message = req.Message
	cmd.Phone = req.Phone
	cmd.Remark = req.Remark
	return
}

type deleteRecipientDTO struct {
	RecipientId string `gorm:"column:recipient_id" json:"recipient_id"`
}

func (req *deleteRecipientDTO) toCmd() (cmd app.CmdToDeleteRecipient, err error) {
	cmd.RecipientId = req.RecipientId
	return
}
