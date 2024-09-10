package infrastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageSubscribeAdapter struct {
	mock.Mock
}

func (m *MockMessageSubscribeAdapter) GetAllSubsConfig(userName string) ([]MessageSubscribeDAO, error) {
	args := m.Called(userName)
	return args.Get(0).([]MessageSubscribeDAO), args.Error(1)
}

func (m *MockMessageSubscribeAdapter) GetSubsConfig(userName string) ([]MessageSubscribeDAO, int64, error) {
	args := m.Called(userName)
	return args.Get(0).([]MessageSubscribeDAO), args.Get(1).(int64), args.Error(2)
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
	userName := "testUser"
	expectedSubs := []MessageSubscribeDAO{{}}

	mockAdapter.On("GetAllSubsConfig", userName).Return(expectedSubs, nil)

	subs, err := mockAdapter.GetAllSubsConfig(userName)

	assert.NoError(t, err)
	assert.Equal(t, expectedSubs, subs)
	mockAdapter.AssertExpectations(t)
}

func TestGetSubsConfig(t *testing.T) {
	mockAdapter := new(MockMessageSubscribeAdapter)
	userName := "testUser"
	expectedSubs := []MessageSubscribeDAO{{}}
	expectedCount := int64(3)

	mockAdapter.On("GetSubsConfig", userName).Return(expectedSubs, expectedCount, nil)

	subs, count, err := mockAdapter.GetSubsConfig(userName)

	assert.NoError(t, err)
	assert.Equal(t, expectedSubs, subs)
	assert.Equal(t, expectedCount, count)
	mockAdapter.AssertExpectations(t)
}

func TestSaveFilter(t *testing.T) {
	mockAdapter := new(MockMessageSubscribeAdapter)
	cmd := CmdToGetSubscribe{}
	userName := "testUser"

	mockAdapter.On("SaveFilter", cmd, userName).Return(nil)

	err := mockAdapter.SaveFilter(cmd, userName)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}

func TestAddSubsConfig(t *testing.T) {
	mockAdapter := new(MockMessageSubscribeAdapter)
	cmd := CmdToAddSubscribe{}
	userName := "testUser"
	expectedIDs := []uint{1, 2, 3}

	mockAdapter.On("AddSubsConfig", cmd, userName).Return(expectedIDs, nil)

	ids, err := mockAdapter.AddSubsConfig(cmd, userName)

	assert.NoError(t, err)
	assert.Equal(t, expectedIDs, ids)
	mockAdapter.AssertExpectations(t)
}

func TestRemoveSubsConfig(t *testing.T) {
	mockAdapter := new(MockMessageSubscribeAdapter)
	cmd := CmdToDeleteSubscribe{}
	userName := "testUser"

	mockAdapter.On("RemoveSubsConfig", cmd, userName).Return(nil)

	err := mockAdapter.RemoveSubsConfig(cmd, userName)

	assert.NoError(t, err)
	mockAdapter.AssertExpectations(t)
}
