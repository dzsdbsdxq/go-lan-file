package routers

import (
	"github.com/gin-gonic/gin"
	"share.ac.cn/controller"
)

func InitApiRouters(r *gin.RouterGroup) gin.IRoutes {
	apiController := controller.ApiController{}
	router := r.Group("/api")
	{

		router.POST("/newPcConnect", apiController.NewConnect)
		//router.HEAD("/:id", gin.WrapF(handler.HeadFile))
		//router.PATCH("/:id", gin.WrapF(handler.PatchFile))
		//router.GET("/:id", gin.WrapF(handler.GetFile))
	}

	return r
}
