package handlers

import (
	"filemanager/app/config"
	"filemanager/app/helpers/files"
	"filemanager/app/templates"
	"filemanager/app/sessions"
	"net/http"
)

func FileManagerHandler(w http.ResponseWriter, r *http.Request) {
	// Declaring the template variable
	var template templates.Template

	// Importing templates from the index templates/index folder
	template.FromFolder("index").Import("FileManager.html", "FileList.html", "FileIcons.html")

	session := &sessions.Session{}
	s, err := session.FirstRun().RetrieveSession(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	// setting a new value for current_directory
	s.Update("current_directory", config.DefaultDirectory)

	// Getting a current directory from the session
	current_directory, err := s.Get("current_directory")
	if err != nil {
		panic(err)
	}

	d := files.Directory_Paths{
		// Default directory is for checking if we are in another sub-directory
		// If we're in sub-directory, it will make directory to previous directory
		DefaultDirectory: config.DefaultDirectory,

		// A directory where we are at the moment
		// In the root page, it's the default directory
		CurrentDirectory: current_directory,

		// Because this is the root
		// We're going to read the root directory
		PathToRead: config.DefaultDirectory,
	}
	// Reading all files from the directory
	var directory files.Directory = files.Get(&d)
	// Checking if there's any error
	if directory.Error != nil {
		panic(err)
	}
	template.Data = map[string]interface{}{
		"Files": directory.FolderContents,
	}
	err = template.AddWriter(w).Execute()
	if err != nil {
		panic(err)
	}
}