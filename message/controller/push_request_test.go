package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeletePushConfigDTOToCmd(t *testing.T) {
	req := &deletePushConfigDTO{
		SubscribeId: 123,
		RecipientId: 456789,
	}

	cmd, err := req.toCmd()
	assert.NoError(t, err)
	assert.Equal(t, 123, cmd.SubscribeId)
	assert.Equal(t, int64(456789), cmd.RecipientId)
}

func TestNewPushConfigDTOToCmd(t *testing.T) {
	req := &newPushConfigDTO{
		SubscribeId:      321,
		RecipientId:      987654,
		NeedMessage:      true,
		NeedPhone:        false,
		NeedMail:         true,
		NeedInnerMessage: false,
	}

	cmd, err := req.toCmd()
	assert.NoError(t, err)
	assert.Equal(t, 321, cmd.SubscribeId)
	assert.Equal(t, int64(987654), cmd.RecipientId)
	assert.True(t, cmd.NeedMessage)
	assert.False(t, cmd.NeedPhone)
	assert.True(t, cmd.NeedMail)
	assert.False(t, cmd.NeedInnerMessage)
}

func TestUpdatePushConfigDTOToCmd(t *testing.T) {
	req := &updatePushConfigDTO{
		SubscribeId:      []string{"1", "2", "3"},
		RecipientId:      "recipient123",
		NeedMessage:      true,
		NeedPhone:        true,
		NeedMail:         false,
		NeedInnerMessage: true,
	}

	cmd, err := req.toCmd()
	assert.NoError(t, err)
	assert.Equal(t, []string{"1", "2", "3"}, cmd.SubscribeId)
	assert.Equal(t, "recipient123", cmd.RecipientId)
	assert.True(t, cmd.NeedMessage)
	assert.True(t, cmd.NeedPhone)
	assert.False(t, cmd.NeedMail)
	assert.True(t, cmd.NeedInnerMessage)
}
