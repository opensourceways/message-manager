/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package postgresql

import (
	"database/sql"
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	sqlDb *sql.DB
	db    *gorm.DB
)

func Init(cfg *Config) (err error) {
	db, err = gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  cfg.dsn(),
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{},
	)
	if err != nil {
		return
	}

	if sqlDb, err = db.DB(); err != nil {
		return errors.New("db error")
	}

	sqlDb.SetConnMaxLifetime(cfg.getLifeDuration())
	sqlDb.SetMaxOpenConns(cfg.MaxConn)
	sqlDb.SetMaxIdleConns(cfg.MaxIdle)

	return
}

func DB() *gorm.DB {
	return db
}
