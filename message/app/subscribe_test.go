package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/xerrors"
)

type MockMessageSubscribeAdapter struct {
	mock.Mock
}

func (m *MockMessageSubscribeAdapter) GetAllSubsConfig(userName string) ([]MessageSubscribeDTO, error) {
	args := m.Called(userName)
	return args.Get(0).([]MessageSubscribeDTO), args.Error(1)
}

func (m *MockMessageSubscribeAdapter) GetSubsConfig(userName string) ([]MessageSubscribeDTO,
	int64, error) {
	args := m.Called(userName)
	return args.Get(0).([]MessageSubscribeDTO), args.Get(1).(int64), args.Error(2)
}

func (m *MockMessageSubscribeAdapter) SaveFilter(cmd CmdToGetSubscribe,
	userName string) error {
	args := m.Called(cmd, userName)
	return args.Error(0)
}

func (m *MockMessageSubscribeAdapter) AddSubsConfig(cmd CmdToAddSubscribe, userName string) ([]uint, error) {
	args := m.Called(cmd, userName)
	return args.Get(0).([]uint), args.Error(1)
}

func (m *MockMessageSubscribeAdapter) RemoveSubsConfig(cmd CmdToDeleteSubscribe,
	userName string) error {
	args := m.Called(cmd, userName)
	return args.Error(0)
}

func TestGetAllSubsConfig(t *testing.T) {
	mockAdapter := new(MockMessageSubscribeAdapter)
	service := messageSubscribeAppService{messageSubscribeAdapter: mockAdapter}
	userName := "testUser"

	expectedSubs := []MessageSubscribeDTO{{}}

	mockAdapter.On("GetAllSubsConfig", userName).
		Return(expectedSubs, nil)

	data, err := service.GetAllSubsConfig(userName)

	assert.NoError(t, err)
	assert.Equal(t, expectedSubs, data)

	mockAdapter.AssertExpectations(t)
}

func TestGetSubsConfig(t *testing.T) {
	mockAdapter := new(MockMessageSubscribeAdapter)
	service := messageSubscribeAppService{messageSubscribeAdapter: mockAdapter}
	userName := "testUser"

	expectedSubs := []MessageSubscribeDTO{{}}
	expectedCount := int64(0)

	mockAdapter.On("GetSubsConfig", userName).
		Return(expectedSubs, expectedCount, nil)

	data, count, err := service.GetSubsConfig(userName)

	assert.NoError(t, err)
	assert.Equal(t, expectedSubs, data)
	assert.Equal(t, expectedCount, count)

	mockAdapter.AssertExpectations(t)
}

func TestSaveFilter(t *testing.T) {
	mockAdapter := new(MockMessageSubscribeAdapter)
	service := messageSubscribeAppService{messageSubscribeAdapter: mockAdapter}
	cmd := &CmdToGetSubscribe{}
	userName := "testUser"

	mockAdapter.On("SaveFilter", *cmd, userName).Return(nil)

	err := service.SaveFilter(userName, cmd)

	assert.NoError(t, err)

	mockAdapter.AssertExpectations(t)
}

func TestAddSubsConfig(t *testing.T) {
	mockAdapter := new(MockMessageSubscribeAdapter)
	service := messageSubscribeAppService{messageSubscribeAdapter: mockAdapter}
	userName := "testUser"
	cmd := &CmdToAddSubscribe{}

	var expectedIds []uint
	expectedIds = []uint{}

	data, err := service.AddSubsConfig(userName, cmd)

	assert.EqualError(t, err, xerrors.Errorf("必填项不能为空").Error())
	assert.Equal(t, expectedIds, data)

}

func TestRemoveSubsConfig(t *testing.T) {
	mockAdapter := new(MockMessageSubscribeAdapter)
	service := messageSubscribeAppService{messageSubscribeAdapter: mockAdapter}
	cmd := &CmdToDeleteSubscribe{}
	userName := "testUser"

	mockAdapter.On("RemoveSubsConfig", *cmd, userName).Return(nil)

	err := service.RemoveSubsConfig(userName, cmd)

	assert.NoError(t, err)

	mockAdapter.AssertExpectations(t)
}
