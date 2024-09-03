/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package controller the common of controller
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseData is a struct that holds the response data for an API request.
type ResponseData struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// newResponseData return the new response data
func newResponseData(data interface{}) ResponseData {
	return ResponseData{
		Data: data,
	}
}

// newResponseCodeError return the new response data and code
func newResponseCodeMsg(code, msg string) ResponseData {
	return ResponseData{
		Code: code,
		Msg:  msg,
	}
}

// SendBadRequestBody return the 400 about body invalid
func SendBadRequestBody(ctx *gin.Context, err error) {
	if _, ok := err.(errorCode); ok {
		SendError(ctx, err)
	} else {
		_ = ctx.Error(err)
		ctx.JSON(
			http.StatusBadRequest,
			newResponseCodeMsg(errorBadRequestBody, err.Error()),
		)
	}
}

// SendUnauthorized return 401
func SendUnauthorized(ctx *gin.Context, err error) {
	if _, ok := err.(errorCode); ok {
		SendError(ctx, err)
	} else {
		_ = ctx.Error(err)
		ctx.JSON(
			http.StatusUnauthorized,
			newResponseCodeMsg(errorUnauthorized, err.Error()),
		)
	}
}

// SendBadRequestParam return the 400 about param invalid
func SendBadRequestParam(ctx *gin.Context, err error) {
	if _, ok := err.(errorCode); ok {
		SendError(ctx, err)
	} else {
		_ = ctx.Error(err)
		ctx.JSON(
			http.StatusBadRequest,
			newResponseCodeMsg(errorBadRequestParam, err.Error()),
		)
	}
}

// SendRespOfPut return the put request
func SendRespOfPut(ctx *gin.Context, data interface{}) {
	if data == nil {
		ctx.JSON(http.StatusAccepted, newResponseCodeMsg("", "success"))
	} else {
		ctx.JSON(http.StatusAccepted, newResponseData(data))
	}
}

// SendRespOfGet return the get request
func SendRespOfGet(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, newResponseData(data))
}

// SendRespOfPost return the post request
func SendRespOfPost(ctx *gin.Context, data interface{}) {
	if data == nil {
		ctx.JSON(http.StatusCreated, newResponseCodeMsg("", "success"))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(data))
	}
}

// SendRespOfDelete return the delete request
func SendRespOfDelete(ctx *gin.Context) {
	ctx.JSON(http.StatusNoContent, newResponseCodeMsg("", "success"))
}

// SendError return the 400 about param invalid
func SendError(ctx *gin.Context, err error) {
	sc, code := httpError(err)

	_ = ctx.AbortWithError(sc, err)

	ctx.JSON(sc, newResponseCodeMsg(code, err.Error()))
}
