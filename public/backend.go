package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sort"

	"github.com/go-martini/martini"
)

type Directory_Structure struct {
	IsFolder bool
	Name string
	Size string
	Date string
	Ext string
}

// striched folder where all the files will be stored in.
// user shouldn't have access outside the folder.
var public_directory string = "./uploads"
var current_directory string

// calculating the actual size of the file to the human readable format
func CalculateActualSize(FloatSize float64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)
	switch {
		case FloatSize < 100:
			return fmt.Sprintf("%d b", int(FloatSize))
		case FloatSize < 10000:
			return fmt.Sprintf("%.2f kb", float64(FloatSize / kb))
		case FloatSize < 10000000:
			return fmt.Sprintf("%.2f mb", float64(FloatSize / mb))
		case FloatSize < 10000000000:
			return fmt.Sprintf("%.2f gb", float64(FloatSize / gb))
	}
	return ""
}

// gets file extension
func GetFileExt(fileName string) string {
	ext := filepath.Ext(fileName)
	if len(ext) != 0 {
		return ext[1:]
	}
	return ""
}

// getting the file and folder list in the specified directory
func LoadFilesFromDirectory(path string) []Directory_Structure {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	var Structure []Directory_Structure
	for _, f := range files {
		d := Directory_Structure{
			IsFolder: f.IsDir(),
			Name: f.Name(),
			Size: CalculateActualSize(float64(f.Size())),
			Date: f.ModTime().Format("02 Jan"),
			Ext: GetFileExt(f.Name()),
		}
		Structure = append(Structure, d)
	}
	sort.SliceStable(Structure, func(i, j int) bool {
		return Structure[i].IsFolder
	})
	return Structure

}

// Updating the folder name from which one we will load new files.
// it's like travelling in the filesystem.
func UpdateFolder(new_name string) string {
	if new_name == ".." {
		// if new_name parameter will be empty then we will set back the default one.
		// default one is "./" folder
		new_name = public_directory
	} else {
		new_name = fmt.Sprintf("%s/%s", public_directory, new_name)
	}
	// ...now we have to somehow change the folder name and load new list of files.
	current_directory = new_name
	return new_name
}

func Filemanager_Index(w http.ResponseWriter, r *http.Request) {
	current_directory = public_directory
	// getting the files from directory
	var files []Directory_Structure = LoadFilesFromDirectory(public_directory)
	// passing the files to template
	data := map[string]interface{}{
		"Files": files,
	}
	tmpl := template.Must(template.ParseFiles("templates/file-list.html"))
	err := tmpl.Execute(w, data)
	if err != nil {
		// this should be disabled in production
		// ...and all errors should be written in log file.
		panic(err)
	}
}

// when user will click on the file, that file will be downloaded into his computer.
func Filemanager_DownloadFile(w http.ResponseWriter, r *http.Request, param martini.Params){
	var file_name string = filepath.Base(param["name"])
	w.Header().Set("Content-Disposition", "attachment; filename="+file_name)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, fmt.Sprintf("%s/%s", current_directory, file_name))
}

// Changing the folders
func Filemanager_UpdateFolder(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// this should change the folder and update content with the files from new folder.
	// this will prevent from LFI attacks.
	// folder_name := filepath.Base(param["name"])
	folder_name := filepath.Base(r.Form.Get("folder_name"))
	// updating the folder name.
	updated_name := UpdateFolder(folder_name)
	// getting the files from a new filder
	new_files := LoadFilesFromDirectory(updated_name)
	// adding previous folder to the list
	if filepath.Base(current_directory) != filepath.Base(public_directory) {
		// adding a previous folder that will allow to travel back to the previous folder
		previous := Directory_Structure {
			IsFolder: true,
			Name: "..",
		}
		new_files = append(new_files, previous)
	}
	sort.SliceStable(new_files, func(i, j int) bool {
		return new_files[i].IsFolder
	})
	data := map[string]interface{}{
		"Files": new_files,
	}
	tmpl := template.Must(template.ParseFiles("templates/file-list.html"))
	err := tmpl.ExecuteTemplate(w, "files", data)
	if err != nil {
		panic(err)
	}
}

func main(){
	m := martini.Classic()
	// main page where all the files will appear.
	m.Get("/", Filemanager_Index)
	// changing the folder name.
	m.Get("/folder", Filemanager_UpdateFolder)
	// downloading the file from the file manager.
	m.Get("/download/:name", Filemanager_DownloadFile)
	// loading the generated css file.
	m.Get("/assets/css/global.min.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/css/global.min.css")
	})
	// starting the web server
	http.ListenAndServe(":8080", m)
}