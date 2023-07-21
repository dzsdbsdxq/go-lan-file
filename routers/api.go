package routers

import (
	"github.com/gin-gonic/gin"
	"share.ac.cn/controller"
)

func InitApiRouters(r *gin.RouterGroup) gin.IRoutes {
	apiController := controller.NewApiController()
	router := r.Group("/api")
	{

		router.POST("/newConnect", apiController.NewConnect)
		//router.HEAD("/:id", gin.WrapF(handler.HeadFile))
		//router.PATCH("/:id", gin.WrapF(handler.PatchFile))
		//router.GET("/:id", gin.WrapF(handler.GetFile))
	}

	return r
}
