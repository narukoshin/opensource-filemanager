package handlers

import (
	"compress/gzip"
	"filemanager/app/sessions"
	"filemanager/app/templates"
	"io"

	"net/http"
	"os"
	"path/filepath"
)

func UploadsHandler(w http.ResponseWriter, r *http.Request) {
	// Checking if the session is valid
	// If session is not valid, user will be redirected to the root page
	session := &sessions.Session{}
	if _, err := session.RetrieveSession(w, r); err != nil {
		panic(err)
	}


	r.ParseForm()
	// setting the template
	template := templates.Template{}
	template.FromFolder("upload").Import("upload.html")
	template.UseWriter(w).Execute()

}

func FileUploadHandler(w http.ResponseWriter, r *http.Request){
	// Getting session ready
	session := &sessions.Session{}
	s, err := session.RetrieveSession(w, r)
	if err != nil {
		panic(err)
	}

	// Getting current directory
	current_directory, err := s.Get("current_directory")
	if err != nil {
		panic(err)
	}

	// Allowing only 100MB per request.
	err = r.ParseMultipartForm(100 << 20)
	if err != nil {
		panic(err)
	}

	// Parsing additional data about the file
	data := r.Form

	// Starting parsing the blob data
	blob, _, err := r.FormFile("data")
	if err != nil {
		panic(err)
	}
	defer blob.Close()

	// The path of the file where it will be stored
	file_path := filepath.Join(current_directory, data["filename"][0])

	// We need to figure out which blob it is
	// If it's the first one, then we will create a new file
	// If it's the second one, we will update an existing file
	if data["chunks_current"][0] == "1" {
		// Checking if the file already exists
		if _, err := os.Stat(file_path); !os.IsNotExist(err) {
			// If the file exists, we delete it.
			err = os.Remove(file_path)
			if err != nil {
				panic(err)
			}
		}
		// Creating a new file
		file, err := os.Create(file_path)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// compress there
		gzip := gzip.NewWriter(file)
		defer gzip.Close()

		// writing our content to the file
		if _, err := io.Copy(gzip, blob); err != nil {
			panic(err)
		}
	} else {
		// If it's not the first blob anymore
		// In way to create file, we will open it and continue write data.
		file, err := os.OpenFile(file_path, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// compress there
		gzip := gzip.NewWriter(file)
		defer gzip.Close()
		// writing contents to the existing file
		if _, err := io.Copy(gzip, blob); err != nil {
			panic(err)
		}
	}

}