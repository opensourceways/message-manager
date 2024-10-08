// main/main_test.go
package main

import (
	"flag"
	"os"
	"testing"

	"github.com/opensourceways/message-manager/config"
	"github.com/sirupsen/logrus"
)

func TestGatherOptions(t *testing.T) {
	// 测试用例
	tests := []struct {
		args         []string
		expectedFile string
		expectError  error
	}{
		{[]string{"-config-file=config.yaml"}, "config.yaml", nil},
		{[]string{"-invalid-flag"}, "", nil},
		{[]string{}, "", nil}, // 没有提供参数
	}

	for _, test := range tests {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		options, err := gatherOptions(fs, test.args...)

		if nil != test.expectError {
			t.Errorf("expected error: %v, got: %v", test.expectError, err)
		}
		if options.Config != test.expectedFile {
			t.Errorf("expected config file: %s, got: %s", test.expectedFile, options.Config)
		}
	}
}

func TestLoadConfig(t *testing.T) {
	// 创建一个临时配置文件
	tempFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		logrus.Fatalf("failed to create temp file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			return
		}
	}(tempFile.Name())

	_, err = tempFile.WriteString("postgresql:\n  host: localhost\n  port: 5432\n")
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	err = tempFile.Close()
	if err != nil {
		return
	}

	cfg := new(config.Config)

	err = config.LoadConfig(tempFile.Name(), cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.Postgresql.Host != "localhost" || cfg.Postgresql.Port != 5432 {
		t.Errorf("expected host: localhost, port: 5432, got host: %s, port: %d", cfg.Postgresql.Host, cfg.Postgresql.Port)
	}
}
