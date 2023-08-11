package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"share.ac.cn/common"
)

const tmpOnlinePrefix = "acc:tmp:online:"

type TmpOnline struct {
	ShareId  string
	FileName string
	FileKey  string
	FileSize uint64
}

func GetTmpOnlinePrefix() string {
	return tmpOnlinePrefix
}

func getTmpOnlineKey(userId string) (key string) {
	key = fmt.Sprintf("%s%s", tmpOnlinePrefix, userId)
	return
}

func GetTmpOnlineInfo(userId string) (fileOnline *TmpOnline, err error) {
	return _getTmpOnlineInfo(getTmpOnlineKey(userId))
}

func _getTmpOnlineInfo(key string) (fileOnline *TmpOnline, err error) {
	data, err := common.GetClient().Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			common.Log.Errorf("GetFileOnlineInfo: %s, 错误信息: %s", key, err)
			return nil, err
		}
		return nil, err
	}

	fileOnline = &TmpOnline{}
	err = jsoniter.Unmarshal(data, fileOnline)
	if err != nil {
		common.Log.Errorf("获取房间内容在线数据 json Unmarshal: %s, 错误信息: %s", key, err)
		return nil, err
	}
	return
}

// SetTmpOnline 设置文件信息
func SetTmpOnline(userId string, fileOnline *TmpOnline) (err error) {
	key := getTmpOnlineKey(userId)
	valueByte, err := jsoniter.Marshal(fileOnline)
	if err != nil {
		common.Log.Infof("设置文件内容在线数据 : %s, 错误信息: %s", userId, err)
		return
	}
	_, err = common.GetClient().Do("setEx", key, 1800, string(valueByte)).Result()
	if err != nil {
		common.Log.Infof("设置房间内容在线数据 : %s, 错误信息: %s", key, err)
		return
	}
	return
}
