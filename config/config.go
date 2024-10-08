/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package config

import (
	"fmt"
	"os"

	"sigs.k8s.io/yaml"

	"github.com/opensourceways/message-manager/common/cassandra"
	common "github.com/opensourceways/message-manager/common/config"
	"github.com/opensourceways/message-manager/common/postgresql"
	"github.com/opensourceways/message-manager/common/user"
)

type Config struct {
	Postgresql postgresql.Config `yaml:"postgresql"`
	Cassandra  cassandra.Config  `yaml:"cassandra"`
	User       user.Config       `yaml:"user"`
}

func LoadFromYaml(path string, cfg interface{}) error {
	b, err := os.ReadFile(path) // #nosec G304
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, cfg)
}

func LoadConfig(path string, cfg *Config) error {
	if err := LoadFromYaml(path, cfg); err != nil {
		return fmt.Errorf("load from yaml failed, %w", err)
	}

	common.SetDefault(cfg)

	return common.Validate(cfg)
}
