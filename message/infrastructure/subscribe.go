/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package infrastructure

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/opensourceways/message-manager/common/postgresql"
)

func MessageSubscribeAdapter() *messageSubscribeAdapter {
	return &messageSubscribeAdapter{}
}

type messageSubscribeAdapter struct{}

func (ctl *messageSubscribeAdapter) GetAllSubsConfig(userName string) ([]MessageSubscribeDAO, error) {

	var response []MessageSubscribeDAO
	query := postgresql.DB().Table("message_center.subscribe_config").
		Where("user_name = ? OR user_name IS NULL", userName).
		Where(gorm.Expr("subscribe_config.is_deleted = ?", false))

	if result := query.Order("subscribe_config.id").Find(&response); result.Error != nil {
		return []MessageSubscribeDAO{}, xerrors.Errorf("查询失败")
	}

	return response, nil
}

func (ctl *messageSubscribeAdapter) GetSubsConfig(userName string) ([]MessageSubscribeDAO, int64, error) {
	var response []MessageSubscribeDAO

	query := postgresql.DB().Table("message_center.subscribe_config").
		Where(gorm.Expr("subscribe_config.is_deleted = ?", false)).
		Where("subscribe_config.user_name = ?", userName)

	var Count int64
	query.Count(&Count)
	if result := query.Order("subscribe_config.id").Find(&response); result.Error != nil {
		return []MessageSubscribeDAO{}, 0, xerrors.Errorf("查询失败")
	}

	return response, Count, nil
}

func (ctl *messageSubscribeAdapter) SaveFilter(cmd CmdToGetSubscribe, userName string) error {
	if userName == "" {
		return xerrors.Errorf("用户名为空")
	}
	var modeFilter datatypes.JSON
	modeFilter, _ = TransToDbFormat(cmd.Source, cmd.EventType, cmd)
	isDefault := new(bool)
	*isDefault = false
	jsonFilter, err := json.Marshal(cmd)
	if err != nil {
		return xerrors.Errorf("marshal data failed, err:%v", err)
	}

	result := postgresql.DB().Table("message_center.subscribe_config").
		Create(MessageSubscribeDAOWithoutId{
			Source:      cmd.Source,
			EventType:   cmd.EventType,
			SpecVersion: cmd.SpecVersion,
			ModeFilter:  modeFilter,
			ModeName:    cmd.ModeName,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UserName:    userName,
			IsDefault:   isDefault,
			WebFilter:   jsonFilter,
		})
	if result.Error != nil {
		return xerrors.Errorf("保存配置失败")
	} else {
		return nil
	}
}

func (ctl *messageSubscribeAdapter) AddSubsConfig(cmd CmdToAddSubscribe, userName string) ([]uint, error) {
	var existData MessageSubscribeDAO

	if result := postgresql.DB().Table("message_center.subscribe_config").
		Where(gorm.Expr("is_deleted = ?", false)).
		Where("source = ? AND mode_name = ?", cmd.Source, cmd.ModeName).
		Where("user_name = ?", userName).
		Scan(&existData); result.RowsAffected != 0 {
		return []uint{}, xerrors.Errorf("新增配置失败")
	}

	var subscribeIds []uint
	lType := strings.Split(cmd.EventType, ",")
	for _, et := range lType {
		result := postgresql.DB().Table("message_center.subscribe_config").
			Create(MessageSubscribeDAO{
				Source:      cmd.Source,
				EventType:   et,
				SpecVersion: cmd.SpecVersion,
				ModeName:    cmd.ModeName,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				UserName:    userName,
			})
		if result.Error != nil {
			return []uint{}, xerrors.Errorf("新增配置失败")
		}

		var id uint
		postgresql.DB().Table("message_center.subscribe_config").
			Debug().
			Where(gorm.Expr("is_deleted = ?", false)).
			Where("source = ? AND event_type = ? AND mode_name = ? AND user_name = ?",
				cmd.Source, et, cmd.ModeName, userName).Select("id").Scan(&id)
		subscribeIds = append(subscribeIds, id)
	}
	return subscribeIds, nil
}

func (ctl *messageSubscribeAdapter) UpdateSubsConfig(cmd CmdToUpdateSubscribe,
	userName string) error {
	if result := postgresql.DB().Table("message_center.subscribe_config").
		Where(gorm.Expr("is_deleted = ?", false)).
		Where("source = ? AND mode_name = ?", cmd.Source, cmd.OldName).
		Where("user_name = ?", userName).
		Update("mode_name", cmd.NewName); result.Error != nil {
		logrus.Errorf("update subscribe config failed, err:%v", result.Error)
		return xerrors.Errorf("更新配置失败")
	}
	return nil
}

func (ctl *messageSubscribeAdapter) RemoveSubsConfig(cmd CmdToDeleteSubscribe, userName string) error {
	lType := strings.Split(cmd.EventType, ",")

	if result := postgresql.DB().Table("message_center.subscribe_config").
		Where(gorm.Expr("is_deleted = ?", false)).
		Where("source = ? AND event_type IN ? AND mode_name = ?", cmd.Source,
			lType, cmd.ModeName).
		Where("user_name = ?", userName).
		Update("is_deleted", true); result.Error != nil {
		return xerrors.Errorf("删除配置失败")
	}
	return nil
}
