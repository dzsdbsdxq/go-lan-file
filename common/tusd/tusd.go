package tusd

import (
	"fmt"
	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
	"share.ac.cn/config"
	"share.ac.cn/services/websocket/service"
)

var tusdHandler *tusd.UnroutedHandler

func NewTusdServer() {
	var err error
	store := filestore.FileStore{
		Path: config.Conf.File.StoreBasePath,
	}
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	tusdHandler, err = tusd.NewUnroutedHandler(tusd.Config{
		BasePath:                  "/files/",
		StoreComposer:             composer,
		NotifyCompleteUploads:     true,
		MaxSize:                   1024 * 1024 * 1000,
		PreFinishResponseCallback: service.NewFileUploadService().PreFinishResponseCallback,
	})

	if err != nil {
		panic(fmt.Errorf("unable to create handler: %s", err))
	}
	// Start another goroutine for receiving events from the handler whenever
	// an upload is completed. The event will contains details about the upload
	// itself and the relevant HTTP request.
	//go service.NewFileUploadService().Notify(tusdHandler)
}
func GetTusd() *tusd.UnroutedHandler {
	return tusdHandler
}
