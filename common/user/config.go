/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package user

type Config struct {
	AuthorHost string `json:"author_host"       required:"true"`
	Community  string `json:"euler_community" required:"true"`
	AppId      string `json:"euler_app_id" required:"true"`
	AppSecret  string `json:"euler_app_secret" required:"true"`
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
