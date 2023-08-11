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
		router.POST("/file/create/:userId", apiController.CreateUpload)
		router.POST("/callback", apiController.CallBack)
		router.GET("/file/:shareId", apiController.GetFileInfo)
		router.GET("/delete/:shareId", apiController.DeleteObject)
		router.POST("/file/complete/:userId", apiController.GetFileComplete)
	}

	return r
}
