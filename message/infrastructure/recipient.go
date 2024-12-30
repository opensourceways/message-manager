/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package infrastructure

import (
	"time"

	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/opensourceways/message-manager/common/postgresql"
)

func MessageRecipientAdapter() *messageRecipientAdapter {
	return &messageRecipientAdapter{}
}

type messageRecipientAdapter struct{}

type RecipientController struct {
	Name          string    `gorm:"column:recipient_name" json:"recipient_id"`
	Mail          string    `gorm:"column:mail" json:"mail"`
	Message       string    `gorm:"column:message" json:"message"`
	Phone         string    `gorm:"column:phone" json:"phone"`
	Remark        string    `gorm:"column:remark" json:"remark"`
	UserName      string    `gorm:"column:user_id"  json:"user_id"`
	GiteeUserName string    `gorm:"column:gitee_user_name" json:"gitee_user_name"`
	IsDeleted     bool      `gorm:"column:is_deleted" json:"is_deleted"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at" swaggerignore:"true"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at" swaggerignore:"true"`
}

func getTable() *gorm.DB {
	return postgresql.DB().Table("message_center.recipient_config").
		Where(gorm.Expr("is_deleted = ?", false))
}

func (ctl *messageRecipientAdapter) GetRecipientConfig(countPerPage, pageNum int, userName string) (
	[]MessageRecipientDAO, int64, error) {
	var response []MessageRecipientDAO
	var Count int64
	getTable().Where("user_id = ?", userName).Count(&Count)

	offsetNum := (pageNum - 1) * countPerPage

	if result := getTable().Where("user_id = ?", userName).Limit(countPerPage).Offset(offsetNum).
		Order("recipient_config.created_at DESC").
		Find(&response); result.Error != nil {
		return []MessageRecipientDAO{}, 0, xerrors.Errorf("get recipient config failed, err:%v", result.Error.Error())
	}
	return response, Count, nil
}

func (ctl *messageRecipientAdapter) AddRecipientConfig(cmd CmdToAddRecipient,
	userName string) error {
	var existData MessageRecipientDAO
	if result := getTable().Where("recipient_name = ? AND user_id = ?", cmd.Name, userName).
		Scan(&existData); result.RowsAffected != 0 {
		return xerrors.Errorf("接收人姓名不能相同")
	}

	if result := getTable().Create(MessageRecipientDAO{
		Name:      cmd.Name,
		Mail:      cmd.Mail,
		Message:   cmd.Message,
		Phone:     cmd.Phone,
		Remark:    cmd.Remark,
		UserName:  userName,
		IsDeleted: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}); result.Error != nil {
		return xerrors.Errorf("add new recipient failed, err:%v", result.Error)
	}
	return nil
}

func (ctl *messageRecipientAdapter) UpdateRecipientConfig(cmd CmdToUpdateRecipient,
	userName string) error {
	if result := getTable().Where("id = ? AND user_id = ?", cmd.Id, userName).
		Updates(RecipientController{
			Name:      cmd.Name,
			Mail:      cmd.Mail,
			Message:   cmd.Message,
			Phone:     cmd.Phone,
			Remark:    cmd.Remark,
			UpdatedAt: time.Now(),
		}); result.Error != nil {
		return xerrors.Errorf("update recipient config failed, err:%v", result.Error)
	}
	return nil
}

func (ctl *messageRecipientAdapter) RemoveRecipientConfig(cmd CmdToDeleteRecipient,
	userName string) error {
	if result := getTable().Where("id = ? AND user_id = ?", cmd.RecipientId, userName).
		Update("is_deleted", true); result.Error != nil || result.RowsAffected == 0 {
		return xerrors.Errorf("删除配置失败, err:%v", result.Error)
	}
	return nil
}

func (ctl *messageRecipientAdapter) SyncUserInfo(cmd CmdToSyncUserInfo) (uint, error) {
	var oldInfo RecipientController

	if result := getTable().
		Where("user_id = ?", cmd.UserName).
		Scan(&oldInfo); result.RowsAffected != 0 {
		newInfo := &oldInfo
		newInfo.Mail = cmd.Mail
		newInfo.Message = cmd.CountryCode + cmd.Phone
		newInfo.Phone = cmd.CountryCode + cmd.Phone
		newInfo.UserName = cmd.UserName
		newInfo.GiteeUserName = cmd.GiteeUserName
		getTable().Where("user_id = ?", cmd.UserName).Save(&newInfo)
	} else {
		newInfo := RecipientController{
			Mail:          cmd.Mail,
			Message:       cmd.CountryCode + cmd.Phone,
			Phone:         cmd.CountryCode + cmd.Phone,
			UserName:      cmd.UserName,
			GiteeUserName: cmd.GiteeUserName,
		}
		getTable().Create(&newInfo)
	}
	var id uint
	getTable().Where(gorm.Expr("is_deleted = ?", false)).
		Where("user_id = ?", cmd.UserName).Select("id").Scan(&id)

	return id, nil
}
