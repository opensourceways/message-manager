package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryInnerParamsToCmd(t *testing.T) {
	req := &queryInnerParams{
		Source:           "test_source",
		EventType:        "test_event",
		IsRead:           "true",
		KeyWord:          "test_keyword",
		IsBot:            "false",
		GiteeSigs:        "test_sig",
		Repos:            "test_repo",
		CountPerPage:     10,
		PageNum:          1,
		StartTime:        "2024-01-01T00:00:00Z",
		EndTime:          "2024-01-02T00:00:00Z",
		MySig:            "my_sig",
		MyManagement:     "my_management",
		PrState:          "open",
		PrCreator:        "creator",
		PrAssignee:       "assignee",
		IssueState:       "closed",
		IssueCreator:     "issue_creator",
		IssueAssignee:    "issue_assignee",
		NoteType:         "note",
		About:            "about",
		BuildStatus:      "success",
		BuildOwner:       "owner",
		BuildCreator:     "builder",
		BuildEnv:         "production",
		MeetingAction:    "create",
		MeetingSigGroup:  "sig_group",
		MeetingStartTime: "2024-01-01T09:00:00Z",
		MeetingEndTime:   "2024-01-01T10:00:00Z",
		CVEComponent:     "cve_component",
		CVEState:         "open",
		CVEAffected:      "yes",
	}

	cmd, err := req.toCmd()
	assert.NoError(t, err)
	assert.Equal(t, "test_source", cmd.Source)
	assert.Equal(t, "test_event", cmd.EventType)
	assert.Equal(t, "true", cmd.IsRead)
	assert.Equal(t, "test_keyword", cmd.KeyWord)
	assert.Equal(t, 10, cmd.CountPerPage)
	// 继续验证其他字段...
}

func TestQueryInnerParamsQuickToCmd(t *testing.T) {
	req := &queryInnerParamsQuick{
		Source:       "quick_source",
		CountPerPage: 5,
		PageNum:      2,
		ModeName:     "quick_mode",
	}

	cmd, err := req.toCmd()
	assert.NoError(t, err)
	assert.Equal(t, "quick_source", cmd.Source)
	assert.Equal(t, 5, cmd.CountPerPage)
	assert.Equal(t, "quick_mode", cmd.ModeName)
}

func TestMessageStatusToCmd(t *testing.T) {
	req := &messageStatus{
		EventId: "event_123",
	}
	cmd, err := req.toCmd()
	assert.NoError(t, err)
	assert.Equal(t, "event_123", cmd.EventId)
}
