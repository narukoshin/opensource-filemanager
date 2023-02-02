package builder

import (
	"filemanager/app/controllers/files"
	"filemanager/app/sessions"
	"fmt"
	"strings"
)

// Updating the folder name from which one we will load new files.
// it's like travelling in the filesystem.
func FolderName(d *files.Directory_Paths, s *sessions.Session) {
	if d.PathToRead == ".." {
		// if new_name parameter will be empty then we will set back the default one.
		// splitting the path to delete the last element
		splitting := strings.Split(d.CurrentDirectory, "/")
		if len(splitting) != 1 {
			index := len(splitting) - 1
			// deleting the last element
			removing_last := append(splitting[:index], splitting[index+1:]...)
			d.PathToRead = strings.Join(removing_last, "/")
		}
	} else {
		if len(d.CurrentDirectory) == 0 {
			s.Update("current_directory", d.DefaultDirectory)
		}
		if d.PathToRead == "." {
			d.PathToRead = d.DefaultDirectory
		} else {
			d.PathToRead = fmt.Sprintf("%s/%s", d.CurrentDirectory, d.PathToRead)
		}
	}

	// ...now we have to somehow change the folder name and load new list of files.
	// current_directory = new_name
	s.Update("current_directory", d.PathToRead)
	// removing all ".." from the path to fix a bug when user can get outside restricted folder
	d.PathToRead = strings.Replace(d.PathToRead, "..", "", -1)
}