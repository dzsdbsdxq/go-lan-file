package routers

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"share.ac.cn/controller"
	"share.ac.cn/middleware"
)

func InitFileRouters(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	fileController := controller.NewFileController()
	router := r.Group("/files")
	{
		router.POST("/", middleware.AuthMiddleware, fileController.PostFile)
		router.HEAD("/:id", middleware.AuthMiddleware, fileController.HeadFile)
		router.PATCH("/:id", middleware.AuthMiddleware, fileController.PatchFile)
		router.GET("/:id", fileController.GetFile)
	}
	return r
}
