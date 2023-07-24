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
	files := cache.GetFileOnlineAll()
	for _, i2 := range files {
		value := i2
		//判断file是否过期
		if time.Now().After(value.ExpireTime) {
			//文件已过期,删除文件,标记文件
			//删除文件
			err := common.DeleteFile(value.FilePath)
			if err != nil {
				common.Log.Errorf("删除文件出错：文件ID：%s，错误信息%s", value.FileId, err.Error())
				continue
			}
			_ = common.DeleteFile(value.FilePath + ".info")
			cache.DeleteFileOnline(value.ShareId)
		}
	}

	return
}
