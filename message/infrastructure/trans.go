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

func TransEurModeFilterToDbFormat(modeFilter CmdToGetSubscribe) (datatypes.JSON, error) {
	var sOwner, sUser, sStatus, sEnv string
	if modeFilter.BuildOwner != "" {
		lOwner := utils.RemoveEmptyStrings(strings.Split(modeFilter.BuildOwner, ","))
		if len(lOwner) == 1 {
			sOwner = "eq=" + lOwner[0]
		} else if len(lOwner) > 1 {
			sOwner = "oneof=" + strings.Join(lOwner, " ")
		}
	}
	if modeFilter.BuildCreator != "" {
		lUser := utils.RemoveEmptyStrings(strings.Split(modeFilter.BuildCreator, ","))
		if len(lUser) == 1 {
			sUser = "eq=" + lUser[0]
		} else if len(lUser) > 1 {
			sUser = "oneof=" + strings.Join(lUser, " ")
		}
	}
	if modeFilter.BuildStatus != "" {
		lStatus := utils.RemoveEmptyStrings(strings.Split(modeFilter.BuildStatus, ","))
		if len(lStatus) == 1 {
			sStatus = "eq=" + lStatus[0]
		} else if len(lStatus) > 1 {
			sStatus = "oneof=" + strings.Join(lStatus, " ")
		}
	}
	if modeFilter.BuildEnv != "" {
		lEnv := utils.RemoveEmptyStrings(strings.Split(modeFilter.BuildEnv, ","))
		if len(lEnv) == 1 {
			sEnv = "eq=" + lEnv[0]
		} else if len(lEnv) > 1 {
			sEnv = "oneof=" + strings.Join(lEnv, " ")
		}
	}

	startTime := utils.ParseUnixTimestamp(modeFilter.StartTime)
	endTime := utils.ParseUnixTimestamp(modeFilter.EndTime)
	var sEventTime string
	if startTime != nil && endTime != nil {
		sEventTime = fmt.Sprintf("gt=%s,lt=%s", startTime.Format(time.DateTime),
			endTime.Format(time.DateTime))

	}
	dbModeFilter := utils.EurDbFormat{
		Status:    sStatus,
		Owner:     sOwner,
		User:      sUser,
		Env:       sEnv,
		EventTime: sEventTime,
	}

	jsonStr, err := json.Marshal(dbModeFilter)
	if err != nil {
		return datatypes.JSON{}, err
	}

	return jsonStr, nil
}

func TransGiteeModeFilterToDbFormat(eventType string, modeFilter CmdToGetSubscribe) (
	datatypes.JSON, error) {
	repoNames := utils.RemoveEmptyStrings(strings.Split(modeFilter.Repos, ","))
	var lRepoName, lNameSpace []string
	var sRepoName, sNamespace string

	if len(repoNames) != 0 {
		mergeRepoNames := utils.MergePaths(repoNames)
		for _, repoName := range mergeRepoNames {
			if repoName == "*" {
				lRepoName = []string{}
				lNameSpace = []string{}
				break
			} else {
				lName := strings.Split(repoName, "/")
				if lName[1] == "*" {
					lNameSpace = append(lNameSpace, lName[0])
				} else {
					lRepoName = append(lRepoName, repoName)
				}
			}
		}

	}
	lRepoName = utils.RemoveEmptyStrings(lRepoName)
	if len(lRepoName) == 1 {
		sRepoName = "eq=" + lRepoName[0]
	} else if len(lRepoName) > 1 {
		sRepoName = "oneof=" + strings.Join(lRepoName, " ")
	}

	if len(lNameSpace) == 1 {
		sNamespace = "eq=" + lNameSpace[0]
	} else if len(lNameSpace) > 1 {
		sNamespace = "oneof=" + strings.Join(lNameSpace, " ")
	}

	var sMyManagement, sMySig string
	if modeFilter.MyManagement != "" {
		sMyManagement = fmt.Sprintf("oneof=%s", modeFilter.MyManagement)
	}
	if modeFilter.MySig != "" {
		sMySig = fmt.Sprintf("oneof=%s", modeFilter.MySig)
	}
	var sSigGroupName string
	if modeFilter.GiteeSigs != "" {
		sigGroupNames := strings.Split(modeFilter.GiteeSigs, ",")
		sigGroupNames = utils.RemoveEmptyStrings(sigGroupNames)
		if len(sigGroupNames) == 1 {
			sSigGroupName = "eq=" + sigGroupNames[0]
		} else if len(sigGroupNames) > 1 {
			sSigGroupName = "oneof=" + strings.Join(sigGroupNames, " ")
		}
	}
	var sIsBot string
	if modeFilter.IsBot == "true" {
		sIsBot = "oneof=openeuler-ci-bot ci-robot openeuler-sync-bot"
	} else if modeFilter.IsBot == "false" {
		sIsBot = "ne=openeuler-ci-bot,ne=ci-robot,ne=openeuler-sync-bot"
	}

	var sPrState, sAction, sPrMergeStatus string
	if modeFilter.PrState != "" {
		prStates := utils.RemoveEmptyStrings(strings.Split(modeFilter.PrState, ","))
		var lStatus, lAction, lMergeStatus []string
		for _, ps := range prStates {
			if ps == "can_be_merged" || ps == "cannot_be_merged" {
				lMergeStatus = append(lMergeStatus, ps)
			} else if ps == "set_draft" {
				lAction = append(lAction, ps)
			} else {
				lStatus = append(lStatus, ps)
			}
		}

		if len(lStatus) == 1 {
			sPrState = "eq=" + lStatus[0]
		} else if len(lStatus) > 1 {
			sPrState = "oneof=" + strings.Join(lStatus, " ")
		}

		if len(lAction) == 1 {
			sAction = "eq=" + lAction[0]
		} else if len(lAction) > 1 {
			sAction = "oneof=" + strings.Join(lAction, " ")
		}

		if len(lMergeStatus) == 1 {
			sPrMergeStatus = "eq=" + lMergeStatus[0]
		} else if len(lStatus) > 1 {
			sPrMergeStatus = "oneof=" + strings.Join(lMergeStatus, " ")
		}
	}
	var sPrCreator string
	if modeFilter.PrCreator != "" {
		prCreators := utils.RemoveEmptyStrings(strings.Split(modeFilter.
			PrCreator, ","))
		if len(prCreators) == 1 {
			sPrCreator = "eq=" + prCreators[0]
		} else if len(prCreators) > 1 {
			sPrCreator = "oneof=" + strings.Join(prCreators, " ")
		}
	}
	var sPrAssignee string
	if modeFilter.PrAssignee != "" {
		prAssignees := utils.RemoveEmptyStrings(strings.Split(modeFilter.PrAssignee, ","))
		if len(prAssignees) == 1 {
			sPrAssignee = "eq=" + prAssignees[0]
		} else if len(prAssignees) > 1 {
			sPrAssignee = "oneof=" + strings.Join(prAssignees, " ")
		}
	}
	var sIssueState string
	if modeFilter.IssueState != "" {
		issueStates := utils.RemoveEmptyStrings(strings.Split(modeFilter.IssueState, ","))
		if len(issueStates) == 1 {
			sIssueState = "eq=" + issueStates[0]
		} else if len(issueStates) > 1 {
			sIssueState = "oneof=" + strings.Join(issueStates, " ")
		}
	}

	var sIssueCreator string
	if modeFilter.IssueCreator != "" {
		issueCreators := utils.RemoveEmptyStrings(strings.Split(modeFilter.IssueCreator, ","))
		if len(issueCreators) == 1 {
			sIssueCreator = "eq=" + issueCreators[0]
		} else if len(issueCreators) > 1 {
			sIssueCreator = "oneof=" + strings.Join(issueCreators, " ")
		}
	}
	var sIssueAssignee string
	if modeFilter.IssueAssignee != "" {
		issueAssignees := utils.RemoveEmptyStrings(strings.Split(modeFilter.IssueAssignee, ","))
		if len(issueAssignees) == 1 {
			sIssueAssignee = "eq=" + issueAssignees[0]
		} else if len(issueAssignees) > 1 {
			sIssueAssignee = "oneof=" + strings.Join(issueAssignees, " ")
		}
	}
	var sNoteType string
	if modeFilter.NoteType != "" {
		noteTypes := utils.RemoveEmptyStrings(strings.Split(modeFilter.NoteType, ","))
		if len(noteTypes) == 1 {
			sNoteType = "eq=" + noteTypes[0]
		} else if len(noteTypes) > 1 {
			sNoteType = "oneof=" + strings.Join(noteTypes, " ")
		}
	}
	var sAbout string
	if modeFilter.About != "" {
		sAbout = "contains=" + sAbout
	}

	startTime := utils.ParseUnixTimestamp(modeFilter.StartTime)
	endTime := utils.ParseUnixTimestamp(modeFilter.EndTime)
	var sEventTime string
	if startTime != nil && endTime != nil {
		sEventTime = fmt.Sprintf("gt=%s,lt=%s", startTime.Format(time.DateTime),
			endTime.Format(time.DateTime))
	}

	var param interface{}
	switch eventType {
	case "issue":
		param = utils.GiteeIssueDbFormat{
			RepoName:      sRepoName,
			IsBot:         sIsBot,
			Namespace:     sNamespace,
			SigGroupName:  sSigGroupName,
			IssueState:    sIssueState,
			IssueCreator:  sIssueCreator,
			IssueAssignee: sIssueAssignee,
			EventTime:     sEventTime,
			MySig:         sMySig,
			MyManagement:  sMyManagement,
		}
	case "note":
		param = utils.GiteeNoteDbFormat{
			RepoName:     sRepoName,
			IsBot:        sIsBot,
			Namespace:    sNamespace,
			SigGroupName: sSigGroupName,
			NoteType:     sNoteType,
			About:        sAbout,
			EventTime:    sEventTime,
			MySig:        sMySig,
			MyManagement: sMyManagement,
		}
	case "pr":
		param = utils.GiteePullRequestDbFormat{
			RepoName:      sRepoName,
			IsBot:         sIsBot,
			Namespace:     sNamespace,
			SigGroupName:  sSigGroupName,
			PrState:       sPrState,
			PrAction:      sAction,
			PrMergeStatus: sPrMergeStatus,
			PrCreator:     sPrCreator,
			PrAssignee:    sPrAssignee,
			EventTime:     sEventTime,
			MySig:         sMySig,
			MyManagement:  sMyManagement,
		}
	case "push":
		param = utils.GiteePushDbFormat{
			RepoName:     sRepoName,
			IsBot:        sIsBot,
			Namespace:    sNamespace,
			SigGroupName: sSigGroupName,
			EventTime:    sEventTime,
			MySig:        sMySig,
			MyManagement: sMyManagement,
		}
	default:
		return datatypes.JSON{}, xerrors.Errorf("the eventType is wrong")
	}

	jsonStr, err := json.Marshal(param)
	if err != nil {
		return datatypes.JSON{}, err
	}
	return jsonStr, nil
}

func TransMeetingModeFilterToDbFormat(modeFilter CmdToGetSubscribe) (datatypes.JSON,
	error) {
	var sAction, sSigGroup, sStart, sEventTime string
	if modeFilter.MeetingAction != "" {
		lAction := utils.RemoveEmptyStrings(strings.Split(modeFilter.MeetingAction, ","))
		if len(lAction) == 1 {
			sAction = "eq=" + lAction[0]
		} else if len(lAction) > 1 {
			sAction = "oneof=" + strings.Join(lAction, " ")
		}
	}
	if modeFilter.MeetingSigGroup != "" {
		lSig := utils.RemoveEmptyStrings(strings.Split(modeFilter.MeetingSigGroup, ","))
		if len(lSig) == 1 {
			sSigGroup = "eq=" + lSig[0]
		} else if len(lSig) > 1 {
			sSigGroup = "oneof=" + strings.Join(lSig, " ")
		}
	}
	meetingStartTime := utils.ParseUnixTimestamp(modeFilter.StartTime)
	meetingEndTime := utils.ParseUnixTimestamp(modeFilter.EndTime)
	if meetingStartTime != nil && meetingEndTime != nil {
		sStart = fmt.Sprintf("gt=%s,lt=%s", meetingStartTime.Format(time.DateTime),
			meetingEndTime.Format(time.DateTime))

	}

	var sMySig string
	if modeFilter.MySig != "" {
		sMySig = fmt.Sprintf("oneof=%s", modeFilter.MySig)
	}

	startTime := utils.ParseUnixTimestamp(modeFilter.StartTime)
	endTime := utils.ParseUnixTimestamp(modeFilter.EndTime)
	if startTime != nil && endTime != nil {
		sEventTime = fmt.Sprintf("gt=%s,lt=%s", startTime.Format(time.DateTime),
			endTime.Format(time.DateTime))

	}
	dbModeFilter := utils.MeetingDbFormat{
		Action:           sAction,
		SigGroup:         sSigGroup,
		MeetingStartTime: sStart,
		EventTime:        sEventTime,
		MySig:            sMySig,
	}

	jsonStr, err := json.Marshal(dbModeFilter)
	if err != nil {
		return datatypes.JSON{}, err
	}

	return jsonStr, nil
}

func TransCveModeFilterToDbFormat(modeFilter CmdToGetSubscribe) (datatypes.JSON, error) {
	var sComponent, sState, sAffected string

	if modeFilter.CVEComponent != "" {
		lComponent := strings.Split(modeFilter.CVEComponent, ",")
		var template []string
		for _, comp := range lComponent {
			template = append(template, "contains="+comp)
		}
		sComponent = "or " + strings.Join(template, " ")
	}

	if modeFilter.IssueState != "" {
		cveStates := utils.RemoveEmptyStrings(strings.Split(modeFilter.IssueState, ","))
		if len(cveStates) == 1 {
			sState = "eq=" + cveStates[0]
		} else if len(cveStates) > 1 {
			sState = "oneof=" + strings.Join(cveStates, " ")
		}
	}

	if modeFilter.CVEAffected != "" {
		lAffected := strings.Split(modeFilter.CVEAffected, ",")
		var template []string
		for _, affect := range lAffected {
			template = append(template, "contains="+affect)
		}
		sAffected = "or " + strings.Join(template, " ")
	}

	var sMySig string
	if modeFilter.MySig != "" {
		sMySig = fmt.Sprintf("oneof=%s", modeFilter.MySig)
	}

	dbModeFilter := utils.CveDbFormat{
		Component: sComponent,
		State:     sState,
		Affected:  sAffected,
		MySig:     sMySig,
	}

	jsonStr, err := json.Marshal(dbModeFilter)
	if err != nil {
		return datatypes.JSON{}, err
	}

	return jsonStr, nil
}
