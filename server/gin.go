/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/server-common-lib/interrupts"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/opensourceways/message-manager/config"
	"github.com/opensourceways/message-manager/docs"
)

const (
	version         = "development" // program version for this build
	apiDesc         = "message server manager APIs"
	apiTitle        = "message server"
	waitServerStart = 3 // 3s
)

func StartWebServer(cfg *config.Config) {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(logRequest())
	engine.UseRawPath = true

	docs.SwaggerInfo.Title = apiTitle
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Description = apiDesc

	services, err := initServices()
	if err != nil {
		return
	}

	setRouterOfInternal(engine, &services)

	// start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8082),
		Handler: engine,
	}

	defer interrupts.WaitForGracefulShutdown()
	interrupts.ListenAndServe(srv, time.Second*30)

	engine.UseRawPath = true
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()

		errmsg := ""
		for _, ginErr := range c.Errors {
			if errmsg != "" {
				errmsg += ","
			}
			errmsg = fmt.Sprintf("%s%s", errmsg, ginErr.Error())
		}

		if strings.Contains(c.Request.RequestURI, "/swagger/") ||
			strings.Contains(c.Request.RequestURI, "/internal/heartbeat") {
			return
		}

		log := fmt.Sprintf(
			"| %d | %d | %s | %s ",
			c.Writer.Status(),
			endTime.Sub(startTime),
			c.Request.Method,
			c.Request.RequestURI,
		)
		if errmsg != "" {
			log += fmt.Sprintf("| %s ", errmsg)
		}

		logrus.Info(log)
	}
}
