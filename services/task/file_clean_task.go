package task

import (
	"fmt"
	"runtime/debug"
	"share.ac.cn/cache"
	"share.ac.cn/common"
	"time"
)

func CleanExpireFileInit() {
	Timer(3*time.Second, 30*time.Second, cleanExpireFile, "", nil, nil)

}

// 清理超时文件
func cleanExpireFile(param interface{}) (result bool) {
	result = true

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("cleanExpireFile stop", r, string(debug.Stack()))
		}
	}()
	fmt.Println("定时任务，清理过期文件", param)
	common.Log.Info("定时任务，清理过期文件")
	cache.GetFileOnlineAll()

	return
}
