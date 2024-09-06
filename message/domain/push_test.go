package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessagePushAdapter struct {
	mock.Mock
}

func (m *MockMessagePushAdapter) GetPushConfig(subsIds []string, countPerPage, pageNum int, userName string) ([]MessagePushDO, error) {
	args := m.Called(subsIds, countPerPage, pageNum, userName)
	return args.Get(0).([]MessagePushDO), args.Error(1)
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
	subsIds := []string{"sub1", "sub2"}
	countPerPage := 10
	pageNum := 1
	userName := "testUser"
	expectedMessages := []MessagePushDO{{}}

	mockAdapter.On("GetPushConfig", subsIds, countPerPage, pageNum, userName).Return(expectedMessages, nil)

	messages, err := mockAdapter.GetPushConfig(subsIds, countPerPage, pageNum, userName)

	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)
	mockAdapter.AssertExpectations(t)
}

func TestAddPushConfig(t *testing.T) {
	mockAdapter := new(MockMessagePushAdapter)
	cmd := CmdToAddPushConfig{}

	mockAdapter.On("AddPushConfig", cmd).Return(nil)

	err := mockAdapter.AddPushConfig(cmd)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestUpdatePushConfig(t *testing.T) {
	mockAdapter := new(MockMessagePushAdapter)
	cmd := CmdToUpdatePushConfig{}

	mockAdapter.On("UpdatePushConfig", cmd).Return(nil)

	err := mockAdapter.UpdatePushConfig(cmd)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestRemovePushConfig(t *testing.T) {
	mockAdapter := new(MockMessagePushAdapter)
	cmd := CmdToDeletePushConfig{}

	mockAdapter.On("RemovePushConfig", cmd).Return(nil)

	err := mockAdapter.RemovePushConfig(cmd)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}
