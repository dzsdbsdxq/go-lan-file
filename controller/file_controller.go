package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"share.ac.cn/cache"
	"share.ac.cn/common"
	"share.ac.cn/common/rsa"
	"share.ac.cn/common/tusd"
	"share.ac.cn/config"
	"share.ac.cn/model"
	"share.ac.cn/repository"
	"share.ac.cn/response"
	"time"
)

type IFileController interface {
	PostFile(c *gin.Context)
	HeadFile(c *gin.Context)
	PatchFile(c *gin.Context)
	GetFile(c *gin.Context)
}
type FileController struct {
	fileRepository repository.IFileRepository
}

func (f *FileController) PostFile(c *gin.Context) {
	tusd.GetTusd().PostFile(c.Writer, c.Request)
}
func (f *FileController) HeadFile(c *gin.Context) {
	tusd.GetTusd().HeadFile(c.Writer, c.Request)
}
func (f *FileController) PatchFile(c *gin.Context) {
	tusd.GetTusd().PatchFile(c.Writer, c.Request)
}
func (f *FileController) GetFile(c *gin.Context) {

	fileError := response.FileError{
		Name:   "Request Get File",
		Type:   "get",
		Sign:   "",
		Reason: "",
	}

	sign := c.Query("sign")
	if sign == "" {
		fileError.Reason = "sign is not found"
		fileError.HandleFileError(c)
		return
	}

	token, err := f.validateToken(sign)
	if err != nil {
		fileError.Sign = sign
		fileError.Reason = err.Error()
		fileError.HandleFileError(c)
		return
	}
	if !token.Valid {
		fileError.Sign = sign
		fileError.Reason = "token is inValid"
		fileError.HandleFileError(c)
		return
	}

	// 获取 token 中的声明（claims）
	claims := token.Claims.(jwt.MapClaims)
	//查询文件是否存在
	fileModel := &model.Files{}
	fileModel.GetFileByShareId(claims["user_id"].(string))
	if fileModel == nil {
		fileError.Reason = "file is not found"
		fileError.HandleFileError(c)
		return
	}
	//判断file是否过期
	if time.Now().After(fileModel.ExpiredAt) {
		//文件已过期,删除文件,标记文件
		err = fileModel.UpdateColumn("is_del", 2)
		if err != nil {
			common.Log.Errorf("数据库标记文件出错：文件ID：%d，错误信息%s", fileModel.ID, err.Error())
			return
		}
		//删除文件
		err = common.DeleteFile(fileModel.FilePath)
		if err != nil {
			common.Log.Errorf("删除文件出错：文件ID：%d，错误信息%s", fileModel.ID, err.Error())
			return
		}
		_ = common.DeleteFile(fileModel.FilePath + ".info")
		cache.DeleteFileOnline(fileModel.ShareId)
		fileError.Reason = "file is not found"
		fileError.HandleFileError(c)
		return
	}
	fileModel.UpdateDownload()
	tusd.GetTusd().GetFile(c.Writer, c.Request)
}

func (f *FileController) validateToken(tokenString string) (*jwt.Token, error) {
	// 在真实环境中，应该将密钥存储在环境变量或者配置文件中，而不是直接硬编码在代码中
	// 此处仅为示例目的
	secretKey := config.Conf.Jwt.Key
	tokenString = rsa.DecodeStrFromBase64(tokenString)
	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保使用相同的密钥进行解析
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func NewFileController() IFileController {
	return &FileController{
		fileRepository: repository.NewFileRepository(),
	}
}
