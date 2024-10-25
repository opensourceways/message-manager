// utils/utils_test.go
package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseUnixTimestamp(t *testing.T) {
	tests := []struct {
		input    string
		expected *time.Time
	}{
		{"", nil}, // 测试空字符串
		{"1633072800", func() *time.Time {
			t := time.Unix(1633072800, 0)
			return &t
		}()},
		{"invalid", nil}, // 测试无效的时间戳
	}

	for _, test := range tests {
		result := ParseUnixTimestamp(test.input)
		if (result == nil && test.expected != nil) || (result != nil && test.expected == nil) {
			t.Errorf("expected %v, got %v", test.expected, result)
		} else if result != nil && !result.Equal(*test.expected) {
			t.Errorf("expected %v, got %v", *test.expected, *result)
		}
	}
}

func TestIsEurMessage(t *testing.T) {
	if !IsEurMessage(EurSource) {
		t.Errorf("IsEurMessage should return true for EurSource")
	}
	if IsEurMessage("invalid") {
		t.Errorf("IsEurMessage should return false for invalid source")
	}
}

func TestIsGiteeMessage(t *testing.T) {
	if !IsGiteeMessage(GiteeSource) {
		t.Errorf("IsGiteeMessage should return true for GiteeSource")
	}
	if IsGiteeMessage("invalid") {
		t.Errorf("IsGiteeMessage should return false for invalid source")
	}
}

func TestIsMeetingMessage(t *testing.T) {
	if !IsMeetingMessage(MeetingSource) {
		t.Errorf("IsMeetingMessage should return true for MeetingSource")
	}
	if IsMeetingMessage("invalid") {
		t.Errorf("IsMeetingMessage should return false for invalid source")
	}
}

func TestIsCveMessage(t *testing.T) {
	if !IsCveMessage(CveSource) {
		t.Errorf("IsCveMessage should return true for CveSource")
	}
	if IsCveMessage("invalid") {
		t.Errorf("IsCveMessage should return false for invalid source")
	}
}

func TestSortStringList(t *testing.T) {
	input := []string{"b", "a", "*", "c", "d*"}
	expected := []string{"*", "d*", "a", "b", "c"}
	result := sortStringList(input)

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("expected %v, got %v", expected, result)
			break
		}
	}
}

func TestMergePaths(t *testing.T) {
	input := []string{"path1/*", "path2/*", "*", "path3/*", "path1/subpath"}
	expected := []string{"*"}
	result := MergePaths(input)

	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("expected %v, got %v", expected, result)
	}

	input2 := []string{"path1/*", "path2/*", "path1/subpath"}
	expected2 := []string{"path1/*", "path2/*"}
	result2 := MergePaths(input2)

	for i, v := range expected2 {
		if result2[i] != v {
			t.Errorf("expected %v, got %v", expected2, result2)
			break
		}
	}
}

func TestRemoveEmptyStrings(t *testing.T) {
	input := []string{"", "test", "", "example"}
	expected := []string{"test", "example"}
	result := RemoveEmptyStrings(input)

	if len(result) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(result))
	}

	for i, v := range expected {
		if result[i] != v {
			t.Errorf("expected %v, got %v", expected, result)
			break
		}
	}
}

type TestSigInfo struct {
	Sig []string `json:"sig"`
}

type TestGiteeRepo struct {
	FullName string `json:"full_name"`
	Admin    bool   `json:"admin"`
}

// 测试 GetUserSigInfo
func TestGetUserSigInfo(t *testing.T) {
	// 创建一个 HTTP 测试服务器
	handler := func(w http.ResponseWriter, r *http.Request) {
		// 模拟返回的 JSON 数据
		response := TestSigInfo{Sig: []string{}}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			t.Fatalf("could not decode response: %v", err)
		}
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// 测试函数
	sigs, err := GetUserSigInfo("testuser")
	assert.NoError(t, err)
	assert.Equal(t, []string{}, sigs)
}

// 测试 GetUserAdminReposByUsername
func TestGetUserAdminRepos(t *testing.T) {
	// 创建一个 HTTP 测试服务器
	handler := func(w http.ResponseWriter, r *http.Request) {
		// 模拟返回的 JSON 数据
		repos := []TestGiteeRepo{
			{FullName: "testuser/Git-Demo", Admin: true},
			{FullName: "testuser/jhj", Admin: false},
			{FullName: "testuser/testtest", Admin: true},
			{FullName: "testuser/testtest2", Admin: true},
		}
		// 设置 total_count 响应头
		w.Header().Set("total_count", "3")
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(repos)
		if err != nil {
			t.Fatalf("could not decode response: %v", err)
		}
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// 测试函数
	adminRepos, err := GetUserAdminReposByUsername("testuser")
	assert.NoError(t, err)
	assert.Equal(t, []string{"testuser/Git-Demo", "testuser/jhj", "testuser/testtest", "testuser/testtest2"}, adminRepos)
}
