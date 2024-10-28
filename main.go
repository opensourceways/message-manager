/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/message-manager/common/postgresql"
	"github.com/opensourceways/message-manager/common/user"
	"github.com/opensourceways/message-manager/config"
	"github.com/opensourceways/message-manager/server"
	"github.com/opensourceways/message-manager/utils"
)

func gatherOptions(fs *flag.FlagSet, args ...string) (Options, error) {
	var o Options
	o.AddFlags(fs)
	err := fs.Parse(args)

	return o, err
}

type Options struct {
	Config string
}

func (o *Options) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.Config, "config-file", "", "Path to config file.")
}

// @title           Message Manager
// @version         1.0
// @description     This is a Message Manager Server.
func main() {
	o, err := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err != nil {
		logrus.Fatalf("new Options failed, err:%s", err.Error())
		return
	}

	// cfg
	cfg := new(config.Config)

	if err = config.LoadConfig(o.Config, cfg); err != nil {
		logrus.Errorf("load config, err:%s", err.Error())

		return
	}

	// init postgresql
	if err := postgresql.Init(&cfg.Postgresql); err != nil {
		fmt.Println("Postgresql数据库初始化失败, err:", err)
		return
	}

	//if err := cassandra.Init(&cfg.Cassandra); err != nil {
	//	fmt.Println("Cassandra数据库初始化失败")
	//}

	// init user
	user.Init(&cfg.User)

	// init user
	utils.Init(&cfg.Utils)

	server.StartWebServer()
}
