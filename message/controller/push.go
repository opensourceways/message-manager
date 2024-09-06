/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/message-manager/common/user"
	"golang.org/x/xerrors"

	commonctl "github.com/opensourceways/message-manager/common/controller"
	"github.com/opensourceways/message-manager/message/app"
)

func AddRouterForMessagePushController(
	r *gin.Engine,
	s app.MessagePushAppService,
) {
	ctl := messagePushController{
		appService: s,
	}

	v1 := r.Group("/message_center/config")
	v1.GET("/push", ctl.GetPushConfig)
	v1.POST("/push", ctl.AddPushConfig)
	v1.PUT("/push", ctl.UpdatePushConfig)
	v1.DELETE("/push", ctl.RemovePushConfig)
}

type messagePushController struct {
	appService app.MessagePushAppService
}

// GetPushConfig
// @Summary			GetPushConfig
// @Description		get push config
// @Tags			message_push
// @Accept			json
// @Success			202	 {object}  app.MessagePushDTO
// @Failure			500	string system_error  查询失败
// @Router			/message_center/config/push [get]
func (ctl *messagePushController) GetPushConfig(ctx *gin.Context) {
	subsIdsStr := ctx.DefaultQuery("subscribe_id", "")
	subsIds := strings.Split(subsIdsStr, ",")
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}

	countPerPage, err := strconv.Atoi(ctx.Query("count_per_page"))
	if err != nil {
		return
	}
	pageNum, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		return
	}

	if data, err := ctl.appService.GetPushConfig(countPerPage, pageNum, userName,
		subsIds); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data})
	}
}

// AddPushConfig
// @Summary			AddPushConfig
// @Description		add a new push_config
// @Tags			message_push
// @Accept			json
// @Param 			body body newPushConfigDTO true "newPushConfigDTO"
// @Success			202	string Accept  新增配置成功
// @Failure			400	string bad_request  无法解析请求正文
// @Failure			500	string system_error  新增配置失败
// @Router			/message_center/config/push [post]
func (ctl *messagePushController) AddPushConfig(ctx *gin.Context) {
	var req newPushConfigDTO
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to bind params, %w", err))
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %w", err))
		return
	}

	if err := ctl.appService.AddPushConfig(&cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("新增配置失败，err:%v",
			err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "新增配置成功"})
	}
}

// UpdatePushConfig
// @Summary			UpdatePushConfig
// @Description		update a push_config
// @Tags			message_push
// @Param			body body updatePushConfigDTO true "updatePushConfigDTO"
// @Accept			json
// @Success			202	string Accept  更新配置成功
// @Failure			400	string bad_request  无法解析请求正文
// @Failure			500	string system_error  更新配置失败
// @Router			/message_center/config/push [put]
func (ctl *messagePushController) UpdatePushConfig(ctx *gin.Context) {
	var req updatePushConfigDTO
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to bind params, %w", err))
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %w", err))
		return
	}
	if err := ctl.appService.UpdatePushConfig(&cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("更新配置失败,err:%v",
			err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "更新配置成功"})
	}
}

// RemovePushConfig
// @Summary			RemovePushConfig
// @Description		delete a push_config
// @Tags			message_push
// @Accept			json
// @Param			body body deletePushConfigDTO true "deletePushConfigDTO"
// @Success			202 string Accept  删除配置成功
// @Failure         400 string bad_request  无法解析请求正文
// @Failure			500	string system_error  删除配置失败
// @Router			/message_center/config/push [delete]
func (ctl *messagePushController) RemovePushConfig(ctx *gin.Context) {
	var req deletePushConfigDTO
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to bind params, %w", err))
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %w", err))
		return
	}
	if err := ctl.appService.RemovePushConfig(&cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("删除配置失败,err:%v",
			err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "删除配置成功"})
	}
}
