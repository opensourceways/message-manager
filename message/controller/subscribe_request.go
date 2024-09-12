/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package controller

import (
	"github.com/opensourceways/message-manager/message/app"
	"gorm.io/datatypes"
)

type subscribeDTO struct {
	queryInnerParams
	SpecVersion string `json:"spec_version"`
	ModeName    string `json:"mode_name"`
}

func (req *subscribeDTO) toCmd() (cmd app.CmdToGetSubscribe, err error) {
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
	cmd.SpecVersion = req.SpecVersion
	cmd.ModeName = req.ModeName
	return cmd, nil
}

type newSubscribeDTO struct {
	Source      string         `json:"source"`
	EventType   string         `json:"event_type"`
	SpecVersion string         `json:"spec_version"`
	ModeFilter  datatypes.JSON `json:"mode_filter" swaggerignore:"true"`
	ModeName    string         `json:"mode_name"`
}

func (req *newSubscribeDTO) toCmd() (cmd app.CmdToAddSubscribe, err error) {
	cmd.Source = req.Source
	cmd.EventType = req.EventType
	cmd.SpecVersion = req.SpecVersion
	cmd.ModeFilter = req.ModeFilter
	cmd.ModeName = req.ModeName
	return
}

type updateSubscribeDTO struct {
	Source  string `json:"source"`
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

func (req *updateSubscribeDTO) toCmd() (cmd app.CmdToUpdateSubscribe, err error) {
	cmd.Source = req.Source
	cmd.OldName = req.OldName
	cmd.NewName = req.NewName
	return
}

type deleteSubscribeDTO struct {
	Source    string `json:"source"`
	EventType string `json:"event_type"`
	ModeName  string `json:"mode_name"`
}

func (req *deleteSubscribeDTO) toCmd() (cmd app.CmdToDeleteSubscribe, err error) {
	cmd.Source = req.Source
	cmd.EventType = req.EventType
	cmd.ModeName = req.ModeName
	return
}
