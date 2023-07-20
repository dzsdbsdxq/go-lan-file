package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"share.ac.cn/common"
	"share.ac.cn/model"
)

const (
	roomOnlinePrefix    = "acc:room:online:" // 房间在线状态
	roomOnlineCacheTime = 24 * 60 * 60
)

func GetRoomOnlinePrefix() string {
	return roomOnlinePrefix
}

func getRoomOnlineKey(roomId string) (key string) {
	key = fmt.Sprintf("%s%s", roomOnlinePrefix, roomId)
	return
}
func GetRoomOnlineInfo(roomId string) (roomOnline *model.RoomOnline, err error) {
	key := getRoomOnlineKey(roomId)
	data, err := common.GetClient().Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("GetRoomOnlineInfo", roomId, err)
			return
		}
		return
	}

	roomOnline = &model.RoomOnline{}
	err = json.Unmarshal(data, roomOnline)
	if err != nil {
		fmt.Println("获取房间在线数据 json Unmarshal", roomId, err)

		return
	}
	return
}

// SetRoomInfo 设置房间信息
func SetRoomInfo(roomId string, roomOnline *model.RoomOnline) (err error) {
	key := getRoomOnlineKey(roomId)
	valueByte, err := json.Marshal(roomOnline)

	if err != nil {
		fmt.Println("设置房间在线数据 json Marshal", key, err)

		return
	}
	_, err = common.GetClient().Do("setEx", key, roomOnlineCacheTime, string(valueByte)).Result()
	if err != nil {
		fmt.Println("设置房间在线数据 ", key, err)

		return
	}
	return
}
