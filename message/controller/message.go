/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/message-manager/common/user"
	"golang.org/x/xerrors"

	commonctl "github.com/opensourceways/message-manager/common/controller"
	"github.com/opensourceways/message-manager/message/app"
)

func AddRouterForMessageListController(
	r *gin.Engine,
	s app.MessageListAppService,
) {
	ctl := messageListController{
		appService: s,
	}

	v1 := r.Group("/message_center")
	v1.GET("/inner", ctl.GetInnerMessage)
	v1.GET("/inner_quick", ctl.GetInnerMessageQuick)
	v1.GET("/inner/count", ctl.CountAllUnReadMessage)
	v1.PUT("/inner", ctl.SetMessageIsRead)
	v1.DELETE("/inner", ctl.RemoveMessage)
}

type messageListController struct {
	appService app.MessageListAppService
}

// GetInnerMessageQuick
// @Summary			GetInnerMessageQuick
// @Description		get inner message by filter
// @Tags			message_center
// @Accept			json
// @Success			202	 {object}  app.MessageListDTO
// @Failure			500	string system_error  查询失败
// @Failure         400 string bad_request  请求参数错误
// @Router			/message_center/inner_quick [get]
func (ctl *messageListController) GetInnerMessageQuick(ctx *gin.Context) {
	var params queryInnerParamsQuick
	if err := ctx.ShouldBindQuery(&params); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to bind params, %w", err))
		return
	}

	cmd, err := params.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %w", err))

		return
	}
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if data, count, err := ctl.appService.GetInnerMessageQuick(userName, &cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetInnerMessage
// @Summary			GetInnerMessage
// @Description		get inner message
// @Tags			message_center
// @Accept			json
// @Success			202	 {object}  app.MessageListDTO
// @Failure			500	string system_error  查询失败
// @Failure         400 string bad_request  无法解析请求正文
// @Router			/message_center/inner [get]
func (ctl *messageListController) GetInnerMessage(ctx *gin.Context) {
	var params queryInnerParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to bind params, %w", err))
		return
	}

	cmd, err := params.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %w", err))
		return
	}
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if data, count, err := ctl.appService.GetInnerMessage(userName, &cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// CountAllUnReadMessage
// @Summary			CountAllUnReadMessage
// @Description		get unread inner message count
// @Tags			message_center
// @Accept			json
// @Success			202 {object} map[string]interface{} "成功响应"
// @Failure			401 {object} string "未授权"
// @Failure			500 {object} string "系统错误"
// @Router			/message_center/inner/count [get]
func (ctl *messageListController) CountAllUnReadMessage(ctx *gin.Context) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if data, err := ctl.appService.CountAllUnReadMessage(userName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("获取失败，"+
			"err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"count": data})
	}
}

// SetMessageIsRead
// @Summary			SetMessageIsRead
// @Description		set message read
// @Tags			message_center
// @Param			body body messageStatus true "messageStatus"
// @Accept			json
// @Success			202	string accepted 设置已读成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  设置已读失败
// @Router			/message_center/inner [put]
func (ctl *messageListController) SetMessageIsRead(ctx *gin.Context) {
	var messages []messageStatus
	if err := ctx.BindJSON(&messages); err != nil {
		ctx.JSON(http.StatusBadRequest, "无法解析请求正文")
		return
	}
	for _, msg := range messages {
		cmd, err := msg.toCmd()
		if err != nil {
			commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %w", err))
			return
		}
		if err := ctl.appService.SetMessageIsRead(&cmd); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf(
				"设置已读失败，err:%v", err)})
			return
		}
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": "设置已读成功"})
}

// RemoveMessage
// @Summary			RemoveMessage
// @Description		remove message
// @Tags			message_center
// @Param			body body messageStatus true "messageStatus"
// @Accept			json
// @Success			202	string accepted 消息删除成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  消息删除失败
// @Router			/message_center/inner [delete]
func (ctl *messageListController) RemoveMessage(ctx *gin.Context) {
	var messages []messageStatus

	if err := ctx.BindJSON(&messages); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("无法解析请求正文"))
		return
	}

	for _, msg := range messages {
		cmd, err := msg.toCmd()
		if err != nil {
			commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %w", err))
			return
		}
		if err := ctl.appService.RemoveMessage(&cmd); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("消息删除失败，"+
				"err:%v", err)})
			return
		}
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": "消息删除成功"})
}
