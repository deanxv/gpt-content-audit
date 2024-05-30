package router

import (
	"github.com/gin-gonic/gin"
	"gpt-content-audit/common/config"
	"gpt-content-audit/controller"
	"gpt-content-audit/middleware"
)

func SetApiRouter(router *gin.Engine) {

	router.Use(middleware.CORS())

	if config.Enable == 1 {
		router.Use(middleware.RequestRateLimit())
		v1Router := router.Group("/v1")
		v1Router.Use(middleware.Auth())
		v1Router.POST("/chat/completions", controller.ChatForOpenAI)
		v1Router.OPTIONS("/chat/completions", controller.ChatForOpenAI)
		v1Router.POST("/images/generations", controller.ImagesForOpenAI)
	}

	router.NoRoute(func(c *gin.Context) {
		middleware.ForwardTo(c, config.BaseUrl)
	})
}
