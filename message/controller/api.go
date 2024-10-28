package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/message-manager/utils"
	"github.com/sirupsen/logrus"
)

func AddWebRouter(r *gin.Engine) {
	v1 := r.Group("/message_center")
	v1.POST("/get_to_do_pulls", GetTodoPullRequest)
}

type RequestData struct {
	UserName string `json:"user_name"`
}

func GetTodoPullRequest(ctx *gin.Context) {
	var req RequestData
	if err := ctx.BindJSON(&req); err != nil {
		logrus.Errorf("bad request data, err:%v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	if req.UserName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}

	data, count, err := utils.GetTodoPulls(req.UserName)
	if err != nil {
		logrus.Errorf("get to-do pulls failed, err:%v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "get to-do pulls failed"})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"total_count": count, "to_do_pulls": data})
}
