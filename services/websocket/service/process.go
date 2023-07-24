package service

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"share.ac.cn/common"
	rqs "share.ac.cn/request"
	"share.ac.cn/response"
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
	//defer func() {
	//	if r := recover(); r != nil {
	//		fmt.Println("处理数据 stop", r)
	//	}
	//}()
	request := &rqs.Request{}

	err := jsoniter.Unmarshal(message, request)
	if err != nil {
		fmt.Println("数据处理 json Unmarshal:", err)
		client.Message <- []byte("数据不合法")
		return
	}

	requestData, err := jsoniter.Marshal(request.Data)
	if err != nil {
		fmt.Println("数据处理 json Marshal", err)
		client.Message <- []byte("处理数据失败")
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
	responseHead := response.NewResponseHead(client.Id, client.Group, seq, cmd, code, msg, data)
	headByte, err := jsoniter.Marshal(responseHead)
	if err != nil {
		fmt.Println("处理数据 json Marshal", err)
		return
	}
	client.Message <- headByte
	return
}
