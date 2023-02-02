package handlers

import (
	"filemanager/app/sessions"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-martini/martini"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request, params martini.Params) {
	session := &sessions.Session{}
	s, err := session.RetrieveSession(w, r)
	if err != nil {
		panic(err)
	}

	base_name := filepath.Base(params["name"])
	current_directory, err := s.Get("current_directory")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+base_name)
	
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, fmt.Sprintf("%s/%s", current_directory, base_name))
}