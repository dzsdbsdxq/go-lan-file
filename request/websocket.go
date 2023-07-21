package request

// Request Websocket通用请求数据格式
type Request struct {
	Seq  string      `json:"seq"`            //消息唯一ID
	Cmd  string      `json:"cmd"`            //请求操作指令
	Data interface{} `json:"data,omitempty"` //数据json
}

// Login 登录请求的数据
type Login struct {
	UserId string `json:"user_id"` //通过NewConnect接口获取到的用户ID
}

// HeartBeat 心跳请求数据
type HeartBeat struct {
	UserId string `json:"user_id,omitempty"`
}
