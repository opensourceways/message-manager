/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package infrastructure

import (
	"fmt"

	"github.com/opensourceways/message-manager/common/postgresql"
	"github.com/opensourceways/message-manager/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

func MessageListAdapter() *messageAdapter {
	return &messageAdapter{}
}

type messageAdapter struct{}

func (s *messageAdapter) CountAllUnReadMessage(userName string) ([]CountDAO, error) {
	var CountData []CountDAO
	sqlCount := `SELECT inner_message.source, COUNT(*) FROM message_center.inner_message 
    JOIN message_center.cloud_event_message ON inner_message.event_id = cloud_event_message.event_id 
         AND inner_message.source = cloud_event_message.source 
	JOIN message_center.recipient_config ON 
		cast(inner_message.recipient_id AS BIGINT) = recipient_config.id 
	WHERE is_read = ? AND recipient_config.user_id = ? 
	AND inner_message.is_deleted = ? 
	AND recipient_config.is_deleted = ?
	GROUP BY inner_message.source`
	if result := postgresql.DB().Raw(sqlCount, false, userName, false, false).
		Scan(&CountData); result.Error != nil {
		return []CountDAO{}, xerrors.Errorf("get count failed, err:%v", result.Error)
	}
	return CountData, nil
}

func (s *messageAdapter) SetMessageIsRead(source, eventId string) error {
	if result := postgresql.DB().Table("message_center.inner_message").
		Where("inner_message.source = ? AND inner_message.event_id = ?", source,
			eventId).Where("inner_message.is_deleted = ?", false).
		Update("is_read", true); result.Error != nil {
		return xerrors.Errorf("set message is_read failed, err:%v", result.Error.Error())
	}
	return nil
}

func (s *messageAdapter) RemoveMessage(source, eventId string) error {
	if result := postgresql.DB().Table("message_center.inner_message").
		Where("inner_message.source = ? AND inner_message."+
			"event_id = ?", source, eventId).
		Update("is_deleted", true); result.Error != nil {
		return xerrors.Errorf("remove inner message failed, err:%v", result.Error.Error())
	}
	return nil
}

func pagination(messages []MessageListDAO, pageNum, countPerPage int) []MessageListDAO {
	if countPerPage == 0 {
		return messages
	}
	start := (pageNum - 1) * countPerPage
	end := start + countPerPage
	if start > len(messages) {
		return []MessageListDAO{}
	}
	if end > len(messages) {
		return messages[start:]
	}
	return messages[start:end]
}

func (s *messageAdapter) GetAllToDoMessage(userName string, giteeUsername string, isDone bool,
	pageNum, countPerPage int, startTime string) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO

	issueTodo, issueCount, err := s.GetIssueToDoMessage(userName, giteeUsername, isDone,
		0, 0, startTime)
	if err != nil {
		return []MessageListDAO{}, 0, err
	}
	prTodo, prCount, err := s.GetPullRequestToDoMessage(userName, giteeUsername, isDone,
		0, 0, startTime)
	if err != nil {
		return []MessageListDAO{}, 0, err
	}
	cveTodo, cveCount, err := s.GetCVEToDoMessage(userName, giteeUsername, isDone,
		0, 0, startTime)
	response = append(response, issueTodo...)
	response = append(response, prTodo...)
	response = append(response, cveTodo...)
	return pagination(response, pageNum, countPerPage), issueCount + prCount + cveCount, nil
}

func (s *messageAdapter) GetAllAboutMessage(userName string, giteeUsername string, isBot *bool,
	pageNum, countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	giteeAbout, giteeCount, err := s.GetGiteeAboutMessage(userName, giteeUsername, isBot,
		0, 0, startTime, isRead)
	if err != nil {
		return []MessageListDAO{}, 0, err
	}
	forumAbout, forumCount, err := s.GetForumAboutMessage(userName, isBot,
		0, 0, startTime, isRead)
	if err != nil {
		return []MessageListDAO{}, 0, err
	}
	response = append(response, giteeAbout...)
	response = append(response, forumAbout...)

	return pagination(response, pageNum, countPerPage), giteeCount + forumCount, nil
}

func (s *messageAdapter) GetAllWatchMessage(userName string, giteeUsername string, pageNum,
	countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	forumMsg, forumCount, err := s.GetForumSystemMessage(userName,
		0, 0, startTime, isRead)
	if err != nil {
		return []MessageListDAO{}, 0, err
	}
	cveMsg, cveCount, err := s.GetCVEMessage(userName, giteeUsername,
		0, 0, startTime, isRead)
	if err != nil {
		return []MessageListDAO{}, 0, err
	}
	giteeMsg, giteeCount, err := s.GetGiteeMessage(userName, giteeUsername,
		0, 0, startTime, isRead)
	if err != nil {
		return []MessageListDAO{}, 0, err
	}
	eurMsg, eurCount, err := s.GetEurMessage(userName, 0, 0, startTime, isRead)
	if err != nil {
		return []MessageListDAO{}, 0, err
	}
	response = append(response, forumMsg...)
	response = append(response, cveMsg...)
	response = append(response, giteeMsg...)
	response = append(response, eurMsg...)
	return pagination(response, pageNum, countPerPage), forumCount + cveCount + giteeCount + eurCount, nil
}

func (s *messageAdapter) GetForumSystemMessage(userName string, pageNum,
	countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO

	query := `select cem.*, im.is_read from message_center.cloud_event_message cem
		join message_center.inner_message im on im.event_id = cem.event_id
		join message_center.recipient_config rc on rc.id = im.recipient_id
		where im.is_deleted = false and rc.is_deleted = false and cem.source = 'forum'
		  and rc.user_id = ? and cem.type IN ('12','24','37')`

	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	if isRead != nil && *isRead == false {
		query += ` and im.is_read = false`
	}

	query += ` order by time desc`

	if result := postgresql.DB().Raw(query, userName).Scan(&response); result.Error != nil {
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetForumAboutMessage(userName string, isBot *bool, pageNum,
	countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	query := `select cem.*, im.is_read from message_center.cloud_event_message cem
		join message_center.inner_message im on im.event_id = cem.event_id
		join message_center.recipient_config rc on rc.id = im.recipient_id
		where im.is_deleted = false and rc.is_deleted = false and cem.source = 'forum'
		  and rc.user_id = ? and cem.type NOT IN ('12','24','37') `
	if isBot != nil {
		if *isBot {
			query += `and jsonb_extract_path_text(cem.data_json::jsonb,
		'Data', 'OriginalUsername') = 'system'`
		} else {
			query += `and jsonb_extract_path_text(cem.data_json::jsonb,
		'Data', 'OriginalUsername') <> 'system'`
		}
	}
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	if isRead != nil && *isRead == false {
		query += ` and im.is_read = false`
	}
	query += ` order by time desc`
	if result := postgresql.DB().Raw(query, userName).Scan(&response); result.Error != nil {
		logrus.Errorf("get inner message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetMeetingToDoMessage(userName string, filter int,
	pageNum, countPerPage int) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	query := `select *
		from (select distinct on (cem.source_url) cem.*, im.is_read
		      from cloud_event_message cem
		               join message_center.inner_message im on cem.event_id = im.event_id
		               join message_center.recipient_config rc on im.recipient_id = rc.id
		      where cem.type = 'meeting'
		        and im.is_deleted = false and rc.is_deleted = false
		        and cem.source = 'https://www.openEuler.org/meeting'
		        and (rc.user_id = ?)
		        and (cem.data_json #>> '{Action}') <> 'delete_meeting'
		        order by cem.source_url, updated_at desc) a`

	if filter == 1 {
		query += ` where NOW() <= to_timestamp(a.data_json ->> 'MeetingEndTime', 
'YYYY-MM-DDHH24:MI')`
	} else if filter == 2 {
		query += ` where NOW() > to_timestamp(a.data_json ->> 'MeetingEndTime', 
'YYYY-MM-DDHH24:MI')`
	}
	query += ` order by updated_at`
	if result := postgresql.DB().Debug().Raw(query, userName).Scan(&response); result.Error != nil {
		logrus.Errorf("get inner message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetCVEToDoMessage(userName, giteeUsername string, isDone bool, pageNum,
	countPerPage int, startTime string) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	if giteeUsername == "" {
		return []MessageListDAO{}, 0, nil
	}
	query := `select *
		from (select distinct on (cem.source_url) cem.*, im.is_read
		      from cloud_event_message cem
		               join message_center.inner_message im on cem.event_id = im.event_id
		               join message_center.recipient_config rc on im.recipient_id = rc.id
		      where cem.type = 'issue'
		        and cem.source = 'cve'
		        and im.is_deleted = false and rc.is_deleted = false
		        and (rc.gitee_user_name = ? or rc.user_id = ?)
		        and (cem.data_json #>> '{IssueEvent,Issue,Assignee,UserName}') = ?
		      order by cem.source_url, cem.updated_at desc) a`
	if isDone {
		query += ` where (a.data_json #>> '{IssueEvent,Issue,State}') IN ('rejected','closed')`
	} else {
		query += ` where (a.data_json #>> '{IssueEvent,Issue,State}') NOT IN ('rejected','closed')`
	}
	if startTime != "" {
		query += fmt.Sprintf(` and a.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	query += ` order by updated_at desc`

	if result := postgresql.DB().Raw(query, giteeUsername, userName,
		giteeUsername).Scan(&response); result.Error != nil {
		logrus.Errorf("get inner message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetCVEMessage(userName, giteeUsername string, pageNum, countPerPage int,
	startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	if giteeUsername == "" {
		return []MessageListDAO{}, 0, nil
	}
	query := `select cem.*, im.is_read
		from cloud_event_message cem
		         join message_center.inner_message im on cem.event_id = im.event_id
		         join message_center.recipient_config rc on im.recipient_id = rc.id
		where cem.type = 'issue'
		  and cem.source = 'cve'
		  and im.is_deleted = false and rc.is_deleted = false
		  and (rc.gitee_user_name = ? or rc.user_id = ?)`
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	if isRead != nil && *isRead == false {
		query += ` and im.is_read = false`
	}

	query += ` order by cem.updated_at desc`
	if result := postgresql.DB().Raw(query, giteeUsername, userName).Scan(&response); result.
		Error != nil {
		logrus.Errorf("get inner message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetIssueToDoMessage(userName, giteeUsername string, isDone bool,
	pageNum, countPerPage int, startTime string) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	if giteeUsername == "" {
		return []MessageListDAO{}, 0, nil
	}
	query := `select *
		from (select distinct on (cem.source_url) cem.*, im.is_read
		      from cloud_event_message cem
		               join message_center.inner_message im on cem.event_id = im.event_id
		               join message_center.recipient_config rc on im.recipient_id = rc.id
		      where cem.type = 'issue'
		        and cem.source = 'https://gitee.com'
		        and im.is_deleted = false and rc.is_deleted = false
		        and (rc.gitee_user_name = ? or rc.user_id = ?)
		        and (cem.data_json #>> '{IssueEvent,Issue,Assignee,UserName}') = ?
 			  order by cem.source_url, cem.updated_at desc) a`
	if isDone {
		query += ` where (a.data_json #>> '{IssueEvent,Issue,State}') IN ('rejected','closed')`
	} else {
		query += ` where (a.data_json #>> '{IssueEvent,Issue,State}') NOT IN ('rejected',
'closed')`
	}
	if startTime != "" {
		query += fmt.Sprintf(` and a.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	query += ` order by updated_at desc`
	if result := postgresql.DB().Raw(query, giteeUsername, userName, giteeUsername).
		Scan(&response); result.Error != nil {
		logrus.Errorf("get inner message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetPullRequestToDoMessage(userName, giteeUsername string, isDone bool,
	pageNum, countPerPage int, startTime string) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	if giteeUsername == "" {
		return []MessageListDAO{}, 0, nil
	}
	query := `select *
		from (select distinct on (cem.source_url) cem.*, im.is_read
		      from cloud_event_message cem
		               join message_center.inner_message im on cem.event_id = im.event_id
		               join message_center.recipient_config rc on im.recipient_id = rc.id
		          and cem.type = 'pr'
		          and cem.source = 'https://gitee.com'
		          and im.is_deleted = false and rc.is_deleted = false
		          and (rc.gitee_user_name = ? or rc.user_id = ?)
		          and (cem.data_json ->> 'Assignees') :: text like ?
		      order by cem.source_url, cem.updated_at desc) a`
	if isDone {
		query += ` where (a.data_json #>> '{PullRequestEvent,State}') IN ('closed', 'merged')`
	} else {
		query += ` where (a.data_json #>> '{PullRequestEvent,State}') NOT IN ('closed', 'merged')`
	}
	if startTime != "" {
		query += fmt.Sprintf(` and a.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}

	query += ` order by updated_at desc`
	if result := postgresql.DB().Raw(query, giteeUsername, userName,
		"%"+giteeUsername+"%").Scan(&response); result.Error != nil {
		logrus.Errorf("get inner message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetGiteeAboutMessage(userName, giteeUsername string, isBot *bool,
	pageNum, countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	if giteeUsername == "" {
		return []MessageListDAO{}, 0, nil
	}
	query := `select cem.*, im.is_read
		from cloud_event_message cem
		         join message_center.inner_message im on cem.event_id = im.event_id
		         join message_center.recipient_config rc on im.recipient_id = rc.id
		    and cem.type = 'note'
		    and cem.source = 'https://gitee.com'
		    and im.is_deleted = false and rc.is_deleted = false
		    and (rc.gitee_user_name = ? or rc.user_id = ?)
		    and (cem.data_json #>> '{NoteEvent,Issue,User,UserName}' = ?
		        or cem.data_json #>> '{NoteEvent,PullRequest,User,UserName}' = ?
		        or cem.data_json #>> '{NoteEvent,Comment,Body}' like ?)`
	if isBot != nil {
		if *isBot {
			query += ` and cem."user" IN ('openeuler-ci-bot','ci-robot','openeuler-sync-bot') `
		} else {
			query += ` and cem."user" NOT IN ('openeuler-ci-bot','ci-robot','openeuler-sync-bot') `
		}
	}
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	if isRead != nil && *isRead == false {
		query += ` and im.is_read = false`
	}
	query += ` order by cem.updated_at desc`
	if result := postgresql.DB().Raw(query, giteeUsername, userName,
		giteeUsername, giteeUsername, "%"+giteeUsername+"%").Scan(&response); result.Error != nil {
		logrus.Errorf("get inner message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetGiteeMessage(userName, giteeUsername string, pageNum,
	countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	if giteeUsername == "" {
		return []MessageListDAO{}, 0, nil
	}
	query := `select cem.*, im.is_read
		from cloud_event_message cem
		         join message_center.inner_message im on cem.event_id = im.event_id
		         join message_center.recipient_config rc on im.recipient_id = rc.id
		    and cem.source = 'https://gitee.com'
			and im.is_deleted = false and rc.is_deleted = false
		    and (rc.gitee_user_name = ? or rc.user_id = ?)
		and cem."user" NOT IN ('openeuler-ci-bot','ci-robot','openeuler-sync-bot')`
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	if isRead != nil && *isRead == false {
		query += ` and im.is_read = false`
	}
	query += ` order by cem.updated_at desc`
	if result := postgresql.DB().Raw(query, giteeUsername, userName).Scan(&response); result.
		Error != nil {
		logrus.Errorf("get inner message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetEurMessage(userName string, pageNum,
	countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	query := `select cem.*, im.is_read
		from cloud_event_message cem
		         join message_center.inner_message im on cem.event_id = im.event_id
		         join message_center.recipient_config rc on im.recipient_id = rc.id
		    and cem.source = 'https://eur.openeuler.openatom.cn'
		    and im.is_deleted = false and rc.is_deleted = false
		    and rc.user_id = ?
			and (cem.data_json #>> '{Body,User}' = ?
		         or cem.data_json #>> '{Body,Owner}' = ?)`
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	if isRead != nil && *isRead == false {
		query += ` and im.is_read = false`
	}
	query += ` order by cem.updated_at desc`
	if result := postgresql.DB().Raw(query, userName, userName, userName).Scan(&response); result.
		Error != nil {
		logrus.Errorf("get inner message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) CountAllMessage(userName, giteeUserName string) (CountDataDAO, error) {
	isRead := false
	_, todoCountNotDone, _ := s.GetAllToDoMessage(userName, giteeUserName, false, 1, 0, "")

	_, aboutCount, _ := s.GetAllAboutMessage(userName, giteeUserName, nil, 1, 0, "", &isRead)

	_, watchCount, _ := s.GetAllWatchMessage(userName, giteeUserName, 1, 0, "", &isRead)

	_, meetingCount, _ := s.GetMeetingToDoMessage(userName, 1, 1, 0)
	return CountDataDAO{
		TodoCount:    todoCountNotDone,
		AboutCount:   aboutCount,
		WatchCount:   watchCount,
		MeetingCount: meetingCount,
	}, nil
}
