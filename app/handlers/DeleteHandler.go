package handlers

import (
	"encoding/json"
	"filemanager/app/sessions"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-martini/martini"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request, param martini.Params) {
	w.Header().Set("Content-Type", "application/json")

	session := &sessions.Session{}
	s, err := session.RetrieveSession(w, r)
	if err != nil {
		panic(err)
	}
	
	current_directory, err := s.Get("current_directory")
	if err != nil {
		panic(err)
	}


	base_name := filepath.Base(param["name"])
	delete_path := fmt.Sprintf("%s/%s", current_directory, base_name)

	if t, err := os.Stat(delete_path); !os.IsNotExist(err) {
		// If it's a folder, then we are deleting it with all the files.
		// We're deleting all files to avoid error about that directory is not empty.
		err = os.RemoveAll(delete_path)
		if err != nil {
			panic(err)
		}

		// Let's check it's folder or file that we just deleted
		var strtype string
		if t.IsDir() {
			strtype = "folder"
		} else {
			strtype = "file"
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"message":fmt.Sprintf("%s %s was successfully deleted", strtype, base_name)})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"message":fmt.Sprintf("%s does not exist", base_name)})
	}

}