/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package user

type Config struct {
	AuthorHost string `json:"author_host"       required:"true"`
	Community  string `json:"community" required:"true"`
	AppId      string `json:"app_id" required:"true"`
	AppSecret  string `json:"app_secret" required:"true"`
}

var config Config

func Init(cfg *Config) {
	config = Config{
		AuthorHost: cfg.AuthorHost,
		Community:  cfg.Community,
		AppId:      cfg.AppId,
		AppSecret:  cfg.AppSecret,
	}
}
