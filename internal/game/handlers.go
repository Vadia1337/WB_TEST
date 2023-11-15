package game

import "github.com/gin-gonic/gin"

type Handlers interface {
	UserTasks(c *gin.Context)
	CreateTasks(c *gin.Context)
	UserInfo(c *gin.Context)
	StartGame(c *gin.Context)
}
