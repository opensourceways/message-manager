package allerror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// 测试创建新错误
	err := New("custom_error", "This is a custom error")
	assert.Equal(t, "This is a custom error", err.Error())
	assert.Equal(t, "custom_error", err.ErrorCode())

	// 测试创建新错误时消息为空
	err = New("another_error", "")
	assert.Equal(t, "another error", err.Error()) // 下划线替换为空格
}

func TestNewNotFound(t *testing.T) {
	err := NewNotFound(ErrorCodeRepoNotFound, "Repository not found")
	assert.Equal(t, "Repository not found", err.Error())
	assert.Equal(t, ErrorCodeRepoNotFound, err.ErrorCode())
}

func TestIsNotFound(t *testing.T) {
	err := NewNotFound(ErrorCodeRepoNotFound, "Repository not found")
	assert.True(t, IsNotFound(err))

	normalErr := errors.New("some other error")
	assert.False(t, IsNotFound(normalErr))
}

func TestNewNoPermission(t *testing.T) {
	err := NewNoPermission("Access denied")
	assert.Equal(t, "Access denied", err.Error())
	assert.Equal(t, errorCodeNoPermission, err.ErrorCode())
}

func TestIsNoPermission(t *testing.T) {
	err := NewNoPermission("Access denied")
	assert.True(t, IsNoPermission(err))

	normalErr := errors.New("some other error")
	assert.False(t, IsNoPermission(normalErr))
}

func TestNewInvalidParam(t *testing.T) {
	err := NewInvalidParam("Invalid parameter provided")
	assert.Equal(t, "Invalid parameter provided", err.Error())
	assert.Equal(t, errorCodeInvalidParam, err.ErrorCode())
}

func TestNewOverLimit(t *testing.T) {
	err := NewOverLimit("rate_limit_exceeded", "Rate limit exceeded")
	assert.Equal(t, "Rate limit exceeded", err.Error())
	assert.Equal(t, "rate_limit_exceeded", err.ErrorCode())
}

func TestIsErrorCodeEmptyRepo(t *testing.T) {
	err := New(ErrorCodeEmptyRepo, "The repository is empty")
	assert.True(t, IsErrorCodeEmptyRepo(err))

	err = New("some_other_error", "Some other error occurred")
	assert.False(t, IsErrorCodeEmptyRepo(err))
}
