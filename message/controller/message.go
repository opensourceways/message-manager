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
	v1.POST("/inner", ctl.GetInnerMessage)
	v1.GET("/inner_quick", ctl.GetInnerMessageQuick)
	v1.GET("/inner/count", ctl.CountAllUnReadMessage)
	v1.PUT("/inner", ctl.SetMessageIsRead)
	v1.DELETE("/inner", ctl.RemoveMessage)

	//release-openeuler-summit
	v1.GET("/inner/todo", ctl.GetAllTodoMessage)
	v1.GET("/inner/about", ctl.GetAllAboutMessage)
	v1.GET("/inner/watch", ctl.GetAllWatchMessage)

	v1.GET("/inner/forum/system", ctl.GetForumSystemMessage)
	v1.GET("/inner/forum/about", ctl.GetForumAboutMessage)
	v1.GET("/inner/meeting/todo", ctl.GetMeetingToDoMessage)
	v1.GET("/inner/cve/todo", ctl.GetCVEToDoMessage)
	v1.GET("/inner/cve", ctl.GetCVEMessage)
	v1.GET("/inner/gitee/issue/todo", ctl.GetIssueToDoMessage)
	v1.GET("/inner/gitee/pr/todo", ctl.GetPullRequestToDoMessage)
	v1.GET("/inner/gitee/about", ctl.GetGiteeAboutMessage)
	v1.GET("/inner/gitee", ctl.GetGiteeMessage)
	v1.GET("/inner/eur", ctl.GetEurMessage)
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
// @Id		getInnerMessageQuick
func (ctl *messageListController) GetInnerMessageQuick(ctx *gin.Context) {
	var params queryInnerParamsQuick
	if err := ctx.ShouldBindQuery(&params); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to bind params, %v", err))
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
// @Param           Params query queryInnerParams true "Query InnerParams"
// @Accept			json
// @Success			202	 {object}  app.MessageListDTO
// @Failure			500	string system_error  查询失败
// @Failure         400 string bad_request  无法解析请求正文
// @Router			/message_center/inner [post]
// @Id		getInnerMessage
func (ctl *messageListController) GetInnerMessage(ctx *gin.Context) {
	var params queryInnerParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
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
// @Id		countAllUnReadMessage
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
// @Id		setMessageIsRead
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
// @Id	    removeMessage
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

// GetForumSystemMessage get form system message
func (ctl *messageListController) GetForumSystemMessage(ctx *gin.Context) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if data, count, err := ctl.appService.GetForumSystemMessage(userName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetForumAboutMessage get form about message
func (ctl *messageListController) GetForumAboutMessage(ctx *gin.Context) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if data, count, err := ctl.appService.GetForumAboutMessage(userName, params.IsBot); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetMeetingToDoMessage get meeting to do message
func (ctl *messageListController) GetMeetingToDoMessage(ctx *gin.Context) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if data, count, err := ctl.appService.GetMeetingToDoMessage(userName, params.GiteeUserName,
		params.Filter); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetCVEToDoMessage get cve to do message
func (ctl *messageListController) GetCVEToDoMessage(ctx *gin.Context) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data, count, err := ctl.appService.GetCVEToDoMessage(userName, params.GiteeUserName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetCVEMessage get cve message
func (ctl *messageListController) GetCVEMessage(ctx *gin.Context) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data, count, err := ctl.appService.GetCVEMessage(userName, params.GiteeUserName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetIssueToDoMessage get issue to do message
func (ctl *messageListController) GetIssueToDoMessage(ctx *gin.Context) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data, count, err := ctl.appService.GetIssueToDoMessage(userName, params.GiteeUserName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetPullRequestToDoMessage get pull request to do message
func (ctl *messageListController) GetPullRequestToDoMessage(ctx *gin.Context) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data, count, err := ctl.appService.GetPullRequestToDoMessage(userName, params.GiteeUserName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetGiteeAboutMessage get gitee about message
func (ctl *messageListController) GetGiteeAboutMessage(ctx *gin.Context) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if data, count, err := ctl.appService.GetGiteeAboutMessage(userName, params.GiteeUserName,
		params.IsBot); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetGiteeMessage get gitee message
func (ctl *messageListController) GetGiteeMessage(ctx *gin.Context) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data, count, err := ctl.appService.GetGiteeMessage(userName, params.GiteeUserName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetEurMessage get eur message
func (ctl *messageListController) GetEurMessage(ctx *gin.Context) {
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if data, count, err := ctl.appService.GetEurMessage(userName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetAllTodoMessage get alltodo message
func (ctl *messageListController) GetAllTodoMessage(ctx *gin.Context) {
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if data, count, err := ctl.appService.GetAllToDoMessage(userName, params.GiteeUserName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetAllAboutMessage get all about message
func (ctl *messageListController) GetAllAboutMessage(ctx *gin.Context) {
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if data, count, err := ctl.appService.GetAllAboutMessage(userName, params.GiteeUserName,
		params.IsBot); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetAllWatchMessage get all watch message
func (ctl *messageListController) GetAllWatchMessage(ctx *gin.Context) {
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userName, err := user.GetEulerUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if data, count, err := ctl.appService.GetAllWatchMessage(userName,
		params.GiteeUserName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}
