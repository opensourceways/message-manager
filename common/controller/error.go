/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package controller the public ability of controller domain .
package controller

import (
	"net/http"
)

const (
	errorSystemError     = "system_error"
	errorBadRequestParam = "bad_request_param"
	errorUnauthorized    = "unauthorized"
)

type errorCode interface {
	ErrorCode() string
}

type errorNoPermission interface {
	errorCode

	NoPermission()
}

func httpError(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}

	sc := http.StatusInternalServerError
	code := errorSystemError

	if v, ok := err.(errorCode); ok {
		code = v.ErrorCode()

		if _, ok := err.(errorNoPermission); ok {
			sc = http.StatusForbidden
		}
	}

	return sc, code
}
