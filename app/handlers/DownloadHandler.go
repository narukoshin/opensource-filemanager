package handlers

import (
	"compress/gzip"
	"encoding/json"
	"filemanager/app/sessions"
	"fmt"
	"io"
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

	// Reading header of the file
	// To check if the file is compressed
	header := make([]byte, 2)
	f.Read(header)

	// After we read the first 2 bytes of the file
	// Setting it back to the beginning.
	f.Seek(0, 0)

	// Retrieving some information about the file
	// fs, _ := os.Stat(file_path)

	// Setting headers
	w.Header().Set("Content-Disposition", "attachment; filename="+base_name)
	w.Header().Set("Content-Type", "application/octet-stream")
	// w.Header().Set("Content-Length", strconv.FormatInt(fs.Size(), 10))

	if header[0] == 0x1F && header[1] == 0x8B  {
		// File is compressed
		gzip, err := gzip.NewReader(f)
		if err != nil {
			panic(err)
		}
		defer gzip.Close()

		_, err = io.Copy(w, gzip)
		if err != nil {
			panic(err)
		}
	} else {
		// File is not compressed
		_, err = io.Copy(w, f)
		if err != nil {
			panic(err)
		}
	}
}