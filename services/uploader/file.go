package uploader

import (
	"fmt"
	tusd "github.com/tus/tusd/pkg/handler"
)

type IFileUploadService interface {
	Notify(handler *tusd.UnroutedHandler)
}
type FileUploadService struct {
}

func (f *FileUploadService) Notify(handler *tusd.UnroutedHandler) {
	for {
		select {
		case event := <-handler.CompleteUploads:
			fmt.Println(event.Upload)
			fmt.Println("haha:", event.HTTPRequest.Header.Get("Sec-Ch-Ua-Platform"))
			fmt.Printf("Upload %s finished\n", event.Upload.ID)
		}

	}
}
func NewFileUploadService() IFileUploadService {
	return &FileUploadService{}
}
