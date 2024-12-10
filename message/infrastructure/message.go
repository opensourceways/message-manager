/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package infrastructure

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"gorm.io/gorm"

	"github.com/opensourceways/message-manager/common/postgresql"
	"github.com/opensourceways/message-manager/utils"
)

func MessageListAdapter() *messageAdapter {
	return &messageAdapter{}
}

type messageAdapter struct{}

// 单值过滤
func applySingleValueFilter(query *gorm.DB, column string, value string) *gorm.DB {
	if value != "" {
		query = query.Where(column+" = ?", value)
	}
	return query
}

// 关键字过滤
func applyKeyWordFilter(query *gorm.DB, field string, keyWord string) *gorm.DB {
	if keyWord != "" {
		query = query.Where(field+" ILIKE ?", "%"+keyWord+"%")
	}
	return query
}

// 仓库过滤
func applyRepoFilter(query *gorm.DB, myManagement string, repos string) *gorm.DB {
	var lRepo []string
	if myManagement != "" {
		lRepo, _ = utils.GetUserAdminRepos(myManagement)
	}
	if repos != "" {
		lRepo = append(lRepo, strings.Split(repos, ",")...)
	}
	if len(lRepo) != 0 {
		query = query.Where("cloud_event_message.source_group = ANY(?)", fmt.Sprintf("{%s}",
			strings.Join(lRepo, ",")))
	}
	return query
}

// 时间过滤
func applyTimeFilter(query *gorm.DB, startTime string, endTime string) *gorm.DB {
	start := utils.ParseUnixTimestamp(startTime)
	end := utils.ParseUnixTimestamp(endTime)
	if start != nil && end != nil {
		query = query.Where("cloud_event_message.time BETWEEN ? AND ?", *start, *end)
	} else if start != nil {
		query = query.Where("cloud_event_message.time >= ?", *start)
	} else if end != nil {
		query = query.Where("cloud_event_message.time <= ?", *end)
	}
	return query
}

// 处理机器人过滤条件
func applyBotFilter(query *gorm.DB, isBot string, eventType string) *gorm.DB {
	botNames := []string{"ci-robot", "openeuler-ci-bot", "openeuler-sync-bot"}
	condition := func(event string) string {
		return fmt.Sprintf(`jsonb_extract_path_text(cloud_event_message.data_json,
'%sEvent', 'Sender', 'Name')`, event)
	}

	generateConditions := func(operator string) string {
		var suffix string
		if operator == "=" {
			suffix = " ANY(%s)"
		} else {
			suffix = " ALL(?)"
		}
		conditions := []string{
			condition("Issue") + " " + operator + suffix,
			condition("PullRequest") + " " + operator + suffix,
			condition("Note") + " " + operator + suffix,
		}
		return strings.Join(conditions, " OR ")
	}
	defaultSuffix := fmt.Sprintf("{%s}", strings.Join(botNames, ","))

	var event string
	if eventType == "pr" {
		event = "PullRequest"
	} else if eventType == "issue" {
		event = "Issue"
	} else {
		event = "Note"
	}

	if isBot == "true" {
		if eventType != "" {
			query = query.Where(condition(event)+" = ANY(?)", defaultSuffix)
		} else {
			query = query.Where(generateConditions("="),
				defaultSuffix, defaultSuffix, defaultSuffix)
		}
	} else if isBot == "false" {
		if eventType != "" {
			query = query.Where(condition(event)+" <> ALL(?)", defaultSuffix)
		} else {
			query = query.Where(generateConditions("<>"),
				defaultSuffix, defaultSuffix, defaultSuffix)
		}
	}

	return query
}

// sig组过滤
func applySigGroupFilter(query *gorm.DB, mySig string, giteeSigs string) *gorm.DB {
	var lSig []string
	// 获取我的sig组
	if mySig != "" {
		sigs, err := utils.GetUserSigInfo(mySig)
		if err == nil {
			lSig = append(lSig, sigs...)
		}
	}
	// 添加 Gitee 仓库所属sig
	if giteeSigs != "" {
		lSig = append(lSig, strings.Split(giteeSigs, ",")...)
	}
	// 如果有sig，则添加过滤条件
	if len(lSig) > 0 {
		query = query.Where("jsonb_extract_path_text(cloud_event_message.data_json, "+
			"'SigGroupName') = ANY(?)",
			fmt.Sprintf("{%s}", strings.Join(lSig, ",")))
	}

	return query
}

func applyPrAssigneeFilter(query *gorm.DB, assignee string) *gorm.DB {
	if assignee != "" {
		query = query.Where("jsonb_extract_path_text(cloud_event_message.data_json, "+
			"'Assignees') ILIKE ?", "%"+assignee+"%")
	}
	return query
}

// 复合过滤，处理PullRequest和Issue
func applyCompositeFilters(query *gorm.DB, eventType string, state string, creator string,
	assignee string) *gorm.DB {
	if eventType == "IssueEvent" {
		query = applySingleValueFilter(query, fmt.Sprintf("jsonb_extract_path_text("+
			"cloud_event_message.data_json, '%s', 'Issue', 'State')", eventType), state)
		query = applySingleValueFilter(query, fmt.Sprintf("jsonb_extract_path_text("+
			"cloud_event_message.data_json, '%s', 'Assignee', 'Login')", eventType), assignee)
	} else if eventType == "PullRequestEvent" {
		query = applySingleValueFilter(query, fmt.Sprintf("jsonb_extract_path_text("+
			"cloud_event_message.data_json, '%s', 'State')", eventType), state)
		query = applyPrAssigneeFilter(query, assignee)
	}
	query = applySingleValueFilter(query, fmt.Sprintf("jsonb_extract_path_text("+
		"cloud_event_message.data_json, '%s', 'User', 'Login')", eventType), creator)

	return query
}

// @某人消息过滤
func applyAboutFilter(query *gorm.DB, about string) *gorm.DB {
	if about != "" {
		query = query.Where("jsonb_extract_path_text(cloud_event_message.data_json, 'NoteEvent', "+
			"'Comment', 'Body') LIKE ?", "%"+about+"%")
	}
	return query
}

// Build相关过滤
func applyBuildFilters(query *gorm.DB, buildStatus string, buildOwner string,
	buildCreator string, buildEnv string) *gorm.DB {
	query = applySingleValueFilter(query, "jsonb_extract_path_text(cloud_event_message.data_json,"+
		" 'Body', 'Status')", buildStatus)
	query = applySingleValueFilter(query, "cloud_event_message.user", buildOwner)
	query = applySingleValueFilter(query, "jsonb_extract_path_text(cloud_event_message.data_json,"+
		" 'Body', 'User')", buildCreator)
	query = applySingleValueFilter(query, "jsonb_extract_path_text(cloud_event_message.data_json,"+
		" 'Body', 'Chroot')", buildEnv)
	return query
}

// 会议相关过滤
func applyMeetingFilters(query *gorm.DB, meetingAction string, meetingSigGroup string,
	meetingStartTime string) *gorm.DB {
	query = applySingleValueFilter(query, "jsonb_extract_path_text(cloud_event_message.data_json,"+
		" 'Action')", meetingAction)
	query = applySingleValueFilter(query, "jsonb_extract_path_text(cloud_event_message.data_json,"+
		" 'Msg', 'GroupName')", meetingSigGroup)

	if meetingStartTime != "" {
		start := utils.ParseUnixTimestamp(meetingStartTime)
		if start != nil {
			logrus.Infof("the time is %v, the time is %v, the date is %v", meetingStartTime,
				start.Format(time.DateTime), start.Format(time.DateOnly))
			query = query.Where("jsonb_extract_path_text(cloud_event_message.data_json,"+
				" 'Msg', 'Date') = ?", start.Format(time.DateOnly))
		}
	}
	return query
}

// CVE相关过滤
func applyCVEFilters(query *gorm.DB, cveComponent string, cveState string, cveAffected string) *gorm.DB {
	if cveComponent != "" {
		lComponent := strings.Split(cveComponent, ",")
		var sql []string
		for _, comp := range lComponent {
			sql = append(sql, "%"+comp+"%")
		}
		query = query.Where("jsonb_extract_path_text(cloud_event_message.data_json, "+
			"'CVEComponent') LIKE ANY (?)", fmt.Sprintf("{%s}", strings.Join(sql, ",")))
	}

	if cveState != "" {
		query = query.Where("jsonb_extract_path_text(cloud_event_message.data_json, 'IssueEvent',"+
			" 'Issue', 'State') = ANY (?)", fmt.Sprintf("{%s}", cveState))
	}

	if cveAffected != "" {
		lAffected := strings.Split(cveAffected, ",")
		var sql []string
		for _, affect := range lAffected {
			sql = append(sql, "%"+affect+"%")
		}
		query = query.Where("jsonb_extract_path_text(cloud_event_message.data_json, "+
			"'CVEAffectVersion') ILIKE ANY (?)", fmt.Sprintf("{%s}", strings.Join(sql, ",")))
	}
	return query
}

func GenQuery(query *gorm.DB, params CmdToGetInnerMessage) *gorm.DB {
	// 简单过滤
	query = applySingleValueFilter(query, "inner_message.source", params.Source)
	query = applySingleValueFilter(query, "is_read", params.IsRead)
	query = applySingleValueFilter(query, "cloud_event_message.type", params.EventType)
	query = applySingleValueFilter(query, "jsonb_extract_path_text(cloud_event_message."+
		"data_json, 'NoteEvent', 'NoteableType')", params.NoteType)
	query = applyKeyWordFilter(query, "cloud_event_message.source_group", params.KeyWord)
	query = applyRepoFilter(query, params.MyManagement, params.Repos)
	query = applyTimeFilter(query, params.StartTime, params.EndTime)

	// 复杂过滤
	query = applyBotFilter(query, params.IsBot, params.EventType)
	query = applySigGroupFilter(query, params.MySig, params.GiteeSigs)
	query = applyCompositeFilters(query, "PullRequestEvent", params.PrState, params.PrCreator,
		params.PrAssignee)
	query = applyCompositeFilters(query, "IssueEvent", params.IssueState, params.IssueCreator,
		params.IssueAssignee)
	query = applyAboutFilter(query, params.About)
	query = applyBuildFilters(query, params.BuildStatus, params.BuildOwner, params.BuildCreator,
		params.BuildEnv)
	query = applyMeetingFilters(query, params.MeetingAction, params.MeetingSigGroup,
		params.MeetingStartTime)
	query = applyCVEFilters(query, params.CVEComponent, params.CVEState, params.CVEAffected)
	return query
}

func GenQueryQuick(query *gorm.DB, data MessageSubscribeDAO) *gorm.DB {
	var modeFilterMap map[string]interface{}
	err := json.Unmarshal(data.ModeFilter, &modeFilterMap)
	if err != nil {
		logrus.Errorf("unmarshal modefilter failed, err:%v", err)
		return query
	}
	if data.Source != "" {
		query = query.Where("inner_message.source = ?", data.Source)
	}
	for k, v := range modeFilterMap {
		splitK := strings.Split(k, ".")
		vString, ok := v.(string)
		if !ok {
			logrus.Errorf("it's not ok for type string")
			break
		}
		queryString := generateJSONBExtractPath(splitK)

		if strings.Contains(k, "Sender.Name") {
			if strings.Contains(v.(string), "ne=") {
				vString = strings.ReplaceAll(vString, "ne=", "")
				vString = strings.Join(strings.Split(vString, " "), ",")

				query = query.Where(queryString+" <> ALL(?)", fmt.Sprintf("{%s}", vString))
			} else {
				vString = strings.ReplaceAll(vString, "oneof=", "")
				query = query.Where(queryString+" = ANY(?)", fmt.Sprintf("{%s}", vString))
			}
		} else if strings.Contains(k, "NoteEvent.Comment.Body") {
			vString = strings.ReplaceAll(vString, "contains=", "")
			query = query.Where(queryString+" LIKE ?", "%"+vString+"%")
		} else if strings.Contains(k, "MeetingStartTime") {
			// 使用正则表达式提取时间
			re := regexp.MustCompile(`gt=(.*?),lt=(.*?)$`)
			matches := re.FindStringSubmatch(vString)
			query = query.
				Where("jsonb_extract_path_text(cloud_event_message."+
					"data_json,'MeetingStartTime') BETWEEN ? AND ?", matches[1], matches[2])
		} else if strings.Contains(k, "CVEAffectVersion") {
			vString = strings.ReplaceAll(vString, "contains=", "")
			lString := strings.Split(vString, " ")
			var newLString []string
			for _, s := range lString {
				newLString = append(newLString, "%"+s+"%")
			}
			query = query.Where("jsonb_extract_path_text(cloud_event_message.data_json, "+
				"'CVEAffectVersion') ILIKE ANY(?)",
				fmt.Sprintf("{%s}", strings.Join(newLString, ",")))
		} else if strings.Contains(k, "Assignees") {
			vString = strings.ReplaceAll(vString, "contains=", "")
			query = query.Where("jsonb_extract_path_text(cloud_event_message.data_json, "+
				"'Assignees') ILIKE ?", "%"+vString+"%")
		} else {
			if vString != "" {
				vString = strings.ReplaceAll(vString, "oneof=", "")
				vString = strings.ReplaceAll(vString, "eq=", "")
				query = query.Where(queryString+" = ANY(?)", fmt.Sprintf("{%s}", vString))
			}
		}
	}
	return query
}

func generateJSONBExtractPath(fields []string) string {
	var sb strings.Builder
	sb.WriteString("jsonb_extract_path_text(cloud_event_message.data_json")

	for range fields {
		sb.WriteString(", '%s'")
	}

	sb.WriteString(")")

	formatArgs := make([]interface{}, len(fields))
	for i := range fields {
		formatArgs[i] = fields[i]
	}

	return fmt.Sprintf(sb.String(), formatArgs...)
}

func (s *messageAdapter) GetInnerMessageQuick(cmd CmdToGetInnerMessageQuick,
	userName string) ([]MessageListDAO, int64, error) {
	var data []MessageSubscribeDAO
	if result := postgresql.DB().Table("message_center.subscribe_config").
		Where(gorm.Expr("is_deleted = ?", false)).
		Where("user_name = ? OR user_name IS NULL", userName).
		Where("source = ? AND mode_name = ?", cmd.Source, cmd.ModeName).
		Scan(&data); result.Error != nil {
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v", result.Error)
	}

	query := postgresql.DB().Table("message_center.inner_message").
		Joins("JOIN message_center.cloud_event_message ON "+
			"inner_message.event_id = cloud_event_message.event_id").
		Joins("JOIN message_center.recipient_config ON "+
			"inner_message.recipient_id = recipient_config.id").
		Where("inner_message.is_deleted = ? AND recipient_config.is_deleted = ?", false, false).
		Where("recipient_config.user_id = ?", userName)

	offsetNum := (cmd.PageNum - 1) * cmd.CountPerPage
	GenQueryQuick(query, data[0])
	if len(data) != 0 {
		var lType []string
		for _, dt := range data {
			lType = append(lType, dt.EventType)
		}
		query = query.Where("cloud_event_message.type = ANY(?)", fmt.Sprintf("{%s}",
			strings.Join(lType, ",")))
	}
	var Count int64
	query.Count(&Count)

	var response []MessageListDAO
	if result := query.Limit(cmd.CountPerPage).Offset(offsetNum).
		Order("cloud_event_message.time DESC").
		Scan(&response); result.Error != nil {
		return []MessageListDAO{}, 0, xerrors.Errorf("get message failed, err:%v",
			result.Error)
	}
	return response, Count, nil
}

func (s *messageAdapter) GetInnerMessage(cmd CmdToGetInnerMessage,
	userName string) ([]MessageListDAO, int64, error) {
	query := postgresql.DB().Table("message_center.inner_message").
		Joins("JOIN message_center.cloud_event_message ON "+
			"inner_message.event_id = cloud_event_message.event_id").
		Joins("JOIN message_center.recipient_config ON "+
			"inner_message.recipient_id = recipient_config.id").
		Where("inner_message.is_deleted = ? AND recipient_config.is_deleted = ?", false, false).
		Where("recipient_config.user_id = ?", userName)

	GenQuery(query, cmd)

	var Count int64
	query.Count(&Count)

	var response []MessageListDAO
	offsetNum := (cmd.PageNum - 1) * cmd.CountPerPage
	if result := query.Limit(cmd.CountPerPage).Offset(offsetNum).
		Order("cloud_event_message.time DESC").
		Scan(&response); result.Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}

	return response, Count, nil
}

func (s *messageAdapter) CountAllUnReadMessage(userName string) ([]CountDAO, error) {
	var CountData []CountDAO
	query := `SELECT source, SUM(count) AS total_unread_count
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
    JOIN message_center.inner_message im ON cem.event_id = im.event_id
    JOIN message_center.recipient_config rc ON im.recipient_id = rc.id
    WHERE im.is_read = false
      AND im.is_deleted = false
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
	
    	update message_center.inner_message
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

func (s *messageAdapter) RemoveMessage(userName string, eventId string) error {
	query := `
    UPDATE message_center.follow_message
    SET is_deleted = true
    WHERE event_id = ? AND is_deleted = false and recipient_id in (
    	    select id from recipient_config where user_id = ?
    	);

    UPDATE message_center.inner_message
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

func (s *messageAdapter) GetAllToDoMessage(userName string, giteeUsername string, isDone *bool,
	pageNum, countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO

	issueTodo, issueCount, err := s.GetIssueToDoMessage(userName, giteeUsername, isDone,
		0, 0, startTime, isRead)
	if err != nil {
		return []MessageListDAO{}, 0, err
	}
	prTodo, prCount, err := s.GetPullRequestToDoMessage(userName, giteeUsername, isDone,
		0, 0, startTime, isRead)
	if err != nil {
		return []MessageListDAO{}, 0, err
	}
	cveTodo, cveCount, err := s.GetCVEToDoMessage(userName, giteeUsername, isDone,
		0, 0, startTime, isRead)
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

	query := `select cem.*, fm.is_read from follow_message fm
		join cloud_event_message cem on cem.event_id = fm.event_id
		join recipient_config rc on rc.id = fm.recipient_id
		where fm.is_deleted = false and rc.is_deleted = false
		and cem.source = 'forumm' and rc.user_id = ?`

	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	if isRead != nil && *isRead == false {
		query += ` and fm.is_read = false`
	}

	query += ` order by cem.update_at desc`

	if result := postgresql.DB().Raw(query, userName).Scan(&response); result.Error != nil {
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetForumAboutMessage(userName string, isBot *bool, pageNum,
	countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	query := `select cem.*, im.is_read from follow_message im
		join cloud_event_message cem on cem.event_id = im.event_id
		join recipient_config rc on rc.id = im.recipient_id
		where im.is_deleted = false and rc.is_deleted = false
		and cem.source = 'fourm' and rc.user_id = ?`
	if isBot != nil {
		if *isBot {
			query += ` and cem.data_json #>> '{Data, OriginalUsername}' = 'system'`
		} else {
			query += ` and cem.data_json #>> '{Data, OriginalUsername}' <> 'system'`
		}
	}
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	if isRead != nil {
		query += fmt.Sprintf(` and im.is_read = %t`, *isRead)
	}
	query += ` order by time desc`
	if result := postgresql.DB().Raw(query, userName).Scan(&response); result.Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetMeetingToDoMessage(userName string, filter int,
	pageNum, countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	query := `select distinct on (tm.business_id, tm.recipient_id) cem.*, tm.is_read from todo_message tm
		join cloud_event_message cem on cem.event_id = tm.latest_event_id
		join recipient_config rc on rc.id = tm.recipient_id
		where rc.is_deleted = false and tm.is_deleted = false
		and cem.type = 'meeting'
		and rc.user_id = ?`

	if isRead != nil {
		query += fmt.Sprintf(` and tm.is_read = %t`, *isRead)
	}

	if filter == 1 {
		query += ` and tm.is_done = true`
	} else if filter == 2 {
		query += ` and tm.is_done = false`
	}
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	query += ` order by tm.business_id, tm.recipient_id desc`
	if result := postgresql.DB().Debug().Raw(query, userName).
		Scan(&response); result.Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetCVEToDoMessage(userName, giteeUsername string, isDone *bool, pageNum,
	countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	if giteeUsername == "" {
		return []MessageListDAO{}, 0, nil
	}
	query := `select distinct on (tm.business_id, tm.recipient_id) cem.*, 
tm.is_read from todo_message tm
		join cloud_event_message cem on cem.event_id = tm.latest_event_id
		join recipient_config rc on rc.id = tm.recipient_id
		where rc.is_deleted = false and tm.is_deleted = false
		and cem.source = 'cve'
		and (rc.gitee_user_name = ? or rc.user_id = ?)`
	if isDone != nil {
		if *isDone {
			query += ` where (cem.data_json #>> '{IssueEvent,Issue,State}') IN ('rejected',
'closed')`
		} else {
			query += ` where (cem.data_json #>> '{IssueEvent,Issue,State}') NOT IN ('rejected',
'closed')`
		}
	}
	if isRead != nil {
		query += fmt.Sprintf(` and tm.is_read = %t`, *isRead)
	}
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	query += ` order by updated_at desc`

	if result := postgresql.DB().Raw(query, giteeUsername, userName,
		giteeUsername).Scan(&response); result.Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
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
	query := `select cem.*, fm.is_read from follow_message fm
		join cloud_event_message cem on cem.event_id = fm.event_id
		join recipient_config rc on rc.id = fm.recipient_id
		where rc.is_deleted = false and fm.is_deleted = false
		and cem.source = 'cve'
		and (rc.gitee_user_name = ? or rc.user_id = ?)
		`
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	if isRead != nil {
		query += fmt.Sprintf(` and fm.is_read = %t`, *isRead)
	}

	query += ` order by cem.updated_at desc`
	if result := postgresql.DB().Raw(query, giteeUsername, userName).Scan(&response); result.
		Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetIssueToDoMessage(userName, giteeUsername string, isDone *bool,
	pageNum, countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	if giteeUsername == "" {
		return []MessageListDAO{}, 0, nil
	}
	query := `select DISTINCT ON (tm.business_id, tm.recipient_id) cem.*, 
tm.is_read from todo_message tm
		join cloud_event_message cem on cem.event_id = latest_event_id
		join recipient_config rc on rc.id = tm.recipient_id
		where tm.is_deleted = false and rc.is_deleted = false
		and cem.type = 'issue' and cem.source = 'https://gitee.com'
		and (rc.gitee_user_name = ? or rc.user_id = ?)`

	if isDone != nil {
		if *isDone {
			query += ` where (cem.data_json #>> '{IssueEvent,Issue,State}') IN ('rejected','closed')`
		} else {
			query += ` where (cem.data_json #>> '{IssueEvent,Issue,State}') NOT IN ('rejected',
'closed')`
		}
	}
	if isRead != nil {
		query += fmt.Sprintf(` and tm.is_read = %t`, *isRead)
	}
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	query += ` order by tm.business_id, tm.recipient_id, cem.updated_at desc`
	if result := postgresql.DB().Raw(query, giteeUsername, userName, giteeUsername).
		Scan(&response); result.Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetPullRequestToDoMessage(userName, giteeUsername string, isDone *bool,
	pageNum, countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	if giteeUsername == "" {
		return []MessageListDAO{}, 0, nil
	}
	query := `select DISTINCT ON (tm.business_id, tm.recipient_id) cem.*, 
tm.is_read from todo_message tm
		join cloud_event_message cem on cem.event_id = latest_event_id
		join recipient_config rc on rc.id = tm.recipient_id
		where tm.is_deleted = false and rc.is_deleted = false
		and cem.type = 'pr' and (rc.gitee_user_name = ? or rc.user_id = ?)`
	if isDone != nil {
		if *isDone {
			query += ` where (cem.data_json #>> '{PullRequestEvent,State}') IN ('closed', 'merged')`
		} else {
			query += ` where (cem.data_json #>> '{PullRequestEvent,State}') NOT IN ('closed', 'merged')`
		}
	}
	if isRead != nil {
		query += fmt.Sprintf(` and tm.is_read = %t`, *isRead)
	}
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}

	query += ` order by tm.business_id, tm.recipient_id, cem.updated_at desc`
	if result := postgresql.DB().Raw(query, giteeUsername, userName).
		Scan(&response); result.Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
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
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
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
	query := `select cem.*, fm.is_read
		from follow_message fm 
		join cloud_event_message cem on cem.event_id = fm.event_id
		join recipient_config rc on rc.id = fm.recipient_id
		where cem.source = 'https://gitee.com'
		and fm.is_deleted = false and rc.is_deleted = false
		and (rc.user_id = ? or rc.gitee_user_name = ?)
		and cem."user" NOT IN ('openeuler-ci-bot','ci-robot','openeuler-sync-bot')`
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	if isRead != nil {
		query += fmt.Sprintf(` and fm.is_read = %t`, *isRead)
	}
	query += ` order by cem.updated_at desc`
	if result := postgresql.DB().Raw(query, userName, giteeUsername).Scan(&response); result.
		Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetEurMessage(userName string, pageNum,
	countPerPage int, startTime string, isRead *bool) ([]MessageListDAO, int64, error) {
	query := `select * from follow_message fm
    	join cloud_event_message cem on cem.event_id = fm.event_id
    	join recipient_config rc on rc.id = fm.recipient_id
    	where fm.is_deleted = false and rc.is_deleted = false
    	and cem.source = 'https://eur.openeuler.openatom.cn'
		and rc.user_id = ?`
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= '%s'`, *utils.ParseUnixTimestampNew(startTime))
	}
	if isRead != nil && *isRead == false {
		query += ` and im.is_read = false`
	}
	query += ` order by cem.updated_at desc`

	var response []MessageListDAO
	if result := postgresql.DB().Raw(query, userName).
		Scan(&response); result.Error != nil {
		return []MessageListDAO{}, 0, xerrors.Errorf("get message failed, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) CountAllMessage(userName, giteeUserName string) (CountDataDAO, error) {
	isRead := false
	isDone := false
	_, todoCountNotDone, _ :=
		s.GetAllToDoMessage(userName, giteeUserName, &isDone, 1, 0, "", nil)

	_, aboutCount, _ :=
		s.GetAllAboutMessage(userName, giteeUserName, nil, 1, 0, "", &isRead)

	_, watchCount, _ :=
		s.GetAllWatchMessage(userName, giteeUserName, 1, 0, "", &isRead)

	_, meetingCount, _ :=
		s.GetMeetingToDoMessage(userName, 1, 1, 0, "", nil)
	return CountDataDAO{
		TodoCount:    todoCountNotDone,
		AboutCount:   aboutCount,
		WatchCount:   watchCount,
		MeetingCount: meetingCount,
	}, nil
}

func (s *messageAdapter) GetAllMessage(userName string, pageNum, countPerPage int,
	isRead *bool) ([]MessageListDAO, int64, error) {
	query := `SELECT
    fm.is_read,
    cem.*
FROM
    follow_message fm
JOIN
    cloud_event_message cem ON fm.event_id = cem.event_id
JOIN
    recipient_config rc ON fm.recipient_id = rc.id
WHERE
    rc.user_id = ?`
	if isRead != nil {
		query += fmt.Sprintf(" AND fm.is_read = %t", *isRead)
	}
	query += ` 
UNION ALL

SELECT
    tm.is_read,
    cem.*
FROM
    todo_message tm
JOIN
    cloud_event_message cem ON tm.latest_event_id = cem.event_id
JOIN
    recipient_config rc ON tm.recipient_id = rc.id
WHERE
    rc.user_id = ?
`
	if isRead != nil {
		query += fmt.Sprintf(" AND tm.is_read = %t", *isRead)
	}
	query += ` 
UNION ALL

SELECT
    im.is_read,
    cem.*
FROM
    inner_message im
JOIN
    cloud_event_message cem ON im.event_id = cem.event_id
JOIN
    recipient_config rc ON im.recipient_id = rc.id
WHERE
    rc.user_id = ?
`
	if isRead != nil {
		query += fmt.Sprintf(" AND im.is_read = %t", *isRead)
	}
	query += ` 
ORDER BY 
    cem.update_at DESC;
`
	var response []MessageListDAO
	if result := postgresql.DB().Raw(query, userName, userName, userName).Scan(&response); result.
		Error != nil {
		logrus.Errorf("get message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}
