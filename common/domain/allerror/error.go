/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package allerror storage all type of error
package allerror

import (
	"errors"
	"strings"
)

const (
	// errorCodeNoPermission mean no permission
	errorCodeNoPermission = "no_permission"

	ErrorCodeAccessDenied = "unauthorized"

	// ErrorCodeTokenNotFound means token is not found
	ErrorCodeTokenNotFound = "token_not_found"

	// ErrorCodeRepoNotFound means repo is not found
	ErrorCodeRepoNotFound = "repo_not_found"

	// ErrorCodeRepoAlreadyExist means repo is already exist
	ErrorCodeRepoAlreadyExist = "repo_already_exist"

	// ErrorCodeAccessTokenInvalid This error code is for restful api
	ErrorCodeAccessTokenInvalid = "access_token_invalid"

	// ErrorCodeSessionInvalid is const
	ErrorCodeSessionInvalid = "session_invalid"

	// ErrorCodeSessionIdInvalid is const
	ErrorCodeSessionIdInvalid = "session_id_invalid"

	// ErrorCodeSessionIdMissing is const
	ErrorCodeSessionIdMissing = "session_id_missing"

	// ErrorCodeSessionNotFound is const
	ErrorCodeSessionNotFound = "session_not_found"

	// ErrorCodeCSRFTokenMissing is const
	ErrorCodeCSRFTokenMissing = "csrf_token_missing" // #nosec G101

	// ErrorCodeCSRFTokenInvalid is const
	ErrorCodeCSRFTokenInvalid = "csrf_token_invalid" // #nosec G101

	// ErrorCodeCSRFTokenNotFound is const
	ErrorCodeCSRFTokenNotFound = "csrf_token_not_found" // #nosec G101

	// ErrorCodeBranchNotFound means branch is not found
	ErrorCodeBranchNotFound = "branch_not_found"

	// ErrorCodeBranchOrPathNotFound means branch or path is not found
	ErrorCodeBranchOrPathNotFound = "branch_or_path_not_found"

	// ErrorCodeFileNotFound means file is not found
	ErrorCodeFileNotFound = "file_not_found"

	// ErrorCodeNotDir means the target is not a directory
	ErrorCodeNotDir = "not_dir"

	// ErrorCodeNotFile means the target is not a file
	ErrorCodeNotFile = "not_file"

	// ErrorCodeEmptyRepo means the repo is empty
	ErrorCodeEmptyRepo = "empty_repo"

	// ErrorCodeOwnerNotFound means owner is not found
	ErrorCodeOwnerNotFound = "owner_not_found"

	// ErrorCodeModelNotFound means model is not found
	ErrorCodeModelNotFound = "model_not_found"

	// Invalid param
	errorCodeInvalidParam = "invalid_param"

	// Sensitive not pass
	ErrorSensitiveNotPass = "sensitive_not_pass"

	// no files or too many files
	ErrorNoFileOrTooManyFiles = "no_file_or_too_many_files"

	// ErrorEmptyFilePath means file path is empty
	ErrorEmptyFilePath = "empty_file_path"

	ErrorFileContentNoChanged = "file_content_is_not_changed"
)

// errorImpl
type errorImpl struct {
	code string
	msg  string
}

// Error return the errorImpl.msg
func (e errorImpl) Error() string {
	return e.msg
}

// ErrorCode return the errorImpl.code
func (e errorImpl) ErrorCode() string {
	return e.code
}

// New the new errorImpl struct
func New(code string, msg string) errorImpl {
	v := errorImpl{
		code: code,
	}

	if msg == "" {
		v.msg = strings.ReplaceAll(code, "_", " ")
	} else {
		v.msg = msg
	}

	return v
}

// notfoudError not found resource error struct
type notfoudError struct {
	errorImpl
}

// NotFound return empty
func (e notfoudError) NotFound() {}

// NewNotFound new the not found error
func NewNotFound(code string, msg string) notfoudError {
	return notfoudError{errorImpl: New(code, msg)}
}

// IsNotFound checks if an error is of type "notfoundError" and returns true if it is.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	var notfoudError notfoudError
	ok := errors.As(err, &notfoudError)

	return ok
}

// noPermissionError
type noPermissionError struct {
	errorImpl
}

// NoPermission return empty
func (e noPermissionError) NoPermission() {}

// NewNoPermission new the no permission error
func NewNoPermission(msg string) noPermissionError {
	return noPermissionError{errorImpl: New(errorCodeNoPermission, msg)}
}

// IsNoPermission check the error is NoPermission
func IsNoPermission(err error) bool {
	if err == nil {
		return false
	}

	var noPermissionError noPermissionError
	ok := errors.As(err, &noPermissionError)

	return ok
}

// NewInvalidParam new the invalid param
func NewInvalidParam(msg string) errorImpl {
	return New(errorCodeInvalidParam, msg)
}

// limitRateError
type limitRateError struct {
	errorImpl
}

// OverLimit is a marker method for over limit rate error.
func (l limitRateError) OverLimit() {}

// NewOverLimit creates a new over limit error with the specified code and message.
func NewOverLimit(code string, msg string) limitRateError {
	return limitRateError{errorImpl: New(code, msg)}
}

// IsErrorCodeEmptyRepo checks if an error has an error code of ErrorCodeEmptyRepo
func IsErrorCodeEmptyRepo(err error) bool {
	if err == nil {
		return false
	}

	var e errorImpl
	ok := errors.As(err, &e)
	if !ok {
		return false
	}

	return e.ErrorCode() == ErrorCodeEmptyRepo
}
