package routes

import (
	"filemanager/app/config"
	"filemanager/app/handlers"
	"filemanager/app/handlers/api/v1"
	"fmt"
	"github.com/go-martini/martini"
	"net/http"
)

var Martini *martini.ClassicMartini = martini.Classic()

func init() {
	// it will show the first page when user visits the site.
	Martini.Get("/", handlers.FileManagerHandler)

	// Additional routes.
	Martini.Group("/", func(r martini.Router) {
		r.Get("robots.txt", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Disallow: *")
		})
		r.Get("assets/css/global.min.css", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, fmt.Sprintf("%s/css/global.min.css", config.AssetsFolder))
		})
		r.Get("assets/js/main.min.js", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, fmt.Sprintf("%s/js/main.min.js", config.AssetsFolder))
		})
	})

	// API Endpoints
	Martini.Group("/api", func(r martini.Router) {
		r.Group("/v1", func(r martini.Router) {
			r.Any("/?", v1.ShowVersion)
			// Create a new folder
			r.Post("/folder", v1.CreateNewFolder)
			// Delete a folder
			r.Delete("/folder", v1.DeleteFolder)
		})
	})

	// Route when user is changing the folders
	Martini.Get("/folder", handlers.UpdateFolderHandler)

	// Route for deleting the file
	Martini.Get("/delete/:name", handlers.DeleteHandler)

	// Route for downloading the file
	Martini.Get("/download/:name", handlers.DownloadHandler)

	// Time for the upload
	Martini.Get("/upload", handlers.UploadsHandler)
	Martini.Post("/upload/file", handlers.FileUploadHandler)
}
