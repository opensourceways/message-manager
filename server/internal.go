/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package server

import (
	"github.com/gin-gonic/gin"
)

func setRouterOfInternal(engine *gin.Engine, services *allServices) {
	setRouteOfMessage(engine, services)
}
