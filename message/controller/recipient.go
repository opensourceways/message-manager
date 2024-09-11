/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/message-manager/common/user"
	"golang.org/x/xerrors"

	commonctl "github.com/opensourceways/message-manager/common/controller"
	"github.com/opensourceways/message-manager/message/app"
)

func AddRouterForMessageRecipientController(
	r *gin.Engine,
	s app.MessageRecipientAppService,
) {
	ctl := messageRecipientController{
		appService: s,
	}

	v1 := r.Group("/message_center/config")
	v1.GET("/recipient", ctl.GetRecipientConfig)
	v1.POST("/recipient", ctl.AddRecipientConfig)
	v1.POST("/recipient/sync", ctl.SyncUserInfo)
	v1.PUT("/recipient", ctl.UpdateRecipientConfig)
	v1.DELETE("/recipient", ctl.RemoveRecipientConfig)
}

type messageRecipientController struct {
	appService app.MessageRecipientAppService
}

// GetRecipientConfig
// @Summary			GetRecipientConfig
// @Description		get recipient config
// @Tags			recipient
// @Accept			json
// @Success			202	int count
// @Failure			500	string system_error  查询失败
// @Router			/message_center/config/recipient [get]
func (ctl *messageRecipientController) GetRecipientConfig(ctx *gin.Context) {
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
	if data, count, err := ctl.appService.GetRecipientConfig(countPerPage, pageNum, userName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"query_info": data, "count": count})
	}
}

// AddRecipientConfig
// @Summary			AddRecipientConfig
// @Description		add recipient config
// @Tags			recipient
// @Param			body body newRecipientDTO true "newRecipientDTO"
// @Accept			json
// @Success			202	string accepted 新增配置成功
// @Failure			400	string bad_request  无法解析请求正文
// @Failure			500	string server_error  新增配置失败
// @Router			/message_center/config/recipient [post]
func (ctl *messageRecipientController) AddRecipientConfig(ctx *gin.Context) {
	var req newRecipientDTO
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to bind params, %w", err))
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %w", err))
		return
	}
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if err := ctl.appService.AddRecipientConfig(userName, &cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("新增配置失败，err:%v",
			err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "新增配置成功"})
	}
}

// UpdateRecipientConfig
// @Summary			UpdateRecipientConfig
// @Description		update recipient config
// @Tags			recipient
// @Param			body body updateRecipientDTO true "updateRecipientDTO"
// @Accept			json
// @Success			202	string accepted 更新配置成功
// @Failure			400	string bad_request  无法解析请求正文
// @Failure			500	string server_error  更新配置失败
// @Router			/message_center/config/recipient [put]
func (ctl *messageRecipientController) UpdateRecipientConfig(ctx *gin.Context) {
	var req updateRecipientDTO
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to bind params, %w", err))
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %w", err))
		return
	}
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if err := ctl.appService.UpdateRecipientConfig(userName, &cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("更新配置失败，err:%v",
			err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "更新配置成功"})
	}
}

// RemoveRecipientConfig
// @Summary			RemoveRecipientConfig
// @Description		remove recipient config
// @Tags			recipient
// @Param			body body updateRecipientDTO true "updateRecipientDTO"
// @Accept			json
// @Success			202	string accepted 删除配置成功
// @Failure			400	string bad_request  无法解析请求正文
// @Failure			500	string server_error  删除配置失败
// @Router			/message_center/config/recipient [delete]
func (ctl *messageRecipientController) RemoveRecipientConfig(ctx *gin.Context) {
	var req updateRecipientDTO
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to bind params, %w", err))
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %w", err))
		return
	}
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if err := ctl.appService.UpdateRecipientConfig(userName, &cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("删除配置失败，err:%v",
			err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"message": "删除配置成功"})
	}
}

// SyncUserInfo
// @Summary			SyncUserInfo
// @Description		sync user info
// @Tags			recipient
// @Param			body body syncUserInfoDTO true "syncUserInfoDTO"
// @Accept			json
// @Success			202	string accepted 同步用户信息成功
// @Failure			400	string bad_request  无法解析请求正文
// @Failure			500	string server_error  同步用户信息失败
// @Router			/message_center/config//recipient/sync [post]
func (ctl *messageRecipientController) SyncUserInfo(ctx *gin.Context) {
	var req syncUserInfoDTO
	if err := ctx.BindJSON(&req); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to bind params, %w", err))
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %w", err))
		return
	}
	if data, err := ctl.appService.SyncUserInfo(&cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("同步用户信息失败，"+
			"err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"newId": data, "message": "同步用户信息成功"})
	}
}
