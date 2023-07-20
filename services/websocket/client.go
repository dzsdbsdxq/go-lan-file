package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"runtime/debug"
)

const heartbeatExpirationTime = 6 * 60

type login struct {
	RoomId string
	UserId string
	Client *Client
}

func (lg *login) GetKey() string {
	return GetUserKey(lg.RoomId, lg.UserId)
}

type Client struct {
	Addr          string          //客户端地址
	Socket        *websocket.Conn //用户连接
	Send          chan []byte     //待发送的数据
	RoomId        string          //登录的房间
	UserId        string          //用户ID，用户登录以后才有
	FirstTime     uint64          //首次连接时间
	HeartbeatTime uint64          //用户上次心跳时间
	LoginTime     uint64          //登录时间，登录以后才有
}

func NewClient(addr string, socket *websocket.Conn, firstTime uint64) (client *Client) {
	client = &Client{
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, 100),
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
	}
	return
}

// GetKey 获取key
func (c *Client) GetKey() string {
	return GetUserKey(c.RoomId, c.UserId)
}

// 读取客户端数据
func (c *Client) read() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write stop", string(debug.Stack()), r)
		}
	}()
	defer func() {
		fmt.Println("读取客户端数据，关闭send", c)
		close(c.Send)
	}()
	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			fmt.Println("读取客户端数据错误:", c.Addr, err)
			return
		}
		fmt.Println("读取客户端数据 处理:", string(message))

		ProcessData(c, message)
	}
}

// 向客户端写数据
func (c *Client) write() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write stop", string(debug.Stack()), r)

		}
	}()

	defer func() {
		clientManager.Unregister <- c
		c.Socket.Close()
		fmt.Println("Client发送数据 defer", c)
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 发送数据错误 关闭连接
				fmt.Println("Client发送数据 关闭连接", c.Addr, "ok", ok)

				return
			}

			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c *Client) SendMsg(msg []byte) {
	if c == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("SendMsg stop:", r, string(debug.Stack()))
		}
	}()
	c.Send <- msg
}

// 关闭客户端数据
func (c *Client) close() {
	close(c.Send)
}

// Login 用户登录
func (c *Client) Login(roomId string, userId string, loginTime uint64) {
	c.RoomId = roomId
	c.UserId = userId
	c.LoginTime = loginTime
	// 登录成功=心跳一次
	c.Heartbeat(loginTime)
}

// Heartbeat 用户心跳
func (c *Client) Heartbeat(currentTime uint64) {
	c.HeartbeatTime = currentTime

	return
}

// IsHeartbeatTimeout 心跳超时
func (c *Client) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	if c.HeartbeatTime+heartbeatExpirationTime <= currentTime {
		timeout = true
	}

	return
}

// IsLogin 是否登录了
func (c *Client) IsLogin() (isLogin bool) {

	// 用户登录了
	if c.UserId != "" {
		isLogin = true

		return
	}

	return
}
