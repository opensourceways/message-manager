/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package infrastructure

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"

	"github.com/opensourceways/message-manager/common/postgresql"
	"github.com/opensourceways/message-manager/common/user"
	"github.com/opensourceways/message-manager/utils"
)

func MessageListAdapter() *messageAdapter {
	return &messageAdapter{}
}

type messageAdapter struct{}

func (s *messageAdapter) CountAllUnReadMessage(userName string) ([]CountDAO, error) {
	var CountData []CountDAO
	query := `SELECT source, SUM(count) AS count
FROM (
    SELECT cem.source, COUNT(*) AS count
    FROM message_center.cloud_event_message cem
    JOIN message_center.follow_message fm ON cem.event_id = fm.event_id
    JOIN message_center.recipient_config rc ON fm.recipient_id = rc.id
    WHERE fm.is_read = false
      AND fm.is_deleted = false
      AND rc.user_id = ?  -- 替换为实际的用户 ID
    GROUP BY cem.source

    UNION ALL

    SELECT cem.source, COUNT(*) AS count
    FROM message_center.cloud_event_message cem
    JOIN message_center.related_message rm ON cem.event_id = rm.event_id
    JOIN message_center.recipient_config rc ON rm.recipient_id = rc.id
    WHERE rm.is_read = false
      AND rm.is_deleted = false
      AND rc.user_id = ?  -- 替换为实际的用户 ID
    GROUP BY cem.source

    UNION ALL

    SELECT cem.source, COUNT(*) AS count
    FROM message_center.cloud_event_message cem
    JOIN message_center.todo_message tm ON tm.latest_event_id = cem.event_id
    JOIN message_center.recipient_config rc ON tm.recipient_id = rc.id
    WHERE tm.is_read = false
      AND tm.is_deleted = false
      AND rc.user_id = ?  -- 替换为实际的用户 ID
    GROUP BY cem.source
) AS unread_counts
GROUP BY source`
	if result := postgresql.DB().Raw(query, userName, userName, userName).
		Scan(&CountData); result.Error != nil {
		return []CountDAO{}, xerrors.Errorf("get count failed, err:%v", result.Error)
	}
	return CountData, nil
}

func (s *messageAdapter) SetMessageIsRead(userName string, eventId string) error {
	query := `
    	update message_center.follow_message
    	set is_read = true
    	where event_id = ? and is_read = false and recipient_id in (
    	    select id from recipient_config where user_id = ?
    	);
	
    	update message_center.related_message
    	set is_read = true
    	where event_id = ? and is_read = false and recipient_id in (
    	    select id from recipient_config where user_id = ?
    	);
	
    	update message_center.todo_message
    	set is_read = true
    	where latest_event_id = ? and is_read = false and recipient_id in (
    	    select id from recipient_config where user_id = ?
    	);
	`

	if result := postgresql.DB().Exec(query, eventId, userName, eventId, userName, eventId,
		userName); result.Error != nil {
		return xerrors.Errorf("set message is_read failed, err:%v", result.Error.Error())
	}
	return nil
}

func (s *messageAdapter) SetAllMessageIsRead(userName, messageType, giteeUsername,
	startTime string, isRead, isDone, isBot *bool, filter int) error {

	handlers := map[string]func() error{
		"all-todo": func() error {
			return s.makeAllTodoMessageIsRead(userName, giteeUsername, isDone, startTime, isRead)
		},
		"all-about": func() error {
			return s.makeAllAboutMessageIsRead(userName, giteeUsername, isBot, startTime, isRead)
		},
		"all-watch": func() error {
			return s.makeAllWatchMessageIsRead(userName, giteeUsername, startTime, isRead)
		},
		"all-meeting": func() error {
			return s.makeMeetingMessageIsRead(giteeUsername, filter, startTime, isRead)
		},
		"forum-system": func() error {
			return s.makeForumSystemMessageIsRead(userName, startTime, isRead)
		},
		"forum-about": func() error {
			return s.makeForumAboutMessageIsRead(userName, isBot, startTime, isRead)
		},
		"cve-todo": func() error {
			return s.makeCVETodoMessageIsRead(userName, giteeUsername, isDone, startTime, isRead)
		},
		"cve-watch": func() error {
			return s.makeCVEMessageIsRead(userName, giteeUsername, startTime, isRead)
		},
		"issue-todo": func() error {
			return s.makeIssueTodoMessageIsRead(userName, giteeUsername, isDone, startTime, isRead)
		},
		"pr-todo": func() error {
			return s.makePullRequestTodoMessageIsRead(userName, giteeUsername, isDone, startTime, isRead)
		},
		"gitee-about": func() error {
			return s.makeGiteeAboutMessageIsRead(userName, giteeUsername, isBot, startTime, isRead)
		},
		"gitee-watch": func() error {
			return s.makeGiteeMessageIsRead(userName, giteeUsername, startTime, isRead)
		},
		"eur": func() error {
			return s.makeEurMessageIsRead(userName, startTime, isRead)
		},
	}

	handler, exists := handlers[messageType]
	if !exists {
		return fmt.Errorf("unsupported messageType: %s", messageType)
	}
	return handler()
}

func (s *messageAdapter) RemoveMessage(userName string, eventId string) error {
	query := `
    UPDATE message_center.follow_message
    SET is_deleted = true
    WHERE event_id = ? AND is_deleted = false and recipient_id in (
    	    select id from recipient_config where user_id = ?
    	);

    UPDATE message_center.related_message
    SET is_deleted = true
    WHERE event_id = ? AND is_deleted = false and recipient_id in (
    	    select id from recipient_config where user_id = ?
    	);

    UPDATE message_center.todo_message
    SET is_deleted = true
    WHERE latest_event_id = ? AND is_deleted = false and recipient_id in (
    	    select id from recipient_config where user_id = ?
    	);
	`
	if result := postgresql.DB().Exec(query, eventId, userName, eventId, userName,
		eventId, userName); result.Error != nil {
		return xerrors.Errorf("remove inner message failed, err:%v", result.Error.Error())
	}

	return nil
}

func filterTodoSql(query *string, isDone *bool, isRead *bool, startTime string) {
	if isDone != nil {
		*query += fmt.Sprintf(` and is_done=%t`, *isDone)
	}
	if isRead != nil {
		*query += fmt.Sprintf(` and is_read = %t`, *isRead)
	}
	if startTime != "" {
		*query += fmt.Sprintf(` and time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
}

func filterMeetingTodoSql(query *string, isDone *bool, isRead *bool, startTime string) {
	if isDone != nil {
		*query += fmt.Sprintf(` and is_done=%t`, *isDone)
	}
	if isRead != nil {
		*query += fmt.Sprintf(` and is_read = %t`, *isRead)
	}
	if startTime != "" {
		*query += fmt.Sprintf(` and time <= '%s' and time >= NOW()`,
			*utils.ParseUnixTimestampNew(startTime))
	}
}

func filterAboutSql(query *string, isRead *bool, startTime string) {
	if startTime != "" {
		*query += fmt.Sprintf(` and cem.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	if isRead != nil {
		*query += fmt.Sprintf(` and rm.is_read = %t`, *isRead)
	}
}

func filterFollowSql(query *string, isRead *bool, startTime string) {
	if startTime != "" {
		*query += fmt.Sprintf(` and time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	if isRead != nil {
		*query += fmt.Sprintf(` and is_read = %t`, *isRead)
	}
}

func (s *messageAdapter) GetAllToDoMessage(userName string, giteeUsername string, isDone *bool,
	pageNum, countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	query := `with latest_messages as (
    select 
        cem.*,
        tm.is_read,
        tm.is_done
    from
        todo_message tm
    join
        cloud_event_message cem ON cem.event_id = tm.latest_event_id
    join
        recipient_config rc ON rc.id = tm.recipient_id
    where
        tm.is_deleted = false
        and rc.is_deleted = false
        and ((rc.gitee_user_name != '' and rc.gitee_user_name = ?) OR rc.user_id = ?)
        and cem.type <> 'meeting'
	)
	select *, count(*) over () as total_count
	from latest_messages
	where true`
	filterTodoSql(&query, isDone, isRead, startTime)
	query += ` order by updated_at desc limit ? offset ?`

	offset := (pageNum - 1) * countPerPage
	var response []MessageListDAO
	if result := postgresql.DB().Debug().Raw(query, giteeUsername, userName,
		countPerPage, offset).Scan(&response); result.Error != nil {
		return []MessageListDAO{}, 0, xerrors.Errorf("get todo message failed, err:%v",
			result.Error)
	}
	var totalCount int64
	if len(response) != 0 {
		totalCount = response[0].TotalCount
	}
	return response, totalCount, nil
}

func (s *messageAdapter) makeAllTodoMessageIsRead(userName string, giteeUsername string,
	isDone *bool, startTime string, isRead *bool) error {
	query := `with latest_messages as (
    select 
        cem.*,
        tm.is_read,
        tm.is_done,
		tm.recipient_id
    from
        todo_message tm
    join
        cloud_event_message cem ON cem.event_id = tm.latest_event_id
    join
        recipient_config rc ON rc.id = tm.recipient_id
    where
        tm.is_deleted = false
        and rc.is_deleted = false
        and ((rc.gitee_user_name != '' and rc.gitee_user_name = ?) OR rc.user_id = ?)
        and cem.type <> 'meeting'
	)
	select *, count(*) over () as total_count
	from latest_messages
	where true`
	filterTodoSql(&query, isDone, isRead, startTime)

	queryIsRead := fmt.Sprintf(`with tobe_isread as (%s)
        update todo_message set is_read = true where (latest_event_id, 
		recipient_id) in (select event_id, recipient_id from tobe_isread)`, query)
	if result := postgresql.DB().Debug().Raw(queryIsRead, giteeUsername,
		userName); result.Error != nil {
		return xerrors.Errorf("set is message failed, err:%v", result.Error)
	}
	return nil
}

func (s *messageAdapter) GetAllAboutMessage(userName string, giteeUsername string, isBot *bool,
	pageNum, countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	query := `
    select
    	cem.*,
    	rm.is_read,
    	count(*) over () as total_count
	from cloud_event_message cem
	    join message_center.related_message rm on cem.event_id = rm.event_id
	    join message_center.recipient_config rc on rm.recipient_id = rc.id
	where rm.is_deleted = false
	  	and rc.is_deleted = false
	  	and (    
	  	    (cem.type = 'note' and ((rc.gitee_user_name != '' and rc.gitee_user_name = ?) or rc.user_id = ?)`
	if isBot != nil {
		if *isBot {
			query += ` and cem.user IN ('openeuler-ci-bot','ci-robot','openeuler-sync-bot')`
		} else {
			query += ` and cem.user NOT IN ('openeuler-ci-bot','ci-robot','openeuler-sync-bot')`
		}
	}
	query += `)`
	query += `or (cem.source = 'forum' and rc.user_id = ?`
	if isBot != nil {
		if *isBot {
			query += ` and cem.data_json #>> '{Data, OriginalUsername}' = 'system'`
		} else {
			query += ` and cem.data_json #>> '{Data, OriginalUsername}' <> 'system'`
		}
	}
	query += `))`
	filterAboutSql(&query, isRead, startTime)
	query += ` order by updated_at desc limit ? offset ?`
	offset := (pageNum - 1) * countPerPage
	var response []MessageListDAO
	if result := postgresql.DB().Raw(query, giteeUsername, userName, userName, countPerPage,
		offset).Scan(&response); result.Error != nil {
		return []MessageListDAO{}, 0, xerrors.Errorf("get about message failed, err:%v",
			result.Error)
	}
	var totalCount int64
	if len(response) != 0 {
		totalCount = response[0].TotalCount
	}
	return response, totalCount, nil
}

func (s *messageAdapter) makeAllAboutMessageIsRead(userName string, giteeUsername string, isBot *bool,
	startTime string, isRead *bool) error {
	query := `
    select
    	cem.*,
    	rm.is_read,
    	rm.recipient_id
	from cloud_event_message cem
	    join message_center.related_message rm on cem.event_id = rm.event_id
	    join message_center.recipient_config rc on rm.recipient_id = rc.id
	where rm.is_deleted = false
	  	and rc.is_deleted = false
	  	and (    
	  	    (cem.type = 'note' and ((rc.gitee_user_name != '' and rc.gitee_user_name = ?) or rc.user_id = ?)`
	if isBot != nil {
		if *isBot {
			query += ` and cem.user IN ('openeuler-ci-bot','ci-robot','openeuler-sync-bot')`
		} else {
			query += ` and cem.user NOT IN ('openeuler-ci-bot','ci-robot','openeuler-sync-bot')`
		}
	}
	query += `)`
	query += `or (cem.source = 'forum' and rc.user_id = ?`
	if isBot != nil {
		if *isBot {
			query += ` and cem.data_json #>> '{Data, OriginalUsername}' = 'system'`
		} else {
			query += ` and cem.data_json #>> '{Data, OriginalUsername}' <> 'system'`
		}
	}
	query += `))`
	filterAboutSql(&query, isRead, startTime)

	queryIsRead := fmt.Sprintf(`
		with tobe_isread as (%s)
		update related_messsage set is_read = true
		where (event_id, recipient_id) in (select event_id, recipient_id from tobe_isread)`, query)

	if result := postgresql.DB().Debug().Raw(queryIsRead, giteeUsername, userName,
		userName); result.Error != nil {
		return xerrors.Errorf("set is message failed, err:%v", result.Error)
	}
	return nil
}

func (s *messageAdapter) GetAllWatchMessage(userName string, giteeUsername string, pageNum,
	countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	query := `
	with filtered_recipient as (
        select *
        from recipient_config
        where not is_deleted and (user_id = ? or gitee_user_name = ?)
	),
	filtered_messages as (
	    select fm.is_read, cem.*
	    from follow_message fm
	    join cloud_event_message cem on cem.event_id = fm.event_id
	    join filtered_recipient rc on rc.id = fm.recipient_id
	    where not fm.is_deleted
	)
	select *, count(*) over () as total_count
	from filtered_messages 
	where true`
	filterFollowSql(&query, isRead, startTime)
	query += ` order by updated_at desc limit ? offset ?`

	offset := (pageNum - 1) * countPerPage
	var response []MessageListDAO
	if result := postgresql.DB().Debug().Raw(query, userName, giteeUsername, countPerPage,
		offset).Scan(&response); result.Error != nil {
		logrus.Errorf("get watch message failed, err:%v", result.Error)
		return []MessageListDAO{}, 0, xerrors.Errorf("get watch message failed, err:%v", result.Error)
	}
	var totalCount int64
	if len(response) != 0 {
		totalCount = response[0].TotalCount
	}
	return response, totalCount, nil
}

func (s *messageAdapter) makeAllWatchMessageIsRead(userName string, giteeUsername string,
	startTime string, isRead *bool) error {
	query := `
	with filtered_recipient as (
        select *
        from recipient_config
        where not is_deleted and (user_id = ? or gitee_user_name = ?)
	),
	filtered_messages as (
	    select fm.is_read, cem.*, fm.recipient_id
	    from follow_message fm
	    join cloud_event_message cem on cem.event_id = fm.event_id
	    join filtered_recipient rc on rc.id = fm.recipient_id
	    where not fm.is_deleted
	)
	select *
	from filtered_messages 
	where true`
	filterFollowSql(&query, isRead, startTime)

	queryIsRead := fmt.Sprintf(`
		with tobe_isread as (%s)
		update follow_message set is_read = true
		where (event_id, recipient_id) in (select event_id, recipient_id from tobe_isread)`, query)

	if result := postgresql.DB().Debug().Raw(queryIsRead, userName, giteeUsername); result.Error != nil {
		return xerrors.Errorf("set is message failed, err:%v", result.Error)
	}
	return nil
}

func (s *messageAdapter) GetForumSystemMessage(userName string, pageNum,
	countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {

	query := `with filtered_recipient as (
    select *
    from recipient_config
    where not is_deleted and user_id = ?
	),
	filtered_messages as (
	    select fm.is_read, cem.*
	    from follow_message fm
	    join cloud_event_message cem on cem.event_id = fm.event_id
	    join filtered_recipient rc on rc.id = fm.recipient_id
	    where not fm.is_deleted
	)
	select *, count(*) over () as total_count
	from filtered_messages
	where source = 'forum'`
	filterFollowSql(&query, isRead, startTime)

	query += ` order by updated_at desc limit ? offset ?`
	offset := (pageNum - 1) * countPerPage
	var response []MessageListDAO
	if result := postgresql.DB().Raw(query, userName, countPerPage, offset).
		Scan(&response); result.Error != nil {
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	var totalCount int64
	if len(response) != 0 {
		totalCount = response[0].TotalCount
	}
	return response, totalCount, nil
}

func (s *messageAdapter) makeForumSystemMessageIsRead(userName string, startTime string,
	isRead *bool) error {
	query := `with filtered_recipient as (
    select *
    from recipient_config
    where not is_deleted and user_id = ?
	),
	filtered_messages as (
	    select fm.is_read, cem.*, fm.recipient_id
	    from follow_message fm
	    join cloud_event_message cem on cem.event_id = fm.event_id
	    join filtered_recipient rc on rc.id = fm.recipient_id
	    where not fm.is_deleted
	)
	select *
	from filtered_messages
	where source = 'forum'`
	filterFollowSql(&query, isRead, startTime)

	queryIsRead := fmt.Sprintf(`with tobe_isread as (%s)
		update follow_message set is_read = true where (event_id, 
		recipient_id) in (select event_id, recipient_id from tobe_isread)`, query)
	if result := postgresql.DB().Debug().Raw(queryIsRead, userName); result.Error != nil {
		return xerrors.Errorf("set is message failed, err:%v", result.Error)
	}
	return nil
}

func (s *messageAdapter) GetForumAboutMessage(userName string, isBot *bool, pageNum,
	countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	query := `select cem.*, rm.is_read, count(*) over () as total_count
		from related_message rm
		join cloud_event_message cem on cem.event_id = rm.event_id
		join recipient_config rc on rc.id = rm.recipient_id
		where rm.is_deleted = false and rc.is_deleted = false
		and cem.source = 'forum' and rc.user_id = ?`
	if isBot != nil {
		if *isBot {
			query += ` and cem.data_json #>> '{Data, OriginalUsername}' = 'system'`
		} else {
			query += ` and cem.data_json #>> '{Data, OriginalUsername}' <> 'system'`
		}
	}
	filterAboutSql(&query, isRead, startTime)
	query += ` order by time desc limit ? offset ?`

	offset := (pageNum - 1) * countPerPage
	var response []MessageListDAO
	if result := postgresql.DB().Raw(query, userName, countPerPage, offset).
		Scan(&response); result.Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	var totalCount int64
	if len(response) != 0 {
		totalCount = response[0].TotalCount
	}
	return response, totalCount, nil
}

func (s *messageAdapter) makeForumAboutMessageIsRead(userName string, isBot *bool, startTime string,
	isRead *bool) error {
	query := `select cem.*, rm.is_read, rm.recipient_id
		from related_message rm
		join cloud_event_message cem on cem.event_id = rm.event_id
		join recipient_config rc on rc.id = rm.recipient_id
		where rm.is_deleted = false and rc.is_deleted = false
		and cem.source = 'forum' and rc.user_id = ?`
	if isBot != nil {
		if *isBot {
			query += ` and cem.data_json #>> '{Data, OriginalUsername}' = 'system'`
		} else {
			query += ` and cem.data_json #>> '{Data, OriginalUsername}' <> 'system'`
		}
	}
	filterAboutSql(&query, isRead, startTime)

	queryIsRead := fmt.Sprintf(`with tobe_isread as (%s)
		update related_message set is_read = true where (event_id, 
		recipient_id) in (select event_id, recipient_id from tobe_isread`, query)
	if result := postgresql.DB().Debug().Raw(queryIsRead, userName); result.Error != nil {
		return xerrors.Errorf("set is message failed, err:%v", result.Error)
	}
	return nil
}

func (s *messageAdapter) GetMeetingToDoMessage(username string, filter int,
	pageNum, countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	query := `select a.*, count(*) over () as total_count
		from (
		    select distinct on (tm.business_id, tm.recipient_id) tm.is_read, cem.*
		    from todo_message tm
		    join cloud_event_message cem ON cem.event_id = tm.latest_event_id
		    join recipient_config rc ON rc.id = tm.recipient_id
		    where rc.is_deleted = false
		    and tm.is_deleted = false
		    and cem.type = 'meeting'
		    and (rc.gitee_user_name != '' and rc.gitee_user_name = ?)
		    order by tm.business_id, tm.recipient_id, cem.updated_at desc
		) as a where true`

	if filter == 1 {
		query += ` and NOW() <= time`
	} else if filter == 2 {
		query += ` and NOW() > time`
	}
	filterMeetingTodoSql(&query, nil, isRead, startTime)
	query += ` order by time limit ? offset ?`
	giteeUsername, err := user.GetThirdUserName(username)
	if err != nil {
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			xerrors.Errorf("get gitee username failed, err:%v", err))
	}

	offset := (pageNum - 1) * countPerPage
	var response []MessageListDAO
	if result := postgresql.DB().Raw(query, giteeUsername, countPerPage, offset).
		Scan(&response); result.Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	var totalCount int64
	if len(response) != 0 {
		totalCount = response[0].TotalCount
	}
	return response, totalCount, nil
}

func (s *messageAdapter) makeMeetingMessageIsRead(giteeUsername string, filter int,
	startTime string, isRead *bool) error {
	query := `select a.*
		from (
		    select distinct on (tm.business_id, tm.recipient_id) tm.is_read, cem.*, tm.recipient_id
		    from todo_message tm
		    join cloud_event_message cem ON cem.event_id = tm.latest_event_id
		    join recipient_config rc ON rc.id = tm.recipient_id
		    where rc.is_deleted = false
		    and tm.is_deleted = false
		    and cem.type = 'meeting'
		    and (rc.gitee_user_name != '' and rc.gitee_user_name = ?)
		    order by tm.business_id, tm.recipient_id, cem.updated_at desc
		) as a where true`

	if filter == 1 {
		query += ` and NOW() <= time`
	} else if filter == 2 {
		query += ` and NOW() > time`
	}
	filterMeetingTodoSql(&query, nil, isRead, startTime)

	queryIsRead := fmt.Sprintf(`with tobe_isread as (%s)
		update todo_message set is_read = true where (latest_event_id, 
		recipient_id) in (select event_id, recipient_id from tobe_isread)`, query)
	if result := postgresql.DB().Debug().Raw(queryIsRead, giteeUsername); result.Error != nil {
		return xerrors.Errorf("set is message failed, err:%v", result.Error)
	}
	return nil
}

func (s *messageAdapter) GetCVEToDoMessage(userName, giteeUsername string, isDone *bool, pageNum,
	countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	if giteeUsername == "" {
		return []MessageListDAO{}, 0, nil
	}
	query := `select *, count(*) over () as total_count from (
    	select distinct on (tm.business_id, tm.recipient_id) cem.*, 
        	tm.is_read, tm.is_done from todo_message tm
		join cloud_event_message cem on cem.event_id = tm.latest_event_id
		join recipient_config rc on rc.id = tm.recipient_id
		where rc.is_deleted = false and tm.is_deleted = false
		and cem.source = 'cve'
		and ((rc.gitee_user_name != '' and rc.gitee_user_name = ?) or rc.user_id = ?)
		order by tm.business_id, tm.recipient_id, cem.updated_at desc) a where true`
	filterTodoSql(&query, isDone, isRead, startTime)
	query += ` order by updated_at desc limit ? offset ?`

	offset := (pageNum - 1) * countPerPage
	var response []MessageListDAO
	if result := postgresql.DB().Raw(query, giteeUsername, userName, countPerPage, offset).
		Scan(&response); result.Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	var totalCount int64
	if len(response) != 0 {
		totalCount = response[0].TotalCount
	}
	return response, totalCount, nil
}

func (s *messageAdapter) makeCVETodoMessageIsRead(userName, giteeUsername string, isDone *bool,
	startTime string, isRead *bool) error {
	query := `select * from (
    	select distinct on (tm.business_id, tm.recipient_id) cem.*, 
        	tm.is_read, tm.is_done, tm.recipient_id from todo_message tm
		join cloud_event_message cem on cem.event_id = tm.latest_event_id
		join recipient_config rc on rc.id = tm.recipient_id
		where rc.is_deleted = false and tm.is_deleted = false
		and cem.source = 'cve'
		and ((rc.gitee_user_name != '' and rc.gitee_user_name = ?) or rc.user_id = ?)
		order by tm.business_id, tm.recipient_id, cem.updated_at desc) a where true`
	filterTodoSql(&query, isDone, isRead, startTime)

	queryIsRead := fmt.Sprintf(`with tobe_isread as (%s)
		update todo_message set is_read = true where (latest_event_id, 
		recipient_id) in (select event_id, recipient_id from tobe_isread)`, query)
	if result := postgresql.DB().Debug().Raw(queryIsRead, giteeUsername, userName); result.Error != nil {
		return xerrors.Errorf("set is message failed, err:%v", result.Error)
	}
	return nil
}

func (s *messageAdapter) GetCVEMessage(userName, giteeUsername string, pageNum, countPerPage int,
	startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	if giteeUsername == "" {
		return []MessageListDAO{}, 0, nil
	}
	query := `with filtered_recipient as (
    select *
    from recipient_config
    where not is_deleted and ((gitee_user_name != '' and gitee_user_name = ?) or user_id = ?)
	),
	filtered_messages as (
	    select fm.is_read, cem.*
	    from follow_message fm
	    join cloud_event_message cem on cem.event_id = fm.event_id
	    join filtered_recipient rc on rc.id = fm.recipient_id
	    where not fm.is_deleted
	)
	select *, count(*) over () as total_count
	from filtered_messages
	where source = 'cve'`
	filterFollowSql(&query, isRead, startTime)

	query += ` order by updated_at desc limit ? offset ?`
	offset := (pageNum - 1) * countPerPage
	var response []MessageListDAO
	if result := postgresql.DB().Raw(query, giteeUsername, userName, countPerPage, offset).
		Scan(&response); result.Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	var totalCount int64
	if len(response) != 0 {
		totalCount = response[0].TotalCount
	}
	return response, totalCount, nil
}

func (s *messageAdapter) makeCVEMessageIsRead(userName, giteeUsername string, startTime string,
	isRead *bool) error {
	query := `with filtered_recipient as (
    select *
    from recipient_config
    where not is_deleted and ((gitee_user_name != '' and gitee_user_name = ?) or user_id = ?)
	),
	filtered_messages as (
	    select fm.is_read, cem.*, fm.recipient_id
	    from follow_message fm
	    join cloud_event_message cem on cem.event_id = fm.event_id
	    join filtered_recipient rc on rc.id = fm.recipient_id
	    where not fm.is_deleted
	)
	select *
	from filtered_messages
	where source = 'cve'`
	filterFollowSql(&query, isRead, startTime)

	queryIsRead := fmt.Sprintf(`with tobe_isread as (%s)
		update follow_message set is_read = true where (event_id, 
		recipient_id) in (select event_id, recipient_id from tobe_isread)`, query)
	if result := postgresql.DB().Debug().Raw(queryIsRead, giteeUsername, userName); result.Error != nil {
		return xerrors.Errorf("set is message failed, err:%v", result.Error)
	}
	return nil
}

func (s *messageAdapter) GetIssueToDoMessage(userName, giteeUsername string, isDone *bool,
	pageNum, countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	if giteeUsername == "" {
		return []MessageListDAO{}, 0, nil
	}
	query := `select *, count(*) over () as total_count from
        (
		select DISTINCT ON (tm.business_id, tm.recipient_id) cem.*, 
			tm.is_read, tm.is_done from todo_message tm
		join cloud_event_message cem on cem.event_id = latest_event_id
		join recipient_config rc on rc.id = tm.recipient_id
		where tm.is_deleted = false and rc.is_deleted = false
		and cem.type = 'issue' and cem.source = 'https://gitee.com'
		and ((rc.gitee_user_name != '' and rc.gitee_user_name = ?) or rc.user_id = ?)
		order by tm.business_id, tm.recipient_id, cem.updated_at desc) a where true`

	filterTodoSql(&query, isDone, isRead, startTime)
	query += ` order by updated_at desc limit ? offset ?`

	offset := (pageNum - 1) * countPerPage
	var response []MessageListDAO
	if result := postgresql.DB().Raw(query, giteeUsername, userName, countPerPage, offset).
		Scan(&response); result.Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	var totalCount int64
	if len(response) != 0 {
		totalCount = response[0].TotalCount
	}
	return response, totalCount, nil
}

func (s *messageAdapter) makeIssueTodoMessageIsRead(userName, giteeUsername string, isDone *bool,
	startTime string, isRead *bool) error {
	query := `select * from
        (
		select DISTINCT ON (tm.business_id, tm.recipient_id) cem.*, 
			tm.is_read, tm.is_done, tm.recipient_id from todo_message tm
		join cloud_event_message cem on cem.event_id = latest_event_id
		join recipient_config rc on rc.id = tm.recipient_id
		where tm.is_deleted = false and rc.is_deleted = false
		and cem.type = 'issue' and cem.source = 'https://gitee.com'
		and ((rc.gitee_user_name != '' and rc.gitee_user_name = ?) or rc.user_id = ?)
		order by tm.business_id, tm.recipient_id, cem.updated_at desc) a where true`
	filterTodoSql(&query, isDone, isRead, startTime)

	queryIsRead := fmt.Sprintf(`with tobe_isread as (%s)
		update todo_message set is_read = true where (latest_event_id, 
		recipient_id) in (select event_id, recipient_id from tobe_isread)`, query)
	if result := postgresql.DB().Debug().Raw(queryIsRead, giteeUsername, userName); result.Error != nil {
		return xerrors.Errorf("set is message failed, err:%v", result.Error)
	}
	return nil
}

func (s *messageAdapter) GetPullRequestToDoMessage(userName, giteeUsername string, isDone *bool,
	pageNum, countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	if giteeUsername == "" {
		return []MessageListDAO{}, 0, nil
	}
	query := `select *, count(*) over () as total_count from
        (
		select DISTINCT ON (tm.business_id, tm.recipient_id) cem.*, 
			tm.is_read, tm.is_done from todo_message tm
		join cloud_event_message cem on cem.event_id = latest_event_id
		join recipient_config rc on rc.id = tm.recipient_id
		where tm.is_deleted = false and rc.is_deleted = false
		and cem.type = 'pr' and ((rc.gitee_user_name != '' and rc.gitee_user_name = ?) or rc.user_id = ?)
		order by tm.business_id, tm.recipient_id, cem.updated_at desc) a where true`

	filterTodoSql(&query, isDone, isRead, startTime)

	query += ` order by updated_at desc limit ? offset ?`

	offset := (pageNum - 1) * countPerPage
	var response []MessageListDAO
	if result := postgresql.DB().Raw(query, giteeUsername, userName, countPerPage, offset).
		Scan(&response); result.Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	var totalCount int64
	if len(response) != 0 {
		totalCount = response[0].TotalCount
	}
	return response, totalCount, nil
}

func (s *messageAdapter) makePullRequestTodoMessageIsRead(userName, giteeUsername string, isDone *bool,
	startTime string, isRead *bool) error {
	query := `select * from
        (
		select DISTINCT ON (tm.business_id, tm.recipient_id) cem.*, 
			tm.is_read, tm.is_done, tm.recipient_id from todo_message tm
		join cloud_event_message cem on cem.event_id = latest_event_id
		join recipient_config rc on rc.id = tm.recipient_id
		where tm.is_deleted = false and rc.is_deleted = false
		and cem.type = 'pr' and ((rc.gitee_user_name != '' and rc.gitee_user_name = ?) or rc.user_id = ?)
		order by tm.business_id, tm.recipient_id, cem.updated_at desc) a where true`
	filterTodoSql(&query, isDone, isRead, startTime)

	queryIsRead := fmt.Sprintf(`with tobe_isread as (%s)
		update todo_message set is_read = true where (latest_event_id, 
		recipient_id) in (select event_id, recipient_id from tobe_isread)`, query)
	if result := postgresql.DB().Debug().Raw(queryIsRead, giteeUsername, userName); result.Error != nil {
		return xerrors.Errorf("set is message failed, err:%v", result.Error)
	}
	return nil
}

func (s *messageAdapter) GetGiteeAboutMessage(userName, giteeUsername string, isBot *bool,
	pageNum, countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	if giteeUsername == "" {
		return []MessageListDAO{}, 0, nil
	}
	query := `select cem.*, rm.is_read, count(*) over () as total_count
		from cloud_event_message cem
			join message_center.related_message rm on cem.event_id = rm.event_id
			join message_center.recipient_config rc on rm.recipient_id = rc.id
		where cem.type = 'note'
		and cem.source = 'https://gitee.com'
		and rm.is_deleted = false and rc.is_deleted = false
		and ((rc.gitee_user_name != '' and rc.gitee_user_name = ?) or rc.user_id = ?)`
	if isBot != nil {
		if *isBot {
			query += ` and cem."user" IN ('openeuler-ci-bot','ci-robot','openeuler-sync-bot') `
		} else {
			query += ` and cem."user" NOT IN ('openeuler-ci-bot','ci-robot','openeuler-sync-bot') `
		}
	}
	filterAboutSql(&query, isRead, startTime)
	query += ` order by cem.updated_at desc limit ? offset ?`

	offset := (pageNum - 1) * countPerPage
	var response []MessageListDAO
	if result := postgresql.DB().Raw(query, giteeUsername, userName, countPerPage, offset).
		Scan(&response); result.Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	var totalCount int64
	if len(response) != 0 {
		totalCount = response[0].TotalCount
	}
	return response, totalCount, nil
}

func (s *messageAdapter) makeGiteeAboutMessageIsRead(userName, giteeUsername string, isBot *bool,
	startTime string, isRead *bool) error {
	query := `select cem.*, rm.is_read, rm.recipient_id
		from cloud_event_message cem
			join message_center.related_message rm on cem.event_id = rm.event_id
			join message_center.recipient_config rc on rm.recipient_id = rc.id
		where cem.type = 'note'
		and cem.source = 'https://gitee.com'
		and rm.is_deleted = false and rc.is_deleted = false
		and ((rc.gitee_user_name != '' and rc.gitee_user_name = ?) or rc.user_id = ?)`
	if isBot != nil {
		if *isBot {
			query += ` and cem."user" IN ('openeuler-ci-bot','ci-robot','openeuler-sync-bot') `
		} else {
			query += ` and cem."user" NOT IN ('openeuler-ci-bot','ci-robot','openeuler-sync-bot') `
		}
	}
	filterAboutSql(&query, isRead, startTime)
	queryIsRead := fmt.Sprintf(`with tobe_isread as (%s)
		update related_message set is_read = true where (event_id, 
		recipient_id) in (select event_id, recipient_id from tobe_isread)`, query)
	if result := postgresql.DB().Debug().Raw(queryIsRead, giteeUsername, userName); result.Error != nil {
		return xerrors.Errorf("set is message failed, err:%v", result.Error)
	}
	return nil
}

func (s *messageAdapter) GetGiteeMessage(userName, giteeUsername string, pageNum,
	countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	query := `with filtered_recipient as (
    select *
    from recipient_config
    where not is_deleted and ((gitee_user_name != '' and gitee_user_name = ?) or user_id = ?)
	),
	filtered_messages as (
	    select fm.is_read, cem.*
	    from follow_message fm
	    join cloud_event_message cem on cem.event_id = fm.event_id
	    join filtered_recipient rc on rc.id = fm.recipient_id
	    where not fm.is_deleted
	)
	select *, count(*) over () as total_count
	from filtered_messages
	where source = 'https://gitee.com'`
	filterFollowSql(&query, isRead, startTime)
	query += ` order by updated_at desc limit ? offset ?`

	offset := (pageNum - 1) * countPerPage
	var response []MessageListDAO
	if result := postgresql.DB().Raw(query, giteeUsername, userName, countPerPage, offset).
		Scan(&response); result.Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	var totalCount int64
	if len(response) != 0 {
		totalCount = response[0].TotalCount
	}
	return response, totalCount, nil
}

func (s *messageAdapter) makeGiteeMessageIsRead(userName, giteeUsername string, startTime string,
	isRead *bool) error {
	query := `with filtered_recipient as (
    select *
    from recipient_config
    where not is_deleted and ((gitee_user_name != '' and gitee_user_name = ?) or user_id = ?)
	),
	filtered_messages as (
	    select fm.is_read, cem.*, fm.recipient_id
	    from follow_message fm
	    join cloud_event_message cem on cem.event_id = fm.event_id
	    join filtered_recipient rc on rc.id = fm.recipient_id
	    where not fm.is_deleted
	)
	select *
	from filtered_messages
	where source = 'https://gitee.com'`
	filterFollowSql(&query, isRead, startTime)
	queryIsRead := fmt.Sprintf(`with tobe_isread as (%s)
		update follow_message set is_read = true where (event_id, 
		recipient_id) in (select event_id, recipient_id from tobe_isread)`, query)
	if result := postgresql.DB().Debug().Raw(queryIsRead, giteeUsername, userName); result.Error != nil {
		return xerrors.Errorf("set is message failed, err:%v", result.Error)
	}
	return nil
}

func (s *messageAdapter) GetEurMessage(userName string, pageNum,
	countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	query := `with filtered_recipient as (
    select *
    from recipient_config
    where not is_deleted and user_id = ?
	),
	filtered_messages as (
	    select fm.is_read, cem.*
	    from follow_message fm
	    join cloud_event_message cem on cem.event_id = fm.event_id
	    join filtered_recipient rc on rc.id = fm.recipient_id
	    where not fm.is_deleted
	)
	select *, count(*) over () as total_count
	from filtered_messages
	where source = 'https://eur.openeuler.openatom.cn'`
	filterFollowSql(&query, isRead, startTime)
	query += ` order by updated_at desc limit ? offset ?`

	offset := (pageNum - 1) * countPerPage
	var response []MessageListDAO
	if result := postgresql.DB().Raw(query, userName, countPerPage, offset).
		Scan(&response); result.Error != nil {
		return []MessageListDAO{}, 0, xerrors.Errorf("get message failed, err:%v",
			result.Error)
	}
	var totalCount int64
	if len(response) != 0 {
		totalCount = response[0].TotalCount
	}
	return response, totalCount, nil
}

func (s *messageAdapter) makeEurMessageIsRead(userName string, startTime string,
	isRead *bool) error {
	query := `with filtered_recipient as (
    select *
    from recipient_config
    where not is_deleted and user_id = ?
	),
	filtered_messages as (
	    select fm.is_read, cem.*
	    from follow_message fm
	    join cloud_event_message cem on cem.event_id = fm.event_id
	    join filtered_recipient rc on rc.id = fm.recipient_id
	    where not fm.is_deleted
	)
	select *
	from filtered_messages
	where source = 'https://eur.openeuler.openatom.cn'`
	filterFollowSql(&query, isRead, startTime)
	queryIsRead := fmt.Sprintf(`with tobe_isread as (%s)
		update follow_message set is_read = true where (event_id, 
		recipient_id) in (select event_id, recipient_id from tobe_isread)`, query)
	if result := postgresql.DB().Debug().Raw(queryIsRead, userName); result.Error != nil {
		return xerrors.Errorf("set is message failed, err:%v", result.Error)
	}
	return nil
}

func (s *messageAdapter) CountAllMessage(userName string, giteeUserName string) (CountDataDAO, error) {

	response := CountDataDAO{}
	query := `
WITH params AS (SELECT ? AS user_id, ? AS gitee_user_name)
SELECT (SELECT count(*)
        FROM message_center.follow_message fm
                 JOIN recipient_config rc ON fm.recipient_id = rc.id
        WHERE (rc.user_id = params.user_id
            OR (rc.gitee_user_name != '' and rc.gitee_user_name = params.gitee_user_name))
          AND rc.is_deleted IS false
          AND fm.is_deleted IS false
          AND fm.is_read IS false
          AND fm.source in ('forum', 'https://eur.openeuler.openatom.cn', 'cve', 'https://gitee.com'))
		AS watch_count,

       (SELECT count(*)
        FROM message_center.related_message rm
                 JOIN recipient_config rc ON rm.recipient_id = rc.id
        WHERE (rc.user_id = params.user_id
            OR (rc.gitee_user_name != '' and rc.gitee_user_name = params.gitee_user_name))
          AND rc.is_deleted IS false
          AND rm.is_deleted IS false
          AND rm.is_read IS false
          AND rm.source in ('forum', 'https://gitee.com')) AS about_count,

       (SELECT count(*)
        FROM message_center.todo_message tm
                 JOIN recipient_config rc ON tm.recipient_id = rc.id
                 JOIN cloud_event_message cem ON tm.latest_event_id = cem.event_id
        WHERE (rc.gitee_user_name != '' and rc.gitee_user_name = params.gitee_user_name)
          AND rc.is_deleted IS false
          AND tm.is_deleted IS false
          AND tm.is_done IS false
          AND tm.source = 'https://www.openEuler.org/meeting'
          AND cem.time >= current_timestamp) AS meeting_count,

       (SELECT count(*)
        FROM message_center.todo_message tm
                 JOIN recipient_config rc ON tm.recipient_id = rc.id
        WHERE (rc.user_id = params.user_id
            OR (rc.gitee_user_name != '' and rc.gitee_user_name = params.gitee_user_name))
          AND rc.is_deleted IS false
          AND tm.is_deleted IS false
          AND tm.is_done IS false
          AND tm.source in ('forum', 'cve', 'https://gitee.com')) AS todo_count
FROM params;
`
	if result := postgresql.DB().Raw(query, userName, giteeUserName).Scan(&response); result.Error != nil {
		logrus.Errorf("get count failed, err:%v", result.Error.Error())
		return CountDataDAO{}, xerrors.Errorf("查询失败, err:%v", result.Error)
	}
	return response, nil
}

func (s *messageAdapter) GetAllMessage(userName string, pageNum, countPerPage int,
	isRead *bool) ([]MessageListDAO, int64, error) {
	query := `with filtered_recipient as (
            select *
            from recipient_config
            where not is_deleted and user_id = ?
		),
		all_messages as (
		    select fm.is_read, cem.*
		    from follow_message fm
		             join cloud_event_message cem on cem.event_id = fm.event_id
		             join filtered_recipient rc on rc.id = fm.recipient_id
		    where fm.is_deleted = false
		union all
		    select tm.is_read, cem.*
		    from todo_message tm
		             join cloud_event_message cem on cem.event_id = tm.latest_event_id
		             join filtered_recipient rc on rc.id = tm.recipient_id
		    where tm.is_deleted = false
		union all   
		    select rm.is_read, cem.*
		    from related_message rm
		             join cloud_event_message cem on cem.event_id = rm.event_id
		             join filtered_recipient rc on rc.id = rm.recipient_id
		    where rm.is_deleted = false
		)
	select *, count(*) over () as total_count
	from all_messages`
	if isRead != nil {
		query += fmt.Sprintf(" where is_read = %t", *isRead)
	}
	query += ` order by updated_at desc limit ? offset ?`

	offset := (pageNum - 1) * countPerPage
	var response []MessageListDAO
	if result := postgresql.DB().Raw(query, userName, countPerPage, offset).
		Scan(&response); result.Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	var totalCount int64
	if len(response) != 0 {
		totalCount = response[0].TotalCount
	}
	return response, totalCount, nil
}
