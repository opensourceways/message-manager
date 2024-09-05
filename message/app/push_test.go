package app

import (
	"github.com/stretchr/testify/mock"
)

type MockMessagePushAdapter struct {
	mock.Mock
}

func (m *MockMessagePushAdapter) GetPushConfig(subsIds []string, countPerPage, pageNum int,
	userName string) ([]MessagePushDTO, error) {
	args := m.Called(subsIds, countPerPage, pageNum, userName)
	return args.Get(0).([]MessagePushDTO), args.Error(1)
}

func (m *MockMessagePushAdapter) AddPushConfig(cmd CmdToAddPushConfig) error {
	args := m.Called(cmd)
	return args.Error(0)
}

func (m *MockMessagePushAdapter) UpdatePushConfig(cmd CmdToUpdatePushConfig) error {
	args := m.Called(cmd)
	return args.Error(0)
}

func (m *MockMessagePushAdapter) RemovePushConfig(cmd CmdToDeletePushConfig) error {
	args := m.Called(cmd)
	return args.Error(0)
}
