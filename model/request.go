package model

// Request 通用请求数据格式
type Request struct {
	Seq  string      `json:"seq"`            //消息唯一ID
	Cmd  string      `json:"cmd"`            //请求操作指令
	Data interface{} `json:"data,omitempty"` //数据json
}

// Login 登录请求数据
type Login struct {
	Token  string `json:"token"` //验证用户是否登录
	RoomId string `json:"room_id"`
	UserId string `json:"user_id"`
}

// HeartBeat 心跳请求数据
type HeartBeat struct {
	UserId string `json:"user_id,omitempty"`
}
