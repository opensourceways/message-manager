package app

import (
	"testing"

	"golang.org/x/xerrors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/datatypes"
)

// MockMessageSubscribeAdapter 是 MessageSubscribeAdapter 的模拟实现
type MockMessageSubscribeAdapter struct {
	mock.Mock
}

func (m *MockMessageSubscribeAdapter) GetAllSubsConfig(userName string) ([]MessageSubscribeDTO, error) {
	args := m.Called(userName)
	return args.Get(0).([]MessageSubscribeDTO), args.Error(1)
}

func (m *MockMessageSubscribeAdapter) GetSubsConfig(userName string) ([]MessageSubscribeDTO, int64, error) {
	args := m.Called(userName)
	return args.Get(0).([]MessageSubscribeDTO), args.Get(1).(int64), args.Error(2)
}

func (m *MockMessageSubscribeAdapter) SaveFilter(cmd CmdToGetSubscribe, userName string) error {
	args := m.Called(cmd, userName)
	return args.Error(0)
}

func (m *MockMessageSubscribeAdapter) AddSubsConfig(cmd CmdToAddSubscribe, userName string) ([]uint, error) {
	args := m.Called(cmd, userName)
	return args.Get(0).([]uint), args.Error(1)
}

func (m *MockMessageSubscribeAdapter) RemoveSubsConfig(cmd CmdToDeleteSubscribe, userName string) error {
	args := m.Called(cmd, userName)
	return args.Error(0)
}

func TestGetAllSubsConfig(t *testing.T) {
	mockAdapter := new(MockMessageSubscribeAdapter)
	service := NewMessageSubscribeAppService(mockAdapter)

	userName := "testUser"
	mockData := []MessageSubscribeDTO{
		{ModeName: "mode1"},
		{ModeName: "mode2"},
	}

	mockAdapter.On("GetAllSubsConfig", userName).Return(mockData, nil)

	data, err := service.GetAllSubsConfig(userName)

	assert.NoError(t, err)
	assert.Equal(t, mockData, data)
	mockAdapter.AssertExpectations(t)

	mockAdapter.On("GetAllSubsConfig", "").Return([]MessageSubscribeDTO{},
		xerrors.Errorf("查询失败"))

	data1, err1 := service.GetAllSubsConfig("")

	assert.ErrorContains(t, err1, "查询失败")
	assert.Equal(t, []MessageSubscribeDTO{}, data1)
	mockAdapter.AssertExpectations(t)
}

func TestGetSubsConfig(t *testing.T) {
	mockAdapter := new(MockMessageSubscribeAdapter)
	service := NewMessageSubscribeAppService(mockAdapter)

	userName := "testUser"
	mockData := []MessageSubscribeDTO{
		{ModeName: "mode1"},
	}
	mockCount := int64(len(mockData))

	mockAdapter.On("GetSubsConfig", userName).Return(mockData, mockCount, nil)

	data, count, err := service.GetSubsConfig(userName)

	assert.NoError(t, err)
	assert.Equal(t, mockData, data)
	assert.Equal(t, mockCount, count)
	mockAdapter.AssertExpectations(t)

	mockAdapter.On("GetSubsConfig", "").Return([]MessageSubscribeDTO{},
		int64(0), xerrors.Errorf("查询失败"))

	data1, count1, err1 := service.GetSubsConfig("")

	assert.ErrorContains(t, err1, "查询失败")
	assert.Equal(t, []MessageSubscribeDTO{}, data1)
	assert.Equal(t, int64(0), count1)
	mockAdapter.AssertExpectations(t)
}

func TestSaveFilter(t *testing.T) {
	mockAdapter := new(MockMessageSubscribeAdapter)
	service := NewMessageSubscribeAppService(mockAdapter)

	userName := "testUser"
	cmd := CmdToGetSubscribe{Source: "source1", EventType: "event_type", IsRead: "false"}
	mockAdapter.On("SaveFilter", cmd, userName).Return(nil)

	err := service.SaveFilter(userName, &cmd)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)

	cmd1 := CmdToGetSubscribe{}
	mockAdapter.On("SaveFilter", cmd1, "").Return(xerrors.Errorf("用户名为空"))

	err1 := service.SaveFilter("", &cmd1)

	assert.ErrorContains(t, err1, "用户名为空")
	mockAdapter.AssertExpectations(t)
}

func TestAddSubsConfig(t *testing.T) {
	mockAdapter := new(MockMessageSubscribeAdapter)
	service := NewMessageSubscribeAppService(mockAdapter)

	userName := "testUser"
	cmd := CmdToAddSubscribe{
		Source:      "source1",
		EventType:   "event_type",
		SpecVersion: "1.0",
		ModeName:    "mode1",
		ModeFilter:  datatypes.JSON(`{"key": "value"}`),
	}
	mockData := []uint{1, 2, 3}

	mockAdapter.On("AddSubsConfig", cmd, userName).Return(mockData, nil)

	data, err := service.AddSubsConfig(userName, &cmd)

	assert.NoError(t, err)
	assert.Equal(t, mockData, data)
	mockAdapter.AssertExpectations(t)

	cmd1 := CmdToAddSubscribe{
		Source:      "https://gitee.com",
		EventType:   "note",
		SpecVersion: "1.0",
		ModeName:    "我提的issue的评论",
		ModeFilter:  datatypes.JSON(`{"NoteEvent.Issue.User.Login": "eq=MaoMao19970922"}`),
	}
	mockAdapter.On("AddSubsConfig", cmd1, "hourunze97").Return([]uint{},
		xerrors.Errorf("新增配置失败"))

	data1, err1 := service.AddSubsConfig("hourunze97", &cmd1)
	assert.ErrorContains(t, err1, "新增配置失败")
	assert.Equal(t, []uint{}, data1)
	mockAdapter.AssertExpectations(t)
}

func TestAddSubsConfig_Error(t *testing.T) {
	mockAdapter := new(MockMessageSubscribeAdapter)
	service := NewMessageSubscribeAppService(mockAdapter)

	userName := "testUser"
	cmd := CmdToAddSubscribe{} // 模式名称和过滤器为空

	_, err := service.AddSubsConfig(userName, &cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "必填项不能为空")
}

func TestRemoveSubsConfig(t *testing.T) {
	mockAdapter := new(MockMessageSubscribeAdapter)
	service := NewMessageSubscribeAppService(mockAdapter)

	userName := "testUser"
	cmd := CmdToDeleteSubscribe{
		Source:    "source1",
		EventType: "event_type",
		ModeName:  "mode1",
	}
	mockAdapter.On("RemoveSubsConfig", cmd, userName).Return(nil)

	err := service.RemoveSubsConfig(userName, &cmd)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestRemoveSubsConfig_Error(t *testing.T) {
	mockAdapter := new(MockMessageSubscribeAdapter)
	service := NewMessageSubscribeAppService(mockAdapter)

	userName := "testUser"
	cmd := CmdToDeleteSubscribe{
		Source:    "source1",
		EventType: "event_type",
		ModeName:  "mode1",
	}
	mockAdapter.On("RemoveSubsConfig", cmd, userName).Return(xerrors.New("error"))

	err := service.RemoveSubsConfig(userName, &cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "remove subs failed")
}
