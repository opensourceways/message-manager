/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package app

import (
	"regexp"

	"golang.org/x/xerrors"

	"github.com/opensourceways/message-manager/message/domain"
)

type MessageRecipientAppService interface {
	GetRecipientConfig(countPerPage, pageNum int, userName string) ([]MessageRecipientDTO, int64,
		error)
	AddRecipientConfig(userName string, cmd *CmdToAddRecipient) error
	UpdateRecipientConfig(userName string, cmd *CmdToUpdateRecipient) error
	RemoveRecipientConfig(userName string, cmd *CmdToDeleteRecipient) error
	SyncUserInfo(cmd *CmdToSyncUserInfo) (uint, error)
}

func NewMessageRecipientAppService(
	messageRecipientAdapter domain.MessageRecipientAdapter,
) MessageRecipientAppService {
	return &messageRecipientAppService{
		messageRecipientAdapter: messageRecipientAdapter,
	}
}

type messageRecipientAppService struct {
	messageRecipientAdapter domain.MessageRecipientAdapter
}

const (
	EmailRegexp = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	EmailMaxLen = 254
	PhoneRegexp = `^\+861[3-9]\d{9}$`
	PhoneLen    = 14
)

func isValidEmail(email string) bool {
	// 简单邮箱正则表达式，可根据需要调整
	emailRegex := regexp.MustCompile(EmailRegexp)
	return emailRegex.MatchString(email) && len(email) <= EmailMaxLen
}

func isValidPhoneNumber(phoneNumber string) bool {
	// 中国大陆手机号码的简单正则表达式，可能根据情况调整
	phoneRegex := regexp.MustCompile(PhoneRegexp)
	return phoneRegex.MatchString(phoneNumber) && len(phoneNumber) == PhoneLen
}

func validateData(email string, phoneNumber string) error {
	if !isValidEmail(email) {
		return xerrors.Errorf("the email is invalid, email:%s", email)
	}
	if !isValidPhoneNumber(phoneNumber) {
		return xerrors.Errorf("the phone number is invalid, phone:%s", phoneNumber)
	}

	return nil
}

func (s *messageRecipientAppService) GetRecipientConfig(countPerPage, pageNum int, userName string) (
	[]MessageRecipientDTO, int64, error) {

	data, count, err := s.messageRecipientAdapter.GetRecipientConfig(countPerPage, pageNum,
		userName)
	if err != nil {
		return []MessageRecipientDTO{}, 0, err
	}
	return data, count, nil
}

func (s *messageRecipientAppService) AddRecipientConfig(userName string,
	cmd *CmdToAddRecipient) error {
	if cmd.Name == "" {
		return xerrors.Errorf("the recipient is null")
	}

	if err := validateData(cmd.Mail, cmd.Phone); err != nil {
		return xerrors.Errorf("data is invalid, err:%v", err.Error())
	}
	err := s.messageRecipientAdapter.AddRecipientConfig(*cmd, userName)
	if err != nil {
		return err
	}
	return nil
}

func (s *messageRecipientAppService) UpdateRecipientConfig(userName string,
	cmd *CmdToUpdateRecipient) error {
	if err := validateData(cmd.Mail, cmd.Phone); err != nil {
		return xerrors.Errorf("data is invalid, err:%v", err.Error())
	}

	err := s.messageRecipientAdapter.UpdateRecipientConfig(*cmd, userName)
	if err != nil {
		return err
	}
	return nil
}

func (s *messageRecipientAppService) RemoveRecipientConfig(userName string,
	cmd *CmdToDeleteRecipient) error {
	err := s.messageRecipientAdapter.RemoveRecipientConfig(*cmd, userName)
	if err != nil {
		return err
	}
	return nil
}

func (s *messageRecipientAppService) SyncUserInfo(cmd *CmdToSyncUserInfo) (uint, error) {
	data, err := s.messageRecipientAdapter.SyncUserInfo(*cmd)
	if err != nil {
		return 0, xerrors.Errorf("sync user info failed, err:%v", err)
	}
	return data, nil
}
