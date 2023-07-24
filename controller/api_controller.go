package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"share.ac.cn/cache"
	"share.ac.cn/common"
	"share.ac.cn/common/rsa"
	"share.ac.cn/config"
	"share.ac.cn/model"
	"share.ac.cn/response"
	"time"
)

type IApiController interface {
	NewConnect(c *gin.Context)
	GetFileInfo(c *gin.Context)
	GetFileComplete(c *gin.Context)
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

func (api *ApiController) GetFileInfo(c *gin.Context) {
	// 获取路径中的activityId
	shareId := c.Param("shareId")
	if shareId == "" {
		response.Fail(c, nil, "ID不正确")
		return
	}
	fileModel := &model.Files{}
	//查询文件是否存在
	fileModel.GetFileByShareId(shareId)
	fmt.Println(fileModel)
	//info, err := cache.GetFileOnlineInfo(shareId)
	if fileModel == nil {
		response.Fail(c, nil, "文件不存在")
		return
	}
	fileModel.UpdateViews()
	token, _ := api.generateToken(shareId)
	//生成token info.ExpireTime.Unix()
	url := `http://localhost:10000/files/` + fileModel.FileId + "?sign=" + token
	//
	response.Success(c, map[string]interface{}{
		"file_id":     fileModel.FileId,
		"uid":         fileModel.Uid,
		"file_name":   fileModel.FileName,
		"file_size":   fileModel.FileSize,
		"views":       fileModel.Views,
		"downloads":   fileModel.Downloads,
		"expired_at":  fileModel.ExpiredAt,
		"expire":      fileModel.ExpiredAt.Unix(),
		"is_uploaded": fileModel.IsUploaded,
		"file_ext":    fileModel.FileExt,
		"url":         url,
	}, "success")

	return

}

func (api *ApiController) GetFileComplete(c *gin.Context) {
	// 获取路径中的fileId
	fileId := c.Param("fileId")
	if fileId == "" {
		response.Fail(c, nil, "ID不正确")
		return
	}
	//文件是否存在
	fileModel := model.Files{}
	file, err := fileModel.GetFileByFileId(fileId)
	if err != nil {
		response.Fail(c, nil, "文件不存在")
		return
	}
	//不让用户查看文件路径
	file.FilePath = ""
	response.Success(c,
		map[string]interface{}{
			"file_id":      file.FileId,
			"create_at":    file.CreatedAt,
			"expired_at":   file.ExpiredAt,
			"expired_time": file.ExpiredAt.Unix(),
			"uid":          file.Uid,
			"file_size":    file.FileSize,
			"shareUrl":     config.Conf.System.HttpBaseWeb + file.ShareId,
			"complete":     true,
		}, "success")
	return
}
func (api *ApiController) generateToken(shareId string) (string, error) {
	// 设置过期时间，这里设置为1小时
	expirationTime := time.Now().Add(24 * time.Hour)

	// 创建一个新的令牌对象
	token := jwt.New(jwt.SigningMethodHS256)

	// 设置声明（claims），这里可以加入你的自定义信息
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = shareId
	claims["exp"] = expirationTime.Unix()

	// 在这里替换为你自己的密钥，确保它是一个长且随机的字符串
	// 在真实环境中，应该将密钥存储在环境变量或者配置文件中，而不是直接硬编码在代码中
	// 此处仅为示例目的
	secretKey := config.Conf.Jwt.Key

	// 签署 token，并获取完整的 token 字符串
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return rsa.EncodeStr2Base64(tokenString), nil
}
func NewApiController() IApiController {
	return &ApiController{}
}
