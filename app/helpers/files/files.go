package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Directory_Structure struct {
	IsFolder bool
	Name 	 string
	Size 	 string
	Date 	 string
	Ext 	 string
}

type Directory_Paths struct {
	DefaultDirectory  string
	CurrentDirectory  string
	PathToRead		  string
}

type Directory struct {
	FolderContents  []Directory_Structure
	Error			error
}

// Loading the folders and files from the directory
func Get(d *Directory_Paths) Directory{
	// Checking if the folder exists
	if _, err := os.Stat(d.PathToRead); os.IsNotExist(err) {
		// if the folder doesn't exist.
		return Directory{
			FolderContents: []Directory_Structure{},
			Error: err,
		}
	}
	files, err := ioutil.ReadDir(d.PathToRead)
	if err != nil {
		return Directory{
			FolderContents: []Directory_Structure{},
			Error: err,
		}
	}

	// Today's date
	current_time := time.Now()

	var Structure []Directory_Structure
	for _, f := range files {
		var date string
		diff := int((current_time.Sub(f.ModTime()).Hours() / 24))
		switch {
			case diff == 0:
				date = "today"
			case diff == 1:
				date = "yesterday"
			case diff > 365:
				date = f.ModTime().Format("02 Jan 2006")
			default:
				date = f.ModTime().Format("02 Jan")
		}
		d := Directory_Structure{
			IsFolder: f.IsDir(),
			Name: f.Name(),
			Size: CalculateActualSize(float64(f.Size())),
			Date: date,
			Ext: GetFileExt(f.Name()),
		}
		Structure = append(Structure, d)
	}
	// Adding a folder go go back to the previous folder
	if filepath.Base(d.CurrentDirectory) != filepath.Base(d.DefaultDirectory) {
		goBack := Directory_Structure {
			IsFolder: true,
			Name: "..",
		}
		Structure = append(Structure, goBack)
	}
	// sorting by type, if it's a folder, then all folders should be at the top and after all folders follows files.
	sort.SliceStable(Structure, func(i, j int) bool {
		return Structure[i].IsFolder
	})
	return Directory{
		FolderContents: Structure,
		Error: nil,
	}
}

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
		case FloatSize < 1000000:
			return fmt.Sprintf("%.2f kb", float64(FloatSize / kb))
		case FloatSize < 1000000000:
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