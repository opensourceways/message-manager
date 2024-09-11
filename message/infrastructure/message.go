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

func GenQuery(query *gorm.DB, params CmdToGetInnerMessage) *gorm.DB {
	// 消息源
	if params.Source != "" {
		query = query.Where("inner_message.source = ?", params.Source)
	}

	// 消息是否已读
	if params.IsRead != "" {
		query = query.Where("is_read = ?", params.IsRead)
	}

	//是否为机器人消息
	if params.IsBot != "" {
		lIsBot := strings.Split(params.IsBot, ",")
		if len(lIsBot) == 1 {
			if lIsBot[0] == "true" {
				switch params.EventType {
				case "pr":
					query = query.
						Where(`jsonb_extract_path_text(cloud_event_message.data_json, 
'PullRequestEvent','Sender','Name') =  ANY(?)`,
							"{ci-robot, openeuler-ci-bot, openeuler-sync-bot}")
				case "issue":
					query = query.
						Where(`jsonb_extract_path_text(cloud_event_message.data_json, 
'IssueEvent','Sender','Name') = ANY(?)`,
							"{ci-robot, openeuler-ci-bot, openeuler-sync-bot}")
				case "note":
					query = query.
						Where(`jsonb_extract_path_text(cloud_event_message.data_json, 
'NoteEvent','Sender','Name') = ANY(?)`,
							"{ci-robot, openeuler-ci-bot, openeuler-sync-bot}")
				}
			} else {
				switch params.EventType {
				case "pr":
					query = query.
						Where(`jsonb_extract_path_text(cloud_event_message.data_json, 
'PullRequestEvent','Sender','Name') <> ALL(?)`,
							"{ci-robot, openeuler-ci-bot, openeuler-sync-bot}")
				case "issue":
					query = query.
						Where(`jsonb_extract_path_text(cloud_event_message.data_json, 
'IssueEvent','Sender','Name') <> ALL(?)`,
							"{ci-robot, openeuler-ci-bot, openeuler-sync-bot}")
				case "note":
					query = query.
						Where(`jsonb_extract_path_text(cloud_event_message.data_json, 
'NoteEvent','Sender','Name') <> ALL(?)`,
							"{ci-robot, openeuler-ci-bot, openeuler-sync-bot}")
				}

			}
		}
	}

	// gitee事件类型
	if params.EventType != "" {
		query = query.Where("cloud_event_message.type = ANY(?)", fmt.Sprintf("{%s}",
			params.EventType))
	}

	// 关键字模糊搜索
	if params.KeyWord != "" {
		query = query.Where("cloud_event_message.source_group ILIKE ?", "%"+params.KeyWord+"%")
	}

	var lSig []string
	// 我的sig组
	if params.MySig != "" {
		lSig, _ = utils.GetUserSigInfo(params.MySig)
	}

	// gitee仓库所属sig组
	if params.GiteeSigs != "" {
		lSig = append(lSig, strings.Split(params.GiteeSigs, ",")...)
	}

	if len(lSig) != 0 {
		query = query.
			Where("jsonb_extract_path_text(cloud_event_message.data_json, 'SigGroupName') = ANY(?)",
				fmt.Sprintf("{%s}", strings.Join(lSig, ",")))
	}

	var lRepo []string
	// 我管理的仓库
	if params.MyManagement != "" {
		lRepo, _ = utils.GetUserAdminRepos(params.MyManagement)
	}

	// 按仓库
	if params.Repos != "" {
		lRepo = append(lRepo, strings.Split(params.Repos, ",")...)
	}

	if len(lRepo) != 0 {
		query = query.Where("cloud_event_message.source_group = ANY(?)",
			fmt.Sprintf("{%s}", strings.Join(lSig, ",")))
	}

	// 按时间
	startTime := utils.ParseUnixTimestamp(params.StartTime)
	endTime := utils.ParseUnixTimestamp(params.EndTime)
	if startTime != nil && endTime != nil {
		query = query.Where("cloud_event_message.time BETWEEN ? AND ?", *startTime, *endTime)
	} else if startTime != nil {
		query = query.Where("cloud_event_message.time >= ?", *startTime)
	} else if endTime != nil {
		query = query.Where("cloud_event_message.time <= ?", *endTime)
	}

	// PullRequest事件的状态
	if params.PrState != "" {
		query = query.
			Where("jsonb_extract_path_text(cloud_event_message.data_json, 'PullRequestEvent',"+
				"'PullRequest', 'State') = ANY(?)", fmt.Sprintf("{%s}", params.PrState))
	}

	// PullRequest的创建者
	if params.PrCreator != "" {
		query = query.
			Where("jsonb_extract_path_text(cloud_event_message.data_json,"+
				"'PullRequestEvent','PullRequest','User','Login') = ANY(?)",
				fmt.Sprintf("{%s}", params.PrCreator))
	}

	// PullRequest指派人
	if params.PrAssignee != "" {
		query = query.Where("jsonb_extract_path_text(cloud_event_message.data_json,"+
			"'PullRequestEvent','PullRequest','Assignee','Login') = ANY(?)",
			fmt.Sprintf("{%s}", params.PrCreator))
	}

	// Issue事件状态
	if params.IssueState != "" {
		query = query.
			Where("jsonb_extract_path_text(cloud_event_message.data_json,"+
				"'IssueEvent', 'Issue', 'State') = ANY(?)",
				fmt.Sprintf("{%s}", params.IssueState))
	}

	// Issue创建者
	if params.IssueCreator != "" {
		query = query.
			Where("jsonb_extract_path_text(cloud_event_message.data_json,"+
				"'IssueEvent','Issue','User','Login') = ANY(?)",
				fmt.Sprintf("{%s}", params.IssueCreator))
	}

	// Issue指派人
	if params.IssueAssignee != "" {
		query = query.
			Where("jsonb_extract_path_text(cloud_event_message.data_json,"+
				"'IssueEvent','Issue','Assignee','Login') = ANY(?)",
				fmt.Sprintf("{%s}", params.IssueAssignee))
	}

	// 评论类型
	if params.NoteType != "" {
		query = query.
			Where("jsonb_extract_path_text(cloud_event_message.data_json,"+
				"'NoteEvent','NoteableType') = ANY(?)", fmt.Sprintf("{%s}", params.NoteType))
	}

	// @某人消息
	if params.About != "" {
		query = query.Where("jsonb_extract_path_text(cloud_event_message."+
			"data_json,'NoteEvent','Comment','Body') LIKE ?", "%"+params.About+"%")
	}

	// eur构建任务状态
	if params.BuildStatus != "" {
		query = query.
			Where("jsonb_extract_path_text(cloud_event_message.data_json,"+
				"'Body','Status') = ANY(?)", fmt.Sprintf("{%s}", params.BuildStatus))
	}

	// eur仓库owner
	if params.BuildOwner != "" {
		query = query.Where("cloud_event_message.user = ANY(?)", fmt.Sprintf("{%s}",
			params.BuildOwner))
	}

	// eur构建任务创建者
	if params.BuildCreator != "" {
		query = query.
			Where("jsonb_extract_path_text(cloud_event_message.data_json,"+
				"'Body','User') = ANY(?)", fmt.Sprintf("{%s}", params.BuildCreator))
	}

	if params.BuildEnv != "" {
		query = query.Where("jsonb_extract_path_text(cloud_event_message.data_json,"+
			"'Body','Chroot') = ANY(?)", fmt.Sprintf("{%s}", params.BuildEnv))
	}

	// 会议状态
	if params.MeetingAction != "" {
		query = query.
			Where("jsonb_extract_path_text(cloud_event_message.data_json,'Action') = ANY(?)",
				fmt.Sprintf("{%s}", params.MeetingAction))
	}

	// 会议所属sig组
	if params.MeetingSigGroup != "" {
		query = query.
			Where("jsonb_extract_path_text(cloud_event_message.data_json,"+
				"'Msg','GroupName') = ANY(?)", fmt.Sprintf("{%s}", params.MeetingSigGroup))
	}

	//会议时间
	if params.MeetingStartTime != "" && params.MeetingEndTime != "" {
		meetingStartTime := utils.ParseUnixTimestamp(params.MeetingStartTime)
		meetingEndTime := utils.ParseUnixTimestamp(params.MeetingEndTime)

		query = query.
			Where("jsonb_extract_path_text(cloud_event_message."+
				"data_json,'MeetingStartTime') BETWEEN ? AND ?",
				meetingStartTime.Format(time.DateTime), meetingEndTime.Format(time.DateTime))
	}

	if params.CVEComponent != "" {
		lComponent := strings.Split(params.CVEComponent, ",")
		var sql []string
		for _, comp := range lComponent {
			sql = append(sql, "%"+comp+"%")
		}
		query = query.
			Where("jsonb_extract_path_text(cloud_event_message.data_json,"+
				"'CVEComponent') LIKE ANY (?)", fmt.Sprintf("{%s}", strings.Join(sql, ",")))
	}

	if params.CVEState != "" {
		query = query.
			Where("jsonb_extract_path_text(cloud_event_message.data_json,"+
				"'IssueEvent','Issue','StateName') = ANY(?)", fmt.Sprintf("{%s}", params.CVEState))
	}

	if params.CVEAffected != "" {
		lAffected := strings.Split(params.CVEAffected, ",")
		var sql []string
		for _, affect := range lAffected {
			sql = append(sql, "%"+affect+"%")
		}
		query = query.Where("jsonb_extract_path_text(cloud_event_message.data_json,"+
			"'CVEAffectVersion') LIKE ANY (?)", fmt.Sprintf("{%s}", strings.Join(sql, ",")))
	}

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
	if err := postgresql.DB().Raw(sqlCount, false, userName, false, false).
		Scan(&CountData); err != nil {
		return []CountDAO{}, xerrors.Errorf("get count failed, err:%v", err)
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
