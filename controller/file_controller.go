package controller

import (
	"github.com/gin-gonic/gin"
	"share.ac.cn/cache"
	"share.ac.cn/common"
	"share.ac.cn/model"
	"share.ac.cn/repository"
	"share.ac.cn/response"
	"time"
)

type IFileController interface {
	UploadFile(c *gin.Context)
	DownloadFile(c *gin.Context)
}
type FileController struct {
	fileRepository repository.IFileRepository
}

func (f *FileController) UploadFile(c *gin.Context) {
	//自动生成临时唯一房间号
	shareId := common.RandPass(4)
	err := cache.SetFileOnline(shareId, &model.FileOnline{
		CreateTime:   time.Time{},
		ExpireTime:   time.Time{},
		FileSize:     0,
		FileViews:    0,
		FileDowns:    0,
		ShareId:      shareId,
		FileId:       "111",
		FileName:     "222",
		FilePath:     "333",
		FileExt:      "444",
		FileHash:     "555",
		FileHashName: "666",
	})
	if err != nil {
		return
	}
}

func (f *FileController) DownloadFile(c *gin.Context) {
	// 获取路径中的shareId
	shareId := c.Param("shareId")
	if shareId == "" {
		response.Fail(c, nil, "分享ID不正确")
		return
	}
	file, err := f.fileRepository.GetFileInfoByShareId(shareId)
	if err != nil {
		response.Fail(c, nil, err.Error())
		return
	}

	common.Log.Info(file)
}

func NewFileController() IFileController {
	return &FileController{
		fileRepository: repository.NewFileRepository(),
	}
}