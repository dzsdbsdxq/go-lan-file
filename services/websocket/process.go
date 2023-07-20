package websocket

import (
	"fmt"
	"github.com/goccy/go-json"
	"share.ac.cn/common"
	"share.ac.cn/model"
	"sync"
)

type DisPoseFunc func(client *Client, seq string, message []byte) (code uint32, msg string, data interface{})

var (
	handlers        = make(map[string]DisPoseFunc)
	handlersRWMutex sync.RWMutex
)

func Register(key string, value DisPoseFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value
}

func getHandlers(key string) (value DisPoseFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()
	value, ok = handlers[key]
	return
}

func ProcessData(client *Client, message []byte) {
	fmt.Println("处理数据:", client.Addr, string(message))
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("处理数据 stop", r)
		}
	}()
	request := &model.Request{}

	err := json.Unmarshal(message, request)
	if err != nil {
		fmt.Println("数据处理 json Unmarshal:", err)
		client.SendMsg([]byte("数据不合法"))
		return
	}

	requestData, err := json.Marshal(request.Data)
	if err != nil {
		fmt.Println("数据处理 json Marshal", err)
		client.SendMsg([]byte("处理数据失败"))
		return
	}
	seq := request.Seq
	cmd := request.Cmd

	var (
		code uint32
		msg  string
		data interface{}
	)

	fmt.Println("request", cmd, client.Addr)
	//采用map注册方式
	if value, ok := getHandlers(cmd); ok {
		code, msg, data = value(client, seq, requestData)
	} else {
		code = common.RoutingNotExist
		fmt.Println("处理数据，路由不存在", client.Addr, "cmd", cmd)
	}

	msg = common.GetErrorMessage(code, msg)
	responseHead := model.NewResponseHead(seq, cmd, code, msg, data)
	headByte, err := json.Marshal(responseHead)
	if err != nil {
		fmt.Println("处理数据 json Marshal", err)

		return
	}

	client.SendMsg(headByte)

	fmt.Println("acc_response send", client.Addr, client.RoomId, client.UserId, "cmd", cmd, "code", code)

	return
}
