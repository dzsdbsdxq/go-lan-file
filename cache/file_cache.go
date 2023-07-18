package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"share.ac.cn/common"
	"share.ac.cn/model"
)

const (
	fileOnlinePrefix = "acc:file:online:" // 文件在线状态
)

func GetFileOnlinePrefix() string {
	return fileOnlinePrefix
}

func getFileOnlineKey(shareId string) (key string) {
	key = fmt.Sprintf("%s%s", fileOnlinePrefix, shareId)
	return
}
func GetFileOnlineInfo(shareId string) (fileOnline *model.FileOnline, err error) {
	return _getFileOnlineInfo(getFileOnlineKey(shareId))
}

func _getFileOnlineInfo(key string) (fileOnline *model.FileOnline, err error) {
	data, err := common.GetClient().Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			common.Log.Errorf("GetFileOnlineInfo: %s, 错误信息: %s", key, err)
			return nil, err
		}
		return nil, err
	}

	fileOnline = &model.FileOnline{}
	err = json.Unmarshal(data, fileOnline)
	if err != nil {
		common.Log.Errorf("获取房间内容在线数据 json Unmarshal: %s, 错误信息: %s", key, err)
		return nil, err
	}
	return
}

// SetFileOnline 设置文件信息
func SetFileOnline(shareId string, fileOnline *model.FileOnline) (err error) {
	key := getFileOnlineKey(shareId)
	valueByte, err := json.Marshal(fileOnline)
	if err != nil {
		common.Log.Infof("设置文件内容在线数据 : %s, 错误信息: %s", shareId, err)
		return
	}
	_, err = common.GetClient().Set(key, string(valueByte), 0).Result()
	if err != nil {
		common.Log.Infof("设置房间内容在线数据 : %s, 错误信息: %s", key, err)
		return
	}
	return
}
func GetFileOnlineAll() {
	var files = make([]*model.FileOnline, 0)
	key := GetFileOnlinePrefix()
	iter := common.GetClient().Scan(0, key+"*", 0).Iterator()
	for iter.Next() {
		info, err := _getFileOnlineInfo(iter.Val())
		if err != nil {
			continue
		}
		files = append(files, info)
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	return

}
