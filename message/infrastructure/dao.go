/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package infrastructure

import (
	"time"

	"gorm.io/datatypes"
)

type MessageListDAO struct {
	Title           string    `gorm:"column:title" json:"title"`
	Summary         string    `gorm:"column:summary" json:"summary"`
	Source          string    `gorm:"column:source" json:"source"`
	Type            string    `gorm:"column:type" json:"type"`
	EventId         string    `gorm:"column:event_id" json:"event_id"`
	DataContentType string    `gorm:"column:data_content_type" json:"data_content_type"`
	DataSchema      string    `gorm:"column:data_schema" json:"data_schema"`
	SpecVersion     string    `gorm:"column:spec_version" json:"spec_version"`
	EventTime       time.Time `gorm:"column:time" json:"time"`
	User            string    `gorm:"column:user" json:"user"`
	SourceUrl       string    `gorm:"column:source_url" json:"source_url"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at" swaggerignore:"true"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at" swaggerignore:"true"`
	IsRead          bool      `gorm:"column:is_read" json:"is_read"`
	SourceGroup     string    `gorm:"column:source_group" json:"source_group"`
}

type MessagePushDAO struct {
	SubscribeId      int       `gorm:"column:subscribe_id" json:"subscribe_id"`
	RecipientId      int64     `gorm:"column:recipient_id" json:"recipient_id"`
	NeedMessage      *bool     `gorm:"column:need_message" json:"need_message"`
	NeedPhone        *bool     `gorm:"column:need_phone" json:"need_phone"`
	NeedMail         *bool     `gorm:"column:need_mail" json:"need_mail"`
	NeedInnerMessage *bool     `gorm:"column:need_inner_message" json:"need_inner_message"`
	IsDeleted        bool      `gorm:"column:is_deleted" json:"is_deleted" swaggerignore:"true"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at" swaggerignore:"true"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at" swaggerignore:"true"`
}

type MessageRecipientDAO struct {
	Id            string    `json:"id,omitempty"`
	Name          string    `gorm:"column:recipient_name" json:"recipient_id"`
	Mail          string    `gorm:"column:mail" json:"mail"`
	Message       string    `gorm:"column:message" json:"message"`
	Phone         string    `gorm:"column:phone" json:"phone"`
	Remark        string    `gorm:"column:remark" json:"remark"`
	UserName      string    `gorm:"column:user_id"  json:"user_id"`
	GiteeUserName string    `gorm:"column:gitee_user_name" json:"gitee_user_name"`
	IsDeleted     bool      `gorm:"column:is_deleted" json:"is_deleted"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at" swaggerignore:"true"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at" swaggerignore:"true"`
}

type MessageSubscribeDAO struct {
	Id          uint           `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
	Source      string         `gorm:"column:source"        json:"source"`
	EventType   string         `gorm:"column:event_type"    json:"event_type"`
	SpecVersion string         `gorm:"column:spec_version"  json:"spec_version"`
	ModeName    string         `gorm:"column:mode_name"     json:"mode_name"`
	ModeFilter  datatypes.JSON `gorm:"column:mode_filter"   json:"mode_filter" swaggerignore:"true"`
	CreatedAt   time.Time      `gorm:"column:created_at"    json:"created_at"  swaggerignore:"true"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"    json:"updated_at"  swaggerignore:"true"`
	UserName    string         `gorm:"column:user_name"     json:"user_name"`
	IsDefault   *bool          `gorm:"column:is_default"    json:"is_default"`
	WebFilter   datatypes.JSON `gorm:"column:web_filter"    json:"web_filter"  swaggerignore:"true"`
}

type MessageSubscribeDAOWithoutId struct {
	Source      string         `gorm:"column:source"        json:"source"`
	EventType   string         `gorm:"column:event_type"    json:"event_type"`
	SpecVersion string         `gorm:"column:spec_version"  json:"spec_version"`
	ModeName    string         `gorm:"column:mode_name"     json:"mode_name"`
	ModeFilter  datatypes.JSON `gorm:"column:mode_filter"   json:"mode_filter" swaggerignore:"true"`
	CreatedAt   time.Time      `gorm:"column:created_at"    json:"created_at"  swaggerignore:"true"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"    json:"updated_at"  swaggerignore:"true"`
	UserName    string         `gorm:"column:user_name"     json:"user_name"`
	IsDefault   *bool          `gorm:"column:is_default"    json:"is_default"`
}

type CountDAO struct {
	Source string `json:"source"`
	Count  int    `json:"count"`
}

type CmdToGetInnerMessageQuick struct {
	Source       string `json:"source"`
	CountPerPage int    `json:"count_per_page"`
	PageNum      int    `json:"page"`
	ModeName     string `json:"mode_name"`
}

type CmdToGetInnerMessage struct {
	Source           string `json:"source"`
	EventType        string `json:"event_type"`
	IsRead           string `json:"is_read"`
	KeyWord          string `json:"key_word"`
	IsBot            string `json:"is_bot"`
	GiteeSigs        string `json:"sig"`
	Repos            string `json:"repos"`
	CountPerPage     int    `json:"count_per_page"`
	PageNum          int    `json:"page"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	MySig            string `json:"my_sig"`
	MyManagement     string `json:"my_management"`
	PrState          string `json:"pr_state"`
	PrCreator        string `json:"pr_creator"`
	PrAssignee       string `json:"pr_assignee"`
	IssueState       string `json:"issue_state"`
	IssueCreator     string `json:"issue_creator"`
	IssueAssignee    string `json:"issue_assignee"`
	NoteType         string `json:"note_type"`
	About            string `json:"about"`
	BuildStatus      string `json:"build_status"`
	BuildOwner       string `json:"build_owner"`
	BuildCreator     string `json:"build_creator"`
	BuildEnv         string `json:"build_env"`
	MeetingAction    string `json:"meeting_action"`
	MeetingSigGroup  string `json:"meeting_sig"`
	MeetingStartTime string `json:"meeting_start_time"`
	MeetingEndTime   string `json:"meeting_end_time"`
	CVEComponent     string `json:"cve_component"`
	CVEState         string `json:"cve_state"`
	CVEAffected      string `json:"cve_affected"`
}

type CmdToSetIsRead struct {
	Source  string `json:"source"`
	EventId string `json:"event_id"`
}

type CmdToAddPushConfig struct {
	SubscribeId      int   `json:"subscribe_id"`
	RecipientId      int64 `json:"recipient_id"`
	NeedMessage      bool  `json:"need_message"`
	NeedPhone        bool  `json:"need_phone"`
	NeedMail         bool  `json:"need_mail"`
	NeedInnerMessage bool  `json:"need_inner_message"`
}

type CmdToUpdatePushConfig struct {
	SubscribeId      []int  `json:"subscribe_id"`
	RecipientId      string `json:"recipient_id"`
	NeedMessage      bool   `json:"need_message"`
	NeedPhone        bool   `json:"need_phone"`
	NeedMail         bool   `json:"need_mail"`
	NeedInnerMessage bool   `json:"need_inner_message"`
}

type CmdToDeletePushConfig struct {
	SubscribeId int   `json:"subscribe_id"`
	RecipientId int64 `json:"recipient_id"`
}

type CmdToAddRecipient struct {
	Name    string `json:"recipient_id"`
	Mail    string `json:"mail"`
	Message string `json:"message"`
	Phone   string `json:"phone"`
	Remark  string `json:"remark"`
}

type CmdToUpdateRecipient struct {
	Id      string `json:"id"`
	Name    string `json:"recipient_id"`
	Mail    string `json:"mail"`
	Message string `json:"message"`
	Phone   string `json:"phone"`
	Remark  string `json:"remark"`
}

type CmdToDeleteRecipient struct {
	RecipientId string `json:"recipient_id"`
}

type CmdToSyncUserInfo struct {
	Mail          string `json:"mail"`
	Phone         string `json:"phone"`
	CountryCode   string `json:"country_code"`
	UserName      string `json:"user_name"`
	GiteeUserName string `json:"gitee_user_name"`
}

type CmdToGetSubscribe struct {
	Source           string `json:"source"`
	EventType        string `json:"event_type"`
	IsRead           string `json:"is_read"`
	KeyWord          string `json:"key_word"`
	IsBot            string `json:"is_bot"`
	GiteeSigs        string `json:"sig"`
	Repos            string `json:"repos"`
	CountPerPage     int    `json:"count_per_page"`
	PageNum          int    `json:"page"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	MySig            string `json:"my_sig"`
	MyManagement     string `json:"my_management"`
	PrState          string `json:"pr_state"`
	PrCreator        string `json:"pr_creator"`
	PrAssignee       string `json:"pr_assignee"`
	IssueState       string `json:"issue_state"`
	IssueCreator     string `json:"issue_creator"`
	IssueAssignee    string `json:"issue_assignee"`
	NoteType         string `json:"note_type"`
	About            string `json:"about"`
	BuildStatus      string `json:"build_status"`
	BuildOwner       string `json:"build_owner"`
	BuildCreator     string `json:"build_creator"`
	BuildEnv         string `json:"build_env"`
	MeetingAction    string `json:"meeting_action"`
	MeetingSigGroup  string `json:"meeting_sig"`
	MeetingStartTime string `json:"meeting_start_time"`
	MeetingEndTime   string `json:"meeting_end_time"`
	CVEComponent     string `json:"cve_component"`
	CVEState         string `json:"cve_state"`
	CVEAffected      string `json:"cve_affected"`
	SpecVersion      string `json:"spec_version"`
	ModeName         string `json:"mode_name"`
}

type CmdToAddSubscribe struct {
	Source      string         `json:"source"`
	EventType   string         `json:"event_type"`
	SpecVersion string         `json:"spec_version"`
	ModeFilter  datatypes.JSON `json:"mode_filter" swaggerignore:"true"`
	ModeName    string         `json:"mode_name"`
}

type CmdToUpdateSubscribe struct {
	Source  string `json:"source"`
	OldName string `json:"old_name"`
	NewName string `json:"new_name"`
}

type CmdToDeleteSubscribe struct {
	Source    string `json:"source"`
	EventType string `json:"event_type"`
	ModeName  string `json:"mode_name"`
}
