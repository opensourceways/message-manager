/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package domain

type MessageSubscribeAdapter interface {
	GetAllSubsConfig(userName string) ([]MessageSubscribeDO, error)
	GetSubsConfig(userName string) ([]MessageSubscribeDOWithPushConfig, int64, error)
	AddSubsConfig(cmd CmdToAddSubscribe, userName string) ([]uint, error)
	UpdateSubsConfig(cmd CmdToUpdateSubscribe, userName string) error
	RemoveSubsConfig(cmd CmdToDeleteSubscribe, userName string) error
}
