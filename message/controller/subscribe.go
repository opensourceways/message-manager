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

func AddRouterForMessageSubscribeController(
	r *gin.Engine,
	s app.MessageSubscribeAppService,
) {
	ctl := messageSubscribeController{
		appService: s,
	}
	v1 := r.Group("/message_center/config")
	v1.GET("/subs", ctl.GetSubsConfig)
	v1.GET("/subs/all", ctl.GetAllSubsConfig)
	v1.POST("/subs", ctl.AddSubsConfig)
	v1.PUT("/subs", ctl.UpdateSubsConfig)
	v1.DELETE("/subs", ctl.RemoveSubsConfig)
}

type messageSubscribeController struct {
	appService app.MessageSubscribeAppService
}

// GetAllSubsConfig
// @Summary			GetAllSubsConfig
// @Description		get all subscribe_config
// @Tags			message_subscribe
// @Accept			json
// @Success			202	 {object}  app.MessageSubscribeDTO
// @Failure			401	string unauthorized  用户未授权
// @Failure			500	string system_error  查询失败
// @Router			/message_center/config/subs/all [get]
// @Id		getAllSubsConfig
func (ctl *messageSubscribeController) GetAllSubsConfig(ctx *gin.Context) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	data, err := ctl.appService.GetAllSubsConfig(userName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data})
	}
}

// GetSubsConfig
// @Summary			GetSubsConfig
// @Description		get subscribe_config
// @Tags			message_subscribe
// @Accept			json
// @Success			202	 {object}  app.MessageSubscribeDTO
// @Failure			401	string unauthorized  用户未授权
// @Failure			500	string system_error  查询失败
// @Router			/message_center/config/subs [get]
// @Id		getSubsConfig
func (ctl *messageSubscribeController) GetSubsConfig(ctx *gin.Context) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if data, count, err := ctl.appService.GetSubsConfig(userName); err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// AddSubsConfig
// @Summary			AddSubsConfig
// @Description		add a subscribe_config
// @Tags			message_subscribe
// @Param			body body newSubscribeDTO true "newSubscribeDTO"
// @Accept			json
// @Success			202	string Accept  新增配置成功
// @Failure			400	string bad_request  无法解析请求正文
// @Failure			401	string unauthorized  用户未授权
// @Failure			500	string system_error  新增配置失败
// @Router			/message_center/config/subs [post]
// @Id		addSubsConfig
func (ctl *messageSubscribeController) AddSubsConfig(ctx *gin.Context) {
	var req newSubscribeDTO
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to bind params, %w", err))
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx,
			xerrors.Errorf("failed to convert req to cmd, %w", err))
		return
	}
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	data, err := ctl.appService.AddSubsConfig(userName, &cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": xerrors.Errorf("新增配置失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"newId": data, "message": "新增配置成功"})
	}
}

// UpdateSubsConfig
// @Summary			UpdateSubsConfig
// @Description		update a subscribe_config
// @Tags			message_subscribe
// @Param			body body updateSubscribeDTO true "updateSubscribeDTO"
// @Accept			json
// @Success			202	string Accept  更新配置成功
// @Failure			400	string bad_request  无法解析请求正文
// @Failure			401	string unauthorized  用户未授权
// @Failure			500	string system_error  更新配置成功
// @Router			/message_center/config/subs [put]
// @Id		updateSubsConfig
func (ctl *messageSubscribeController) UpdateSubsConfig(ctx *gin.Context) {
	var req updateSubscribeDTO
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to bind params, %w", err))
		return
	}
	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx,
			xerrors.Errorf("failed to convert req to cmd, %w", err))
		return
	}
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	err = ctl.appService.UpdateSubsConfig(userName, &cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "更新配置成功", "error": err})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "更新配置成功"})
	}
}

// RemoveSubsConfig
// @Summary			RemoveSubsConfig
// @Description		delete a subscribe_config by source and type
// @Tags			message_subscribe
// @Accept			json
// @Param 			body body deleteSubscribeDTO true "deleteSubscribeDTO"
// @Success			202 string Accept  删除配置成功
// @Failure			400	string bad_request  无法解析请求正文
// @Failure			401	string unauthorized  用户未授权
// @Failure			500	string system_error  删除配置失败
// @Router			/message_center/config/subs [delete]
// @Id		removeSubsConfig
func (ctl *messageSubscribeController) RemoveSubsConfig(ctx *gin.Context) {
	var req deleteSubscribeDTO
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to bind params, %w", err))
		return
	}
	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx,
			xerrors.Errorf("failed to convert req to cmd, %w", err))
		return
	}
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	err = ctl.appService.RemoveSubsConfig(userName, &cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": xerrors.Errorf("删除配置失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "删除配置成功"})
	}
}
