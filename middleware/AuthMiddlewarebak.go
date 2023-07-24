package middleware

import (
	"github.com/gin-gonic/gin"
	"share.ac.cn/services/websocket/service"
)

func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.Abort()
		return
	}
	client, ok := service.LoginUserManager.GetUserClient(token)
	if !ok {
		c.Abort()
		return
	}
	//用户未登陆
	if !client.User.OnLine {
		c.Abort()
		return
	}
	c.Next()
}
