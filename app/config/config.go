package config

import "path/filepath"

// Declaring some global variables
var DefaultDirectory string = "uploads"

var AppFolder		 string = "app"
var AssetsFolder 	 string = filepath.Join(AppFolder, "assets")
var TemplatesFolder  string	= filepath.Join(AppFolder, "templates")
var LogsFolder		 string = filepath.Join(AppFolder, "logs")