package service

import (
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"runtime/debug"
	"share.ac.cn/common"
	"time"
)

// Client 单个 websocket 信息
type Client struct {
	Id            string
	Group         string
	Addr          string    //客户端地址
	FirstTime     time.Time //首次连接时间
	HeartBeatTime time.Time //用户上次心跳时间
	LoginTime     time.Time //登录时间，登录以后才有
	Socket        *websocket.Conn
	Message       chan []byte
}

func NewClient(group string, addr string, socket *websocket.Conn) *Client {
	return &Client{
		Id:            uuid.NewV4().String(),
		Group:         group,
		Addr:          addr,
		FirstTime:     time.Now(),
		HeartBeatTime: time.Now(),
		LoginTime:     time.Now(),
		Socket:        socket,
		Message:       make(chan []byte, 1024),
	}
}

// 读信息，从 websocket 连接直接读取数据
func (c *Client) Read() {
	defer func() {
		if r := recover(); r != nil {
			common.Log.Info("write stop", string(debug.Stack()), r)
		}
	}()
	defer func() {
		clientManager.UnRegister <- c
		common.Log.Infof("client [%s] disconnect", c.Id)
		if err := c.Socket.Close(); err != nil {
			common.Log.Error("client [%s] disconnect err: %s", c.Id, err)
		}
	}()

	for {
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			common.Log.Error("读取客户端数据错误:", c.Addr, err)
			break
		}
		common.Log.Infof("client [%s] receive message: %s", c.Id, string(message))
		c.Message <- message
	}
}

// 写信息，从 channel 变量 Send 中读取数据写入 websocket 连接
func (c *Client) Write() {
	defer func() {
		if r := recover(); r != nil {
			common.Log.Error("write stop", string(debug.Stack()), r)
		}
	}()
	defer func() {
		common.Log.Info("client [%s] disconnect", c.Id)
		clientManager.UnRegister <- c
		if err := c.Socket.Close(); err != nil {
			common.Log.Info("client [%s] disconnect err: %s", c.Id, err)
		}
	}()

	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			common.Log.Info("client [%s] write message: %s", c.Id, string(message))
			err := c.Socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				common.Log.Error("client [%s] write message err: %s", c.Id, err)
			}
		}
	}
}
