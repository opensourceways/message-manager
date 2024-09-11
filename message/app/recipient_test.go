package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/xerrors"
)

// MockMessageRecipientAdapter 是 MessageRecipientAdapter 的模拟实现
type MockMessageRecipientAdapter struct {
	mock.Mock
}

func (m *MockMessageRecipientAdapter) GetRecipientConfig(countPerPage, pageNum int, userName string) ([]MessageRecipientDTO, int64, error) {
	args := m.Called(countPerPage, pageNum, userName)
	if args.Get(0) == nil {
		return []MessageRecipientDTO{}, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]MessageRecipientDTO), args.Get(1).(int64), args.Error(2)
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
	service := NewMessageRecipientAppService(mockAdapter)

	userName := "testUser"
	countPerPage := 10
	pageNum := 1

	mockData := []MessageRecipientDTO{
		{Id: "1", Name: "Recipient 1", Mail: "recipient1@example.com", Phone: "+8613800138000"},
	}

	mockAdapter.On("GetRecipientConfig", countPerPage, pageNum, userName).Return(mockData, int64(len(mockData)), nil)

	data, count, err := service.GetRecipientConfig(countPerPage, pageNum, userName)

	assert.NoError(t, err)
	assert.Equal(t, mockData, data)
	assert.Equal(t, int64(len(mockData)), count)
	mockAdapter.AssertExpectations(t)
}

func TestGetRecipientConfig_Error(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)

	userName := "testUser"
	countPerPage := 10
	pageNum := 1

	mockAdapter.On("GetRecipientConfig", countPerPage, pageNum, userName).Return(nil, int64(0), xerrors.New("error"))

	data, count, err := service.GetRecipientConfig(countPerPage, pageNum, userName)

	assert.Error(t, err)
	assert.Empty(t, data)
	assert.Equal(t, int64(0), count)
	mockAdapter.AssertExpectations(t)
}

func TestAddRecipientConfig(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)

	cmd := &CmdToAddRecipient{
		Name:    "Recipient 1",
		Mail:    "recipient1@example.com",
		Phone:   "+8613800138000",
		Message: "Hello",
		Remark:  "Test recipient",
	}
	userName := "testUser"

	mockAdapter.On("AddRecipientConfig", *cmd, userName).Return(nil)

	err := service.AddRecipientConfig(userName, cmd)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestAddRecipientConfig_Error_NullName(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)

	cmd := &CmdToAddRecipient{
		Name:    "",
		Mail:    "recipient1@example.com",
		Phone:   "+8613800138000",
		Message: "Hello",
		Remark:  "Test recipient",
	}
	userName := "testUser"

	err := service.AddRecipientConfig(userName, cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "the recipient is null")
}

func TestAddRecipientConfig_Error_InvalidData(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)

	cmd := &CmdToAddRecipient{
		Name:    "Recipient 1",
		Mail:    "invalid-email",
		Phone:   "+8613800138000",
		Message: "Hello",
		Remark:  "Test recipient",
	}
	userName := "testUser"

	err := service.AddRecipientConfig(userName, cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data is invalid")
}

func TestUpdateRecipientConfig(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)

	cmd := &CmdToUpdateRecipient{
		Id:      "1",
		Name:    "Updated Recipient",
		Mail:    "updated@example.com",
		Phone:   "+8613800138000",
		Message: "Hello",
		Remark:  "Updated recipient",
	}
	userName := "testUser"

	mockAdapter.On("UpdateRecipientConfig", *cmd, userName).Return(nil)

	err := service.UpdateRecipientConfig(userName, cmd)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestUpdateRecipientConfig_Error_InvalidData(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)

	cmd := &CmdToUpdateRecipient{
		Id:      "1",
		Name:    "Updated Recipient",
		Mail:    "invalid-email",
		Phone:   "+8613800138000",
		Message: "Hello",
		Remark:  "Updated recipient",
	}
	userName := "testUser"

	err := service.UpdateRecipientConfig(userName, cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data is invalid")
}

func TestRemoveRecipientConfig(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)

	cmd := &CmdToDeleteRecipient{
		RecipientId: "1",
	}
	userName := "testUser"

	mockAdapter.On("RemoveRecipientConfig", *cmd, userName).Return(nil)

	err := service.RemoveRecipientConfig(userName, cmd)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestRemoveRecipientConfig_Error(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)

	cmd := &CmdToDeleteRecipient{
		RecipientId: "1",
	}
	userName := "testUser"

	mockAdapter.On("RemoveRecipientConfig", *cmd, userName).Return(xerrors.New("error"))

	err := service.RemoveRecipientConfig(userName, cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error")
}

func TestSyncUserInfo(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)

	cmd := &CmdToSyncUserInfo{
		Mail:          "user@example.com",
		Phone:         "+8613800138000",
		CountryCode:   "86",
		UserName:      "testUser",
		GiteeUserName: "giteeUser",
	}

	mockAdapter.On("SyncUserInfo", *cmd).Return(uint(1), nil)

	data, err := service.SyncUserInfo(cmd)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), data)
	mockAdapter.AssertExpectations(t)
}

func TestSyncUserInfo_Error(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)

	cmd := &CmdToSyncUserInfo{
		Mail:          "user@example.com",
		Phone:         "+8613800138000",
		CountryCode:   "86",
		UserName:      "testUser",
		GiteeUserName: "giteeUser",
	}

	mockAdapter.On("SyncUserInfo", *cmd).Return(uint(0), xerrors.New("error"))

	data, err := service.SyncUserInfo(cmd)

	assert.Error(t, err)
	assert.Equal(t, uint(0), data)
	mockAdapter.AssertExpectations(t)
}
