package service

import (
	"fmt"
	tusd "github.com/tus/tusd/pkg/handler"
	"math/rand"
	"strings"
	"time"
)

type IFileUploadService interface {
	Notify(handler *tusd.UnroutedHandler)
	PreFinishResponseCallback(hook tusd.HookEvent) error
}
type FileUploadService struct {
}

func (f *FileUploadService) RandPass(lenNum int) string {
	var chars = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	str := strings.Builder{}
	length := len(chars)
	rand.Seed(time.Now().UnixNano()) //重新播种，否则值不会变
	for i := 0; i < lenNum; i++ {
		str.WriteString(chars[rand.Intn(length)])

	}
	return str.String()
}

func (f *FileUploadService) PreFinishResponseCallback(hook tusd.HookEvent) error {
	//userId := hook.HTTPRequest.Header.Get("Authorization")
	//shareId := f.RandPass(4)
	//expireTime := time.Now().Add(24 * time.Hour)
	//fileInfo := &model.Files{
	//	Uid:        userId,
	//	IsUploaded: 1,
	//	IsDel:      1,
	//	Views:      0,
	//	Downloads:  0,
	//	FileSize:   int(hook.Upload.Size),
	//	ExpiredAt:  expireTime,
	//	FileName:   hook.Upload.MetaData["filename"],
	//	FileId:     hook.Upload.ID,
	//	ShareId:    shareId,
	//	FilePath:   hook.Upload.Storage["Path"],
	//	FileExt:    hook.Upload.MetaData["filetype"],
	//}
	//_, _ = fileInfo.AddFiles()
	//fileOnline := &model.FileOnline{
	//	CreateTime:   time.Now(),
	//	ExpireTime:   expireTime,
	//	FileSize:     uint64(hook.Upload.Size),
	//	FileViews:    0,
	//	FileDowns:    0,
	//	ShareId:      shareId,
	//	FileId:       hook.Upload.ID,
	//	FileName:     hook.Upload.MetaData["filename"],
	//	FilePath:     hook.Upload.Storage["Path"],
	//	FileExt:      hook.Upload.MetaData["filetype"],
	//	FileHash:     "",
	//	FileHashName: "",
	//}
	//_ = cache.SetFileOnline(shareId, fileOnline)
	return nil
}

func (f *FileUploadService) Notify(handler *tusd.UnroutedHandler) {
	for {
		select {
		case event := <-handler.CompleteUploads:
			fmt.Println(event.Upload.ID)
			//获取auth
			//userId := event.HTTPRequest.Header.Get("Authorization")
			//shareId := f.RandPass(4)
			//fileInfo := &model.Files{
			//	Uid:        userId,
			//	IsUploaded: 1,
			//	IsDel:      1,
			//	Views:      0,
			//	Downloads:  0,
			//	FileSize:   int(event.Upload.Size),
			//	ExpiredAt:  time.Now().Add(24 * time.Hour),
			//	FileName:   event.Upload.MetaData["filename"],
			//	FileId:     event.Upload.ID,
			//	ShareId:    shareId,
			//	FilePath:   event.Upload.Storage["Path"],
			//	FileExt:    event.Upload.MetaData["filetype"],
			//}
			//_, _ = fileInfo.AddFiles()
			//fileOnline := &model.FileOnline{
			//	CreateTime:   time.Now(),
			//	ExpireTime:   time.Now().Add(24 * time.Hour),
			//	FileSize:     uint64(event.Upload.Size),
			//	FileViews:    0,
			//	FileDowns:    0,
			//	ShareId:      shareId,
			//	FileId:       event.Upload.ID,
			//	FileName:     event.Upload.MetaData["filename"],
			//	FilePath:     event.Upload.Storage["Path"],
			//	FileExt:      event.Upload.MetaData["filetype"],
			//	FileHash:     "",
			//	FileHashName: "",
			//}
			//_ = cache.SetFileOnline(shareId, fileOnline)
			//service.NewLoginUsers()
			//fmt.Println(userId)
			//
			////获取登陆客户端
			//client := loginUserManager.GetUserClient(userId)
			//fmt.Println(client)
			//client.SendMsg(map[string]interface{}{
			//	"to":       userId,
			//	"fileName": event.Upload.MetaData["filename"],
			//	"fileId":   event.Upload.ID,
			//	"complete": true,
			//	"url":      "http://localhost:10000/" + shareId,
			//})

		}

	}
}
func NewFileUploadService() IFileUploadService {
	return &FileUploadService{}
}
