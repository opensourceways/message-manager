package app

import (
	"github.com/stretchr/testify/mock"
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
