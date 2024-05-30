package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"gpt-content-audit/common/config"
	"gpt-content-audit/model"
	"net/http"
	"strings"
)

func isValidSecret(secret string) bool {
	return config.Authorization != "" && !lo.Contains(config.Authorizations, secret)
}

func authHelper(c *gin.Context) {
	secret := c.Request.Header.Get("Authorization")
	secret = strings.Replace(secret, "Bearer ", "", 1)
	if isValidSecret(secret) {
		c.JSON(http.StatusUnauthorized, model.OpenAIErrorResponse{
			OpenAIError: model.OpenAIError{
				Message: "Auth failed",
				Type:    "invalid_request_error",
				Code:    "INVALID_AUTHORIZATION",
			},
		})
		c.Abort()
		return
	}
	c.Next()
	return
}

func Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelper(c)
	}
}
