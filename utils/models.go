/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package utils

type SigData struct {
	Sig  []string `json:"sig"`
	Type []string `json:"type"`
}

type SigInfo struct {
	SigData `json:"data"`
}

type GiteeRepo struct {
	FullName   string `json:"full_name"`
	Permission `json:"permission"`
}

type Permission struct {
	Pull  bool `json:"pull"`
	Push  bool `json:"push"`
	Admin bool `json:"admin"`
}

type EurDbFormat struct {
	Status    string `json:"Body.Status,omitempty"`
	Owner     string `json:"Body.Owner,omitempty"`
	User      string `json:"Body.User,omitempty"`
	Env       string `json:"Body.Chroot,omitempty"`
	EventTime string `json:"EventTime,omitempty"`
}

type GiteeIssueDbFormat struct {
	RepoName        string `json:"IssueEvent.Repository.FullName,omitempty"`
	IsBot           string `json:"IssueEvent.Sender.Name,omitempty"`
	Namespace       string `json:"IssueEvent.Repository.Namespace,omitempty"`
	SigGroupName    string `json:"SigGroupName,omitempty"`
	IssueState      string `json:"IssueEvent.Issue.State,omitempty"`
	IssueCreator    string `json:"IssueEvent.Issue.User.Login,omitempty"`
	IssueAssignee   string `json:"IssueEvent.Issue.Assignee.Login,omitempty"`
	EventTime       string `json:"EventTime,omitempty"`
	MyManagement    string `json:"RepoAdmins,omitempty"`
	OtherManagement string `json:"RepoAdmins,omitempty"`
	MySig           string `json:"SigMaintainers,omitempty"`
	OtherSig        string `json:"SigMaintainers,omitempty"`
}
type GiteeNoteDbFormat struct {
	RepoName        string `json:"NoteEvent.Repository.FullName,omitempty"`
	IsBot           string `json:"NoteEvent.Sender.Name,omitempty"`
	Namespace       string `json:"NoteEvent.Repository.Namespace,omitempty"`
	SigGroupName    string `json:"SigGroupName,omitempty"`
	NoteType        string `json:"NoteEvent.NoteableType,omitempty"`
	About           string `json:"NoteEvent.Comment.Body,omitempty"`
	EventTime       string `json:"EventTime,omitempty"`
	MyManagement    string `json:"RepoAdmins,omitempty"`
	OtherManagement string `json:"RepoAdmins,omitempty"`
	MySig           string `json:"SigMaintainers,omitempty"`
	OtherSig        string `json:"SigMaintainers,omitempty"`
}
type GiteePullRequestDbFormat struct {
	RepoName        string `json:"PullRequestEvent.Repository.FullName,omitempty"`
	IsBot           string `json:"PullRequestEvent.Sender.Name,omitempty"`
	Namespace       string `json:"PullRequestEvent.Repository.Namespace,omitempty"`
	SigGroupName    string `json:"SigGroupName,omitempty"`
	PrState         string `json:"PullRequestEvent.PullRequest.State,omitempty"`
	PrAction        string `json:"PullRequestEvent.Action,omitempty"`
	PrMergeStatus   string `json:"PullRequestEvent.MergeStatus,omitempty"`
	PrCreator       string `json:"PullRequestEvent.PullRequest.User.Login,omitempty"`
	PrAssignee      string `json:"PullRequestEvent.PullRequest.Assignee.Login,omitempty"`
	EventTime       string `json:"EventTime,omitempty"`
	MyManagement    string `json:"RepoAdmins,omitempty"`
	OtherManagement string `json:"RepoAdmins,omitempty"`
	MySig           string `json:"SigMaintainers,omitempty"`
	OtherSig        string `json:"SigMaintainers,omitempty"`
}

type GiteeNoTypeDbFormat struct {
}

type MeetingDbFormat struct {
	Action   string `json:"Action,omitempty"`
	SigGroup string `json:"Msg.GroupName,omitempty"`
	Date     string `json:"Msg.Date,omitempty"`
	MySig    string `json:"SigMaintainers,omitempty"`
	OtherSig string `json:"SigMaintainers,omitempty"`
}

type CveDbFormat struct {
	Component    string `json:"CVEComponent,omitempty"`
	State        string `json:"IssueEvent.Issue.State,omitempty"`
	SigGroupName string `json:"SigGroupName,omitempty"`
	Affected     string `json:"CVEAffectVersion,omitempty"`
	MySig        string `json:"SigMaintainers,omitempty"`
	OtherSig     string `json:"SigMaintainers,omitempty"`
}
