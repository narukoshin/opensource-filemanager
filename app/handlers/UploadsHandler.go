package handlers

import (
	"filemanager/app/sessions"
	"filemanager/app/templates"
	"fmt"
	"io"

	"net/http"
	"os"
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

	// Allowing file upload till 1GB file size.
	r.ParseMultipartForm(1024 << 20)
	file, header, err := r.FormFile("myFile")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Retrieving the file name
	file_name := filepath.Base(header.Filename)

	// Creating a new file, where we'll upload a new file
	f, err := os.Create(fmt.Sprintf("%s/%s", current_directory, file_name))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	
	// Copying contents from the memory to the file we just created
	_, err = io.Copy(f, file)
	if err != nil {
		panic(err)
	}
}