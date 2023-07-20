package routers

import (
	"github.com/gin-gonic/gin"
	"share.ac.cn/common"
)

func InitTusdRouters(r *gin.RouterGroup) gin.IRoutes {
	// 启用全局跨域中间件
	handler := common.GetTusd()
	router := r.Group("/files")
	{

		router.POST("/", gin.WrapF(handler.PostFile))
		router.HEAD("/:id", gin.WrapF(handler.HeadFile))
		router.PATCH("/:id", gin.WrapF(handler.PatchFile))
		router.GET("/:id", gin.WrapF(handler.GetFile))
	}

	return r
}
