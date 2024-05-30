package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gpt-content-audit/check"
	"gpt-content-audit/common"
	"gpt-content-audit/common/config"
	logger "gpt-content-audit/common/loggger"
	"gpt-content-audit/middleware"
	"gpt-content-audit/router"
	"os"
	"strconv"

	"github.com/gin-contrib/sessions/cookie"
)

func main() {
	logger.SetupLogger()
	logger.SysLog(fmt.Sprintf("GPT Content Audit %s started", common.Version))

	check.CheckEnvVariable()

	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	if config.DebugEnabled {
		logger.SysLog("running in debug mode")
	}
	var err error

	server := gin.New()
	server.Use(gin.Recovery())
	server.Use(middleware.RequestId())
	middleware.SetUpLogger(server)
	store := cookie.NewStore([]byte(config.SessionSecret))
	server.Use(sessions.Sessions("session", store))

	router.SetRouter(server)
	var port = os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(*common.Port)
	}
	err = server.Run(":" + port)
	if err != nil {
		logger.FatalLog("failed to start HTTP server: " + err.Error())
	}
}
