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
	mockData := []MessageRecipientDTO{
		{Name: "Recipient 1", Mail: "test1@example.com", Phone: "+8612345678901"},
		{Name: "Recipient 2", Mail: "test2@example.com", Phone: "+8612345678902"},
	}
	mockAdapter.On("GetRecipientConfig", 10, 1, userName).Return(mockData, int64(len(mockData)), nil)

	data, count, err := service.GetRecipientConfig(10, 1, userName)

	assert.NoError(t, err)
	assert.Equal(t, mockData, data)
	assert.Equal(t, int64(len(mockData)), count)

	mockAdapter.AssertExpectations(t)

	mockAdapter.On("GetRecipientConfig", 10, 1, "").Return([]MessageRecipientDTO{{}}, int64(0), nil)
	data, count, err = service.GetRecipientConfig(10, 1, "")
	assert.NoError(t, err)
	assert.Equal(t, []MessageRecipientDTO{{}}, data)
	assert.Equal(t, int64(0), count)

	mockAdapter.AssertExpectations(t)
}

func TestAddRecipientConfig(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)
	userName := "testUser"

	cmd1 := &CmdToAddRecipient{Name: "", Mail: "test@example.com", Phone: "+8612345678901"}
	err1 := service.AddRecipientConfig(userName, cmd1)
	assert.ErrorContains(t, err1, "the recipient is null")

	cmd3 := &CmdToAddRecipient{Name: "MaoMao19970922", Mail: "1043170898@qq.com", Phone: "+8615315420821"}
	mockAdapter.On("AddRecipientConfig", *cmd3, "hourunze97").Return(xerrors.Errorf("接收人姓名不能相同"))
	err3 := service.AddRecipientConfig("hourunze97", cmd3)
	assert.ErrorContains(t, err3, "接收人姓名不能相同")
	mockAdapter.AssertExpectations(t)
}

func TestAddRecipientConfig_InvalidEmail(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)

	cmd := &CmdToAddRecipient{Name: "Recipient 1", Mail: "invalid_email", Phone: "+8612345678901"}
	userName := "testUser"

	err := service.AddRecipientConfig(userName, cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "the email is invalid")
}

func TestAddRecipientConfig_InvalidPhone(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)

	cmd := &CmdToAddRecipient{Name: "Recipient 1", Mail: "test@example.com", Phone: "invalid_phone"}
	userName := "testUser"

	err := service.AddRecipientConfig(userName, cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "the phone number is invalid")
}

func TestUpdateRecipientConfig(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)
	userName := "testUser"

	cmd := &CmdToUpdateRecipient{Mail: "test@example.com", Phone: "+8615315420821"}
	mockAdapter.On("UpdateRecipientConfig", *cmd, userName).Return(xerrors.Errorf(
		"update recipient config failed"))

	err := service.UpdateRecipientConfig(userName, cmd)
	assert.ErrorContains(t, err, "update recipient config failed")
	mockAdapter.AssertExpectations(t)
}
func TestUpdateRecipientConfig_InvalidEmail(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)

	cmd := &CmdToUpdateRecipient{Mail: "invalid_email", Phone: "+8612345678901"}
	userName := "testUser"

	err := service.UpdateRecipientConfig(userName, cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "the email is invalid")
}

func TestUpdateRecipientConfig_InvalidPhone(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)

	cmd := &CmdToUpdateRecipient{Mail: "test@example.com", Phone: "invalid_phone"}
	userName := "testUser"

	err := service.UpdateRecipientConfig(userName, cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "the phone number is invalid")
}

func TestRemoveRecipientConfig(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)

	cmd := CmdToDeleteRecipient{RecipientId: "0"}
	userName := "testUser"

	mockAdapter.On("RemoveRecipientConfig", cmd, userName).Return(nil)

	err := service.RemoveRecipientConfig(userName, &cmd)

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

func TestSyncUserInfo_Error(t *testing.T) {
	mockAdapter := new(MockMessageRecipientAdapter)
	service := NewMessageRecipientAppService(mockAdapter)

	cmd := CmdToSyncUserInfo{UserName: "testUser"}
	mockAdapter.On("SyncUserInfo", cmd).Return(uint(0), xerrors.New("error"))

	_, err := service.SyncUserInfo(&cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sync user info failed")
}
