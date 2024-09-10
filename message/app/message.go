/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package app

import (
	"golang.org/x/xerrors"

	"github.com/opensourceways/message-manager/message/domain"
)

type MessageListAppService interface {
	GetInnerMessageQuick(userName string, cmd *CmdToGetInnerMessageQuick) ([]MessageListDTO,
		int64, error)
	GetInnerMessage(userName string, cmd *CmdToGetInnerMessage) ([]MessageListDTO, int64, error)
	CountAllUnReadMessage(userName string) ([]CountDTO, error)
	SetMessageIsRead(cmd *CmdToSetIsRead) error
	RemoveMessage(cmd *CmdToSetIsRead) error
}

func NewMessageListAppService(
	messageListAdapter domain.MessageListAdapter,
) MessageListAppService {
	return &messageListAppService{
		messageListAdapter: messageListAdapter,
	}
}

type messageListAppService struct {
	messageListAdapter domain.MessageListAdapter
}

func (s *messageListAppService) GetInnerMessageQuick(userName string,
	cmd *CmdToGetInnerMessageQuick) ([]MessageListDTO, int64, error) {
	response, count, err := s.messageListAdapter.GetInnerMessageQuick(*cmd, userName)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) GetInnerMessage(userName string,
	cmd *CmdToGetInnerMessage) ([]MessageListDTO, int64, error) {
	response, count, err := s.messageListAdapter.GetInnerMessage(*cmd, userName)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) CountAllUnReadMessage(userName string) ([]CountDTO, error) {
	count, err := s.messageListAdapter.CountAllUnReadMessage(userName)
	if err != nil {
		return []CountDTO{}, err
	}
	return count, nil
}

func (s *messageListAppService) SetMessageIsRead(cmd *CmdToSetIsRead) error {
	if err := s.messageListAdapter.SetMessageIsRead(cmd.Source, cmd.EventId); err != nil {
		return xerrors.Errorf("set message is_read failed, err:%v", err.Error())
	}
	return nil
}

func (s *messageListAppService) RemoveMessage(cmd *CmdToSetIsRead) error {
	if err := s.messageListAdapter.RemoveMessage(cmd.Source, cmd.EventId); err != nil {
		return xerrors.Errorf("set message is_read failed, err:%v", err.Error())
	}
	return nil
}
