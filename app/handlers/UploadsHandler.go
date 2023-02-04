package handlers

import (
	"filemanager/app/sessions"
	"filemanager/app/templates"
	"io"

	"net/http"
	"os"
	"path/filepath"
)

func UploadsHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// setting the template
	template := templates.Template{}
	template.FromFolder("upload").Import("upload.html")
	template.UseWriter(w).Execute()

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
	file, _, err := r.FormFile("myFile")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Retrieving the file name
	file_name := filepath.Base(r.FormValue("filename"))

	// Building the file path
	file_path := filepath.Join(current_directory, file_name)

	// Checking if the file exists
	var f *os.File
	if _, err := os.Stat(file_path); os.IsNotExist(err){
		// Creating a new file, where we'll upload a new file
		f, err = os.Create(file_path)
		if err != nil {
			panic(err)
		}
	} else {
		// If the file already exists,writing to the file
		f, err = os.OpenFile(file_path, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}
	}
	defer f.Close()
	
	// Copying contents from the memory to the file we just created
	_, err = io.Copy(f, file)
	if err != nil {
		panic(err)
	}

	/******************** MEMORY LEAK SOLUTION ********************/
	// Copying contents from the memory to the file we just created
	// Add buffer to copy the contents to avoid memory leak
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := f.Write(buf[:n]); err != nil {
			panic(err)
		}
	}
}