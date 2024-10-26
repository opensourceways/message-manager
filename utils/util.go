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
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

const (
	EurSource     = "https://eur.openeuler.openatom.cn"
	GiteeSource   = "https://gitee.com"
	MeetingSource = "https://www.openEuler.org/meeting"
	CveSource     = "cve"

	eulerUserSigUrl   = "https://dsapi.osinfra.cn/query/user/ownertype?community=openeuler&user=%s"
	giteeUserReposUrl = "https://gitee.com/api/v5/users/%s/repos?type=all&sort=full_name&page=%d" +
		"&per_page=%d"
	giteeGetPullsUrl = "https://gitee.com/api/v5/repos/%s/%s/pulls?&state=open" +
		"&sort=created&direction=desc&page=%d&per_page=%d&assignee=%s"
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
	url := fmt.Sprintf(eulerUserSigUrl, userName)
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
	repos, err := getRepos(giteeUserReposUrl, userName)
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

type PullRequest struct {
	PullRequestUrl string `json:"url"`
}

func getPulls(url, owner, repoName, username string) ([]PullRequest, error) {
	var pulls []PullRequest
	page := 1
	perPage := 100

	var totalCount int

	for {
		url := fmt.Sprintf(url, owner, repoName, page, perPage, username)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return []PullRequest{}, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return []PullRequest{}, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []PullRequest{}, err
		}

		var members []PullRequest
		err = json.Unmarshal(body, &members)
		logrus.Infof("the data is %v", string(body))
		if err != nil {
			return []PullRequest{}, err
		}

		pulls = append(pulls, members...)

		if totalCount == 0 {
			totalCount, err = strconv.Atoi(resp.Header.Get("total_count"))
			if err != nil {
				return []PullRequest{}, xerrors.Errorf("trans to int failed, err:%v", err)
			}
		}

		if len(members) < perPage {
			break
		}
		page++
		err = resp.Body.Close()
		if err != nil {
			return []PullRequest{}, xerrors.Errorf("close body failed, err :%v", err)
		}
	}
	return pulls, nil
}

func GetTodoPulls(userName string) ([]PullRequest, error) {
	repos, err := GetUserAdminReposByUsername(userName)
	if err != nil {
		return nil, err
	}

	// Step 2: 获取待评审的 PR 列表
	var PullRequests []PullRequest
	var wg sync.WaitGroup

	for _, repo := range repos {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()
			lRepo := strings.Split(repo, "/")
			owner, repoName := lRepo[0], lRepo[1]
			pulls, err := getPulls(giteeGetPullsUrl, owner, repoName, userName)
			if err != nil {
				return
			}
			PullRequests = append(PullRequests, pulls...)
		}(repo)
	}
	wg.Wait()
	return PullRequests, nil
}
