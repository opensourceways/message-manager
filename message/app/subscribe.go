/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package app

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"

	"github.com/opensourceways/message-manager/common/user"
	"github.com/opensourceways/message-manager/message/domain"
)

type MessageSubscribeAppService interface {
	GetAllSubsConfig(ctx *gin.Context) ([]MessageSubscribeDTO, error)
	GetSubsConfig(ctx *gin.Context) ([]MessageSubscribeDTO, int64, error)
	SaveFilter(ctx *gin.Context, cmd *CmdToGetSubscribe) error
	AddSubsConfig(ctx *gin.Context, cmd *CmdToAddSubscribe) ([]uint, error)
	RemoveSubsConfig(ctx *gin.Context, cmd *CmdToDeleteSubscribe) error
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

func (s *messageSubscribeAppService) GetAllSubsConfig(ctx *gin.Context) ([]MessageSubscribeDTO, error) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		return []MessageSubscribeDTO{}, xerrors.Errorf("get username failed, err:%v", err.Error())
	}
	response, err := s.messageSubscribeAdapter.GetAllSubsConfig(userName)
	if err != nil {
		return []MessageSubscribeDTO{}, err
	}
	return response, nil
}

func (s *messageSubscribeAppService) GetSubsConfig(ctx *gin.Context) ([]MessageSubscribeDTO,
	int64, error) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		return []MessageSubscribeDTO{}, 0, xerrors.Errorf("get username failed, err:%v", err.Error())
	}
	response, count, err := s.messageSubscribeAdapter.GetSubsConfig(userName)
	if err != nil {
		return []MessageSubscribeDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageSubscribeAppService) SaveFilter(ctx *gin.Context, cmd *CmdToGetSubscribe) error {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		return xerrors.Errorf("get username failed, err:%v", err.Error())
	}
	err = s.messageSubscribeAdapter.SaveFilter(*cmd, userName)
	if err != nil {
		return err
	}
	return nil
}

func (s *messageSubscribeAppService) AddSubsConfig(ctx *gin.Context,
	cmd *CmdToAddSubscribe) ([]uint, error) {

	if cmd.ModeName == "" || cmd.ModeFilter == nil || len(cmd.ModeFilter) == 0 {
		return []uint{}, xerrors.Errorf("必填项不能为空")
	}

	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		return []uint{}, xerrors.Errorf("get username failed, err:%v", err.Error())
	}

	data, err := s.messageSubscribeAdapter.AddSubsConfig(*cmd, userName)
	if err != nil {
		return []uint{}, xerrors.Errorf("add subs failed, err:%v", err)
	} else {
		return data, nil
	}
}

func (s *messageSubscribeAppService) RemoveSubsConfig(ctx *gin.Context, cmd *CmdToDeleteSubscribe) error {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		return xerrors.Errorf("get username failed, err:%v", err.Error())
	}

	err = s.messageSubscribeAdapter.RemoveSubsConfig(*cmd, userName)
	if err != nil {
		return xerrors.Errorf("remove subs failed, err:%v", err)
	} else {
		return nil
	}
}
