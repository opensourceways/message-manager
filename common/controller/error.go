/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package controller the public ability of controller domain .
package controller

import (
	"net/http"

	"github.com/opensourceways/message-manager/common/domain/allerror"
)

const (
	errorSystemError     = "system_error"
	errorBadRequestBody  = "bad_request_body"
	errorBadRequestParam = "bad_request_param"
	errorUnauthorized    = "unauthorized"
)

type errorCode interface {
	ErrorCode() string
}

type errorNotFound interface {
	errorCode

	NotFound()
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

		if _, ok := err.(errorNotFound); ok {
			sc = http.StatusNotFound

		} else if _, ok := err.(errorNoPermission); ok {
			sc = http.StatusForbidden

		} else {
			switch code {
			case allerror.ErrorCodeAccessTokenInvalid:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeSessionIdMissing:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeSessionIdInvalid:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeSessionNotFound:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeSessionInvalid:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeCSRFTokenMissing:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeCSRFTokenInvalid:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeCSRFTokenNotFound:
				sc = http.StatusUnauthorized

			case allerror.ErrorCodeAccessDenied:
				sc = http.StatusUnauthorized

			default:
				sc = http.StatusBadRequest
			}
		}
	}

	return sc, code
}
