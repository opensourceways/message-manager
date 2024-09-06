package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessagePushAdapter struct {
	mock.Mock
}

func (m *MockMessagePushAdapter) GetPushConfig(subsIds []string, countPerPage, pageNum int,
	userName string) ([]MessagePushDTO, error) {
	args := m.Called(subsIds, countPerPage, pageNum, userName)
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
	service := messagePushAppService{messagePushAdapter: mockAdapter}
	countPerPage, pageNum := 10, 0
	userName := "testUser"
	var subsIds []string

	expectedData := []MessagePushDTO{{}}

	mockAdapter.On("GetPushConfig", subsIds, countPerPage, pageNum, userName).Return(expectedData, nil)

	data, err := service.GetPushConfig(countPerPage, pageNum, userName, subsIds)

	assert.NoError(t, err)
	assert.Equal(t, data, expectedData)

	mockAdapter.AssertExpectations(t)
}

func TestAddPushConfig(t *testing.T) {
	mockAdapter := new(MockMessagePushAdapter)
	service := messagePushAppService{messagePushAdapter: mockAdapter}
	cmd := &CmdToAddPushConfig{}

	mockAdapter.On("AddPushConfig", *cmd).Return(nil)

	err := service.AddPushConfig(cmd)

	assert.NoError(t, err)

	mockAdapter.AssertExpectations(t)
}

func TestUpdatePushConfig(t *testing.T) {
	mockAdapter := new(MockMessagePushAdapter)
	service := messagePushAppService{messagePushAdapter: mockAdapter}
	cmd := &CmdToUpdatePushConfig{}

	mockAdapter.On("UpdatePushConfig", *cmd).Return(nil)

	err := service.UpdatePushConfig(cmd)

	assert.NoError(t, err)

	mockAdapter.AssertExpectations(t)
}

func TestRemoveConfig(t *testing.T) {
	mockAdapter := new(MockMessagePushAdapter)
	service := messagePushAppService{messagePushAdapter: mockAdapter}
	cmd := &CmdToDeletePushConfig{}

	mockAdapter.On("RemovePushConfig", *cmd).Return(nil)

	err := service.RemovePushConfig(cmd)

	assert.NoError(t, err)

	mockAdapter.AssertExpectations(t)
}
