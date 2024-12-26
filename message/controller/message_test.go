package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/xerrors"

	"github.com/opensourceways/message-manager/message/app"
)

// MockMessageListAppService 是 MessageListAppService 的模拟实现
type MockMessageListAppService struct {
	mock.Mock
}

func (m *MockMessageListAppService) GetAllToDoMessage(userName string, giteeUsername string, isDone bool, pageNum, countPerPage int, startTime string) ([]app.MessageListDTO, int64, error) {
	panic("implement me")
}
func (m *MockMessageListAppService) GetAllAboutMessage(userName string, giteeUsername string, isBot *bool, pageNum, countPerPage int, startTime string, isRead *bool) ([]app.MessageListDTO, int64, error) {
	panic("implement me")
}
func (m *MockMessageListAppService) GetAllWatchMessage(userName string, giteeUsername string, pageNum, countPerPage int, startTime string, isRead *bool) ([]app.MessageListDTO, int64, error) {
	panic("implement me")
}
func (m *MockMessageListAppService) CountAllMessage(userName string, giteeUsername string) (app.CountDataDTO, error) {
	panic("implement me")
}
func (m *MockMessageListAppService) GetForumSystemMessage(userName string, pageNum, countPerPage int, startTime string, isRead *bool) ([]app.MessageListDTO, int64, error) {
	panic("implement me")
}
func (m *MockMessageListAppService) GetForumAboutMessage(userName string, isBot *bool, pageNum, countPerPage int, startTime string, isRead *bool) ([]app.MessageListDTO, int64, error) {
	panic("implement me")
}
func (m *MockMessageListAppService) GetMeetingToDoMessage(userName string, filter int, pageNum, countPerPage int) ([]app.MessageListDTO, int64, error) {
	panic("implement me")
}
func (m *MockMessageListAppService) GetCVEToDoMessage(userName string, giteeUsername string, isDone bool, pageNum, countPerPage int, startTime string) ([]app.MessageListDTO, int64, error) {
	panic("implement me")
}
func (m *MockMessageListAppService) GetCVEMessage(userName string, giteeUsername string, pageNum, countPerPage int, startTime string, isRead *bool) ([]app.MessageListDTO, int64, error) {
	panic("implement me")
}
func (m *MockMessageListAppService) GetIssueToDoMessage(userName string, giteeUsername string, isDone bool, pageNum, countPerPage int, startTime string) ([]app.MessageListDTO, int64, error) {
	panic("implement me")
}
func (m *MockMessageListAppService) GetPullRequestToDoMessage(userName string, giteeUsername string, isDone bool, pageNum, countPerPage int, startTime string) ([]app.MessageListDTO, int64, error) {
	panic("implement me")
}
func (m *MockMessageListAppService) GetGiteeAboutMessage(userName string, giteeUsername string, isBot *bool, pageNum, countPerPage int, startTime string, isRead *bool) ([]app.MessageListDTO, int64, error) {
	panic("implement me")
}
func (m *MockMessageListAppService) GetGiteeMessage(userName string, giteeUsername string, pageNum, countPerPage int, startTime string, isRead *bool) ([]app.MessageListDTO, int64, error) {
	panic("implement me")
}
func (m *MockMessageListAppService) GetEurMessage(userName string, pageNum, countPerPage int, startTime string, isRead *bool) ([]app.MessageListDTO, int64, error) {
	panic("implement me")
}

func (m *MockMessageListAppService) CountAllUnReadMessage(userName string) ([]app.CountDTO, error) {
	args := m.Called(userName)
	return args.Get(0).([]app.CountDTO), args.Error(1)
}

func (m *MockMessageListAppService) SetMessageIsRead(cmd *app.CmdToSetIsRead) error {
	args := m.Called(cmd)
	return args.Error(0)
}

func (m *MockMessageListAppService) RemoveMessage(cmd *app.CmdToSetIsRead) error {
	args := m.Called(cmd)
	return args.Error(0)
}

func TestCountAllUnReadMessage_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	req, err := http.NewRequest("GET", "/message_center/inner/count", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCountAllUnReadMessage_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	mockService.On("CountAllUnReadMessage", "testUser").
		Return(0, xerrors.New("error"))

	req, err := http.NewRequest("GET", "/message_center/inner/count", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSetMessageIsRead_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	messages := []messageStatus{
		{Source: "source1", EventId: "event1"},
	}
	body, err := json.Marshal(messages)
	if err != nil {
		t.Fatal("Failed to marshal messages", err)
	}
	mockService.On("SetMessageIsRead", mock.Anything).Return(nil)

	req, err := http.NewRequest("PUT", "/message_center/inner", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	mockService.AssertExpectations(t)
}

func TestSetMessageIsRead_BindError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	req, err := http.NewRequest("PUT", "/message_center/inner", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSetMessageIsRead_ConvertError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	messages := []messageStatus{
		{Source: "source1", EventId: "event1"},
	}
	body, err := json.Marshal(messages)
	if err != nil {
		t.Fatal("Failed to marshal messages:", err)
	}
	mockService.On("SetMessageIsRead", mock.Anything).Return(xerrors.New("error"))

	req, err := http.NewRequest("PUT", "/message_center/inner", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRemoveMessage_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	messages := []messageStatus{
		{Source: "source1", EventId: "event1"},
	}
	body, err := json.Marshal(messages)
	if err != nil {
		t.Fatal("Failed to marshal messages:", err)
	}
	mockService.On("RemoveMessage", mock.Anything).Return(nil)

	req, err := http.NewRequest("DELETE", "/message_center/inner", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
	mockService.AssertExpectations(t)
}

func TestRemoveMessage_BindError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	req, err := http.NewRequest("DELETE", "/message_center/inner", nil)
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRemoveMessage_ConvertError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	mockService := new(MockMessageListAppService)
	AddRouterForMessageListController(r, mockService)

	messages := []messageStatus{
		{Source: "source1", EventId: "event1"},
	}
	body, err := json.Marshal(messages)
	if err != nil {
		t.Fatal("Failed to marshal messages:", err)
	}
	mockService.On("RemoveMessage", mock.Anything).Return(xerrors.New("error"))

	req, err := http.NewRequest("DELETE", "/message_center/inner", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal("Failed to create request:", err)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
