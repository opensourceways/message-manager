/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package domain

type MessageListAdapter interface {
	GetInnerMessageQuick(cmd CmdToGetInnerMessageQuick, serName string) (
		[]MessageListDO, int64, error)
	GetInnerMessage(cmd CmdToGetInnerMessage, userName string) ([]MessageListDO, int64, error)
	CountAllUnReadMessage(userName string) ([]CountDO, error)
	SetMessageIsRead(source, eventId string) error
	RemoveMessage(source, eventId string) error
	GetForumSystemMessage(userName string) ([]MessageListDO, int64, error)
	GetForumAboutMessage(userName string) ([]MessageListDO, int64, error)
	GetMeetingToDoMessage(userName string, giteeUsername string) ([]MessageListDO, int64, error)
	GetCVEToDoMessage(userName, giteeUsername string) ([]MessageListDO, int64, error)
	GetCVEMessage(userName, giteeUsername string) ([]MessageListDO, int64, error)
	GetIssueToDoMessage(userName, giteeUsername string) ([]MessageListDO, int64, error)
	GetPullRequestToDoMessage(userName, giteeUsername string) ([]MessageListDO, int64, error)
	GetGiteeAboutMessage(userName, giteeUsername string) ([]MessageListDO, int64, error)
	GetGiteeMessage(userName, giteeUsername string) ([]MessageListDO, int64, error)
	GetEurMessage(userName string) ([]MessageListDO, int64, error)
}
