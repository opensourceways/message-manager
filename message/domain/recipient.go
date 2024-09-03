/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package domain

type MessageRecipientAdapter interface {
	GetRecipientConfig(countPerPage, pageNum int, userName string) ([]MessageRecipientDO, int64,
		error)
	AddRecipientConfig(cmd CmdToAddRecipient, userName string) error
	UpdateRecipientConfig(cmd CmdToUpdateRecipient, userName string) error
	RemoveRecipientConfig(cmd CmdToDeleteRecipient, userName string) error
	SyncUserInfo(cmd CmdToSyncUserInfo) (uint, error)
}
