package utils

import (
	"WB_TEST/config"
	"github.com/gin-gonic/gin"
)

func ReadBody(ctx *gin.Context, model interface{}) error {
	err := ctx.Bind(model)
	if err != nil {
		return err
	}

	return nil
}

func SendJsonResponse(c *gin.Context, code int, message interface{}) {
	if code == 200 || code == 201 {
		message = map[string]interface{}{"result": message}
	} else {
		message = map[string]interface{}{"error": message}
	}

	c.JSON(code, message)
}

func CreateSessionCookie(cfg *config.Config, session string) (
	name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool) {

	return "refreshToken", session, cfg.RefreshTokenExpire, "/", "", false, true
}
