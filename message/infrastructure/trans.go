/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package infrastructure

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"golang.org/x/xerrors"
	"gorm.io/datatypes"

	"github.com/opensourceways/message-manager/utils"
)

func TransToDbFormat(source string, eventType string, request CmdToGetSubscribe) (datatypes.JSON,
	error) {
	if utils.IsEurMessage(source) {
		return TransEurModeFilterToDbFormat(request)
	} else if utils.IsGiteeMessage(source) {
		return TransGiteeModeFilterToDbFormat(eventType, request)
	} else if utils.IsMeetingMessage(source) {
		return TransMeetingModeFilterToDbFormat(request)
	} else if utils.IsCveMessage(source) {
		return TransCveModeFilterToDbFormat(request)
	}
	return nil, xerrors.Errorf("can not trans, source : %v", source)
}

func buildStringFilter(values []string) string {
	values = utils.RemoveEmptyStrings(values)
	if len(values) == 1 {
		return "eq=" + values[0]
	} else if len(values) > 1 {
		return "oneof=" + strings.Join(values, " ")
	}
	return ""
}

func buildTimeFilter(startTime, endTime *time.Time) string {
	if startTime != nil && endTime != nil {
		return fmt.Sprintf("gt=%s,lt=%s", startTime.Format(time.DateTime),
			endTime.Format(time.DateTime))
	}
	return ""
}

func marshalToJson(param interface{}) (datatypes.JSON, error) {
	jsonStr, err := json.Marshal(param)
	if err != nil {
		return datatypes.JSON{}, err
	}
	return jsonStr, nil
}

// TransEurModeFilterToDbFormat 处理 Eur 过滤器
func TransEurModeFilterToDbFormat(modeFilter CmdToGetSubscribe) (datatypes.JSON, error) {
	dbModeFilter := utils.EurDbFormat{
		Owner:  buildStringFilter(strings.Split(modeFilter.BuildOwner, ",")),
		User:   buildStringFilter(strings.Split(modeFilter.BuildCreator, ",")),
		Status: buildStringFilter(strings.Split(modeFilter.BuildStatus, ",")),
		Env:    buildStringFilter(strings.Split(modeFilter.BuildEnv, ",")),
		EventTime: buildTimeFilter(utils.ParseUnixTimestamp(modeFilter.StartTime),
			utils.ParseUnixTimestamp(modeFilter.EndTime)),
	}
	return marshalToJson(dbModeFilter)
}

// 获取机器人过滤器
func getBotFilter(isBot string) string {
	if isBot == "true" {
		return "oneof=openeuler-ci-bot ci-robot openeuler-sync-bot"
	} else if isBot == "false" {
		return "ne=openeuler-ci-bot,ne=ci-robot,ne=openeuler-sync-bot"
	}
	return ""
}

func getStateFilter(prState string) (state, action, mergeStatus string) {
	// 处理 PR 状态、动作和合并状态
	if prState != "" {
		prStates := utils.RemoveEmptyStrings(strings.Split(prState, ","))
		var lStatus, lAction, lMergeStatus []string

		for _, ps := range prStates {
			switch ps {
			case "can_be_merged", "cannot_be_merged":
				lMergeStatus = append(lMergeStatus, ps)
			case "set_draft":
				lAction = append(lAction, ps)
			default:
				lStatus = append(lStatus, ps)
			}
		}

		if len(lStatus) == 1 {
			state = "eq=" + lStatus[0]
		} else if len(lStatus) > 1 {
			state = "oneof=" + strings.Join(lStatus, " ")
		}

		if len(lAction) == 1 {
			action = "eq=" + lAction[0]
		} else if len(lAction) > 1 {
			action = "oneof=" + strings.Join(lAction, " ")
		}

		if len(lMergeStatus) == 1 {
			mergeStatus = "eq=" + lMergeStatus[0]
		} else if len(lMergeStatus) > 1 {
			mergeStatus = "oneof=" + strings.Join(lMergeStatus, " ")
		}
	}
	return
}

func getOneOfFilter(value string) string {
	if value != "" {
		return "oneof=" + value
	}
	return ""
}

func getNotOneOfFilter(value string) string {
	if value != "" {
		return "oneof!=" + value
	}
	return ""
}

// 获取 Gitee 数据库格式
func getGiteeDbFormat(eventType string, modeFilter CmdToGetSubscribe, lRepoName,
	lNameSpace []string) interface{} {
	sRepoName := buildStringFilter(lRepoName)
	sNamespace := buildStringFilter(lNameSpace)
	sMyManagement := getOneOfFilter(modeFilter.MyManagement)
	sOtherManagement := getNotOneOfFilter(modeFilter.OtherManagement)
	sMySig := getOneOfFilter(modeFilter.MySig)
	sOtherSig := getNotOneOfFilter(modeFilter.OtherSig)
	sSigGroupName := buildStringFilter(strings.Split(modeFilter.GiteeSigs, ","))
	sIsBot := getBotFilter(modeFilter.IsBot)
	eventTime := buildTimeFilter(utils.ParseUnixTimestamp(modeFilter.StartTime),
		utils.ParseUnixTimestamp(modeFilter.EndTime))

	sState, sAction, sMergeStatus := getStateFilter(modeFilter.PrState)

	switch eventType {
	case "issue":
		return utils.GiteeIssueDbFormat{
			RepoName:        sRepoName,
			IsBot:           sIsBot,
			Namespace:       sNamespace,
			SigGroupName:    sSigGroupName,
			IssueState:      buildStringFilter(strings.Split(modeFilter.IssueState, ",")),
			IssueCreator:    buildStringFilter(strings.Split(modeFilter.IssueCreator, ",")),
			IssueAssignee:   buildStringFilter(strings.Split(modeFilter.IssueAssignee, ",")),
			EventTime:       eventTime,
			MySig:           sMySig,
			MyManagement:    sMyManagement,
			OtherSig:        sOtherSig,
			OtherManagement: sOtherManagement,
		}
	case "note":
		return utils.GiteeNoteDbFormat{
			RepoName:        sRepoName,
			IsBot:           sIsBot,
			Namespace:       sNamespace,
			SigGroupName:    sSigGroupName,
			NoteType:        buildStringFilter(strings.Split(modeFilter.NoteType, ",")),
			About:           buildStringFilter([]string{modeFilter.About}),
			EventTime:       eventTime,
			MySig:           sMySig,
			MyManagement:    sMyManagement,
			OtherSig:        sOtherSig,
			OtherManagement: sOtherManagement,
		}
	case "pr":
		return utils.GiteePullRequestDbFormat{
			RepoName:        sRepoName,
			IsBot:           sIsBot,
			Namespace:       sNamespace,
			SigGroupName:    sSigGroupName,
			PrState:         sState,
			PrAction:        sAction,
			PrMergeStatus:   sMergeStatus,
			PrCreator:       buildStringFilter(strings.Split(modeFilter.PrCreator, ",")),
			PrAssignee:      buildStringFilter(strings.Split(modeFilter.PrAssignee, ",")),
			EventTime:       eventTime,
			MySig:           sMySig,
			MyManagement:    sMyManagement,
			OtherSig:        sOtherSig,
			OtherManagement: sOtherManagement,
		}
	case "push":
		return utils.GiteePushDbFormat{
			RepoName:        sRepoName,
			IsBot:           sIsBot,
			Namespace:       sNamespace,
			SigGroupName:    sSigGroupName,
			EventTime:       eventTime,
			MySig:           sMySig,
			MyManagement:    sMyManagement,
			OtherSig:        sOtherSig,
			OtherManagement: sOtherManagement,
		}
	default:
		return nil
	}
}

// TransGiteeModeFilterToDbFormat 处理 GiteeMode 过滤器
func TransGiteeModeFilterToDbFormat(eventType string, modeFilter CmdToGetSubscribe) (
	datatypes.JSON, error) {
	repoNames := utils.RemoveEmptyStrings(strings.Split(modeFilter.Repos, ","))
	var lRepoName, lNameSpace []string

	// 处理仓库名和命名空间
	for _, repoName := range utils.MergePaths(repoNames) {
		if repoName == "*" {
			lRepoName, lNameSpace = []string{}, []string{}
			break
		}
		lName := strings.Split(repoName, "/")
		if lName[1] == "*" {
			lNameSpace = append(lNameSpace, lName[0])
		} else {
			lRepoName = append(lRepoName, repoName)
		}
	}
	dbModeFilter := getGiteeDbFormat(eventType, modeFilter, lRepoName, lNameSpace)
	return marshalToJson(dbModeFilter)
}

// TransMeetingModeFilterToDbFormat 处理 MeetingMode 过滤器
func TransMeetingModeFilterToDbFormat(modeFilter CmdToGetSubscribe) (datatypes.JSON, error) {
	dbModeFilter := utils.MeetingDbFormat{
		Action:   buildStringFilter(strings.Split(modeFilter.MeetingAction, ",")),
		SigGroup: buildStringFilter(strings.Split(modeFilter.MeetingSigGroup, ",")),
		MeetingStartTime: buildTimeFilter(utils.ParseUnixTimestamp(modeFilter.StartTime),
			utils.ParseUnixTimestamp(modeFilter.EndTime)),
		EventTime: buildTimeFilter(utils.ParseUnixTimestamp(modeFilter.StartTime),
			utils.ParseUnixTimestamp(modeFilter.EndTime)),
		MySig:    getOneOfFilter(modeFilter.MySig),
		OtherSig: getNotOneOfFilter(modeFilter.OtherSig),
	}
	return marshalToJson(dbModeFilter)
}

// 处理 CVE 组件过滤器
func buildCveComponentFilter(component string) string {
	if component == "" {
		return ""
	}
	lComponent := strings.Split(component, ",")
	var template []string
	for _, comp := range lComponent {
		template = append(template, "contains="+comp)
	}
	return "or " + strings.Join(template, " ")
}

// 处理 CVE 影响过滤器
func buildCveAffectedFilter(affected string) string {
	if affected == "" {
		return ""
	}
	lAffected := strings.Split(affected, ",")
	var template []string
	for _, aff := range lAffected {
		template = append(template, "contains="+aff)
	}
	if len(template) == 1 {
		return template[0]
	} else if len(template) > 1 {
		return "or " + strings.Join(template, " ")
	}
	return ""
}

// TransCveModeFilterToDbFormat 处理 CVE 过滤器
func TransCveModeFilterToDbFormat(modeFilter CmdToGetSubscribe) (datatypes.JSON, error) {
	dbModeFilter := utils.CveDbFormat{
		Component:    buildCveComponentFilter(modeFilter.CVEComponent),
		State:        buildStringFilter(strings.Split(modeFilter.IssueState, ",")),
		Affected:     buildCveAffectedFilter(modeFilter.CVEAffected),
		SigGroupName: buildStringFilter(strings.Split(modeFilter.GiteeSigs, ",")),
		MySig:        getOneOfFilter(modeFilter.MySig),
		OtherSig:     getNotOneOfFilter(modeFilter.OtherSig),
	}
	return marshalToJson(dbModeFilter)
}
