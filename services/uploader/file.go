package uploader

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
	"path/filepath"
	"share.ac.cn/common/unique"
	"share.ac.cn/config"
	"share.ac.cn/model"
	"share.ac.cn/request"
	"share.ac.cn/response"
	"strconv"
)

type IFileUploadService interface {
	Upload(c *gin.Context, req *request.FileUploadRequest) (*response.FileUploadResponse, error)
	saveChunkToLocalFromMultiPartForm(c *gin.Context, tempFileName string, currentChunkSize int) (err error)
	UniqueFilename(filename string) string
}
type FileUploadService struct{}
type FileUploadInfo struct {
	Id               int    `json:"id"`
	ChunkNumber      int    `json:"chunkNumber"`      // 当前是第几片
	ChunkSize        int    `json:"chunkSize"`        // 每片的大小
	CurrentChunkSize int    `json:"currentChunkSize"` // 当前分片的大小
	TotalSize        int    `json:"totalSize"`        // 文件的总大小
	TotalChunks      int    `json:"totalChunks"`      // 总分片数
	FileName         string `json:"fileName"`
	Identifier       string `json:"identifier"` // fileMd5值
	HasBeenUploaded  string `json:"hasBeenUploaded"`
}

func NewFileUploadService() IFileUploadService {
	return &FileUploadService{}
}

func (f *FileUploadService) Upload(c *gin.Context, req *request.FileUploadRequest) (*response.FileUploadResponse, error) {
	// 返回给前端的对象
	resultInfo := response.FileUploadResponse{}
	files := &model.Files{}
	fmt.Println(req)
	tempFileName := req.Identifier + strconv.Itoa(req.ChunkNumber) + filepath.Ext(req.FileName)
	targetFileName := req.Identifier + filepath.Ext(req.FileName)
	//分片校验，上传前的预检请求
	if c.Request.Method == "GET" {
		isExit := files.ExistFileByFileId(req.Identifier)
		// 查询不到记录，说明未曾这是第一次上传
		if !isExit {
			fileDetail := &model.Files{
				Uid:             0,
				IsUploaded:      0,
				TotalChunks:     req.TotalChunks,
				IsDel:           1,
				Views:           0,
				Downloads:       0,
				FileName:        req.FileName,
				FileId:          req.Identifier,
				UploadId:        req.UploadId,
				ShareId:         req.ShareId,
				FileSize:        req.FileSize,
				FilePath:        "",
				FileExt:         filepath.Ext(req.FileName),
				HasBeenUploaded: "",
			}
			id, err := files.AddFiles(fileDetail)
			if err != nil || id < 0 {
				return nil, errors.New(fmt.Sprintf("fail to insert upload detail record fineName:[%v] md5:[%v] err:[%v]", req.FileName, req.Identifier, err))
			}
			// 如果单个chunk就可以完成文件的上传，直接告诉前端可以合并了
			// 前端单线程顺序发送chunk，故下面的条件衡成立
			if req.ChunkNumber == req.TotalChunks {
				resultInfo.IsUploaded = 0
				resultInfo.Merge = 1
				return &resultInfo, nil
			}
			resultInfo.IsUploaded = 0
			resultInfo.Merge = 0
			resultInfo.HasBeenUploaded = ""
			return &resultInfo, nil

		}
		// 校验一下，当前文件是否曾经上传过，并且完整的上传完成了
		detail, err := files.GetFileByFileId(req.Identifier)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("post: fail to findUploadDetailByFileId fileId:[%v] targetFileName:[%v]", req.Identifier, targetFileName))
		}
		// 已经上传过了
		if detail.IsUploaded == 1 {
			resultInfo.IsUploaded = 1
			resultInfo.Merge = 0
			return &resultInfo, nil
		}
		// 处理断点续传chunk
		resultInfo.IsUploaded = 0
		resultInfo.Merge = 0
		resultInfo.HasBeenUploaded = detail.HasBeenUploaded // 将曾经上传过的记录发送给前端
		return &resultInfo, nil
	}
	//POST 请求
	detail, err := files.GetFileByFileId(req.Identifier)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("post: fail to findUploadDetailByFileId fileId:[%v] targetFileName:[%v]", req.Identifier, targetFileName))
	}
	// 判断是否已经上传成了～
	if detail.IsUploaded == 1 {

		// todo 判断是否曾经merge过

		fmt.Printf("post: has been uploaded the file fileName:[%v] md5:[%s]", targetFileName, req.Identifier)
		resultInfo.IsUploaded = 1
		resultInfo.Merge = 0
		resultInfo.HttpCode = 201
		// todo 改变状态码，非200
		return &resultInfo, nil
	}
	// todo 当前currentSize为空， 需要特殊处理一下

	// 保存当前chunk
	err = f.saveChunkToLocalFromMultiPartForm(c, tempFileName, req.CurrentChunkSize)
	if err != nil {
		fmt.Printf("post fail to save chunk fileName:[%v] md5:[%v] chunkNumber:[%v]", req.FileName, req.Identifier, req.ChunkNumber)
		// 告诉前端重传
		resultInfo.HttpCode = 500
		resultInfo.IsUploaded = 1
		resultInfo.Merge = 0
		return &resultInfo, nil
	}
	// 多协程并发修改单行数据会产生覆盖，但是前端会单线程访问，故下面的操作安全
	if detail.TotalChunks == req.ChunkNumber {
		detail.HasBeenUploaded = detail.HasBeenUploaded + strconv.Itoa(req.ChunkNumber)
	} else {
		detail.HasBeenUploaded = detail.HasBeenUploaded + strconv.Itoa(req.ChunkNumber) + ":"
	}
	err = detail.UpdateColumn("has_been_uploaded", detail.HasBeenUploaded)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("fail to updateColumn err:[%v]", err))
	}

	// 如果前端是并发上传的化，这个顺序不能保证～，所以下面的条件应该是比对 totalChunks == chunkNumber不能保证全部正确
	// 目前前端会采用单条线程顺序发送chunk，故下面的条件衡成立
	// 服务端判断，当前chunkNumber == totalChunks时, 先保存当前chunk 再向前端发送特殊响应，前端接收到后会发送merge请求
	if req.ChunkNumber == req.TotalChunks {
		// 更新数据库 标记文件全部上传完成
		detail.IsUploaded = 1
		err = detail.UpdateColumn("is_uploaded", 1)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("fail to updateColumn err:[%v]", err))
		}
		resultInfo.IsUploaded = 0
		// 告诉前端可以merge了
		resultInfo.Merge = 1
		return &resultInfo, nil
	}
	resultInfo.IsUploaded = 0
	resultInfo.Merge = 0
	resultInfo.Msg = "ok"
	return &resultInfo, nil
}

/**
 *   desc：
 *       将chunk中的数据暂时存在本地
 *
 *	 params：
 *	 	 tempFileName: 当前分片使用的文件名
 *		 currentChunkSize: 当前分片的大小
 */
func (f *FileUploadService) saveChunkToLocalFromMultiPartForm(c *gin.Context, tempFileName string, currentChunkSize int) (err error) {
	// 创建文件夹
	path, err := os.Getwd()
	if config.Conf.File.StoreBasePath != "" {
		path = config.Conf.File.StoreBasePath
	}
	folderPath := path + "/upload/"
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		fmt.Println("创建文件夹")
		// 必须分成两步
		// 先创建文件夹
		err = os.MkdirAll(folderPath, 0777)
		if err != nil {
			return err
		}
		// 再修改权限
		err = os.Chmod(folderPath, 0777)
		if err != nil {
			return err
		}
	}
	// 保存本次上传的分片 /upload/username/fileMd5_chunkNumber.suffix
	// todo 为什么第一次是空的 (因为第一次是get)
	if c.Request.MultipartForm == nil {
		return
	}
	fileHeader := c.Request.MultipartForm.File["file"][0]
	if fileHeader == nil {
		err = errors.New("fileHeader 为空")
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		fmt.Printf("error : %v", err)
		return
	}
	//// 关闭文件
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)
	// 在本地创建文件，如果没有就创建，如果有就打开
	myFile, err := os.Create(folderPath + "/" + tempFileName)
	if err != nil {
		fmt.Printf("error : %v", err)
		return
	}
	defer func(myFile *os.File) {
		_ = myFile.Close()
	}(myFile)

	// 循环读取客户端发送过来的文件
	buf := make([]byte, currentChunkSize)
	num, err := file.Read(buf)
	if err != nil {
		fmt.Printf("error : %v", err)
		return
	}
	fmt.Printf("本轮读取了 num=[%v] byte", num)
	// 保存文件
	num, err = myFile.Write(buf)
	if err != nil {
		fmt.Printf("error : %v", err)
		return
	}
	fmt.Printf("本次保存分片Size为 num == [%v] ", num)
	return
}

// UniqueFilename 生成唯一文件名
func (f *FileUploadService) UniqueFilename(filename string) string {
	ext := filepath.Ext(filename)
	return unique.String().String() + ext
}
