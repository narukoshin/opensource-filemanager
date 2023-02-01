package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-martini/martini"
	// I should use sessions to store current_directory value
	// It seems to be impossible to manage this by using only non-global or struct variables
	// ...so I should to try to do this by using session variables... that can fix my problem and I will be a step closer to run it on  production.
	// docs: https://gowebexamples.com/sessions/
	"github.com/gorilla/sessions"
)

type Directory_Structure struct {
	IsFolder bool
	Name string
	Size string
	Date string
	Ext string
}

// restricted folder where all the files will be stored in.
// user shouldn't have access outside the folder.
var public_directory string = "uploads"
// var current_directory string

var Store *sessions.CookieStore = sessions.NewCookieStore([]byte("secret-password"))

// calculating the actual size of the file to the human readable format
func CalculateActualSize(FloatSize float64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)
	switch {
		case FloatSize < 1000:
			return fmt.Sprintf("%d b", int(FloatSize))
		case FloatSize < 10000:
			return fmt.Sprintf("%.2f kb", float64(FloatSize / kb))
		case FloatSize < 100000000:
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
func LoadFilesFromDirectory(path string, session *sessions.Session) ([]Directory_Structure,  bool) {
	// checking if the folder exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// if the folder doesn't exist
		return []Directory_Structure{}, false
	}
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
	// Adding a folder go go back to the previous folder
	if filepath.Base(session.Values["current_directory"].(string)) != filepath.Base(public_directory) {
		previous := Directory_Structure {
			IsFolder: true,
			Name: "..",
		}
		Structure = append(Structure, previous)
	}
	// sorting by type, if it's a folder, then all folders should be at the top and after all folders follows files.
	sort.SliceStable(Structure, func(i, j int) bool {
		return Structure[i].IsFolder
	})
	return Structure, true

}

// Updating the folder name from which one we will load new files.
// it's like travelling in the filesystem.
func UpdateFolder(new_name string, session *sessions.Session) string {
	if new_name == ".." {
		// if new_name parameter will be empty then we will set back the default one.
		// splitting the path to delete the last element
		splitting := strings.Split(session.Values["current_directory"].(string), "/")
		if len(splitting) != 1 {
			index := len(splitting) - 1
			// deleting the last element
			removing_last := append(splitting[:index], splitting[index+1:]...)
			new_name = strings.Join(removing_last, "/")
		}
	} else {
		if len(session.Values["current_directory"].(string)) == 0 {
			session.Values["current_directory"] = public_directory
		}
		if new_name == "." {
			new_name = public_directory
		} else {
			new_name = fmt.Sprintf("%s/%s", session.Values["current_directory"], new_name)
		}
	}

	// ...now we have to somehow change the folder name and load new list of files.
	// current_directory = new_name
	session.Values["current_directory"] = new_name
	// removing all ".." from the path to fix a bug when user can get outside restricted folder
	new_name = strings.Replace(new_name, "..", "", -1)
	fmt.Println(session.Values["current_directory"].(string))
	return new_name
}

func Filemanager_Index(w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, "file-manager")
	if err != nil {
		panic(err)
	}
	session.Values["current_directory"] = public_directory
	session.Save(r, w)

	// getting the files from directory
	files, _ := LoadFilesFromDirectory(public_directory, session)
	
	// passing the files to template
	data := map[string]interface{}{
		"Files": files,
	}
	tmpl := template.Must(template.ParseFiles("templates/index/FileManager.html"))
	tmpl.ParseFiles("templates/index/FileList.html", "templates/index/FileIcons.html")
	err = tmpl.Execute(w, data)
	if err != nil {
		// this should be disabled in production
		// ...and all errors should be written in log file.
		panic(err)
	}
}

// when user will click on the file, that file will be downloaded into his computer.
func Filemanager_DownloadFile(w http.ResponseWriter, r *http.Request, param martini.Params){
	session, err := Store.Get(r, "file-manager")
	if err != nil {
		panic(err)
	}
	var file_name string = filepath.Base(param["name"])
	w.Header().Set("Content-Disposition", "attachment; filename="+file_name)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, fmt.Sprintf("%s/%s", session.Values["current_directory"], file_name))
}

// Changing the folders
func Filemanager_UpdateFolder(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	session, err := Store.Get(r, "file-manager")
	if err != nil {
		panic(err)
	}
	if session.Values["current_directory"] == nil {
		session.Values["current_directory"] = ""
		session.Save(r, w)
	}
	// this should change the folder and update content with the files from new folder.
	// this will prevent from LFI attacks.
	folder_name := filepath.Base(r.Form.Get("folder_name"))
	// updating the folder name.
	updated_name := UpdateFolder(folder_name, session)
	// getting the files from a new filder
	new_files, exist := LoadFilesFromDirectory(updated_name, session)
	if !exist {
		session.Values["current_directory"] = public_directory
		new_files, _ = LoadFilesFromDirectory(public_directory, session)
	}
	session.Save(r, w)
	// adding previous folder to the list
	data := map[string]interface{}{
		"Files": new_files,
	}
	tmpl := template.Must(template.ParseFiles("templates/index/FileManager.html"))
	tmpl.ParseFiles("templates/index/FileList.html", "templates/index/FileIcons.html")
	err = tmpl.ExecuteTemplate(w, "files", data)
	if err != nil {
		panic(err)
	}
}

func Filemanager_DeleteFile(w http.ResponseWriter, r *http.Request, param martini.Params){
	session, err := Store.Get(r, "file-manager")
	if err != nil {
		panic(err)
	}
	var file_name string = filepath.Base(param["name"])
	file_path := fmt.Sprintf("%s/%s", session.Values["current_directory"], file_name)

	if _, err := os.Stat(file_path); !os.IsNotExist(err) {
		// deleting the file.
		os.Remove(file_path)
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
	// deleting the file from the file manager.
	m.Get("/delete/:name", Filemanager_DeleteFile)
	// loading the generated css file.
	m.Get("/assets/css/global.min.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/css/global.min.css")
	})
	// javascript code
	m.Get("/assets/js/main.min.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "assets/js/main.min.js")
	})
	// Reading the robots.txt file
	m.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadFile("robots.txt")
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, string(body))
	})
	// starting the web server
	http.ListenAndServe(":8080", m)
}