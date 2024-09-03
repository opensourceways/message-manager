package app

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"

	"github.com/opensourceways/message-manager/common/user"
	"github.com/opensourceways/message-manager/message/domain"
)

type MessageListAppService interface {
	GetInnerMessageQuick(ctx *gin.Context, cmd *CmdToGetInnerMessageQuick) ([]MessageListDTO,
		int64, error)
	GetInnerMessage(ctx *gin.Context, cmd *CmdToGetInnerMessage) ([]MessageListDTO, int64, error)
	CountAllUnReadMessage(ctx *gin.Context) (int64, error)
	SetMessageIsRead(ctx *gin.Context, cmd *CmdToSetIsRead) error
	RemoveMessage(ctx *gin.Context, cmd *CmdToSetIsRead) error
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

func (s *messageListAppService) GetInnerMessageQuick(ctx *gin.Context,
	cmd *CmdToGetInnerMessageQuick) ([]MessageListDTO, int64, error) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		return []MessageListDTO{}, 0, xerrors.Errorf("get username failed, err:%v", err.Error())
	}

	response, count, err := s.messageListAdapter.GetInnerMessageQuick(*cmd, userName)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) GetInnerMessage(ctx *gin.Context,
	cmd *CmdToGetInnerMessage) ([]MessageListDTO, int64, error) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		return []MessageListDTO{}, 0, xerrors.Errorf("get username failed, err:%v", err.Error())
	}

	response, count, err := s.messageListAdapter.GetInnerMessage(*cmd, userName)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) CountAllUnReadMessage(c *gin.Context) (int64, error) {
	userName, err := user.GetEulerUserName(c)
	if err != nil {
		return 0, xerrors.Errorf("get username failed, err:%v", err.Error())
	}
	count, err := s.messageListAdapter.CountAllUnReadMessage(userName)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *messageListAppService) SetMessageIsRead(ctx *gin.Context, cmd *CmdToSetIsRead) error {
	if err := s.messageListAdapter.SetMessageIsRead(cmd.Source, cmd.EventId); err != nil {
		return xerrors.Errorf("set message is_read failed, err:%v", err.Error())
	}
	return nil
}

func (s *messageListAppService) RemoveMessage(ctx *gin.Context, cmd *CmdToSetIsRead) error {
	if err := s.messageListAdapter.RemoveMessage(cmd.Source, cmd.EventId); err != nil {
		return xerrors.Errorf("set message is_read failed, err:%v", err.Error())
	}
	return nil
}
