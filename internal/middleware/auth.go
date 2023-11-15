package middleware

import (
	"WB_TEST/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (mc *MiddlewareChief) AuthMiddleware(c *gin.Context) {
	bearerHeader := c.Request.Header.Get("Authorization")
	if bearerHeader != "" {

		headerParts := strings.Split(bearerHeader, " ")
		if len(headerParts) != 2 {
			mc.logger.Error("AuthMiddleware: ", "len(headerParts) != 2")
			utils.SendJsonResponse(c, http.StatusUnauthorized, "Your request contains incorrect data")
			c.Abort()
			return
		}

		jwt := headerParts[1]

		userID, err := utils.JWTValidate(jwt, mc.cfg)
		if err == utils.JWTIsDied {
			if c.Request.RequestURI == "/auth/refresh" {
				c.Next()
				return
			}
		}
		if err != nil {
			mc.logger.Error("AuthMiddleware: ", "jwt is died", err.Error())
			utils.SendJsonResponse(c, http.StatusUnauthorized, "Your jwt is done")
			c.Abort()
			return
		}

		c.Set("userID", userID)

		return
	}

	utils.SendJsonResponse(c, http.StatusUnauthorized, "where is jwt")
	c.Abort()
}
