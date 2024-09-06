package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageRecipientAdapter struct {
	mock.Mock
}

func (m *MockMessageRecipientAdapter) GetRecipientConfig(countPerPage, pageNum int, userName string) ([]MessageRecipientDO, int64, error) {
	args := m.Called(countPerPage, pageNum, userName)
	return args.Get(0).([]MessageRecipientDO), args.Get(1).(int64), args.Error(2)
}

func (m *MockMessageRecipientAdapter) AddRecipientConfig(cmd CmdToAddRecipient, userName string) error {
	args := m.Called(cmd, userName)
	return args.Error(0)
}

func (m *MockMessageRecipientAdapter) UpdateRecipientConfig(cmd CmdToUpdateRecipient, userName string) error {
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
	countPerPage := 10
	pageNum := 1
	userName := "testUser"
	expectedRecipients := []MessageRecipientDO{{}}
	expectedCount := int64(5)

	mockAdapter.On("GetRecipientConfig", countPerPage, pageNum, userName).Return(expectedRecipients, expectedCount, nil)

	recipients, count, err := mockAdapter.GetRecipientConfig(countPerPage, pageNum, userName)

	assert.NoError(t, err)
	assert.Equal(t, expectedRecipients, recipients)
	assert.Equal(t, expectedCount, count)
	mockAdapter.AssertExpectations(t)
}

func TestAddRecipientConfig(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	cmd := CmdToAddRecipient{}
	userName := "testUser"

	mockAdapter.On("AddRecipientConfig", cmd, userName).Return(nil)

	err := mockAdapter.AddRecipientConfig(cmd, userName)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestUpdateRecipientConfig(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	cmd := CmdToUpdateRecipient{}
	userName := "testUser"

	mockAdapter.On("UpdateRecipientConfig", cmd, userName).Return(nil)

	err := mockAdapter.UpdateRecipientConfig(cmd, userName)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestRemoveRecipientConfig(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	cmd := CmdToDeleteRecipient{}
	userName := "testUser"

	mockAdapter.On("RemoveRecipientConfig", cmd, userName).Return(nil)

	err := mockAdapter.RemoveRecipientConfig(cmd, userName)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestSyncUserInfo(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	cmd := CmdToSyncUserInfo{}
	expectedUserID := uint(123)

	mockAdapter.On("SyncUserInfo", cmd).Return(expectedUserID, nil)

	userID, err := mockAdapter.SyncUserInfo(cmd)

	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, userID)
	mockAdapter.AssertExpectations(t)
}
