/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package controller

import "github.com/opensourceways/message-manager/message/app"

type queryInnerParams struct {
	Source           string `json:"source"`     // 消息源
	EventType        string `json:"event_type"` // 事件类型
	IsRead           string `json:"is_read"`    // 是否已读
	KeyWord          string `json:"key_word"`   // 关键字模糊搜索
	IsBot            string `json:"is_bot"`     // 是否机器人
	GiteeSigs        string `json:"sig"`        // sig组筛选
	Repos            string `json:"repos"`      // 仓库筛选
	CountPerPage     int    `json:"count_per_page"`
	PageNum          int    `json:"page"`
	StartTime        string `json:"start_time"`         // 起始时间
	EndTime          string `json:"end_time"`           // 结束时间
	MySig            string `json:"my_sig"`             // 我的sig组
	MyManagement     string `json:"my_management"`      // 我管理的仓库
	PrState          string `json:"pr_state"`           // pr事件状态
	PrCreator        string `json:"pr_creator"`         // pr提交者
	PrAssignee       string `json:"pr_assignee"`        // pr指派者
	IssueState       string `json:"issue_state"`        // issue事件状态
	IssueCreator     string `json:"issue_creator"`      // issue提交者
	IssueAssignee    string `json:"issue_assignee"`     // issue指派者
	NoteType         string `json:"note_type"`          // 评论类型
	About            string `json:"about"`              // @我的
	BuildStatus      string `json:"build_status"`       // eur构建状态
	BuildOwner       string `json:"build_owner"`        // eur我的项目
	BuildCreator     string `json:"build_creator"`      // eur我执行的
	BuildEnv         string `json:"build_env"`          // eur构建环境
	MeetingAction    string `json:"meeting_action"`     // 会议操作
	MeetingSigGroup  string `json:"meeting_sig"`        // 会议所属sig
	MeetingStartTime string `json:"meeting_start_time"` // 会议开始时间

	MeetingEndTime string `json:"meeting_end_time"` // 会议结束时间
	CVEComponent   string `json:"cve_component"`    // cve组件仓库
	CVEState       string `json:"cve_state"`        // cve漏洞状态
	CVEAffected    string `json:"cve_affected"`     // cve影响系统版本
}

func (req *queryInnerParams) toCmd() (cmd app.CmdToGetInnerMessage, err error) {
	cmd.Source = req.Source
	cmd.EventType = req.EventType
	cmd.IsRead = req.IsRead
	cmd.KeyWord = req.KeyWord
	cmd.IsBot = req.IsBot
	cmd.GiteeSigs = req.GiteeSigs
	cmd.Repos = req.Repos
	cmd.CountPerPage = req.CountPerPage
	cmd.PageNum = req.PageNum
	cmd.StartTime = req.StartTime
	cmd.EndTime = req.EndTime
	cmd.MySig = req.MySig
	cmd.MyManagement = req.MyManagement
	cmd.PrState = req.PrState
	cmd.PrCreator = req.PrCreator
	cmd.PrAssignee = req.PrAssignee
	cmd.IssueState = req.IssueState
	cmd.IssueCreator = req.IssueCreator
	cmd.IssueAssignee = req.IssueAssignee
	cmd.NoteType = req.NoteType
	cmd.About = req.About
	cmd.BuildStatus = req.BuildStatus
	cmd.BuildOwner = req.BuildOwner
	cmd.BuildCreator = req.BuildCreator
	cmd.BuildEnv = req.BuildEnv
	cmd.MeetingAction = req.MeetingAction
	cmd.MeetingSigGroup = req.MeetingSigGroup
	cmd.MeetingStartTime = req.MeetingStartTime
	cmd.MeetingEndTime = req.MeetingEndTime
	cmd.CVEComponent = req.CVEComponent
	cmd.CVEState = req.CVEState
	cmd.CVEAffected = req.CVEAffected
	return cmd, nil
}

type queryInnerParamsQuick struct {
	Source       string `form:"source" json:"source"`
	CountPerPage int    `form:"count_per_page" json:"count_per_page"`
	PageNum      int    `form:"page" json:"page"`
	ModeName     string `form:"mode_name" json:"mode_name"`
}

func (req *queryInnerParamsQuick) toCmd() (cmd app.CmdToGetInnerMessageQuick, err error) {
	cmd.Source = req.Source
	cmd.CountPerPage = req.CountPerPage
	cmd.ModeName = req.ModeName
	cmd.PageNum = req.PageNum
	return cmd, nil
}

type messageStatus struct {
	Source  string `json:"source"`
	EventId string `json:"event_id"`
}

func (req *messageStatus) toCmd() (cmd app.CmdToSetIsRead, err error) {
	cmd.EventId = req.EventId
	cmd.Source = req.Source
	return cmd, nil
}

type QueryParams struct {
	GiteeUserName string `form:"gitee_user_name"`
	IsBot         bool   `form:"is_bot"`
	Filter        int    `form:"filter"`
	IsDone        bool   `form:"is_done"`
}
