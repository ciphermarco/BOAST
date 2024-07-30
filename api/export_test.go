package api

import (
	"net/http"

	app "github.com/ciphermarco/BOAST"
	"github.com/ciphermarco/BOAST/log"
)

type ExportAPI struct {
	http.Handler
}

func NewTestAPI(statusPath string, strg app.Storage) *ExportAPI {
	handler, err := api("", statusPath, strg)
	if err != nil {
		log.Fatalln(err)
	}
	return &ExportAPI{
		Handler: handler,
	}
}
