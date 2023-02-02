package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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