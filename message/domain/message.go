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

	GetAllToDoMessage(userName, giteeUsername string, isDone bool, pageNum,
		countPerPage int, startTime string) ([]MessageListDO, int64, error)
	GetAllAboutMessage(userName, giteeUsername string, isBot bool, pageNum,
		countPerPage int, startTime string, isRead *bool) ([]MessageListDO, int64, error)
	GetAllWatchMessage(userName, giteeUsername string, pageNum, countPerPage int,
		startTime string, isRead *bool) ([]MessageListDO, int64, error)

	GetForumSystemMessage(userName string, pageNum, countPerPage int,
		startTime string, isRead *bool) ([]MessageListDO, int64, error)
	GetForumAboutMessage(userName string, isBot bool, pageNum,
		countPerPage int, startTime string, isRead *bool) ([]MessageListDO, int64, error)
	GetMeetingToDoMessage(userName string, giteeUsername string, filter int, pageNum,
		countPerPage int) ([]MessageListDO, int64, error)
	GetCVEToDoMessage(userName, giteeUsername string, isDone bool, pageNum,
		countPerPage int, startTime string) ([]MessageListDO, int64, error)
	GetCVEMessage(userName, giteeUsername string, pageNum, countPerPage int,
		startTime string, isRead *bool) ([]MessageListDO, int64, error)
	GetIssueToDoMessage(userName, giteeUsername string, isDone bool, pageNum,
		countPerPage int, startTime string) ([]MessageListDO, int64, error)
	GetPullRequestToDoMessage(userName, giteeUsername string, isDone bool, pageNum,
		countPerPage int, startTime string) ([]MessageListDO, int64, error)
	GetGiteeAboutMessage(userName, giteeUsername string, isBot bool,
		pageNum, countPerPage int, startTime string, isRead *bool) ([]MessageListDO, int64, error)
	GetGiteeMessage(userName, giteeUsername string, pageNum, countPerPage int,
		startTime string, isRead *bool) ([]MessageListDO, int64, error)
	GetEurMessage(userName string, pageNum, countPerPage int, startTime string,
		isRead *bool) ([]MessageListDO, int64, error)
	CountAllMessage(username, giteeUsername string) (CountDataDO, error)
}
