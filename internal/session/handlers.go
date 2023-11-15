package session

import "github.com/gin-gonic/gin"

type Handlers interface {
	RefreshTokens(c *gin.Context)
}
