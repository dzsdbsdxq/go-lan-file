package service

//本模块源码参考https://blog.csdn.net/qq_34857250/article/details/105122272/
//主要是对代码进行了拆分和增加心跳检测逻辑代码
import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"share.ac.cn/common"
	"sync"
)

// 用户登录管理
var loginUserManager = NewLoginUsers()

// Manager 管理所有websocket的信息
type Manager struct {
	Clients    map[*Client]bool // 全部的连接
	Group      map[string]map[string]*Client
	Lock       sync.Mutex
	Register   chan *Client
	UnRegister chan *Client
	Login      chan *Client //带用户信息的登录
}

type Head struct {
	Seq     string // 消息的Id
	Cmd     string // 消息的cmd 动作
	Message interface{}
}

// MessageData 单个发送数据信息
type MessageData struct {
	*Head
	Id    string
	Group string
}

// GroupMessageData 组广播数据信息
type GroupMessageData struct {
	*Head
	Group string
}

// BroadCastMessageData 广播发送数据信息
type BroadCastMessageData struct {
	*Head
}

// InClient 判断客户端是否存在
func (manager *Manager) InClient(client *Client) (ok bool) {
	manager.Lock.Lock()
	defer manager.Lock.Unlock()
	// 连接存在，在添加
	_, ok = manager.Clients[client]
	return
}

// AddClients 添加客户端
func (manager *Manager) AddClients(client *Client) {
	manager.Clients[client] = true
}

// DelClients 删除客户端
func (manager *Manager) DelClients(client *Client) {

	if _, ok := manager.Clients[client]; ok {
		delete(manager.Clients, client)
	}
}

// Start 启动 websocket 管理器
func (manager *Manager) start() {
	common.Log.Info("websocket manage start:监听websocket进行中...")
	for {
		select {
		// 注册
		case client := <-manager.Register:
			manager.EventRegister(client)
		//登录
		case client := <-manager.Login:
			manager.EventLogin(client)

		// 注销
		case client := <-manager.UnRegister:
			manager.EventUnregister(client)

			// 发送广播数据到某个组的 channel 变量 Send 中
			//case data := <-manager.boardCast:
			//	if groupMap, ok := manager.wsGroup[data.GroupId]; ok {
			//		for _, conn := range groupMap {
			//			conn.Send <- data.Data
			//		}
			//	}
		}
	}
}

// EventRegister 用户建立连接事件
func (manager *Manager) EventRegister(client *Client) {
	common.Log.Infof("client [%s] connect", client.Id)
	common.Log.Infof("register client [%s] to group [%s]", client.Id, client.Group)

	manager.Lock.Lock()
	defer manager.Lock.Unlock()

	if manager.Group[client.Group] == nil {
		manager.Group[client.Group] = make(map[string]*Client, 1000)
	}
	manager.Group[client.Group][client.Id] = client

	//添加客户端
	manager.AddClients(client)

	manager.Send(client.Id, client.Group, "enter", "hello~")
}

// EventUnregister 用户断开连接
func (manager *Manager) EventUnregister(client *Client) {
	common.Log.Infof("EventUnregister 用户断开连接 client [%s] from group [%s]", client.Id, client.Group)
	manager.Lock.Lock()
	defer manager.Lock.Unlock()
	//删除客户端
	manager.DelClients(client)
	if client.User != nil {
		//删除登录用户连接
		loginUserManager.DelUsers(client)
	}

	if _, ok := manager.Group[client.Group]; ok {
		if _, ok := manager.Group[client.Group][client.Id]; ok {
			close(client.Message)
			delete(manager.Group[client.Group], client.Id)
			if len(manager.Group[client.Group]) == 0 {
				delete(manager.Group, client.Group)
			}
		}
	}
	//Todo 给client所在的group发送离开消息
	manager.SendGroup(client.Group, "exit", "bye~")
}

// EventLogin 用户登录
func (manager *Manager) EventLogin(client *Client) {
	//如果存在连接，则添加登录用户
	if manager.InClient(client) {
		loginUserManager.AddUsers(client)
		common.Log.Info("EventLogin 用户登录", client.Addr, client.User.UserId)
		//发送登录成功信息
		manager.Send(client.Id, client.Group, "login", []byte("success"))
	}

}

// SendService 处理单个 client 发送数据
func (manager *Manager) sendService(data *MessageData) {
	common.Log.Info("处理单个 client 发送数据...")
	if groupMap, ok := manager.Group[data.Group]; ok {
		if conn, ok := groupMap[data.Id]; ok {
			fmt.Println("send:", groupMap[data.Id], manager.Group[data.Group])
			marshal, _ := jsoniter.Marshal(data)
			conn.Message <- marshal
		}
	}
}

// SendGroupService 处理 group 广播数据
func (manager *Manager) sendGroupService(data *GroupMessageData) {
	if groupMap, ok := manager.Group[data.Group]; ok {
		for _, conn := range groupMap {
			marshal, _ := jsoniter.Marshal(data)
			conn.Message <- marshal
		}
	}
}

// SendAllService 处理广播数据
func (manager *Manager) sendAllService(data *BroadCastMessageData) {
	for _, v := range manager.Group {
		for _, conn := range v {
			marshal, _ := jsoniter.Marshal(data)
			conn.Message <- marshal
		}
	}
}

// Send 向指定的 client 发送数据
func (manager *Manager) Send(id string, group string, cmd string, message interface{}) {
	data := &MessageData{
		Id:    id,
		Group: group,
		Head: &Head{
			Seq:     common.GetRandomId(11),
			Cmd:     cmd,
			Message: message,
		},
	}
	fmt.Println(data)
	manager.sendService(data)
}

// SendGroup 向指定的 Group 广播
func (manager *Manager) SendGroup(group string, cmd string, message interface{}) {
	data := &GroupMessageData{
		Group: group,
		Head: &Head{
			Seq:     common.GetRandomId(11),
			Cmd:     cmd,
			Message: message,
		},
	}
	manager.sendGroupService(data)
}

// SendAll 广播
func (manager *Manager) SendAll(cmd string, message interface{}) {
	data := &BroadCastMessageData{
		Head: &Head{
			Seq:     common.GetRandomId(11),
			Cmd:     cmd,
			Message: message,
		},
	}
	manager.sendAllService(data)
}

// RegisterClient 注册
func (manager *Manager) RegisterClient(client *Client) {
	manager.Register <- client
}

// UnRegisterClient 注销
func (manager *Manager) UnRegisterClient(client *Client) {
	manager.UnRegister <- client
}

// LenGroup 当前组个数
func (manager *Manager) LenGroup() uint {
	return uint(len(manager.Group))
}

// LenClient 当前连接个数
func (manager *Manager) LenClient() uint {
	return uint(len(manager.Clients))
}

// Info 获取 wsManager 管理器信息
func (manager *Manager) Info() map[string]interface{} {
	managerInfo := make(map[string]interface{})
	managerInfo["groupLen"] = manager.LenGroup()
	managerInfo["clientLen"] = manager.LenClient()
	managerInfo["chanRegisterLen"] = len(manager.Register)
	managerInfo["chanUnregisterLen"] = len(manager.UnRegister)
	return managerInfo
}

// NewWebsocketManager  初始化 wsManager 管理器
func NewWebsocketManager() (clientManager *Manager) {
	clientManager = &Manager{
		Clients:    make(map[*Client]bool, 1000),
		Group:      make(map[string]map[string]*Client, 1000),
		Login:      make(chan *Client, 1000),
		Register:   make(chan *Client, 1000),
		UnRegister: make(chan *Client, 1000),
	}
	return
}
