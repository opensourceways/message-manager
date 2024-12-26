/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package app

import (
	"golang.org/x/xerrors"

	"github.com/opensourceways/message-manager/message/domain"
)

type MessageSubscribeAppService interface {
	GetAllSubsConfig(userName string) ([]MessageSubscribeDTO, error)
	GetSubsConfig(userName string) ([]MessageSubscribeDTOWithPushConfig, int64, error)
	AddSubsConfig(userName string, cmd *CmdToAddSubscribe) ([]uint, error)
	UpdateSubsConfig(userName string, cmd *CmdToUpdateSubscribe) error
	RemoveSubsConfig(userName string, cmd *CmdToDeleteSubscribe) error
}

func NewMessageSubscribeAppService(
	messageSubscribeAdapter domain.MessageSubscribeAdapter,
) MessageSubscribeAppService {
	return &messageSubscribeAppService{
		messageSubscribeAdapter: messageSubscribeAdapter,
	}
}

type messageSubscribeAppService struct {
	messageSubscribeAdapter domain.MessageSubscribeAdapter
}

func (s *messageSubscribeAppService) GetAllSubsConfig(userName string) ([]MessageSubscribeDTO, error) {
	response, err := s.messageSubscribeAdapter.GetAllSubsConfig(userName)
	if err != nil {
		return []MessageSubscribeDTO{}, err
	}
	return response, nil
}

func (s *messageSubscribeAppService) GetSubsConfig(userName string) ([]MessageSubscribeDTOWithPushConfig,
	int64, error) {
	response, count, err := s.messageSubscribeAdapter.GetSubsConfig(userName)
	if err != nil {
		return []MessageSubscribeDTOWithPushConfig{}, 0, err
	}
	return response, count, nil
}

func (s *messageSubscribeAppService) AddSubsConfig(userName string,
	cmd *CmdToAddSubscribe) ([]uint, error) {

	if cmd.ModeName == "" || cmd.ModeFilter == nil || len(cmd.ModeFilter) == 0 {
		return []uint{}, xerrors.Errorf("必填项不能为空")
	}

	data, err := s.messageSubscribeAdapter.AddSubsConfig(*cmd, userName)
	if err != nil {
		return []uint{}, xerrors.Errorf("add subs failed, err:%v", err)
	} else {
		return data, nil
	}
}

func (s *messageSubscribeAppService) UpdateSubsConfig(userName string,
	cmd *CmdToUpdateSubscribe) error {
	err := s.messageSubscribeAdapter.UpdateSubsConfig(*cmd, userName)
	if err != nil {
		return xerrors.Errorf("update subs failed, err:%v", err)
	} else {
		return nil
	}
}

func (s *messageSubscribeAppService) RemoveSubsConfig(userName string,
	cmd *CmdToDeleteSubscribe) error {
	err := s.messageSubscribeAdapter.RemoveSubsConfig(*cmd, userName)
	if err != nil {
		return xerrors.Errorf("remove subs failed, err:%v", err)
	} else {
		return nil
	}
}
