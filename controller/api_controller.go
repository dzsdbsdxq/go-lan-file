package controller

import (
	"github.com/gin-gonic/gin"
	"share.ac.cn/cache"
	"share.ac.cn/common"
	"share.ac.cn/model"
	"share.ac.cn/response"
	"time"
)

type IApiController interface {
	NewConnect(c *gin.Context)
}
type ApiController struct {
}

// NewConnect 创建连接，返回上传通道（userId）
func (api *ApiController) NewConnect(c *gin.Context) {
	//自动创建用户ID
	userId := common.GetUserIdByRandom()
	//检查redis是否安装启动
	_, err := common.GetClient().Ping().Result()
	if err != nil {
		response.Fail(c, nil, "系统错误")
		return
	}
	//将用户注册到Redis系统中
	userOnline := &model.UserOnline{
		UserId:        userId,
		LoginTime:     time.Now(),
		HeartbeatTime: time.Now(),
		LogOutTime:    time.Now(),
		Qua:           "",
		DeviceInfo:    c.Request.Header.Get("User-Agent"),
		ShareId:       "",
		OnLine:        false,
	}
	err = cache.SetUserOnlineInfo(userId, userOnline)
	if err != nil {
		return
	}
	response.Success(c, userOnline.UserId, "success")
}

func NewApiController() IApiController {
	return &ApiController{}
}
