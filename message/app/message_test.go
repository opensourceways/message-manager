package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageListAdapter struct {
	mock.Mock
}

func (m *MockMessageListAdapter) GetInnerMessageQuick(cmd CmdToGetInnerMessageQuick,
	userName string) ([]MessageListDTO, int64, error) {
	args := m.Called(cmd, userName)
	return args.Get(0).([]MessageListDTO), args.Get(1).(int64), args.Error(2)
}

func (m *MockMessageListAdapter) GetInnerMessage(cmd CmdToGetInnerMessage,
	userName string) ([]MessageListDTO, int64, error) {
	args := m.Called(cmd, userName)
	return args.Get(0).([]MessageListDTO), args.Get(1).(int64), args.Error(2)
}

func (m *MockMessageListAdapter) CountAllUnReadMessage(userName string) ([]CountDTO, error) {
	args := m.Called(userName)
	return args.Get(0).([]CountDTO), args.Error(1)
}

func (m *MockMessageListAdapter) SetMessageIsRead(source, eventId string) error {
	args := m.Called(source, eventId)
	return args.Error(0)
}

func (m *MockMessageListAdapter) RemoveMessage(source, eventId string) error {
	args := m.Called(source, eventId)
	return args.Error(0)
}

func TestGetInnerMessageQuick(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	service := messageListAppService{messageListAdapter: mockAdapter}
	cmd := &CmdToGetInnerMessageQuick{}
	userName := "testUser"

	expectedMessages := []MessageListDTO{{}}
	expectedCount := int64(0)

	mockAdapter.On("GetInnerMessageQuick", *cmd, userName).Return(expectedMessages,
		expectedCount, nil)

	messages, count, err := service.GetInnerMessageQuick(userName, cmd)

	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)
	assert.Equal(t, expectedCount, count)

	mockAdapter.AssertExpectations(t)
}

func TestGetInnerMessage(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	service := messageListAppService{messageListAdapter: mockAdapter}
	cmd := &CmdToGetInnerMessage{}
	userName := "testUser"

	expectedMessages := []MessageListDTO{{}}
	expectedCount := int64(0)

	mockAdapter.On("GetInnerMessage", *cmd, userName).Return(expectedMessages,
		expectedCount, nil)

	messages, count, err := service.GetInnerMessage(userName, cmd)

	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)
	assert.Equal(t, expectedCount, count)

	mockAdapter.AssertExpectations(t)
}

func TestCountAllUnReadMessage(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	service := messageListAppService{messageListAdapter: mockAdapter}
	userName := "testUser"

	expectedData := []CountDTO{{}}
	mockAdapter.On("CountAllUnReadMessage", userName).Return(expectedData, nil)

	data, err := service.CountAllUnReadMessage(userName)

	assert.NoError(t, err)
	assert.Equal(t, expectedData, data)

	mockAdapter.AssertExpectations(t)
}

func TestSetMessageIsRead(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	service := messageListAppService{messageListAdapter: mockAdapter}
	cmd := &CmdToSetIsRead{}

	mockAdapter.On("SetMessageIsRead", cmd.Source, cmd.EventId).Return(nil)

	err := service.SetMessageIsRead(cmd)

	assert.NoError(t, err)

	mockAdapter.AssertExpectations(t)
}

func TestRemoveMessage(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	service := messageListAppService{messageListAdapter: mockAdapter}
	cmd := &CmdToSetIsRead{}

	mockAdapter.On("RemoveMessage", cmd.Source, cmd.EventId).Return(nil)

	err := service.RemoveMessage(cmd)

	assert.NoError(t, err)

	mockAdapter.AssertExpectations(t)
}
