package routers

import (
	"github.com/gin-gonic/gin"
	"share.ac.cn/controller"
)

func InitWebRouters(r *gin.Engine) *gin.Engine {
	r.Static("css", "web/css")
	r.Static("js", "web/js")
	r.Static("fonts", "web/fonts")
	r.Static("img", "web/img")
	r.LoadHTMLGlob("./web/*.htm")
	r.GET("/", controller.IndexPage)
	r.GET("/:code", controller.IndexPage)
	return r
}
