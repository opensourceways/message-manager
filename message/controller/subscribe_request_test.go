package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/datatypes"
)

func TestSubscribeDTOToCmd(t *testing.T) {
	req := &subscribeDTO{
		queryInnerParams: queryInnerParams{
			Source:           "test_source",
			EventType:        "test_event",
			IsRead:           "true",
			KeyWord:          "keyword",
			IsBot:            "false",
			GiteeSigs:        "sig",
			Repos:            "repo",
			CountPerPage:     10,
			PageNum:          1,
			StartTime:        "2024-01-01T00:00:00Z",
			EndTime:          "2024-01-02T00:00:00Z",
			MySig:            "my_sig",
			MyManagement:     "management",
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
		},
		SpecVersion: "1.0",
		ModeName:    "test_mode",
	}

	cmd, err := req.toCmd()
	assert.NoError(t, err)
	assert.Equal(t, "test_source", cmd.Source)
	assert.Equal(t, "test_event", cmd.EventType)
	assert.Equal(t, "1.0", cmd.SpecVersion)
	assert.Equal(t, "test_mode", cmd.ModeName)
	// 继续验证其他字段...
}

func TestNewSubscribeDTOToCmd(t *testing.T) {
	req := &newSubscribeDTO{
		Source:      "new_source",
		EventType:   "new_event",
		SpecVersion: "1.1",
		ModeFilter:  datatypes.JSON([]byte(`{"filter": "value"}`)),
		ModeName:    "new_mode",
	}

	cmd, err := req.toCmd()
	assert.NoError(t, err)
	assert.Equal(t, "new_source", cmd.Source)
	assert.Equal(t, "new_event", cmd.EventType)
	assert.Equal(t, "1.1", cmd.SpecVersion)
	assert.JSONEq(t, `{"filter": "value"}`, string(cmd.ModeFilter))
	assert.Equal(t, "new_mode", cmd.ModeName)
}

func TestDeleteSubscribeDTOToCmd(t *testing.T) {
	req := &deleteSubscribeDTO{
		Source:   "delete_source",
		ModeName: "delete_mode",
	}

	cmd, err := req.toCmd()
	assert.NoError(t, err)
	assert.Equal(t, "delete_source", cmd.Source)
	assert.Equal(t, "delete_mode", cmd.ModeName)
}
