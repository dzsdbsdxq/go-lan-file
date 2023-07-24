package common

import (
	"errors"
	"fmt"
	"os"
)

func DeleteFile(filePath string) error {
	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("文件不存在")
		}
		return errors.New(fmt.Sprintf("无法获取文件信息：%s\n", err))
	}
	// 检查文件权限
	fileInfo, err := os.Stat(filePath)
	if err != nil {

		return errors.New(fmt.Sprintf("无法获取文件信息：%s\n", err))
	}
	if fileInfo.Mode().IsDir() {
		return errors.New(fmt.Sprintf("无法删除目录，请使用 os.Remove 或 os.RemoveAll 来删除目录"))
	}
	// 删除文件
	err = os.Remove(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("无法删除文件：%s\n", err))
	}
	return nil
}
