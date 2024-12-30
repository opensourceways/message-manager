/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package user

type Config struct {
	AuthorHost     string `json:"author_host"       required:"true"`
	EulerCommunity string `json:"euler_community" required:"true"`
	EulerAppId     string `json:"euler_app_id" required:"true"`
	EulerAppSecret string `json:"euler_app_secret" required:"true"`
}

var config Config

func Init(cfg *Config) {
	config = Config{
		AuthorHost:     cfg.AuthorHost,
		EulerCommunity: cfg.EulerCommunity,
		EulerAppId:     cfg.EulerAppId,
		EulerAppSecret: cfg.EulerAppSecret,
	}
}
