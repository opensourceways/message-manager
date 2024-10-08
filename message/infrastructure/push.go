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

func MessagePushAdapter() *messagePushAdapter {
	return &messagePushAdapter{}
}

type messagePushAdapter struct{}

func (s *messagePushAdapter) GetPushConfig(subsIds []string, countPerPage, pageNum int,
	userName string) ([]MessagePushDAO, error) {

	query := postgresql.DB().Table("message_center.push_config").
		Where(gorm.Expr("push_config.is_deleted = ?", false))
	offsetNum := (pageNum - 1) * countPerPage
	query = query.Select("push_config.*").
		Joins("JOIN message_center.recipient_config ON recipient_config.id = push_config."+
			"recipient_id").
		Where(gorm.Expr("recipient_config.is_deleted = ?", false)).
		Where("recipient_config.user_id = ?", userName)

	var response []MessagePushDAO
	if result := query.Limit(countPerPage).Offset(offsetNum).
		Find(&response, "subscribe_id IN ?", subsIds); result.Error != nil {
		return []MessagePushDAO{}, xerrors.Errorf("查询失败")
	}

	return response, nil
}

func (s *messagePushAdapter) AddPushConfig(cmd CmdToAddPushConfig) error {

	var existData MessagePushDAO
	if result := postgresql.DB().Table("message_center.push_config").
		Where(gorm.Expr("is_deleted = ?", false)).
		Where("subscribe_id = ? AND recipient_id = ?", cmd.SubscribeId, cmd.RecipientId).
		Scan(&existData); result.RowsAffected != 0 {
		return xerrors.Errorf("新增配置失败，配置已存在")
	}

	if result := postgresql.DB().Table("message_center.push_config").
		Create(MessagePushDAO{
			SubscribeId:      cmd.SubscribeId,
			RecipientId:      cmd.RecipientId,
			NeedMessage:      &cmd.NeedMessage,
			NeedPhone:        &cmd.NeedPhone,
			NeedMail:         &cmd.NeedMail,
			NeedInnerMessage: &cmd.NeedInnerMessage,
			IsDeleted:        false,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}); result.Error != nil {
		return xerrors.Errorf("新增配置失败，err:%v", result.Error)
	}
	return nil
}

func (s *messagePushAdapter) UpdatePushConfig(cmd CmdToUpdatePushConfig) error {
	if result := postgresql.DB().Table("message_center.push_config").Debug().
		Where("is_deleted = ?", false).
		Where("subscribe_id IN ? AND recipient_id = ?", cmd.SubscribeId, cmd.RecipientId).
		Updates(&MessagePushDAO{
			NeedMessage:      &cmd.NeedMessage,
			NeedPhone:        &cmd.NeedPhone,
			NeedMail:         &cmd.NeedMail,
			NeedInnerMessage: &cmd.NeedInnerMessage,
			UpdatedAt:        time.Now(),
		}); result.Error != nil {
		return xerrors.Errorf("更新配置失败，err:%v", result.Error)
	}
	return nil
}

func (s *messagePushAdapter) RemovePushConfig(cmd CmdToDeletePushConfig) error {
	if result := postgresql.DB().Table("message_center.push_config").
		Where(gorm.Expr("is_deleted IS NULL OR is_deleted = ?", false)).
		Where("subscribe_id = ? AND recipient_id = ?", cmd.SubscribeId, cmd.RecipientId).
		Update("is_deleted", true); result.Error != nil {
		return xerrors.Errorf("删除配置失败，err:%v", result.Error)
	}
	return nil
}
