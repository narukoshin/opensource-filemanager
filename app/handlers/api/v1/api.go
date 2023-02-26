package v1

import (
	"encoding/json"
	"filemanager/app/config"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var Version string = "v1.0.0"

type FolderJSON struct {
	FolderName string `json:"folder_name"`
	Path       string `json:"path"`
}

type CreateFolderJSON FolderJSON
type DeleteFolderJSON FolderJSON

type ErrorResponse struct {
	Error bool   `json:"error"`
	Code  string `json:"error_code"`
	MessageResponse
}

var DefaultTimeFormat string = "2006-01-02 15:04:05"

type MessageResponse struct {
	Message                string `json:"message"`
	AdditionalDataResponse `json:"additional_data"`
}

type AdditionalDataResponse struct {
	Path string `json:"path"`
	Time string `json:"time"`
}

// Removing all malicous characters that can be used to gain access to the files out-of-box
func PathValidation(path string) string {
	var new_path string
	// Removing all double dots (..)
	path = strings.ReplaceAll(path, ".", "")
	// Removing all slashes from both sides
	new_path = strings.TrimLeft(path, "/")
	new_path = strings.TrimRight(new_path, "/")
	return new_path
}

func ShowVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"api_version": Version,
			"routes": map[string]interface{}{
				"/folder": map[string]interface{}{
					"description": "Create or delete the folder",
					"methods":     []string{"POST", "DELETE"},
					"args": map[string]interface{}{
						"folder_name": map[string]interface{}{
							"description": "Name of the folder that will be created or deleted",
							"type":        []string{"string"},
							"required":    true,
						},
						"path": map[string]interface{}{
							"description": "Path to directory in which will be created or deleted the folder",
							"type":        []string{"string"},
							"required":    false,
						},
					},
				},
			},
			"author": map[string]interface{}{
				"name":   "Naru Koshin",
				"role":   "Software Engineer",
				"github": "https://github.com/narukoshin",
			},
		},
	)
}

// This will create a new folder
func CreateNewFolder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// Sending 500 status code
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(
			ErrorResponse{
				Error: true,
				Code:  "panic",
				MessageResponse: MessageResponse{
					Message: err.Error(),
					AdditionalDataResponse: AdditionalDataResponse{
						Time: time.Now().Format(DefaultTimeFormat),
					},
				},
			},
		)
	}
	var data CreateFolderJSON
	err = json.Unmarshal(body, &data)
	if err != nil {
		// Sending 500 status code
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(
			ErrorResponse{
				Error: true,
				Code:  "panic",
				MessageResponse: MessageResponse{
					Message: err.Error(),
					AdditionalDataResponse: AdditionalDataResponse{
						Time: time.Now().Format(DefaultTimeFormat),
					},
				},
			},
		)
	}
	// Checking if json values are not empty
	if data.FolderName == "" {
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: true,
			Code:  "empty_folder_name",
			MessageResponse: MessageResponse{
				"Parameter 'folder_name' is mandatory in order to create a new folder",
				AdditionalDataResponse{
					Path: data.Path,
					Time: time.Now().Format(DefaultTimeFormat),
				},
			},
		})
		return
	}
	// Updating folder name
	data.FolderName = filepath.Base(data.FolderName)
	// Updating path by adding the uploads directory and filtered one.
	data.Path = filepath.Join(config.DefaultDirectory, data.Path)

	if _, err := os.Stat(data.Path); os.IsNotExist(err) {
		// Sending 404 status code
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: true,
			Code:  "directory_not_found",
			MessageResponse: MessageResponse{
				fmt.Sprintf("Folder '%s' does not exist\nCannot create a folder in the directory which does not exist", filepath.Base(data.Path)),
				AdditionalDataResponse{
					Path: data.Path,
					Time: time.Now().Format(DefaultTimeFormat),
				},
			},
		})
		return
	}
	fullPath := filepath.Join(data.Path, data.FolderName)
	// Checking if the folder already exists
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		// sendin 409 status code
		w.WriteHeader(409)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: true,
			Code:  "folder_already_exists",
			MessageResponse: MessageResponse{
				fmt.Sprintf("Folder '%s' already exists", data.FolderName),
				AdditionalDataResponse{
					Path: data.Path,
					Time: time.Now().Format(DefaultTimeFormat),
				},
			},
		})
		return
	}
	// After all the tests, we can securely create a new folder
	err = os.Mkdir(fullPath, 644)
	if err != nil {
		// Sending 500 status code
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(
			ErrorResponse{
				Error: true,
				Code:  "panic",
				MessageResponse: MessageResponse{
					Message: err.Error(),
					AdditionalDataResponse: AdditionalDataResponse{
						Time: time.Now().Format(DefaultTimeFormat),
					},
				},
			},
		)
	}
	json.NewEncoder(w).Encode(MessageResponse{
		Message: fmt.Sprintf("Folder '%s' was successfully created", data.FolderName),
		AdditionalDataResponse: AdditionalDataResponse{
			Path: data.Path,
			Time: time.Now().Format(DefaultTimeFormat),
		},
	})
}

// Delete a folder
func DeleteFolder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// Sending 500 status code
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(
			ErrorResponse{
				Error: true,
				Code:  "panic",
				MessageResponse: MessageResponse{
					Message: err.Error(),
					AdditionalDataResponse: AdditionalDataResponse{
						Time: time.Now().Format(DefaultTimeFormat),
					},
				},
			},
		)
	}
	var data DeleteFolderJSON
	err = json.Unmarshal(body, &data)
	if err != nil {
		// Sending 500 status code
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(
			ErrorResponse{
				Error: true,
				Code:  "panic",
				MessageResponse: MessageResponse{
					Message: err.Error(),
					AdditionalDataResponse: AdditionalDataResponse{
						Time: time.Now().Format(DefaultTimeFormat),
					},
				},
			},
		)
	}
	// Checking if json values are not empty
	if data.FolderName == "" {
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: true,
			Code:  "empty_folder_name",
			MessageResponse: MessageResponse{
				"Parameter 'folder_name' is mandatory in order to create a new folder",
				AdditionalDataResponse{
					Path: data.Path,
					Time: time.Now().Format(DefaultTimeFormat),
				},
			},
		})
		return
	}
	// Updating folder name
	data.FolderName = filepath.Base(data.FolderName)
	// Updating path by adding the uploads directory and filtered one.
	data.Path = filepath.Join(config.DefaultDirectory, data.Path)

	if _, err := os.Stat(data.Path); os.IsNotExist(err) {
		// Sending 404 status code
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: true,
			Code:  "directory_not_found",
			MessageResponse: MessageResponse{
				fmt.Sprintf("Folder '%s' does not exist\nCannot delete a folder in the directory which does not exist", filepath.Base(data.Path)),
				AdditionalDataResponse{
					Path: data.Path,
					Time: time.Now().Format(DefaultTimeFormat),
				},
			},
		})
		return
	}
	fullPath := filepath.Join(data.Path, data.FolderName)
	// Checking if the folder already exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// Sending 404 status code
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: true,
			Code:  "folder_not_found",
			MessageResponse: MessageResponse{
				fmt.Sprintf("Folder '%s' does not exist", data.FolderName),
				AdditionalDataResponse{
					Path: data.Path,
					Time: time.Now().Format(DefaultTimeFormat),
				},
			},
		})
		return
	}
	// After all the tests, we can securely create a new folder
	err = os.RemoveAll(fullPath)
	if err != nil {
		// Sending 500 status code
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(
			ErrorResponse{
				Error: true,
				Code:  "panic",
				MessageResponse: MessageResponse{
					Message: err.Error(),
					AdditionalDataResponse: AdditionalDataResponse{
						Time: time.Now().Format(DefaultTimeFormat),
					},
				},
			},
		)
	}
	json.NewEncoder(w).Encode(MessageResponse{
		Message: fmt.Sprintf("Folder '%s' was successfully deleted", data.FolderName),
		AdditionalDataResponse: AdditionalDataResponse{
			Path: data.Path,
			Time: time.Now().Format(DefaultTimeFormat),
		},
	})
}
