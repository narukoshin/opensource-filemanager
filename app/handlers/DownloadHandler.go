package handlers

import (
	"encoding/json"
	"compress/gzip"
	"filemanager/app/sessions"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
	// Reading and uncompressing file
	file_path := filepath.Join(current_directory, base_name)
	if _, err := os.Stat(file_path); os.IsNotExist(err) {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": fmt.Sprintf("File %s doesn't exist", base_name)})
		return
	}
	f, err := os.Open(file_path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	

	gzip, err := gzip.NewReader(f)
	if err != nil {
		panic(err)
	}
	defer gzip.Close()

	uncompresed, err := ioutil.ReadAll(gzip)
	if err != nil {
		panic(err)
	}


	// Setting headers
	w.Header().Set("Content-Disposition", "attachment; filename="+base_name)
	w.Header().Set("Content-Type", "application/octet-stream")

	// Sending the uncompressed file
	w.Write(uncompresed)
}