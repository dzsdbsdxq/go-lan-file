package websocket

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"share.ac.cn/cache"
	"share.ac.cn/model"
	"time"
)

// UserList 查询所有用户
func UserList(roomId string) (userList []string) {

	userList = make([]string, 0)
	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerAll(currentTime)
	if err != nil {
		fmt.Println("给全体用户发消息", err)

		return
	}

	for _, server := range servers {
		var (
			list []string
		)
		if IsLocal(server) {
			list = GetUserList(roomId)
		}
		userList = append(userList, list...)
	}

	return
}

// CheckUserOnline 查询用户是否在线
func CheckUserOnline(roomId string, userId string) (online bool) {
	// 全平台查询
	if roomId == "" {
		for _, roomId := range GetRoomIds() {
			online, _ = checkUserOnline(roomId, userId)
			if online == true {
				break
			}
		}
	} else {
		online, _ = checkUserOnline(roomId, userId)
	}

	return
}

// 查询用户 是否在线
func checkUserOnline(roomId string, userId string) (online bool, err error) {
	key := GetUserKey(roomId, userId)
	userOnline, err := cache.GetUserOnlineInfo(key)
	if err != nil {
		if err == redis.Nil {
			fmt.Println("GetUserOnlineInfo", roomId, userId, err)

			return false, nil
		}

		fmt.Println("GetUserOnlineInfo", roomId, userId, err)

		return
	}

	online = userOnline.IsOnline()

	return
}

// SendUserMessage 给用户发送消息
func SendUserMessage(roomId string, userId string, msgId, message string) (sendResults bool, err error) {

	data := model.GetTextMsgData(userId, msgId, message)

	client := GetUserClient(roomId, userId)

	if client != nil {
		// 在本机发送
		sendResults, err = SendUserMessageLocal(roomId, userId, data)
		if err != nil {
			fmt.Println("给用户发送消息", roomId, userId, err)
		}

		return
	}
	sendResults = true

	return
}

// SendUserMessageLocal 给本机用户发送消息
func SendUserMessageLocal(roomId string, userId string, data string) (sendResults bool, err error) {

	client := GetUserClient(roomId, userId)
	if client == nil {
		err = errors.New("用户不在线")

		return
	}

	// 发送消息
	client.SendMsg([]byte(data))
	sendResults = true

	return
}

// SendUserMessageAll 给全体用户发消息
func SendUserMessageAll(roomId string, userId string, msgId, cmd, message string) (sendResults bool, err error) {
	sendResults = true

	currentTime := uint64(time.Now().Unix())
	servers, err := cache.GetServerAll(currentTime)
	if err != nil {
		fmt.Println("给全体用户发消息", err)
		return
	}
	for _, server := range servers {
		if IsLocal(server) {
			data := model.GetMsgData(userId, msgId, cmd, message)
			AllSendMessages(roomId, userId, data)
		}
	}

	return
}
