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

// 处理布尔值的过滤条件
func applyBotFilter(query *gorm.DB, isBot string, eventType string) *gorm.DB {
	if isBot == "true" {
		botNames := "{ci-robot, openeuler-ci-bot, openeuler-sync-bot}"
		condition := fmt.Sprintf(`jsonb_extract_path_text(cloud_event_message.data_json, 
'%sEvent', 'Sender', 'Name')`, eventType)
		query = query.Where(condition+" = ANY(?)", botNames)
	} else if isBot == "false" {
		botNames := "{ci-robot, openeuler-ci-bot, openeuler-sync-bot}"
		condition := fmt.Sprintf(`jsonb_extract_path_text(cloud_event_message.data_json, 
'%sEvent', 'Sender', 'Name')`, eventType)
		query = query.Where(condition+" <> ALL(?)", botNames)
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

// 复合过滤，处理PullRequest和Issue
func applyCompositeFilters(query *gorm.DB, eventType string, state string, creator string,
	assignee string) *gorm.DB {
	query = applySingleValueFilter(query, fmt.Sprintf("jsonb_extract_path_text("+
		"cloud_event_message.data_json, '%s', 'State')", eventType), state)
	query = applySingleValueFilter(query, fmt.Sprintf("jsonb_extract_path_text("+
		"cloud_event_message.data_json, '%s', 'User', 'Login')", eventType), creator)
	query = applySingleValueFilter(query, fmt.Sprintf("jsonb_extract_path_text("+
		"cloud_event_message.data_json, '%s', 'Assignee', 'Login')", eventType), assignee)
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
	meetingStartTime string, meetingEndTime string) *gorm.DB {
	query = applySingleValueFilter(query, "jsonb_extract_path_text(cloud_event_message.data_json,"+
		" 'Action')", meetingAction)
	query = applySingleValueFilter(query, "jsonb_extract_path_text(cloud_event_message.data_json,"+
		" 'Msg', 'GroupName')", meetingSigGroup)

	if meetingStartTime != "" && meetingEndTime != "" {
		start := utils.ParseUnixTimestamp(meetingStartTime)
		end := utils.ParseUnixTimestamp(meetingEndTime)
		if start != nil && end != nil {
			query = query.Where("jsonb_extract_path_text(cloud_event_message.data_json, "+
				"'MeetingStartTime') BETWEEN ? AND ?", start.Format(time.DateTime), end.Format(time.DateTime))
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
	query = applySingleValueFilter(query, "jsonb_extract_path_text(cloud_event_message.data_json,"+
		" 'IssueEvent', 'Issue', 'StateName')", cveState)

	if cveAffected != "" {
		lAffected := strings.Split(cveAffected, ",")
		var sql []string
		for _, affect := range lAffected {
			sql = append(sql, "%"+affect+"%")
		}
		query = query.Where("jsonb_extract_path_text(cloud_event_message.data_json, "+
			"'CVEAffectVersion') LIKE ANY (?)", fmt.Sprintf("{%s}", strings.Join(sql, ",")))
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
		params.MeetingStartTime, params.MeetingEndTime)
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
	query = query.Where("inner_message.source = ?", data.Source)
	query = query.Where("cloud_event_message.type = ?", data.EventType)

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
		} else {
			vString = strings.ReplaceAll(vString, "oneof=", "")
			vString = strings.ReplaceAll(vString, "eq=", "")
			query = query.Where(queryString+" = ANY(?)", fmt.Sprintf("{%s}", vString))
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
	var data MessageSubscribeDAO
	if result := postgresql.DB().Table("message_center.subscribe_config").
		Where(gorm.Expr("is_deleted = ?", false)).
		Where("user_name = ? OR user_name IS NULL", userName).
		Where("source = ? AND mode_name = ?", cmd.Source, cmd.ModeName).
		Scan(&data); result.Error != nil {
		return []MessageListDAO{}, 0, xerrors.Errorf("查询失败, err:%v", result.Error)
	}

	query := postgresql.DB().Table("message_center.inner_message").
		Joins("JOIN message_center.cloud_event_message ON "+
			"inner_message.event_id = cloud_event_message.event_id AND"+
			" inner_message.source = cloud_event_message.source").
		Joins("JOIN message_center.recipient_config ON "+
			"cast(inner_message.recipient_id AS BIGINT) = recipient_config.id").
		Where("inner_message.is_deleted = ? AND recipient_config.is_deleted = ?", false, false).
		Where("recipient_config.user_id = ?", userName)

	offsetNum := (cmd.PageNum - 1) * cmd.CountPerPage

	GenQueryQuick(query, data)

	var Count int64
	query.Count(&Count)

	var response []MessageListDAO
	if result := query.Debug().Limit(cmd.CountPerPage).Offset(offsetNum).
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
			"inner_message.event_id = cloud_event_message.event_id AND"+
			" inner_message.source = cloud_event_message.source").
		Joins("JOIN message_center.recipient_config ON "+
			"cast(inner_message.recipient_id AS BIGINT) = recipient_config.id").
		Where("inner_message.is_deleted = ? AND recipient_config.is_deleted = ?", false, false).
		Where("recipient_config.user_id = ?", userName)

	GenQuery(query, cmd)

	var Count int64
	query.Count(&Count)

	var response []MessageListDAO
	offsetNum := (cmd.PageNum - 1) * cmd.CountPerPage
	if result := query.Debug().Limit(cmd.CountPerPage).Offset(offsetNum).
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
