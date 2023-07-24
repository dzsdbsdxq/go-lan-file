package routers

import (
	"fmt"
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

	// 初始化JWT认证中间件
	authMiddleware, err := middleware.InitAuth()
	if err != nil {
		common.Log.Panicf("初始化JWT中间件失败：%v", err)
		panic(fmt.Sprintf("初始化JWT中间件失败：%v", err))
	}

	InitWebRouters(r)
	// 路由分组
	apiGroup := r.Group("/")
	// 注册路由
	InitFileRouters(apiGroup, authMiddleware) // 注册文件路由
	InitApiRouters(apiGroup, authMiddleware)  //注册api路由
	InitWebSocketRouters(apiGroup)            //注册websocket路由

	common.Log.Info("初始化路由完成！")
	return r
}
