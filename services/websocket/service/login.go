package service

import (
	"fmt"
	"sync"
)

type LoginUsers struct {
	Lock  sync.Mutex
	Users map[string]*Client //登录的用户
}

// AddUsers 添加登录用户
func (login *LoginUsers) AddUsers(client *Client) {
	login.Lock.Lock()
	defer login.Lock.Unlock()
	login.Users[client.User.UserId] = client
}

// DelUsers 删除登录用户
func (login *LoginUsers) DelUsers(client *Client) (result bool) {
	fmt.Println(client.User)
	fmt.Println(login.Users)
	if value, ok := login.Users[client.User.UserId]; ok {
		// 判断是否为相同的用户
		if value.Addr != client.Addr {
			return
		}
		delete(login.Users, client.User.UserId)
		result = true
	}
	return
}
func NewLoginUsers() *LoginUsers {
	return &LoginUsers{
		Users: make(map[string]*Client, 1000),
	}
}
