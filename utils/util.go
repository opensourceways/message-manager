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
)

const (
	EurSource     = "https://eur.openeuler.openatom.cn"
	GiteeSource   = "https://gitee.com"
	MeetingSource = "https://www.openEuler.org/meeting"
	CveSource     = "cve"

	eulerUserSigUrl   = "https://dsapi.osinfra.cn/query/user/ownertype?community=openeuler&user=%s"
	giteeUserReposUrl = "https://gitee.com/api/v5/users/%s/repos?type=all&sort=full_name&page=%d&per_page=%d"
)

func ParseUnixTimestampNew(timestampStr string) *string {
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
	t := time.Unix(seconds, nanoseconds).UTC() // 转换为 UTC 时区

	// 格式化时间为 PostgreSQL 可接受的格式
	formattedTime := t.Format("2006-01-02 15:04:05.999999999 -0700")
	return &formattedTime
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
func GetUserAdminRepos(userName string) ([]string, error) {
	var repos []GiteeRepo
	page := 1
	perPage := 100
	for {
		members, err := fetchUserRepos(userName, page, perPage)
		if err != nil {
			return nil, err // 直接返回 nil 而不是空切片
		}
		repos = append(repos, members...)
		// 如果返回的成员少于每页数量，则表示没有更多数据
		if len(members) < perPage {
			break
		}
		page++
	}
	return filterAdminRepos(repos), nil
}

func fetchUserRepos(userName string, page, perPage int) ([]GiteeRepo, error) {
	url := fmt.Sprintf(giteeUserReposUrl, userName, page, perPage)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // 确保在函数结束时关闭响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var members []GiteeRepo
	err = json.Unmarshal(body, &members)
	if err != nil {
		return nil, err
	}
	return members, nil
}

func filterAdminRepos(repos []GiteeRepo) []string {
	var adminRepos []string
	for _, repo := range repos {
		if repo.Admin {
			adminRepos = append(adminRepos, repo.FullName)
		}
	}
	return adminRepos
}
