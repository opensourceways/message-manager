package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/xerrors"
)

// MockMessagePushAdapter 是 MessagePushAdapter 的模拟实现
type MockMessagePushAdapter struct {
	mock.Mock
}

func (m *MockMessagePushAdapter) GetPushConfig(subsIds []string, countPerPage, pageNum int, userName string) ([]MessagePushDTO, error) {
	args := m.Called(subsIds, countPerPage, pageNum, userName)
	if args.Get(0) == nil {
		return []MessagePushDTO{}, args.Error(1)
	}
	return args.Get(0).([]MessagePushDTO), args.Error(1)
}

func (m *MockMessagePushAdapter) AddPushConfig(cmd CmdToAddPushConfig) error {
	args := m.Called(cmd)
	return args.Error(0)
}

func (m *MockMessagePushAdapter) UpdatePushConfig(cmd CmdToUpdatePushConfig) error {
	args := m.Called(cmd)
	return args.Error(0)
}

func (m *MockMessagePushAdapter) RemovePushConfig(cmd CmdToDeletePushConfig) error {
	args := m.Called(cmd)
	return args.Error(0)
}

func TestGetPushConfig(t *testing.T) {
	mockAdapter := new(MockMessagePushAdapter)
	service := NewMessagePushAppService(mockAdapter)

	userName := "testUser"
	subsIds := []string{"sub1", "sub2"}
	countPerPage := 10
	pageNum := 1

	mockData := []MessagePushDTO{
		{ /* 填充必要字段 */ },
	}

	mockAdapter.On("GetPushConfig", subsIds, countPerPage, pageNum, userName).Return(mockData, nil)

	data, err := service.GetPushConfig(countPerPage, pageNum, userName, subsIds)

	assert.NoError(t, err)
	assert.Equal(t, mockData, data)
	mockAdapter.AssertExpectations(t)
}

func TestGetPushConfig_Error(t *testing.T) {
	mockAdapter := new(MockMessagePushAdapter)
	service := NewMessagePushAppService(mockAdapter)

	userName := "testUser"
	subsIds := []string{"sub1", "sub2"}
	countPerPage := 10
	pageNum := 1

	mockAdapter.On("GetPushConfig", subsIds, countPerPage, pageNum, userName).Return(nil, xerrors.New("error"))

	data, err := service.GetPushConfig(countPerPage, pageNum, userName, subsIds)

	assert.Error(t, err)
	assert.Empty(t, data)
	mockAdapter.AssertExpectations(t)
}

func TestAddPushConfig(t *testing.T) {
	mockAdapter := new(MockMessagePushAdapter)
	service := NewMessagePushAppService(mockAdapter)

	cmd := CmdToAddPushConfig{
		SubscribeId:      1,
		RecipientId:      12345,
		NeedMessage:      true,
		NeedPhone:        false,
		NeedMail:         true,
		NeedInnerMessage: false,
	}
	mockAdapter.On("AddPushConfig", cmd).Return(nil)

	err := service.AddPushConfig(&cmd)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestAddPushConfig_Error(t *testing.T) {
	mockAdapter := new(MockMessagePushAdapter)
	service := NewMessagePushAppService(mockAdapter)

	cmd := CmdToAddPushConfig{
		SubscribeId:      1,
		RecipientId:      12345,
		NeedMessage:      true,
		NeedPhone:        false,
		NeedMail:         true,
		NeedInnerMessage: false,
	}
	mockAdapter.On("AddPushConfig", cmd).Return(xerrors.New("error"))

	err := service.AddPushConfig(&cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "add message push config failed")
}

func TestUpdatePushConfig(t *testing.T) {
	mockAdapter := new(MockMessagePushAdapter)
	service := NewMessagePushAppService(mockAdapter)

	cmd := CmdToUpdatePushConfig{
		SubscribeId:      []int{1, 2},
		RecipientId:      "12345",
		NeedMessage:      true,
		NeedPhone:        false,
		NeedMail:         true,
		NeedInnerMessage: false,
	}
	mockAdapter.On("UpdatePushConfig", cmd).Return(nil)

	err := service.UpdatePushConfig(&cmd)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestUpdatePushConfig_Error(t *testing.T) {
	mockAdapter := new(MockMessagePushAdapter)
	service := NewMessagePushAppService(mockAdapter)

	cmd := CmdToUpdatePushConfig{
		SubscribeId:      []int{1, 2},
		RecipientId:      "12345",
		NeedMessage:      true,
		NeedPhone:        false,
		NeedMail:         true,
		NeedInnerMessage: false,
	}
	mockAdapter.On("UpdatePushConfig", cmd).Return(xerrors.New("error"))

	err := service.UpdatePushConfig(&cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update message push config failed")
}

func TestRemovePushConfig(t *testing.T) {
	mockAdapter := new(MockMessagePushAdapter)
	service := NewMessagePushAppService(mockAdapter)

	cmd := CmdToDeletePushConfig{
		SubscribeId: 1,
		RecipientId: 12345,
	}
	mockAdapter.On("RemovePushConfig", cmd).Return(nil)

	err := service.RemovePushConfig(&cmd)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestRemovePushConfig_Error(t *testing.T) {
	mockAdapter := new(MockMessagePushAdapter)
	service := NewMessagePushAppService(mockAdapter)

	cmd := CmdToDeletePushConfig{
		SubscribeId: 1,
		RecipientId: 12345,
	}
	mockAdapter.On("RemovePushConfig", cmd).Return(xerrors.New("error"))

	err := service.RemovePushConfig(&cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "remove message push config failed")
}
