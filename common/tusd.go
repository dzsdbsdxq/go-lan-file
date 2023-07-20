package common

import (
	"fmt"
	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
	"share.ac.cn/config"
	"share.ac.cn/services/uploader"
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
		BasePath:              "/files/",
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})
	if err != nil {
		panic(fmt.Errorf("unable to create handler: %s", err))
	}
	// Start another goroutine for receiving events from the handler whenever
	// an upload is completed. The event will contains details about the upload
	// itself and the relevant HTTP request.
	go uploader.NewFileUploadService().Notify(tusdHandler)
}
func GetTusd() *tusd.UnroutedHandler {
	return tusdHandler
}
