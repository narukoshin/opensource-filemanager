package handlers

import (
	"filemanager/app/config"
	"filemanager/app/helpers/builder"
	"filemanager/app/helpers/files"
	"filemanager/app/templates"
	"filemanager/app/sessions"
	"net/http"
	"os"
	"path/filepath"
)

func UpdateFolderHandler(w http.ResponseWriter, r *http.Request) {
	// Declaring the template variable
	var template templates.Template

	// Importing templates from the index templates/index folder
	template.FromFolder("index").Import("FileManager.html", "FileList.html", "FileIcons.html")

	r.ParseForm()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	session := &sessions.Session{}
	s, err := session.RetrieveSession(w, r)
	if err != nil {
		panic(err)
	}
	// this should change the folder and update content with the files from new folder.
	// this will prevent from LFI attacks.
	folder_name := filepath.Base(r.Form.Get("folder_name"))

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

		// Setting it to the folder_name
		// later, when we will get the full path
		// We will set it to the full path and load files
		// ..Just to don't repeat the same code
		PathToRead: folder_name,
	}
	// building the path for the folder to load from
	builder.FolderName(&d, session)

	// Updating the current_directory, because we just now changed it.
	// I noticed that we are not changing it anywhere. xd
	err = s.Update("current_directory", d.PathToRead)
	if err != nil {
		panic(err)
	}
	// Updating directory paths.
	d.CurrentDirectory = d.PathToRead

	// getting the files from the path
	var directory files.Directory = files.Get(&d)
	// checking for the errors
	if directory.Error != nil {
		// if there's an error that the folder doesn't exist
		// loading files from the default directory
		if os.IsNotExist(directory.Error){
			// setting the session to the default directory
			session.Update("current_directory", config.DefaultDirectory)
			// loading files from the default directory
			d.PathToRead = config.DefaultDirectory
			directory = files.Get(&d)	
		}
	}

	template.Data = map[string]interface{}{
		"Files": directory.FolderContents,
	}
	err = template.UseWriter(w).WithBlock("files").Execute()
	if err != nil {
		panic(err)
	}
	
}