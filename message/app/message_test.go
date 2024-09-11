package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/xerrors"
)

// MockMessageListAdapter 是 MessageListAdapter 的模拟实现
type MockMessageListAdapter struct {
	mock.Mock
}

func (m *MockMessageListAdapter) GetInnerMessageQuick(cmd CmdToGetInnerMessageQuick, userName string) ([]MessageListDTO, int64, error) {
	args := m.Called(cmd, userName)
	return args.Get(0).([]MessageListDTO), args.Get(1).(int64), args.Error(2)
}

func (m *MockMessageListAdapter) GetInnerMessage(cmd CmdToGetInnerMessage, userName string) ([]MessageListDTO, int64, error) {
	args := m.Called(cmd, userName)
	return args.Get(0).([]MessageListDTO), args.Get(1).(int64), args.Error(2)
}

func (m *MockMessageListAdapter) CountAllUnReadMessage(userName string) ([]CountDTO, error) {
	args := m.Called(userName)
	return args.Get(0).([]CountDTO), args.Error(1)
}

func (m *MockMessageListAdapter) SetMessageIsRead(source string, eventId string) error {
	args := m.Called(source, eventId)
	return args.Error(0)
}

func (m *MockMessageListAdapter) RemoveMessage(source string, eventId string) error {
	args := m.Called(source, eventId)
	return args.Error(0)
}

func TestGetInnerMessageQuick(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	service := NewMessageListAppService(mockAdapter)

	userName := "testUser"
	cmd := CmdToGetInnerMessageQuick{Source: "test_source", CountPerPage: 10, PageNum: 1, ModeName: "test_mode"}
	mockData := []MessageListDTO{
		{
			Title:           "Test Title 1",
			Summary:         "Summary of message 1",
			Source:          "source1",
			Type:            "info",
			EventId:         "event1",
			DataContentType: "application/json",
			DataSchema:      "schema1",
			SpecVersion:     "1.0",
			EventTime:       time.Now(),
			User:            "user1",
			SourceUrl:       "http://example.com/1",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			IsRead:          false,
			SourceGroup:     "group1",
		},
		{
			Title:           "Test Title 2",
			Summary:         "Summary of message 2",
			Source:          "source2",
			Type:            "alert",
			EventId:         "event2",
			DataContentType: "application/json",
			DataSchema:      "schema2",
			SpecVersion:     "1.0",
			EventTime:       time.Now(),
			User:            "user2",
			SourceUrl:       "http://example.com/2",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			IsRead:          true,
			SourceGroup:     "group2",
		},
	}

	mockAdapter.On("GetInnerMessageQuick", cmd, userName).Return(mockData, int64(len(mockData)), nil)

	data, count, err := service.GetInnerMessageQuick(userName, &cmd)

	assert.NoError(t, err)
	assert.Equal(t, mockData, data)
	assert.Equal(t, int64(len(mockData)), count)
	mockAdapter.AssertExpectations(t)

	cmd1 := CmdToGetInnerMessageQuick{Source: "", CountPerPage: 0, PageNum: 1,
		ModeName: "test_mode"}
	mockAdapter.On("GetInnerMessageQuick", cmd1, userName).Return([]MessageListDTO{},
		int64(0), xerrors.Errorf("get inner message failed"))
	data1, count1, err1 := service.GetInnerMessageQuick(userName, &cmd1)
	assert.ErrorContains(t, err1, "get inner message failed")
	assert.Equal(t, []MessageListDTO{}, data1)
	assert.Equal(t, int64(0), count1)
	mockAdapter.AssertExpectations(t)

}

func TestGetInnerMessage(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	service := NewMessageListAppService(mockAdapter)

	userName := "testUser"
	cmd := CmdToGetInnerMessage{Source: "test_source", EventType: "test_event", IsRead: "false", CountPerPage: 10, PageNum: 2}
	mockData := []MessageListDTO{
		{
			Title:           "Test Title 1",
			Summary:         "Summary of message 1",
			Source:          "source1",
			Type:            "info",
			EventId:         "event1",
			DataContentType: "application/json",
			DataSchema:      "schema1",
			SpecVersion:     "1.0",
			EventTime:       time.Now(),
			User:            "user1",
			SourceUrl:       "http://example.com/1",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			IsRead:          false,
			SourceGroup:     "group1",
		},
	}

	mockAdapter.On("GetInnerMessage", cmd, userName).Return(mockData, int64(len(mockData)), nil)

	data, count, err := service.GetInnerMessage(userName, &cmd)

	assert.NoError(t, err)
	assert.Equal(t, mockData, data)
	assert.Equal(t, int64(len(mockData)), count)
	mockAdapter.AssertExpectations(t)

	cmd1 := CmdToGetInnerMessage{Source: "", EventType: "test_event", IsRead: "false",
		CountPerPage: 0, PageNum: 2}
	mockAdapter.On("GetInnerMessage", cmd1, userName).Return([]MessageListDTO{}, int64(0),
		xerrors.Errorf("get inner message failed"))

	data1, count1, err1 := service.GetInnerMessage(userName, &cmd1)
	assert.ErrorContains(t, err1, "get inner message failed")
	assert.Equal(t, []MessageListDTO{}, data1)
	assert.Equal(t, int64(0), count1)
	mockAdapter.AssertExpectations(t)
}

func TestCountAllUnReadMessage(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	service := NewMessageListAppService(mockAdapter)

	userName := "testUser"
	mockData := []CountDTO{
		{Source: "source1", Count: 5},
	}

	mockAdapter.On("CountAllUnReadMessage", userName).Return(mockData, nil)

	data, err := service.CountAllUnReadMessage(userName)

	assert.NoError(t, err)
	assert.Equal(t, mockData, data)
	mockAdapter.AssertExpectations(t)

	mockAdapter.On("CountAllUnReadMessage", "").Return([]CountDTO{},
		xerrors.Errorf("get count failed"))

	data1, err1 := service.CountAllUnReadMessage("")

	assert.ErrorContains(t, err1, "get count failed")
	assert.Equal(t, []CountDTO{}, data1)
	mockAdapter.AssertExpectations(t)

}

func TestSetMessageIsRead(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	service := NewMessageListAppService(mockAdapter)

	cmd := CmdToSetIsRead{Source: "test_source", EventId: "event1"}
	mockAdapter.On("SetMessageIsRead", cmd.Source, cmd.EventId).Return(nil)

	err := service.SetMessageIsRead(&cmd)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestSetMessageIsRead_Error(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	service := NewMessageListAppService(mockAdapter)

	cmd := CmdToSetIsRead{Source: "test_source", EventId: "event1"}
	mockAdapter.On("SetMessageIsRead", cmd.Source, cmd.EventId).Return(xerrors.New("error"))

	err := service.SetMessageIsRead(&cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "set message is_read failed")
}

func TestRemoveMessage(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	service := NewMessageListAppService(mockAdapter)

	cmd := CmdToSetIsRead{Source: "test_source", EventId: "event1"}
	mockAdapter.On("RemoveMessage", cmd.Source, cmd.EventId).Return(nil)

	err := service.RemoveMessage(&cmd)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestRemoveMessage_Error(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	service := NewMessageListAppService(mockAdapter)

	cmd := CmdToSetIsRead{Source: "test_source", EventId: "event1"}
	mockAdapter.On("RemoveMessage", cmd.Source, cmd.EventId).Return(xerrors.New("error"))

	err := service.RemoveMessage(&cmd)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "set message is_read failed")
}
