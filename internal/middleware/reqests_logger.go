package middleware

import (
	"github.com/gin-gonic/gin"
)

func (mc *MiddlewareChief) RequestLogMiddleware(c *gin.Context) {
	c.Next()
	mc.logger.Info("New request: ", "METHOD: ", c.Request.Method, "URI: ",
		c.Request.RequestURI, "BODY: ", c.Request.Body)
}
