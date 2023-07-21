package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"runtime/debug"
	"share.ac.cn/common"
	"share.ac.cn/model"
	"time"
)

const heartbeatExpirationTime = 3 * 60

// Client 单个 websocket 信息
type Client struct {
	Id            string
	Group         string
	Ctx           *gin.Context
	Addr          string    //客户端地址
	FirstTime     time.Time //首次连接时间
	HeartBeatTime time.Time //用户上次心跳时间
	LoginTime     time.Time //登录时间，登录以后才有
	Socket        *websocket.Conn
	Message       chan []byte
	User          *model.UserOnline //用户信息，用户登录以后才有
}

func NewClient(ctx *gin.Context, group string, addr string, socket *websocket.Conn) *Client {
	return &Client{
		Id:            uuid.NewV4().String(),
		Ctx:           ctx,
		Group:         group,
		Addr:          addr,
		FirstTime:     time.Now(),
		HeartBeatTime: time.Now(),
		LoginTime:     time.Now(),
		Socket:        socket,
		Message:       make(chan []byte, 1024),
		User:          nil,
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
			common.Log.Errorf("client [%s] disconnect err: %s", c.Id, err)
		}
	}()

	for {
		messageType, message, err := c.Socket.ReadMessage()
		if err != nil || messageType == websocket.CloseMessage {
			common.Log.Error("读取客户端数据错误:", c.Addr, err)
			break
		}
		common.Log.Infof("client [%s] receive message: %s", c.Id, string(message))
		//c.Message <- message
		//处理消息
		ProcessData(c, message)
	}
}

// 写信息，从 channel 变量 Send 中读取数据写入 websocket 连接
func (c *Client) Write() {
	defer func() {
		if r := recover(); r != nil {
			common.Log.Error("write stop:", string(debug.Stack()), r)
		}
	}()
	defer func() {
		common.Log.Infof("client [%s] disconnect", c.Id)
		clientManager.UnRegister <- c
		if err := c.Socket.Close(); err != nil {
			common.Log.Infof("client [%s] disconnect err: %s", c.Id, err)
		}
	}()

	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			common.Log.Infof("client [%s] write message: %s", c.Id, message)
			err := c.Socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				common.Log.Errorf("client [%s] write message err: %s", c.Id, err)
			}
		}
	}
}

// Login 用户登录
func (c *Client) Login(user *model.UserOnline) {
	c.User = user
	c.LoginTime = time.Now()
	// 登录成功=心跳一次
	c.HeartBeat()
}

// HeartBeat 用户心跳
func (c *Client) HeartBeat() {
	c.HeartBeatTime = time.Now()
}

// IsHeartbeatTimeout 心跳超时
func (c *Client) IsHeartbeatTimeout() bool {
	if c.HeartBeatTime.Before(time.Now().Add(-heartbeatExpirationTime * time.Second)) {
		common.Log.Infof("连接客户端[%s],心跳超时", c.Id)
		return true
	}
	return false
}

// IsLogin 是否登录了
func (c *Client) IsLogin() (isLogin bool) {
	// 用户登录了
	if c.User != nil {
		isLogin = true
		return
	}
	return
}
