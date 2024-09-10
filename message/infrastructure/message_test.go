package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageListAdapter struct {
	mock.Mock
}

func (m *MockMessageListAdapter) GetInnerMessageQuick(cmd CmdToGetInnerMessageQuick,
	userName string) ([]MessageListDAO, int64, error) {
	args := m.Called(cmd, userName)
	return args.Get(0).([]MessageListDAO), args.Get(1).(int64), args.Error(2)
}

func (m *MockMessageListAdapter) GetInnerMessage(cmd CmdToGetInnerMessage, userName string) ([]MessageListDAO, int64, error) {
	args := m.Called(cmd, userName)
	return args.Get(0).([]MessageListDAO), args.Get(1).(int64), args.Error(2)
}

func (m *MockMessageListAdapter) CountAllUnReadMessage(userName string) ([]CountDAO, error) {
	args := m.Called(userName)
	return args.Get(0).([]CountDAO), args.Error(1)
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
	cmd := CmdToGetInnerMessageQuick{}
	userName := "testUser"
	expectedMessages := []MessageListDAO{{}}
	expectedCount := int64(5)

	mockAdapter.On("GetInnerMessageQuick", cmd, userName).Return(expectedMessages, expectedCount, nil)

	messages, count, err := mockAdapter.GetInnerMessageQuick(cmd, userName)

	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)
	assert.Equal(t, expectedCount, count)
	mockAdapter.AssertExpectations(t)
}

func TestGetInnerMessage(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	cmd := CmdToGetInnerMessage{}
	userName := "testUser"
	expectedMessages := []MessageListDAO{{}}
	expectedCount := int64(5)

	mockAdapter.On("GetInnerMessage", cmd, userName).Return(expectedMessages, expectedCount, nil)

	messages, count, err := mockAdapter.GetInnerMessage(cmd, userName)

	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)
	assert.Equal(t, expectedCount, count)
	mockAdapter.AssertExpectations(t)
}

func TestCountAllUnReadMessage(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	userName := "testUser"
	expectedCounts := []CountDAO{{}}

	mockAdapter.On("CountAllUnReadMessage", userName).Return(expectedCounts, nil)

	counts, err := mockAdapter.CountAllUnReadMessage(userName)

	assert.NoError(t, err)
	assert.Equal(t, expectedCounts, counts)
	mockAdapter.AssertExpectations(t)
}

func TestSetMessageIsRead(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	source := "source"
	eventId := "eventId"

	mockAdapter.On("SetMessageIsRead", source, eventId).Return(nil)

	err := mockAdapter.SetMessageIsRead(source, eventId)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestRemoveMessage(t *testing.T) {
	mockAdapter := new(MockMessageListAdapter)
	source := "source"
	eventId := "eventId"

	mockAdapter.On("RemoveMessage", source, eventId).Return(nil)

	err := mockAdapter.RemoveMessage(source, eventId)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}
