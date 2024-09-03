/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package cassandra

type Config struct {
	Host string `json:"host"     required:"true"`
	User string `json:"user"     required:"true"`
	Pwd  string `json:"pwd"      required:"true"`
	Name string `json:"name"     required:"true"`
	Port int    `json:"port"     required:"true"`
}
