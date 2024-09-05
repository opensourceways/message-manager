package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageRecipientAdapter struct {
	mock.Mock
}

func (m *MockMessageRecipientAdapter) GetRecipientConfig(countPerPage,
	pageNum int, userName string) ([]MessageRecipientDTO, int64, error) {
	args := m.Called(countPerPage, pageNum, userName)
	return args.Get(0).([]MessageRecipientDTO), args.Get(1).(int64), args.Error(2)
}

func (m *MockMessageRecipientAdapter) AddRecipientConfig(cmd CmdToAddRecipient,
	userName string) error {
	m.Called(cmd, userName)
	return nil
}

func (m *MockMessageRecipientAdapter) UpdateRecipientConfig(cmd CmdToUpdateRecipient,
	userName string) error {
	args := m.Called(cmd, userName)
	return args.Error(0)
}

func (m *MockMessageRecipientAdapter) RemoveRecipientConfig(cmd CmdToDeleteRecipient, userName string) error {
	args := m.Called(cmd, userName)
	return args.Error(0)
}

func (m *MockMessageRecipientAdapter) SyncUserInfo(cmd CmdToSyncUserInfo) (uint, error) {
	args := m.Called(cmd)
	return args.Get(0).(uint), args.Error(1)
}

func TestGetRecipientConfig(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := messageRecipientAppService{messageRecipientAdapter: mockAdapter}
	countPerPage, pageNum := 10, 0
	userName := "testUser"

	expectedRecipients := []MessageRecipientDTO{{}}
	expectedCount := int64(0)

	mockAdapter.On("GetRecipientConfig", countPerPage, pageNum, userName).
		Return(expectedRecipients, expectedCount, nil)

	messages, count, err := service.GetRecipientConfig(countPerPage, pageNum, userName)

	assert.NoError(t, err)
	assert.Equal(t, expectedRecipients, messages)
	assert.Equal(t, expectedCount, count)

	mockAdapter.AssertExpectations(t)
}

func TestAddRecipientConfig(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := messageRecipientAppService{messageRecipientAdapter: mockAdapter}
	cmd := &CmdToAddRecipient{}
	userName := "testUser"

	err := service.AddRecipientConfig(userName, cmd)
	assert.EqualError(t, err, "the recipient is null")
}

func TestUpdateRecipientConfig(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := messageRecipientAppService{messageRecipientAdapter: mockAdapter}
	cmd := &CmdToUpdateRecipient{}
	userName := "testUser"

	err := service.UpdateRecipientConfig(userName, cmd)

	assert.ErrorContains(t, err, "data is invalid")
}

func TestRemoveRecipientConfig(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := messageRecipientAppService{messageRecipientAdapter: mockAdapter}
	cmd := &CmdToDeleteRecipient{}
	userName := "testUser"

	mockAdapter.On("RemoveRecipientConfig", *cmd, userName).Return(nil)

	err := service.RemoveRecipientConfig(userName, cmd)

	assert.NoError(t, err)

	mockAdapter.AssertExpectations(t)
}

func TestSyncUserInfo(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := messageRecipientAppService{messageRecipientAdapter: mockAdapter}
	cmd := &CmdToSyncUserInfo{}

	expectedId := uint(0)
	mockAdapter.On("SyncUserInfo", *cmd).Return(expectedId, nil)

	data, err := service.SyncUserInfo(cmd)

	assert.NoError(t, err)
	assert.Equal(t, data, expectedId)

	mockAdapter.AssertExpectations(t)
}
