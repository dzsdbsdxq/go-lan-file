package routers

import (
	"github.com/gin-gonic/gin"
	"share.ac.cn/common"
	"share.ac.cn/config"
	"share.ac.cn/middleware"
	"time"
)

func InitRoutes() *gin.Engine {
	//设置模式
	gin.SetMode(config.Conf.System.Mode)

	// 创建带有默认中间件的路由:
	// 日志与恢复中间件
	r := gin.Default()

	// 启用限流中间件
	// 默认每50毫秒填充一个令牌，最多填充200个
	fillInterval := time.Duration(config.Conf.RateLimit.FillInterval)
	capacity := config.Conf.RateLimit.Capacity
	r.Use(middleware.RateLimitMiddleware(time.Millisecond*fillInterval, capacity))

	// 启用全局跨域中间件
	r.Use(middleware.CORSMiddleware())

	// 路由分组
	apiGroup := r.Group("/")

	// 注册路由
	//InitTemplatesRoutes(r)
	//InitBaseRoutes(apiGroup, authMiddleware)              // 注册基础路由, 不需要jwt认证中间件,不需要casbin中间件
	InitFileRouters(apiGroup) // 注册文件路由, casbin鉴权中间件

	common.Log.Info("初始化路由完成！")
	return r
}
