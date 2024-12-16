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
	// basic
	v1.POST("/inner", ctl.GetInnerMessage)
	v1.GET("/inner_quick", ctl.GetInnerMessageQuick)
	v1.GET("/inner/count", ctl.CountAllUnReadMessage)
	v1.PUT("/inner", ctl.SetMessageIsRead)
	v1.DELETE("/inner", ctl.RemoveMessage)

	//release-openeuler-summit
	v1.GET("/inner/todo", ctl.GetAllTodoMessage)
	v1.GET("/inner/about", ctl.GetAllAboutMessage)
	v1.GET("/inner/watch", ctl.GetAllWatchMessage)

	v1.GET("/inner/count_new", ctl.CountAllMessage)
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

	// ubmc
	v1.GET("/all", ctl.GetAllMessage)
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
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %v", err))

		return
	}
	userName, err := user.GetSystemUserName(ctx)
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
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to bind params, %v", err))
		return
	}

	cmd, err := params.toCmd()
	if err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("failed to convert req to cmd, %v", err))
		return
	}
	userName, err := user.GetSystemUserName(ctx)
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
	userName, err := user.GetSystemUserName(ctx)
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
// @Param			eventId body []string true "eventId"
// @Accept			json
// @Success			202	string accepted 设置已读成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  设置已读失败
// @Router			/message_center/inner [put]
// @Id		setMessageIsRead
func (ctl *messageListController) SetMessageIsRead(ctx *gin.Context) {
	var messages []string
	if err := ctx.BindJSON(&messages); err != nil {
		ctx.JSON(http.StatusBadRequest, "无法解析请求正文")
		return
	}
	userName, err := user.GetSystemUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	for _, eventId := range messages {
		if err := ctl.appService.SetMessageIsRead(userName, eventId); err != nil {
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
// @Param			eventId body []string true "eventId"
// @Accept			json
// @Success			202	string accepted 消息删除成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  消息删除失败
// @Router			/message_center/inner [delete]
// @Id	    removeMessage
func (ctl *messageListController) RemoveMessage(ctx *gin.Context) {
	var messages []string

	if err := ctx.BindJSON(&messages); err != nil {
		commonctl.SendBadRequestParam(ctx, xerrors.Errorf("无法解析请求正文"))
		return
	}
	userName, err := user.GetSystemUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	for _, eventId := range messages {
		if err := ctl.appService.RemoveMessage(userName, eventId); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("消息删除失败，"+
				"err:%v", err)})
			return
		}
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": "消息删除成功"})
}

// GetForumSystemMessage get form system message
// @Summary			GetForumSystemMessage
// @Description		get forum system message 获取论坛系统通知消息
// @Tags			message_center_openeuler_summit
// @Param			body body QueryParams true "QueryParams"
// @Accept			json
// @Success			202	string accepted 查询成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  查询失败
// @Router			/message_center/inner/forum/system [get]
// @Id	    getForumSystemMessage
func (ctl *messageListController) GetForumSystemMessage(ctx *gin.Context) {
	userName, err := user.GetSystemUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if data, count, err := ctl.appService.GetForumSystemMessage(userName, params.PageNum,
		params.CountPerPage, params.StartTime, params.IsRead); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetForumAboutMessage get form about message
// @Summary			GetForumAboutMessage
// @Description		get forum about message 获取论坛提到我的消息
// @Tags			message_center_openeuler_summit
// @Param			body body QueryParams true "QueryParams"
// @Accept			json
// @Success			202	string accepted 查询成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  查询失败
// @Router			/message_center/inner/forum/about [get]
// @Id	    getForumAboutMessage
func (ctl *messageListController) GetForumAboutMessage(ctx *gin.Context) {
	userName, err := user.GetSystemUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if data, count, err := ctl.appService.GetForumAboutMessage(userName, params.IsBot,
		params.PageNum, params.CountPerPage, params.StartTime, params.IsRead); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetMeetingToDoMessage get meeting to do message
// @Summary			GetMeetingToDoMessage
// @Description		get meeting to do message 获取待参加的会议消息
// @Tags			message_center_openeuler_summit
// @Param			body body QueryParams true "QueryParams"
// @Accept			json
// @Success			202	string accepted 查询成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  查询失败
// @Router			/message_center/inner/meeting/todo [get]
// @Id	    getMeetingToDoMessage
func (ctl *messageListController) GetMeetingToDoMessage(ctx *gin.Context) {
	userName, err := user.GetSystemUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if data, count, err := ctl.appService.GetMeetingToDoMessage(userName, params.Filter,
		params.PageNum, params.CountPerPage, params.StartTime, params.IsRead); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetCVEToDoMessage get cve to do message
// @Summary			GetCVEToDoMessage
// @Description		get cve to do message 获取待我处理的漏洞消息
// @Tags			message_center_openeuler_summit
// @Param			body body QueryParams true "QueryParams"
// @Accept			json
// @Success			202	string accepted 查询成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  查询失败
// @Router			/message_center/inner/cve/todo [get]
// @Id	    getCVEToDoMessage
func (ctl *messageListController) GetCVEToDoMessage(ctx *gin.Context) {
	userName, err := user.GetSystemUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data, count, err := ctl.appService.GetCVEToDoMessage(userName, params.GiteeUserName,
		params.IsDone, params.PageNum, params.CountPerPage, params.StartTime, params.IsRead); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetCVEMessage get cve message
// @Summary			GetCVEMessage
// @Description		get cve message 获取漏洞关注消息
// @Tags			message_center_openeuler_summit
// @Param			body body QueryParams true "QueryParams"
// @Accept			json
// @Success			202	string accepted 查询成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  查询失败
// @Router			/message_center/inner/cve [get]
// @Id	    getCVEMessage
func (ctl *messageListController) GetCVEMessage(ctx *gin.Context) {
	userName, err := user.GetSystemUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data, count, err := ctl.appService.GetCVEMessage(userName, params.GiteeUserName,
		params.PageNum, params.CountPerPage, params.StartTime, params.IsRead); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetIssueToDoMessage get issue to do message
// @Summary			GetIssueToDoMessage
// @Description		get issue to do message 获取待我处理的issue
// @Tags			message_center_openeuler_summit
// @Param			body body QueryParams true "QueryParams"
// @Accept			json
// @Success			202	string accepted 查询成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  查询失败
// @Router			/message_center/inner/issue/todo [get]
// @Id	    getIssueToDoMessage
func (ctl *messageListController) GetIssueToDoMessage(ctx *gin.Context) {
	userName, err := user.GetSystemUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data, count, err := ctl.appService.GetIssueToDoMessage(userName, params.GiteeUserName,
		params.IsDone, params.PageNum, params.CountPerPage, params.StartTime, params.IsRead); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetPullRequestToDoMessage get pull request to do message
// @Summary			GetPullRequestToDoMessage
// @Description		get pull request to do message 获取待我处理的pr
// @Tags			message_center_openeuler_summit
// @Param			body body QueryParams true "QueryParams"
// @Accept			json
// @Success			202	string accepted 查询成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  查询失败
// @Router			/message_center/inner/pull_request/todo [get]
// @Id	    getPullRequestToDoMessage
func (ctl *messageListController) GetPullRequestToDoMessage(ctx *gin.Context) {
	userName, err := user.GetSystemUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data, count, err := ctl.appService.GetPullRequestToDoMessage(userName,
		params.GiteeUserName, params.IsDone, params.PageNum, params.CountPerPage,
		params.StartTime, params.IsRead); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetGiteeAboutMessage get gitee about message
// @Summary			GetGiteeAboutMessage
// @Description		get gitee about message 获取gitee提到我的
// @Tags			message_center_openeuler_summit
// @Param			body body QueryParams true "QueryParams"
// @Accept			json
// @Success			202	string accepted 查询成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  查询失败
// @Router			/message_center/inner/gitee/about [get]
// @Id	    getGiteeAboutMessage
func (ctl *messageListController) GetGiteeAboutMessage(ctx *gin.Context) {
	userName, err := user.GetSystemUserName(ctx)
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
		params.IsBot, params.PageNum, params.CountPerPage, params.StartTime, params.IsRead); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetGiteeMessage get gitee message
// @Summary			GetGiteeMessage
// @Description		get gitee message 获取gitee 动态消息
// @Tags			message_center_openeuler_summit
// @Param			body body QueryParams true "QueryParams"
// @Accept			json
// @Success			202	string accepted 查询成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  查询失败
// @Router			/message_center/inner/gitee [get]
// @Id	    getGiteeMessage
func (ctl *messageListController) GetGiteeMessage(ctx *gin.Context) {
	userName, err := user.GetSystemUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if data, count, err := ctl.appService.GetGiteeMessage(userName, params.GiteeUserName,
		params.PageNum, params.CountPerPage, params.StartTime, params.IsRead); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetEurMessage get eur message
// @Summary			GetEurMessage
// @Description		get eur message 获取eur关注消息
// @Tags			message_center_openeuler_summit
// @Param			body body QueryParams true "QueryParams"
// @Accept			json
// @Success			202	string accepted 查询成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  查询失败
// @Router			/message_center/inner/eur [get]
// @Id	    getEurMessage
func (ctl *messageListController) GetEurMessage(ctx *gin.Context) {
	userName, err := user.GetSystemUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if data, count, err := ctl.appService.GetEurMessage(userName, params.PageNum,
		params.CountPerPage, params.StartTime, params.IsRead); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetAllTodoMessage get alltodo message
// @Summary			GetAllTodoMessage
// @Description		get all todo message 获取所有待办消息
// @Tags			message_center_openeuler_summit
// @Param			body body QueryParams true "QueryParams"
// @Accept			json
// @Success			202	string accepted 查询成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  查询失败
// @Router			/message_center/inner/todo [get]
// @Id	    getAllTodoMessage
func (ctl *messageListController) GetAllTodoMessage(ctx *gin.Context) {
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userName, err := user.GetSystemUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if data, count, err := ctl.appService.GetAllToDoMessage(userName, params.GiteeUserName,
		params.IsDone, params.PageNum, params.CountPerPage, params.StartTime,
		params.IsRead); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetAllAboutMessage get all about message
// @Summary			GetAllAboutMessage
// @Description		get all about message 获取所有提到我的消息
// @Tags			message_center_openeuler_summit
// @Param			body body QueryParams true "QueryParams"
// @Accept			json
// @Success			202	string accepted 查询成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  查询失败
// @Router			/message_center/inner/about [get]
// @Id	    getAllAboutMessage
func (ctl *messageListController) GetAllAboutMessage(ctx *gin.Context) {
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userName, err := user.GetSystemUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if data, count, err := ctl.appService.GetAllAboutMessage(userName, params.GiteeUserName,
		params.IsBot, params.PageNum, params.CountPerPage, params.StartTime, params.IsRead); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// GetAllWatchMessage get all watch message
// @Summary			GetAllWatchMessage
// @Description		get all watch message 获取所有关注消息
// @Tags			message_center_openeuler_summit
// @Param			body body QueryParams true "QueryParams"
// @Accept			json
// @Success			202	string accepted 查询成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  查询失败
// @Router			/message_center/inner/watch [get]
// @Id	    getAllWatchMessage
func (ctl *messageListController) GetAllWatchMessage(ctx *gin.Context) {
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userName, err := user.GetSystemUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if data, count, err := ctl.appService.GetAllWatchMessage(userName,
		params.GiteeUserName, params.PageNum, params.CountPerPage, params.StartTime, params.IsRead); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}

// CountAllMessage count all message
// @Summary			CountAllMessage
// @Description		count all message 获取所有消息分类数量
// @Tags			message_center_openeuler_summit
// @Param			body body QueryParams true "QueryParams"
// @Accept			json
// @Success			202	string accepted 查询成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  查询失败
// @Router			/message_center/inner/count_new [get]
// @Id	    countAllMessage
func (ctl *messageListController) CountAllMessage(ctx *gin.Context) {
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userName, err := user.GetSystemUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if data, err := ctl.appService.CountAllMessage(userName, params.GiteeUserName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"count": data})
	}
}

// GetAllMessage get all message
// @Summary			GetAllMessage
// @Description		get all message 获取所有消息
// @Tags			message_center_ubmc
// @Param			params body QueryParams true "QueryParams"
// @Accept			json
// @Success			202	string accepted 查询成功
// @Failure         400 string bad_request 无法解析请求正文
// @Failure			500	string system_error  查询失败
// @Router			/message_center/all [get]
// @Id	    getAllMessage
func (ctl *messageListController) GetAllMessage(ctx *gin.Context) {
	var params QueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userName, err := user.GetSystemUserName(ctx)
	if err != nil {
		commonctl.SendUnauthorized(ctx, xerrors.Errorf("get username failed, err:%v", err))
		return
	}
	if data, count, err := ctl.appService.GetAllMessage(userName, params.PageNum, params.CountPerPage, params.IsRead); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": xerrors.Errorf("查询失败，err:%v", err)})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"query_info": data, "count": count})
	}
}
