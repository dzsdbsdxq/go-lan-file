package cache

import (
	"encoding/json"
	"fmt"
	"share.ac.cn/common"
	"share.ac.cn/model"
	"strconv"
)

const (
	serversHashKey       = "acc:hash:servers" // 全部的服务器
	serversHashCacheTime = 2 * 60 * 60        // key过期时间
	serversHashTimeout   = 3 * 60             // 超时时间
)

func getServersHashKey() (key string) {
	key = fmt.Sprintf("%s", serversHashKey)

	return
}

// SetServerInfo 设置服务器信息
func SetServerInfo(server *model.Server, currentTime uint64) (err error) {
	key := getServersHashKey()

	value := fmt.Sprintf("%d", currentTime)

	number, err := common.GetClient().Do("hSet", key, server.String(), value).Int()
	if err != nil {
		fmt.Println("SetServerInfo", key, number, err)

		return
	}
	if number != 1 {

		return
	}
	common.GetClient().Do("Expire", key, serversHashCacheTime)

	return
}

// DelServerInfo 下线服务器信息
func DelServerInfo(server *model.Server) (err error) {
	key := getServersHashKey()
	number, err := common.GetClient().Do("hDel", key, server.String()).Int()
	if err != nil {
		fmt.Println("DelServerInfo", key, number, err)

		return
	}

	if number != 1 {

		return
	}

	common.GetClient().Do("Expire", key, serversHashCacheTime)

	return
}

func GetServerAll(currentTime uint64) (servers []*model.Server, err error) {

	servers = make([]*model.Server, 0)
	key := getServersHashKey()

	val, err := common.GetClient().Do("hGetAll", key).Result()

	valByte, _ := json.Marshal(val)
	fmt.Println("GetServerAll", key, string(valByte))

	serverMap, err := common.GetClient().HGetAll(key).Result()
	fmt.Println("setServerInfo err:", serverMap, err)
	if err != nil {
		fmt.Println("SetServerInfo", key, err)

		return
	}

	for key, value := range serverMap {
		valueUint64, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			fmt.Println("GetServerAll", key, err)

			return nil, err
		}

		// 超时
		if valueUint64+serversHashTimeout <= currentTime {
			continue
		}

		server, err := model.StringToServer(key)
		if err != nil {
			fmt.Println("GetServerAll", key, err)

			return nil, err
		}

		servers = append(servers, server)
	}

	return
}
