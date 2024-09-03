/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package domain

type MessagePushAdapter interface {
	GetPushConfig(subsIds []string, countPerPage, pageNum int,
		userName string) ([]MessagePushDO, error)
	AddPushConfig(cmd CmdToAddPushConfig) error
	UpdatePushConfig(cmd CmdToUpdatePushConfig) error
	RemovePushConfig(cmd CmdToDeletePushConfig) error
}
