package cache

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"share.ac.cn/common"
	"share.ac.cn/model"
)

const (
	roomContentOnlinePrefix    = "acc:room:content:online:" // 房间在线状态
	roomContentOnlineCacheTime = 24 * 60 * 60
)

func GetRoomContentOnlinePrefix() string {
	return roomContentOnlinePrefix
}

func getRoomContentOnlineKey(roomId string) (key string) {
	key = fmt.Sprintf("%s%s", roomContentOnlinePrefix, roomId)
	return
}

func GetRoomContentOnlineInfo(roomId string) (roomContentOnline *model.RoomContentOnline, err error) {
	key := getRoomContentOnlineKey(roomId)
	data, err := common.GetClient().Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("GetRoomContentOnlineInfo", roomId, err)
			return
		}
		return
	}

	roomContentOnline = &model.RoomContentOnline{}
	err = json.Unmarshal(data, roomContentOnline)
	if err != nil {
		fmt.Println("获取房间内容在线数据 json Unmarshal", roomId, err)
		return
	}
	return
}

// SetRoomContent 设置房间信息
func SetRoomContent(roomId string, roomContentOnline *model.RoomContentOnline) (err error) {
	key := getRoomContentOnlineKey(roomId)
	valueByte, err := json.Marshal(roomContentOnline)

	if err != nil {
		fmt.Println("设置房间内容在线数据 json Marshal", key, err)

		return
	}
	_, err = common.GetClient().Do("setEx", key, roomContentOnlineCacheTime, string(valueByte)).Result()
	if err != nil {
		fmt.Println("设置房间内容在线数据 ", key, err)
		return
	}
	return
}
