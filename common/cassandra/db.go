package cassandra

import (
	"github.com/gocql/gocql"
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
		Password: cfg.Pwd,
	}
	sessionInstance, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}

	session = sessionInstance
	return nil
}

func Session() *gocql.Session {
	return session
}
