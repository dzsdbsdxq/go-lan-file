package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"share.ac.cn/common"
	"share.ac.cn/model"
)

const (
	userOnlinePrefix    = "acc:user:online:" // 用户在线状态
	userOnlineCacheTime = 24 * 60 * 60
)

/*********************  查询用户是否在线  ************************/
func getUserOnlineKey(userKey string) (key string) {
	key = fmt.Sprintf("%s%s", userOnlinePrefix, userKey)
	return
}

func GetUserOnlineInfo(userKey string) (userOnline *model.UserOnline, err error) {

	key := getUserOnlineKey(userKey)

	data, err := common.GetClient().Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("GetUserOnlineInfo", userKey, err)

			return
		}

		fmt.Println("GetUserOnlineInfo", userKey, err)

		return
	}

	userOnline = &model.UserOnline{}
	err = json.Unmarshal(data, userOnline)
	if err != nil {
		fmt.Println("获取用户在线数据 json Unmarshal", userKey, err)

		return
	}

	return
}

// SetUserOnlineInfo 设置用户在线数据
func SetUserOnlineInfo(userKey string, userOnline *model.UserOnline) (err error) {

	key := getUserOnlineKey(userKey)

	valueByte, err := json.Marshal(userOnline)
	if err != nil {
		common.Log.Info("设置用户在线数据 json Marshal", key, err)
		return
	}

	_, err = common.GetClient().Do("setEx", key, userOnlineCacheTime, string(valueByte)).Result()
	if err != nil {
		common.Log.Info("设置用户在线数据 ", key, err)
		return
	}

	return
}
