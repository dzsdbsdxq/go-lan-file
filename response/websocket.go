package response

import "encoding/json"

type Head struct {
	Id       string      `json:"Id"`
	Group    string      `json:"Group"`
	Seq      string      `json:"Seq"`      // 消息的Id
	Cmd      string      `json:"Cmd"`      // 消息的cmd 动作
	Response *WSResponse `json:"Response"` // 消息体
}
type WSResponse struct {
	Code    uint32      `json:"code"`
	CodeMsg string      `json:"codeMsg"`
	Data    interface{} `json:"data"` // 数据 json
}

type PushMsg struct {
	Seq  string `json:"seq"`
	Uuid uint64 `json:"uuid"`
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

// NewResponseHead 设置返回消息
func NewResponseHead(id, group, seq, cmd string, code uint32, codeMsg string, data interface{}) *Head {
	response := NewResponse(code, codeMsg, data)

	return &Head{Id: id, Group: group, Seq: seq, Cmd: cmd, Response: response}
}

func (h *Head) String() (headStr string) {
	headBytes, _ := json.Marshal(h)
	headStr = string(headBytes)
	return
}

func NewResponse(code uint32, codeMsg string, data interface{}) *WSResponse {
	return &WSResponse{Code: code, CodeMsg: codeMsg, Data: data}
}
