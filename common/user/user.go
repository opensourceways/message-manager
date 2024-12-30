/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

const OneIdUserCookie = "_Y_G_"

type ManagerTokenRequest struct {
	GrantType string `json:"grant_type"`
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type ManagerTokenResponse struct {
	ManagerToken string `json:"token"`
}

type GetUserInfoResponse struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data `json:"data"`
}

type Data struct {
	UserName string `json:"username"`
	Phone    string `json:"phone"`
	NickName string `json:"nickname"`
}

func JsonMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(t); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func getManagerToken(appId string, appSecret string) (string, error) {
	url := fmt.Sprintf("%s/oneid/manager/token", config.AuthorHost)
	reqBody := ManagerTokenRequest{
		GrantType: "token",
		AppId:     appId,
		AppSecret: appSecret,
	}
	v, err := JsonMarshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(v))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var data ManagerTokenResponse
	if err = json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	return data.ManagerToken, nil
}

func GetEulerUserName(ctx *gin.Context) (string, error) {
	token := ctx.Request.Header.Get("token")
	YGCookie, err := extractYGCookie(ctx.Request.Header.Get("Cookie"))
	if err != nil {
		return "", err
	}
	managerToken, err := getManagerToken(config.EulerAppId, config.EulerAppSecret)
	if err != nil {
		logrus.Errorf("get manager token failed, err:%v", err)
		return "", err
	}
	userName, err := fetchUserName(managerToken, token, YGCookie)
	if err != nil {
		logrus.Errorf("get user name failed, err:%v", err)
		return "", err
	}
	return userName, nil
}

func extractYGCookie(cookieHeader string) (string, error) {
	re := regexp.MustCompile(`_Y_G_=(.*?)(?:;|$)`)
	if re.MatchString(cookieHeader) {
		match := re.FindStringSubmatch(cookieHeader)
		if len(match) > 1 {
			return match[1], nil
		}
	}
	return "", xerrors.Errorf("YG cookie not found")
}

func fetchUserName(managerToken, userToken, YGCookie string) (string, error) {
	url := fmt.Sprintf("%s/oneid/manager/personal/center/user?community=%s", config.AuthorHost, config.EulerCommunity)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("token", managerToken)
	req.Header.Add("user-token", userToken)
	req.Header.Add("Cookie", OneIdUserCookie+"="+YGCookie)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var data GetUserInfoResponse
	if err = json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	if data.UserName == "" {
		return "", xerrors.Errorf("the user name is null")
	}
	return data.UserName, nil
}
