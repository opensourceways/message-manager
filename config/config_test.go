// config/config_test.go
package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Success(t *testing.T) {
	// 创建一个临时 YAML 文件
	tempFile, err := os.CreateTemp("", "test_config.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name()) // 清理临时文件

	// 写入测试配置
	_, err = tempFile.WriteString(`
postgresql:
  user: test_user
  pwd: test_password
  name: test_db
  port: 2345
cassandra:
  host: localhost
  port: 1234
`)
	assert.NoError(t, err)

	// 读取配置
	var cfg Config
	err = LoadConfig(tempFile.Name(), &cfg)
	assert.NoError(t, err)

	// 验证配置内容
	assert.Equal(t, "test_user", cfg.Postgresql.User)
	assert.Equal(t, "test_password", cfg.Postgresql.Pwd)
	assert.Equal(t, "test_db", cfg.Postgresql.Name)
	assert.Equal(t, 2345, cfg.Postgresql.Port)
	assert.Equal(t, "localhost", cfg.Cassandra.Host)
	assert.Equal(t, 1234, cfg.Cassandra.Port)
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	var cfg Config
	err := LoadConfig("non_existent_file.yaml", &cfg)
	assert.Error(t, err) // 应该返回错误
}

func TestLoadConfig_InvalidYaml(t *testing.T) {
	// 创建一个临时 YAML 文件
	tempFile, err := os.CreateTemp("", "invalid_config.yaml")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// 写入无效的 YAML 内容
	_, err = tempFile.WriteString(`invalid_yaml: [}`)
	assert.NoError(t, err)

	// 尝试加载配置
	var cfg Config
	err = LoadConfig(tempFile.Name(), &cfg)
	assert.Error(t, err) // 应该返回错误
}
