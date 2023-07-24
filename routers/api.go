package routers

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"share.ac.cn/controller"
)

func InitApiRouters(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	apiController := controller.NewApiController()
	router := r.Group("/api")
	{

		router.POST("/newConnect", apiController.NewConnect)
		router.GET("/file/:shareId", apiController.GetFileInfo)
		router.GET("/complete/:fileId", apiController.GetFileComplete)
	}

	return r
}
