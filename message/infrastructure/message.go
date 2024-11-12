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
		return []MessageListDAO{}, 0, xerrors.Errorf("get inner message failed, err:%v",
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
		logrus.Errorf("get inner message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}

	return response, Count, nil
}

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

func (s *messageAdapter) GetAllAboutMessage(userName string, giteeUsername string, isBot bool,
	pageNum, countPerPage int, startTime string, isRead bool) ([]MessageListDAO, int64, error) {
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
	countPerPage int, startTime string, isRead bool) ([]MessageListDAO, int64, error) {
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
	countPerPage int, startTime string, isRead bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO

	query := `select * from message_center.cloud_event_message cem
		join message_center.inner_message im on im.event_id = cem.event_id
		join message_center.recipient_config rc on rc.id = im.recipient_id
		where im.is_deleted = false and rc.is_deleted = false and cem.source = 'forum'
		  and rc.user_id = ? and cem.type IN ('12','24','37')`

	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= %s`, startTime)
	}
	if isRead {
		query += ` and im.is_read = true`
	} else {
		query += ` and im.is_read = false`
	}

	if result := postgresql.DB().Raw(query, userName).Scan(&response); result.Error != nil {
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetForumAboutMessage(userName string, isBot bool, pageNum,
	countPerPage int, startTime string, isRead bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	query := `select * from message_center.cloud_event_message cem
		join message_center.inner_message im on im.event_id = cem.event_id
		join message_center.recipient_config rc on rc.id = im.recipient_id
		where im.is_deleted = false and rc.is_deleted = false and cem.source = 'forum'
		  and rc.user_id = ? and cem.type NOT IN ('12','24','37') `
	if isBot {
		query += `and jsonb_extract_path_text(cem.data_json::jsonb,
		'Data', 'OriginalUsername') = 'system'`
	} else {
		query += `and jsonb_extract_path_text(cem.data_json::jsonb,
		'Data', 'OriginalUsername') <> 'system'`
	}
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= %s`, startTime)
	}
	if isRead {
		query += ` and im.is_read = true`
	} else {
		query += ` and im.is_read = false`
	}

	if result := postgresql.DB().Debug().Raw(query, userName).Scan(&response); result.Error != nil {
		logrus.Errorf("get inner message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetMeetingToDoMessage(userName string, giteeUsername string, filter int,
	pageNum, countPerPage int) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	query := `select *
		from (select distinct on (cem.source_url) cem.*
		      from cloud_event_message cem
		               join message_center.inner_message im on cem.event_id = im.event_id
		               join message_center.recipient_config rc on im.recipient_id = rc.id
		      where cem.type = 'meeting'
		        and cem.source = 'https://www.openEuler.org/meeting'
		        and (rc.gitee_user_name = ? or rc.user_id = ?)
		        and (cem.data_json #>> '{Action}') <> 'delete_meeting'`

	if filter == 1 {
		query += ` and NOW() <= to_timestamp(data_json ->> 'MeetingEndTime', 'YYYY-MM-DDHH24:MI')`
	} else if filter == 2 {
		query += ` and NOW() > to_timestamp(data_json ->> 'MeetingEndTime', 'YYYY-MM-DDHH24:MI')`
	}
	query += ` order by cem.source_url, updated_at desc) a
		order by updated_at`
	if result := postgresql.DB().Raw(query, giteeUsername, userName).Scan(&response); result.Error != nil {
		logrus.Errorf("get inner message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetCVEToDoMessage(userName, giteeUsername string, isDone bool, pageNum,
	countPerPage int, startTime string) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO

	query := `select *
		from (select distinct on (cem.source_url) cem.*
		      from cloud_event_message cem
		               join message_center.inner_message im on cem.event_id = im.event_id
		               join message_center.recipient_config rc on im.recipient_id = rc.id
		      where cem.type = 'issue'
		        and cem.source = 'cve'
		        and (rc.gitee_user_name = ? or rc.user_id = ?)
		        and (cem.data_json #>> '{IssueEvent,Issue,Assignee,UserName}') = ?`
	if isDone {
		query += ` and (cem.data_json #>> '{IssueEvent,Issue,State}') = 'closed'`
	} else {
		query += ` and (cem.data_json #>> '{IssueEvent,Issue,State}') = 'open'`
	}
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= %s`, startTime)
	}
	query += ` order by cem.source_url, cem.updated_at desc) a
		order by updated_at desc`

	if result := postgresql.DB().Raw(query, giteeUsername, userName,
		giteeUsername).Scan(&response); result.Error != nil {
		logrus.Errorf("get inner message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetCVEMessage(userName, giteeUsername string, pageNum, countPerPage int,
	startTime string, isRead bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	query := `select *
		from cloud_event_message cem
		         join message_center.inner_message im on cem.event_id = im.event_id
		         join message_center.recipient_config rc on im.recipient_id = rc.id
		where cem.type = 'issue'
		  and cem.source = 'cve'
		  and (rc.gitee_user_name = ? or rc.user_id = ?)`
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= %s`, startTime)
	}
	if isRead {
		query += ` and im.is_read = true`
	} else {
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

	query := `select *
		from (select distinct on (cem.source_url) cem.*
		      from cloud_event_message cem
		               join message_center.inner_message im on cem.event_id = im.event_id
		               join message_center.recipient_config rc on im.recipient_id = rc.id
		      where cem.type = 'issue'
		        and cem.source = 'https://gitee.com'
		        and (rc.gitee_user_name = ? or rc.user_id = ?)
		        and (cem.data_json #>> '{IssueEvent,Issue,Assignee,UserName}') = ?`
	if isDone {
		query += ` and (cem.data_json #>> '{IssueEvent,Issue,State}') = 'closed'`
	} else {
		query += ` and (cem.data_json #>> '{IssueEvent,Issue,State}') = 'open'`
	}
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= %s`, startTime)
	}
	query += ` order by cem.source_url, cem.updated_at desc) a
		order by updated_at desc`
	if result := postgresql.DB().Raw(query, giteeUsername, userName,
		giteeUsername).Scan(&response); result.Error != nil {
		logrus.Errorf("get inner message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetPullRequestToDoMessage(userName, giteeUsername string, isDone bool,
	pageNum, countPerPage int, startTime string) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO

	query := `select *
		from (select distinct on (cem.source_url) cem.*
		      from cloud_event_message cem
		               join message_center.inner_message im on cem.event_id = im.event_id
		               join message_center.recipient_config rc on im.recipient_id = rc.id
		          and cem.type = 'pr'
		          and cem.source = 'https://gitee.com'
		          and (rc.gitee_user_name = ? or rc.user_id = ?)
		          and (cem.data_json ->> 'Assignees') :: text like ?`
	if isDone {
		query += ` and (cem.data_json #>> '{PullRequestEvent,State}') IN ('closed', 'merged')`
	} else {
		query += ` and (cem.data_json #>> '{PullRequestEvent,State}') <> 'closed'
		  and (cem.data_json #>> '{PullRequestEvent,State}') <> 'merged'`
	}
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= %s`, startTime)
	}
	query += ` order by cem.source_url, cem.updated_at desc) a
		order by updated_at desc`
	if result := postgresql.DB().Debug().Raw(query, giteeUsername, userName,
		"%"+giteeUsername+"%").Scan(&response); result.Error != nil {
		logrus.Errorf("get inner message failed, err:%v", result.Error.Error())
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v",
			result.Error)
	}
	return pagination(response, pageNum, countPerPage), int64(len(response)), nil
}

func (s *messageAdapter) GetGiteeAboutMessage(userName, giteeUsername string, isBot bool,
	pageNum, countPerPage int, startTime string, isRead bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	query := `select *
		from cloud_event_message cem
		         join message_center.inner_message im on cem.event_id = im.event_id
		         join message_center.recipient_config rc on im.recipient_id = rc.id
		    and cem.type = 'note'
		    and cem.source = 'https://gitee.com'
		    and (rc.gitee_user_name = ? or rc.user_id = ?)
		    and (cem.data_json #>> '{NoteEvent,Issue,User,UserName}' = ?
		        or cem.data_json #>> '{NoteEvent,PullRequest,User,UserName}' = ?
		        or cem.data_json #>> '{NoteEvent,Comment,Body}' like ?)`
	if isBot {
		query += ` and cem."user" IN ('openeuler-ci-bot','ci-robot','openeuler-sync-bot') `
	} else {
		query += ` and cem."user" NOT IN ('openeuler-ci-bot','ci-robot','openeuler-sync-bot') `
	}
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= %s`, startTime)
	}
	if isRead {
		query += ` and im.is_read = true`
	} else {
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
	countPerPage int, startTime string, isRead bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	query := `select *
		from cloud_event_message cem
		         join message_center.inner_message im on cem.event_id = im.event_id
		         join message_center.recipient_config rc on im.recipient_id = rc.id
		    and cem.source = 'https://gitee.com'
		    and (rc.gitee_user_name = ? or rc.user_id = ?)
		and cem."user" NOT IN ('openeuler-ci-bot','ci-robot','openeuler-sync-bot')`
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= %s`, startTime)
	}
	if isRead {
		query += ` and im.is_read = true`
	} else {
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
	countPerPage int, startTime string, isRead bool) ([]MessageListDAO, int64, error) {
	var response []MessageListDAO
	query := `select *
		from cloud_event_message cem
		         join message_center.inner_message im on cem.event_id = im.event_id
		         join message_center.recipient_config rc on im.recipient_id = rc.id
		    and cem.source = 'https://eur.openeuler.openatom.cn'
		    and rc.user_id = ?
			and (cem.data_json #>> '{Body,User}' = ?
		         or cem.data_json #>> '{Body,Owner}' = ?)`
	if startTime != "" {
		query += fmt.Sprintf(` and cem.time >= %s`, startTime)
	}
	if isRead {
		query += ` and im.is_read = true`
	} else {
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
	_, todoCountNotDone, _ := s.GetAllToDoMessage(userName, giteeUserName, false, 1, 0, "")

	_, aboutCountBot, _ := s.GetAllAboutMessage(userName, giteeUserName, true, 1, 0, "", false)
	_, aboutCountNotBot, _ := s.GetAllAboutMessage(userName, giteeUserName, false, 1, 0, "", false)

	_, watchCount, _ := s.GetAllWatchMessage(userName, giteeUserName, 1, 0, "", false)

	_, meetingCount, _ := s.GetMeetingToDoMessage(userName, giteeUserName, 1, 1, 0)
	return CountDataDAO{
		TodoCount:    todoCountNotDone,
		AboutCount:   aboutCountBot + aboutCountNotBot,
		WatchCount:   watchCount,
		MeetingCount: meetingCount,
	}, nil
}
