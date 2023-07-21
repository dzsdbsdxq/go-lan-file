package model

type Json map[string]interface{}
type Head struct {
	Seq string // 消息的Id
	Cmd string // 消息的cmd 动作
}
