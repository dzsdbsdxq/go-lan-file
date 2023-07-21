package service

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"share.ac.cn/cache"
	"share.ac.cn/common"
	"share.ac.cn/model"
	rqs "share.ac.cn/request"
)

// PingController ping
func PingController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	fmt.Println("webSocket_request ping接口", client.Addr, seq, message)

	data = "pong"

	return
}

// LoginController 用户登录
func LoginController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	request := &rqs.Login{}
	if err := jsoniter.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("用户登录 解析数据失败", seq, err)
		return
	}

	fmt.Println("webSocket_request 用户登录", seq, "ServiceToken", "UserId:", request.UserId)

	// TODO::进行用户权限认证，一般是客户端传入TOKEN，然后检验TOKEN是否合法，通过TOKEN解析出来用户ID
	// 本项目只是演示，所以直接过去客户端传入的用户ID
	//获取缓存中的用户信息
	userInfo, err := cache.GetUserOnlineInfo(request.UserId)
	if err != nil {
		code = common.UnauthorizedUserId
		fmt.Println("用户登录 非法的用户", seq, request.UserId)
		return
	}

	if client.IsLogin() {
		fmt.Println("用户登录 用户已经登录", client.User.UserId, seq)
		code = common.OperationFailure
		return
	}
	userInfo.UserLogin()

	client.Login(userInfo)

	err = cache.SetUserOnlineInfo(request.UserId, userInfo)
	if err != nil {
		code = common.ServerError
		fmt.Println("用户登录 SetUserOnlineInfo", seq, err)

		return
	}

	clientManager.Login <- client

	fmt.Println("用户登录 成功", seq, client.Addr, request.UserId)

	return
}

// HeartbeatController 心跳接口
func HeartbeatController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {

	code = common.OK
	data = model.Json{
		"userId": "",
	}

	request := &model.HeartBeat{}
	if err := jsoniter.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		fmt.Println("心跳接口 解析数据失败", seq, err)
		return
	}

	fmt.Println("webSocket_request 心跳接口", client.User.UserId)

	if !client.IsLogin() {
		fmt.Println("心跳接口 用户未登录", client.User.UserId, seq)
		code = common.NotLoggedIn

		return
	}

	client.HeartBeat()
	client.User.Heartbeat()
	data = model.Json{
		"userId": client.User.UserId,
	}

	return
}
