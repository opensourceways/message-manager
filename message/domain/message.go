/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package domain

type MessageListAdapter interface {
	GetInnerMessageQuick(cmd CmdToGetInnerMessageQuick, serName string) (
		[]MessageListDTO, int64, error)
	GetInnerMessage(cmd CmdToGetInnerMessage, userName string) ([]MessageListDTO, int64, error)
	CountAllUnReadMessage(userName string) (int64, error)
	SetMessageIsRead(source, eventId string) error
	RemoveMessage(source, eventId string) error
}
