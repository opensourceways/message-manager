package app

import (
	"testing"

	"github.com/opensourceways/message-manager/message/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/xerrors"
)

// MockMessageListAdapter 是 MessageListAdapter 的模拟实现
type MockMessageListAdapter struct {
	mock.Mock
}

func (m *MockMessageListAdapter) SetAllMessageIsRead(userName, messageType, giteeUsername, startTime string, isRead, isDone, isBot *bool, filter int) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageListAdapter) GetAllToDoMessage(userName, giteeUsername string, isDone *bool, pageNum, countPerPage int, startTime string, isRead *bool) ([]domain.MessageListDO, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageListAdapter) GetAllAboutMessage(userName, giteeUsername string, isBot *bool, pageNum, countPerPage int, startTime string, isRead *bool) ([]domain.MessageListDO, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageListAdapter) GetAllWatchMessage(userName, giteeUsername string, pageNum, countPerPage int, startTime string, isRead *bool) ([]domain.MessageListDO, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageListAdapter) GetForumSystemMessage(userName string, pageNum, countPerPage int, startTime string, isRead *bool) ([]domain.MessageListDO, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageListAdapter) GetForumAboutMessage(userName string, isBot *bool, pageNum, countPerPage int, startTime string, isRead *bool) ([]domain.MessageListDO, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageListAdapter) GetMeetingToDoMessage(userName string, filter int, pageNum, countPerPage int, startTime string, isRead *bool) ([]domain.MessageListDO, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageListAdapter) GetCVEToDoMessage(userName, giteeUsername string, isDone *bool, pageNum, countPerPage int, startTime string, isRead *bool) ([]domain.MessageListDO, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageListAdapter) GetCVEMessage(userName, giteeUsername string, pageNum, countPerPage int, startTime string, isRead *bool) ([]domain.MessageListDO, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageListAdapter) GetIssueToDoMessage(userName, giteeUsername string, isDone *bool, pageNum, countPerPage int, startTime string, isRead *bool) ([]domain.MessageListDO, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageListAdapter) GetPullRequestToDoMessage(userName, giteeUsername string, isDone *bool, pageNum, countPerPage int, startTime string, isRead *bool) ([]domain.MessageListDO, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageListAdapter) GetGiteeAboutMessage(userName, giteeUsername string, isBot *bool, pageNum, countPerPage int, startTime string, isRead *bool) ([]domain.MessageListDO, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageListAdapter) GetGiteeMessage(userName, giteeUsername string, pageNum, countPerPage int, startTime string, isRead *bool) ([]domain.MessageListDO, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageListAdapter) GetEurMessage(userName string, pageNum, countPerPage int, startTime string, isRead *bool) ([]domain.MessageListDO, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageListAdapter) CountAllMessage(username, giteeUsername string) (domain.CountDataDO, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockMessageListAdapter) GetAllMessage(username string, pageNum, countPerPage int, isRead *bool) ([]domain.MessageListDO, int64, error) {
	//TODO implement me
	panic("implement me")
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
	userName := "testUser"
	eventId := "event1"
	mockAdapter.On("SetMessageIsRead", userName, eventId).Return(nil)

	err := service.SetMessageIsRead(userName, eventId)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestSetMessageIsRead_Error(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	service := NewMessageListAppService(mockAdapter)

	userName := "testUser"
	eventId := "event1"
	mockAdapter.On("SetMessageIsRead", userName, eventId).Return(xerrors.New("error"))

	err := service.SetMessageIsRead(userName, eventId)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "set message is_read failed")
}

func TestRemoveMessage(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	service := NewMessageListAppService(mockAdapter)

	userName := "testUser"
	eventId := "event1"
	mockAdapter.On("RemoveMessage", userName, eventId).Return(nil)

	err := service.RemoveMessage(userName, eventId)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestRemoveMessage_Error(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	service := NewMessageListAppService(mockAdapter)

	userName := "testUser"
	eventId := "event1"
	mockAdapter.On("RemoveMessage", userName, eventId).Return(xerrors.New("error"))

	err := service.RemoveMessage(userName, eventId)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "set message is_read failed")
}
