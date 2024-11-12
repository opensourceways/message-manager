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

	GetAllToDoMessage(userName string, giteeUsername string, isDone bool,
		pageNum, countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error)
	GetAllAboutMessage(userName string, giteeUsername string, isBot bool,
		pageNum, countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error)
	GetAllWatchMessage(userName string, giteeUsername string,
		pageNum, countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error)

	CountAllMessage(userName string, giteeUsername string) (CountDataDTO, error)

	GetForumSystemMessage(userName string, pageNum, countPerPage int,
		startTime string, isRead bool) ([]MessageListDTO, int64, error)
	GetForumAboutMessage(userName string, isBot bool, pageNum,
		countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error)
	GetMeetingToDoMessage(userName string, giteeUsername string, filter int,
		pageNum, countPerPage int, isRead bool) ([]MessageListDTO, int64, error)
	GetCVEToDoMessage(userName string, giteeUsername string, isDone bool,
		pageNum, countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error)
	GetCVEMessage(userName string, giteeUsername string,
		pageNum, countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error)
	GetIssueToDoMessage(userName string, giteeUsername string, isDone bool,
		pageNum, countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error)
	GetPullRequestToDoMessage(userName string, giteeUsername string, isDone bool,
		pageNum, countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error)
	GetGiteeAboutMessage(userName string, giteeUsername string, isBot bool,
		pageNum, countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error)
	GetGiteeMessage(userName string, giteeUsername string, pageNum,
		countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error)
	GetEurMessage(userName string, pageNum, countPerPage int,
		startTime string, isRead bool) ([]MessageListDTO, int64, error)
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

func (s *messageListAppService) GetAllToDoMessage(userName string, giteeUsername string,
	isDone bool, pageNum, countPerPage int, startTime string, isRead bool) (
	[]MessageListDTO, int64, error) {
	response, count, err := s.messageListAdapter.GetAllToDoMessage(userName, giteeUsername,
		isDone, pageNum, countPerPage, startTime, isRead)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) GetAllAboutMessage(userName string, giteeUsername string,
	isBot bool, pageNum, countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error) {
	response, count, err := s.messageListAdapter.GetAllAboutMessage(userName, giteeUsername,
		isBot, pageNum, countPerPage, startTime, isRead)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) GetAllWatchMessage(userName string, giteeUsername string,
	pageNum, countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error) {
	response, count, err := s.messageListAdapter.GetAllWatchMessage(userName, giteeUsername,
		pageNum, countPerPage, startTime, isRead)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) GetForumSystemMessage(userName string, pageNum, countPerPage int,
	startTime string, isRead bool) ([]MessageListDTO, int64, error) {
	response, count, err := s.messageListAdapter.GetForumSystemMessage(userName, pageNum,
		countPerPage, startTime, isRead)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) GetForumAboutMessage(userName string, isBot bool, pageNum,
	countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error) {
	response, count, err := s.messageListAdapter.GetForumAboutMessage(userName, isBot, pageNum,
		countPerPage, startTime, isRead)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) GetMeetingToDoMessage(userName string, giteeUsername string,
	filter int, pageNum, countPerPage int, isRead bool) ([]MessageListDTO, int64, error) {
	response, count, err := s.messageListAdapter.GetMeetingToDoMessage(userName, giteeUsername,
		filter, pageNum, countPerPage, isRead)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) GetCVEToDoMessage(userName string, giteeUsername string,
	isDone bool, pageNum, countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error) {
	response, count, err := s.messageListAdapter.GetCVEToDoMessage(userName, giteeUsername,
		isDone, pageNum, countPerPage, startTime, isRead)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) GetCVEMessage(userName string, giteeUsername string, pageNum,
	countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error) {
	response, count, err := s.messageListAdapter.GetCVEMessage(userName, giteeUsername, pageNum,
		countPerPage, startTime, isRead)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) GetIssueToDoMessage(userName string, giteeUsername string,
	isDone bool, pageNum, countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error) {
	response, count, err := s.messageListAdapter.GetIssueToDoMessage(userName, giteeUsername,
		isDone, pageNum, countPerPage, startTime, isRead)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) GetPullRequestToDoMessage(userName string, giteeUsername string,
	isDone bool, pageNum, countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error) {
	response, count, err := s.messageListAdapter.GetPullRequestToDoMessage(userName,
		giteeUsername, isDone, pageNum, countPerPage, startTime, isRead)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) GetGiteeAboutMessage(userName string, giteeUsername string,
	isBot bool, pageNum, countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error) {
	response, count, err := s.messageListAdapter.GetGiteeAboutMessage(userName, giteeUsername,
		isBot, pageNum, countPerPage, startTime, isRead)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) GetGiteeMessage(userName string, giteeUsername string, pageNum,
	countPerPage int, startTime string, isRead bool) ([]MessageListDTO, int64, error) {
	response, count, err := s.messageListAdapter.GetGiteeMessage(userName, giteeUsername,
		pageNum, countPerPage, startTime, isRead)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) GetEurMessage(userName string, pageNum, countPerPage int,
	startTime string, isRead bool) ([]MessageListDTO, int64, error) {
	response, count, err := s.messageListAdapter.GetEurMessage(userName, pageNum, countPerPage,
		startTime, isRead)
	if err != nil {
		return []MessageListDTO{}, 0, err
	}
	return response, count, nil
}

func (s *messageListAppService) CountAllMessage(userName string, giteeUsername string) (CountDataDTO, error) {
	data, err := s.messageListAdapter.CountAllMessage(userName, giteeUsername)
	if err != nil {
		return CountDataDTO{}, err
	}
	return data, nil
}
