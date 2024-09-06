package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncUserInfoDTOToCmd(t *testing.T) {
	req := &syncUserInfoDTO{
		Mail:          "user@example.com",
		Phone:         "1234567890",
		CountryCode:   "86",
		UserName:      "testuser",
		GiteeUserName: "giteeuser",
	}

	cmd, err := req.toCmd()
	assert.NoError(t, err)
	assert.Equal(t, "user@example.com", cmd.Mail)
	assert.Equal(t, "1234567890", cmd.Phone)
	assert.Equal(t, "86", cmd.CountryCode)
	assert.Equal(t, "testuser", cmd.UserName)
	assert.Equal(t, "giteeuser", cmd.GiteeUserName)
}

func TestNewRecipientDTOToCmd(t *testing.T) {
	req := &newRecipientDTO{
		Name:    "Recipient Name",
		Mail:    "recipient@example.com",
		Message: "Welcome!",
		Phone:   "0987654321",
		Remark:  "Important recipient",
	}

	cmd, err := req.toCmd()
	assert.NoError(t, err)
	assert.Equal(t, "Recipient Name", cmd.Name)
	assert.Equal(t, "recipient@example.com", cmd.Mail)
	assert.Equal(t, "Welcome!", cmd.Message)
	assert.Equal(t, "0987654321", cmd.Phone)
	assert.Equal(t, "Important recipient", cmd.Remark)
}

func TestUpdateRecipientDTOToCmd(t *testing.T) {
	req := &updateRecipientDTO{
		Id:      "recipient-id-123",
		Name:    "Updated Recipient Name",
		Mail:    "updated@example.com",
		Message: "Updated message",
		Phone:   "1122334455",
		Remark:  "Updated remark",
	}

	cmd, err := req.toCmd()
	assert.NoError(t, err)
	assert.Equal(t, "recipient-id-123", cmd.Id)
	assert.Equal(t, "Updated Recipient Name", cmd.Name)
	assert.Equal(t, "updated@example.com", cmd.Mail)
	assert.Equal(t, "Updated message", cmd.Message)
	assert.Equal(t, "1122334455", cmd.Phone)
	assert.Equal(t, "Updated remark", cmd.Remark)
}

func TestDeleteRecipientDTOToCmd(t *testing.T) {
	req := &deleteRecipientDTO{
		RecipientId: "recipient-id-456",
	}

	cmd, err := req.toCmd()
	assert.NoError(t, err)
	assert.Equal(t, "recipient-id-456", cmd.RecipientId)
}
