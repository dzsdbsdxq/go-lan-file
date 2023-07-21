package routers

import (
	"github.com/gin-gonic/gin"
	"share.ac.cn/services/websocket/service"
)

func InitWebSocketRouters(r *gin.RouterGroup) gin.IRoutes {
	r.GET("/acc", service.WebSocketFunc)
	return r
}
