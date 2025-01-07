/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package cassandra

import (
	"github.com/gocql/gocql"
	"golang.org/x/xerrors"
)

var (
	session *gocql.Session
)

func Init(cfg *Config) error {

	cluster := gocql.NewCluster(cfg.Host)
	cluster.Keyspace = cfg.Name
	cluster.Port = cfg.Port
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.User,
		Password: cfg.Password,
	}
	sessionInstance, err := cluster.CreateSession()
	if err != nil {
		return xerrors.Errorf("create session failed, err:%v", err)
	}

	session = sessionInstance
	return nil
}

func Session() *gocql.Session {
	return session
}
