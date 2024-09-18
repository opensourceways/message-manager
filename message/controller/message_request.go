/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package controller

import "github.com/opensourceways/message-manager/message/app"

type queryInnerParams struct {
	Source          string `form:"source" json:"source"`         // 消息源
	EventType       string `form:"event_type" json:"event_type"` // 事件类型
	IsRead          string `form:"is_read" json:"is_read"`       // 是否已读
	KeyWord         string `form:"key_word" json:"key_word"`     // 关键字模糊搜索
	IsBot           string `form:"is_bot" json:"is_bot"`         // 是否机器人
	GiteeSigs       string `form:"sig" json:"sig"`               // sig组筛选
	Repos           string `form:"repos" json:"repos"`           // 仓库筛选
	CountPerPage    int    `form:"count_per_page" json:"count_per_page"`
	PageNum         int    `form:"page" json:"page"`
	StartTime       string `form:"start_time" json:"start_time"`         // 起始时间
	EndTime         string `form:"end_time" json:"end_time"`             // 结束时间
	MySig           string `form:"my_sig" json:"my_sig"`                 // 我的sig组
	MyManagement    string `form:"my_management" json:"my_management"`   // 我管理的仓库
	PrState         string `form:"pr_state" json:"pr_state"`             // pr事件状态
	PrCreator       string `form:"pr_creator" json:"pr_creator"`         // pr提交者
	PrAssignee      string `form:"pr_assignee" json:"pr_assignee"`       // pr指派者
	IssueState      string `form:"issue_state" json:"issue_state"`       // issue事件状态
	IssueCreator    string `form:"issue_creator" json:"issue_creator"`   // issue提交者
	IssueAssignee   string `form:"issue_assignee" json:"issue_assignee"` // issue指派者
	NoteType        string `form:"note_type" json:"note_type"`           // 评论类型
	About           string `form:"about" json:"about"`                   // @我的
	BuildStatus     string `form:"build_status" json:"build_status"`     // eur构建状态
	BuildOwner      string `form:"build_owner" json:"build_owner"`       // eur我的项目
	BuildCreator    string `form:"build_creator" json:"build_creator"`   // eur我执行的
	BuildEnv        string `form:"build_env" json:"build_env"`           // eur构建环境
	MeetingAction   string `form:"meeting_action" json:"meeting_action"` // 会议操作
	MeetingSigGroup string `form:"meeting_sig" json:"meeting_sig"`       // 会议所属sig
	// 会议开始时间
	MeetingStartTime string `form:"meeting_start_time" json:"meeting_start_time"`
	// 会议结束时间
	MeetingEndTime string `form:"meeting_end_time" json:"meeting_end_time"`
	CVEComponent   string `form:"cve_component" json:"cve_component"` // cve组件仓库
	CVEState       string `form:"cve_state" json:"cve_state"`         // cve漏洞状态
	CVEAffected    string `form:"cve_affected" json:"cve_affected"`   // cve影响系统版本
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
