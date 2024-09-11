/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package app

import (
	"golang.org/x/xerrors"

	"github.com/opensourceways/message-manager/message/domain"
)

type MessagePushAppService interface {
	GetPushConfig(countPerPage, pageNum int, userName string,
		subsIds []string) ([]MessagePushDTO, error)
	AddPushConfig(cmd *CmdToAddPushConfig) error
	UpdatePushConfig(cmd *CmdToUpdatePushConfig) error
	RemovePushConfig(cmd *CmdToDeletePushConfig) error
}

func NewMessagePushAppService(
	messagePushAdapter domain.MessagePushAdapter,
) MessagePushAppService {
	return &messagePushAppService{
		messagePushAdapter: messagePushAdapter,
	}
}

type messagePushAppService struct {
	messagePushAdapter domain.MessagePushAdapter
}

func (s *messagePushAppService) GetPushConfig(countPerPage, pageNum int, userName string,
	subsIds []string) ([]MessagePushDTO, error) {

	data, err := s.messagePushAdapter.GetPushConfig(subsIds, countPerPage, pageNum, userName)
	if err != nil {
		return []MessagePushDTO{}, err
	}
	return data, nil
}

func (s *messagePushAppService) AddPushConfig(cmd *CmdToAddPushConfig) error {
	if err := s.messagePushAdapter.AddPushConfig(*cmd); err != nil {
		return xerrors.Errorf("add message push config failed, err:%v", err.Error())
	}
	return nil
}

func (s *messagePushAppService) UpdatePushConfig(cmd *CmdToUpdatePushConfig) error {
	if err := s.messagePushAdapter.UpdatePushConfig(*cmd); err != nil {
		return xerrors.Errorf("update message push config failed, err:%v", err.Error())
	}
	return nil
}

func (s *messagePushAppService) RemovePushConfig(cmd *CmdToDeletePushConfig) error {
	if err := s.messagePushAdapter.RemovePushConfig(*cmd); err != nil {
		return xerrors.Errorf("remove message push config failed, err:%v", err.Error())
	}
	return nil
}
