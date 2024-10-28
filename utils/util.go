/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

const (
	EurSource     = "https://eur.openeuler.openatom.cn"
	GiteeSource   = "https://gitee.com"
	MeetingSource = "https://www.openEuler.org/meeting"
	CveSource     = "cve"
)

func ParseUnixTimestamp(timestampStr string) *time.Time {
	if timestampStr == "" {
		return nil
	}
	// 解析字符串为整数
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return nil
	}
	// 将毫秒转换为秒和纳秒
	seconds := timestamp / 1000
	nanoseconds := (timestamp % 1000) * 1000000
	t := time.Unix(seconds, nanoseconds)
	utcPlus8 := t.Add(8 * time.Hour)
	return &utcPlus8
}

func IsEurMessage(source string) bool {
	return source == EurSource
}

func IsGiteeMessage(source string) bool {
	return source == GiteeSource
}

func IsMeetingMessage(source string) bool {
	return source == MeetingSource
}

func IsCveMessage(source string) bool {
	return source == CveSource
}

func sortStringList(strList []string) []string {
	// 定义一个比较函数,用于排序
	less := func(i, j int) bool {
		// 如果一个字符串包含 "*"，另一个不包含，则把包含 "*" 的排在前面
		if strings.Contains(strList[i], "*") && !strings.Contains(strList[j], "*") {
			return true
		}
		if !strings.Contains(strList[i], "*") && strings.Contains(strList[j], "*") {
			return false
		}
		// 如果两个字符串都包含或都不包含 "*"，则按原来的顺序排列
		return strList[i] < strList[j]
	}

	// 使用自定义的比较函数进行排序
	sort.Slice(strList, less)
	return strList
}

func MergePaths(paths []string) []string {
	pathDict := make(map[string]bool)
	var result []string
	sortPaths := sortStringList(paths)

	// 遍历路径列表
	for _, p := range sortPaths {
		if p == "" {
			continue
		}

		if p == "*" {
			result = []string{"*"}
			break
		}

		lp := strings.Split(p, "/")
		if strings.Contains(p, "*") {
			result = append(result, p)
			pathDict[lp[0]] = true
		} else {
			if pathDict[lp[0]] {
				continue
			} else {
				result = append(result, p)
			}
		}
	}
	return result
}

func RemoveEmptyStrings(input []string) []string {
	var result []string
	for _, str := range input {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}

func GetUserSigInfo(userName string) ([]string, error) {
	url := fmt.Sprintf(config.EulerUserSigUrl, userName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []string{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	var repoSig SigInfo
	err = json.Unmarshal(body, &repoSig)
	if err != nil {
		return []string{}, err
	}
	if repoSig.Sig == nil {
		return []string{}, nil // 确保返回空切片而不是 nil
	}
	return repoSig.Sig, nil
}

func getRepos(url string, param string) ([]GiteeRepo, error) {
	var repos []GiteeRepo
	page := 1
	perPage := 100

	var totalCount int
	for {
		curUrl := fmt.Sprintf(url, param, page, perPage)
		req, err := http.NewRequest("GET", curUrl, nil)
		if err != nil {
			return []GiteeRepo{}, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return []GiteeRepo{}, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []GiteeRepo{}, err
		}

		var members []GiteeRepo
		err = json.Unmarshal(body, &members)
		if err != nil {
			return []GiteeRepo{}, err
		}

		repos = append(repos, members...)

		if totalCount == 0 {
			totalCount, err = strconv.Atoi(resp.Header.Get("total_count"))
			if err != nil {
				return []GiteeRepo{}, xerrors.Errorf("trans to int failed, err:%v", err)
			}
		}

		if len(members) < perPage {
			break
		}
		page++
		err = resp.Body.Close()
		if err != nil {
			return []GiteeRepo{}, xerrors.Errorf("close body failed, err :%v", err)
		}
	}
	return repos, nil
}

func GetUserAdminReposByUsername(userName string) ([]string, error) {
	repos, err := getRepos(config.GiteeUserReposUrl, userName)
	if err != nil {
		return []string{}, err
	}
	var adminRepos []string
	for _, repo := range repos {
		if repo.Admin {
			adminRepos = append(adminRepos, repo.FullName)
		}
	}
	if len(adminRepos) == 0 {
		return []string{}, nil // 确保返回空切片而不是 nil
	}
	return adminRepos, nil
}

type UserInfo struct {
	ID int `json:"id"`
}

func GetUserId(userName string) int {
	curUrl := fmt.Sprintf(config.GiteeGetUserIdUrl, userName, config.GiteeToken)
	req, err := http.NewRequest("GET", curUrl, nil)
	if err != nil {
		return 0
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0
	}

	var giteeUserInfo UserInfo
	err = json.Unmarshal(body, &giteeUserInfo)
	if err != nil {
		return 0
	}

	return giteeUserInfo.ID
}

type PullsData struct {
	TotalCount   int           `json:"total_count"`
	PullRequests []PullRequest `json:"data"`
}

type PullRequest struct {
	IId     int         `json:"iid"`
	Title   string      `json:"title"`
	Project PullProject `json:"project"`
}

type PullProject struct {
	PathWithNamespace string `json:"path_with_namespace"`
}

func getPulls(giteeId int) ([]PullRequest, int, error) {
	var pulls []PullRequest
	page := 1
	perPage := 100
	if giteeId == 0 {
		return []PullRequest{}, 0, xerrors.Errorf("the gitee id is empty")
	}
	var totalCount int
	state := "opened"
	for {
		curUrl := fmt.Sprintf(config.GiteeGetPullsUrl+"&state=%s", config.OpenEulerId,
			config.OpenEulerToken, giteeId, page, perPage, state)
		req, err := http.NewRequest("GET", curUrl, nil)
		if err != nil {
			return []PullRequest{}, 0, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return []PullRequest{}, 0, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []PullRequest{}, 0, err
		}

		var data PullsData
		err = json.Unmarshal(body, &data)
		if err != nil {
			return []PullRequest{}, 0, err
		}

		totalCount = data.TotalCount
		pulls = append(pulls, data.PullRequests...)
		logrus.Infof("the url is %v, the length is %v", curUrl, len(data.PullRequests))

		if len(data.PullRequests) < perPage {
			break
		}
		page++
		err = resp.Body.Close()
		if err != nil {
			return []PullRequest{}, 0, xerrors.Errorf("close body failed, err :%v", err)
		}
	}
	return pulls, totalCount, nil
}

func GetTodoPulls(userName string) ([]string, int, error) {
	giteeUserId := GetUserId(userName)
	PullRequests, totalCount, err := getPulls(giteeUserId)

	if err != nil {
		return nil, 0, err
	}
	var pullUrls []string
	for _, pr := range PullRequests {
		url := fmt.Sprintf("https://gitee.com/%s/pulls/%d", pr.Project.PathWithNamespace, pr.IId)
		pullUrls = append(pullUrls, url)
	}

	return pullUrls, totalCount, nil
}
