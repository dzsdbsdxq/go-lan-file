package model

import (
	"share.ac.cn/common"
	"time"
)

const heartbeatTimeout = 3 * 60 // 用户心跳超时时间

// UserOnline 用户在线状态
type UserOnline struct {
	UserId        string    `json:"userId"`        // 用户Id
	LoginTime     time.Time `json:"loginTime"`     // 用户上次登录时间
	HeartbeatTime time.Time `json:"heartbeatTime"` // 用户上次心跳时间
	LogOutTime    time.Time `json:"logOutTime"`    // 用户退出登录的时间
	Qua           string    `json:"qua"`           // qua
	DeviceInfo    string    `json:"deviceInfo"`    // 设备信息
	ShareId       string    `json:"shareId"`       //文件分享ID
	OnLine        bool      `json:"onLine"`        //是否在线
}

/**********************  数据处理  *********************************/

// UserLogin 用户登录
func (u *UserOnline) UserLogin() {
	u.LoginTime = time.Now()
	u.HeartbeatTime = time.Now()
	u.OnLine = true
}

// Heartbeat 用户心跳
func (u *UserOnline) Heartbeat() {
	u.HeartbeatTime = time.Now()
	u.OnLine = true
	return
}

// LogOut 用户退出登录
func (u *UserOnline) LogOut() {
	u.LogOutTime = time.Now()
	u.OnLine = false
	//Todo 清除redis缓存
}

/**********************  数据操作  *********************************/

// IsOnline 用户是否在线
func (u *UserOnline) IsOnline() (online bool) {
	if u.OnLine {
		return true
	}
	if u.HeartbeatTime.Before(time.Now().Add(-heartbeatTimeout * time.Second)) {
		common.Log.Infof("用户[%s]心跳超时", u.UserId)
		return false
	}
	return true
}
