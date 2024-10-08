package infrastructure

import (
	"encoding/json"
	"testing"

	"github.com/opensourceways/message-manager/utils"
)

func TestTransToDbFormat(t *testing.T) {
	tests := []struct {
		source    string
		eventType string
		filter    CmdToGetSubscribe
		expectErr bool
	}{
		{
			source:    utils.EurSource,
			eventType: "issue",
			filter: CmdToGetSubscribe{
				BuildOwner:   "owner1",
				BuildCreator: "creator1",
				BuildStatus:  "success",
				BuildEnv:     "prod",
				StartTime:    "1622505600",
				EndTime:      "1622592000",
			},
			expectErr: false,
		},
		{
			source:    utils.GiteeSource,
			eventType: "note",
			filter: CmdToGetSubscribe{
				Repos:        "repo1/repo2",
				MyManagement: "management1",
				MySig:        "sig1",
			},
			expectErr: false,
		},
		{
			source:    utils.MeetingSource,
			eventType: "meeting",
			filter: CmdToGetSubscribe{
				MeetingAction:   "start,stop",
				MeetingSigGroup: "group1",
				StartTime:       "1622505600",
				EndTime:         "1622592000",
			},
			expectErr: false,
		},
		{
			source:    utils.CveSource,
			eventType: "cve",
			filter: CmdToGetSubscribe{
				CVEComponent: "component1,component2",
				IssueState:   "open,closed",
			},
			expectErr: false,
		},
		{
			source:    "unknown",
			eventType: "issue",
			filter:    CmdToGetSubscribe{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			result, err := TransToDbFormat(tt.source, tt.eventType, tt.filter)
			if (err != nil) != tt.expectErr {
				t.Errorf("TransToDbFormat() error = %v, wantErr %v", err, tt.expectErr)
			}

			if !tt.expectErr {
				var dbFormat map[string]interface{}
				if err := json.Unmarshal(result, &dbFormat); err != nil {
					t.Fatalf("Failed to unmarshal result: %v", err)
				}
			}
		})
	}
}

func TestTransEurModeFilterToDbFormat(t *testing.T) {
	tests := []struct {
		filter    CmdToGetSubscribe
		expectErr bool
	}{
		{
			filter: CmdToGetSubscribe{
				BuildOwner:   "owner1",
				BuildCreator: "creator1",
				BuildStatus:  "success",
				BuildEnv:     "prod",
				StartTime:    "1622505600",
				EndTime:      "1622592000",
			},
			expectErr: false,
		},
		{
			filter: CmdToGetSubscribe{
				BuildOwner:   "owner2,owner3",
				BuildCreator: "",
				BuildStatus:  "",
				BuildEnv:     "test",
				StartTime:    "1622505600",
				EndTime:      "1622592000",
			},
			expectErr: false,
		},
		{
			filter: CmdToGetSubscribe{
				BuildOwner:   "",
				BuildCreator: "creator3",
				BuildStatus:  "failed",
				BuildEnv:     "",
				StartTime:    "", // 测试空时间
				EndTime:      "",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.filter.BuildOwner, func(t *testing.T) {
			result, err := TransEurModeFilterToDbFormat(tt.filter)
			if (err != nil) != tt.expectErr {
				t.Errorf("TransEurModeFilterToDbFormat() error = %v, wantErr %v", err, tt.expectErr)
			}

			if !tt.expectErr {
				var dbFormat map[string]interface{}
				if err := json.Unmarshal(result, &dbFormat); err != nil {
					t.Fatalf("Failed to unmarshal result: %v", err)
				}
			}
		})
	}
}

func TestTransGiteeModeFilterToDbFormat(t *testing.T) {
	tests := []struct {
		eventType string
		filter    CmdToGetSubscribe
		expectErr bool
	}{
		{
			eventType: "issue",
			filter: CmdToGetSubscribe{
				Repos:        "repo1/repo2",
				MyManagement: "management1",
				MySig:        "sig1",
			},
			expectErr: false,
		},
		{
			eventType: "note",
			filter: CmdToGetSubscribe{
				Repos:        "*", // 特殊情况
				MyManagement: "management2",
				MySig:        "sig2",
			},
			expectErr: false,
		},
		{
			eventType: "pr",
			filter: CmdToGetSubscribe{
				IssueState: "open",
				PrCreator:  "creator1",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.eventType, func(t *testing.T) {
			result, err := TransGiteeModeFilterToDbFormat(tt.eventType, tt.filter)
			if (err != nil) != tt.expectErr {
				t.Errorf("TransGiteeModeFilterToDbFormat() error = %v, wantErr %v", err, tt.expectErr)
			}

			if !tt.expectErr {
				var dbFormat map[string]interface{}
				if err := json.Unmarshal(result, &dbFormat); err != nil {
					t.Fatalf("Failed to unmarshal result: %v", err)
				}
			}
		})
	}
}

func TestTransMeetingModeFilterToDbFormat(t *testing.T) {
	tests := []struct {
		filter    CmdToGetSubscribe
		expectErr bool
	}{
		{
			filter: CmdToGetSubscribe{
				MeetingAction:   "start",
				MeetingSigGroup: "group1",
				StartTime:       "1622505600",
				EndTime:         "1622592000",
			},
			expectErr: false,
		},
		{
			filter: CmdToGetSubscribe{
				MeetingAction:   "stop",
				MeetingSigGroup: "group2",
				StartTime:       "1622505600",
				EndTime:         "1622592000",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.filter.MeetingAction, func(t *testing.T) {
			result, err := TransMeetingModeFilterToDbFormat(tt.filter)
			if (err != nil) != tt.expectErr {
				t.Errorf("TransMeetingModeFilterToDbFormat() error = %v, wantErr %v", err, tt.expectErr)
			}

			if !tt.expectErr {
				var dbFormat map[string]interface{}
				if err := json.Unmarshal(result, &dbFormat); err != nil {
					t.Fatalf("Failed to unmarshal result: %v", err)
				}
			}
		})
	}
}

func TestTransCveModeFilterToDbFormat(t *testing.T) {
	tests := []struct {
		filter    CmdToGetSubscribe
		expectErr bool
	}{
		{
			filter: CmdToGetSubscribe{
				CVEComponent: "component1",
				IssueState:   "open",
				CVEAffected:  "affected1",
			},
			expectErr: false,
		},
		{
			filter: CmdToGetSubscribe{
				CVEComponent: "component2",
				IssueState:   "closed",
				CVEAffected:  "affected2",
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.filter.CVEComponent, func(t *testing.T) {
			result, err := TransCveModeFilterToDbFormat(tt.filter)
			if (err != nil) != tt.expectErr {
				t.Errorf("TransCveModeFilterToDbFormat() error = %v, wantErr %v", err, tt.expectErr)
			}

			if !tt.expectErr {
				var dbFormat map[string]interface{}
				if err := json.Unmarshal(result, &dbFormat); err != nil {
					t.Fatalf("Failed to unmarshal result: %v", err)
				}
			}
		})
	}
}
