package routers

import (
	"github.com/gin-gonic/gin"
	"share.ac.cn/controller"
)

func InitFileRouters(r *gin.RouterGroup) gin.IRoutes {
	fileController := controller.NewFileController()
	router := r.Group("/files")
	{
		router.POST("/", fileController.PostFile)
		router.HEAD("/:id", fileController.HeadFile)
		router.PATCH("/:id", fileController.PatchFile)
		router.GET("/:id", fileController.GetFile)
	}
	return r
}
