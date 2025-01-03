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
	// ErrorCodeRepoNotFound means repo is not found
	ErrorCodeRepoNotFound = "repo_not_found"
	// ErrorCodeEmptyRepo means the repo is empty
	ErrorCodeEmptyRepo = "empty_repo"
	// ErrorCodeModelNotFound means model is not found
	ErrorCodeModelNotFound = "model_not_found"
	// Invalid param
	errorCodeInvalidParam = "invalid_param"
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
