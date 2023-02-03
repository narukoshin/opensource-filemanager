package handlers

import (
	"filemanager/app/sessions"
	"filemanager/app/templates"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func UploadsHandler(w http.ResponseWriter, r *http.Request) {
	// setting the template
	template := templates.Template{}
	template.FromFolder("upload").Import("upload.html")
	template.AddWriter(w).Execute()

}

func FileUploadHandler(w http.ResponseWriter, r *http.Request){
	// Inicializing the session
	session := &sessions.Session{}
	s, err := session.RetrieveSession(w, r)
	if err != nil {
		panic(err)
	}

	// throwing session handler to the black hole
	current_directory, err := s.Get("current_directory")
	if err != nil {
		panic(err)
	}

	file, header, err := r.FormFile("myFile")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file_name := filepath.Base(header.Filename)

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/%s", current_directory, file_name), data, 0644)
	if err != nil {
		panic(err)
	}
}