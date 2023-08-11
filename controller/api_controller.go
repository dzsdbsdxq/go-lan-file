package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"path"
	"share.ac.cn/cache"
	"share.ac.cn/common"
	"share.ac.cn/common/rsa"
	"share.ac.cn/config"
	"share.ac.cn/model"
	"share.ac.cn/response"
	"share.ac.cn/services/uploader"
	"share.ac.cn/services/websocket/service"
	"strconv"
	"time"
)

type IApiController interface {
	NewConnect(c *gin.Context)
	CreateUpload(c *gin.Context)
	GetFileInfo(c *gin.Context)
	GetFileComplete(c *gin.Context)
	GetAccessToken(c *gin.Context)
	CallBack(c *gin.Context)
	DeleteObject(c *gin.Context)
}
type ApiController struct {
}

func (api *ApiController) CreateUpload(c *gin.Context) {
	// 获取路径中的fileId
	userId := c.Param("userId")
	if userId == "" {
		response.Fail(c, nil, "ID不正确")
		return
	}
	//判断userId是否在线
	value, ok := service.LoginUserManager.GetUserClient(userId)
	if !ok {
		response.Fail(c, nil, "ID不在线，无法获取Token")
		return
	}

	fileName, _ := c.GetPostForm("name")

	fileSize, _ := c.GetPostForm("size")
	fileSizeS, _ := strconv.Atoi(fileSize)

	//key := fmt.Sprintf("%s%s%s", time.Now().Format("2006/01/02/"), common.GetRandomId(32), path.Ext(fileName))
	key := fmt.Sprintf("%s%s", time.Now().Format("2006/01/02/"), fileName)
	tmpOnline := &cache.TmpOnline{
		ShareId:  common.RandPass(4),
		FileName: fileName,
		FileKey:  key,
		FileSize: uint64(fileSizeS),
	}
	err := cache.SetTmpOnline(userId, tmpOnline)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}
	token := uploader.GetQiNiuAccessToken()
	response.Success(c, map[string]interface{}{"token": token, "source": key, "info": value.User}, "success")
}

func (api *ApiController) CallBack(c *gin.Context) {

}
func (api *ApiController) DeleteObject(c *gin.Context) {
	shareId := c.Param("shareId")
	info, err := cache.GetFileOnlineInfo(shareId)
	if err != nil {
		return
	}
	err = uploader.DeleteObject(info.FilePath)
	if err != nil {
		return
	}
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
		"url":         uploader.GetQiNiuDownloadUrl(fileModel.FilePath),
	}, "success")

	return

}

func (api *ApiController) GetAccessToken(c *gin.Context) {

}

func (api *ApiController) GetFileComplete(c *gin.Context) {
	// 获取路径中的userId
	userId := c.Param("userId")
	if userId == "" {
		response.Fail(c, nil, "ID不正确")
		return
	}
	fileHash, _ := c.GetPostForm("hash")
	//获取临时上传文件信息
	tmpInfo, _ := cache.GetTmpOnlineInfo(userId)

	expireTime := time.Now().Add(24 * time.Hour)
	fileInfo := &model.Files{
		Uid:        userId,
		IsUploaded: 1,
		IsDel:      1,
		Views:      0,
		Downloads:  0,
		FileSize:   int(tmpInfo.FileSize),
		ExpiredAt:  expireTime,
		FileName:   tmpInfo.FileName,
		FileId:     fileHash,
		ShareId:    tmpInfo.ShareId,
		FilePath:   tmpInfo.FileKey,
		FileExt:    path.Ext(tmpInfo.FileName),
	}
	_, _ = fileInfo.AddFiles()

	fileOnline := &model.FileOnline{
		CreateTime:   time.Now(),
		ExpireTime:   expireTime,
		FileSize:     tmpInfo.FileSize,
		FileViews:    0,
		FileDowns:    0,
		ShareId:      tmpInfo.ShareId,
		FileId:       fileHash,
		FileName:     tmpInfo.FileName,
		FilePath:     tmpInfo.FileKey,
		FileExt:      path.Ext(tmpInfo.FileName),
		FileHash:     fileHash,
		FileHashName: "sha1",
	}
	_ = cache.SetFileOnline(tmpInfo.ShareId, fileOnline)

	//设置七牛云文件有效期
	_ = uploader.DeleteAfterDays(tmpInfo.FileKey)

	//不让用户查看文件路径
	response.Success(c,
		map[string]interface{}{
			"file_id":      fileHash,
			"create_at":    time.Now(),
			"expired_at":   expireTime,
			"expired_time": expireTime.Unix(),
			"uid":          userId,
			"file_size":    tmpInfo.FileSize,
			"shareUrl":     config.Conf.System.HttpBaseWeb + tmpInfo.ShareId,
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
