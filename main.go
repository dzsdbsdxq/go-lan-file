package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"share.ac.cn/common"
	"share.ac.cn/common/tusd"
	"share.ac.cn/config"
	"share.ac.cn/routers"
	"share.ac.cn/services/task"
	"share.ac.cn/services/uploader"
	"share.ac.cn/services/websocket/service"
	"syscall"
	"time"
)

func main() {
	// 加载配置文件到全局配置结构体
	config.InitConfig()
	// 初始化日志
	common.InitLogger()
	// 初始化Validator数据校验
	common.InitValidate()
	//初始化sqlite3数据库
	common.InitDataBase()
	//初始化Redis
	common.InitRedisClient()

	//初始化七牛云
	uploader.InitQiNiuMac()
	// 定时任务
	task.CleanExpireFileInit()
	//初始化Tusd服务
	tusd.NewTusdServer()

	//go websocket.StartWebSocket()
	go service.StartWebSocket()

	// 注册所有路由
	r := routers.InitRoutes()

	host := config.Conf.System.Host
	port := config.Conf.System.Port

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			common.Log.Fatalf("listen: %s\n", err)
		}
	}()

	common.Log.Info(fmt.Sprintf("Server is running at %s:%d", host, port))

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	common.Log.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		common.Log.Fatal("Server forced to shutdown:", err)
	}

	common.Log.Info("Server exiting!")

}
