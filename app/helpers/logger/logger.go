package logger

import (
	"filemanager/app/config"
	"filemanager/app/routes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-martini/martini"
)

var m *martini.ClassicMartini

var logFile *os.File

func CreateLogFile() {	
	m = routes.Martini
	// Create a new log file name for today
	fileName := fmt.Sprintf("logs_%s.log", time.Now().Format("2006_01_02"))

	// Open the new log file
	var err error
	logFile, err = os.OpenFile(filepath.Join(config.LogsFolder, fileName), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Print(err)
	}

	// Set the log output to the new log file
	m.Map(log.New(logFile, "[martini] ", log.LstdFlags))
	m.Use(func(c martini.Context) {
		log.SetPrefix("[martini] ")
	})
	log.SetOutput(logFile)
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), (now.Day()+1), 0, 0, 0, 0, time.Local)
	duration := midnight.Sub(now)
	// Setting the prefix to logger
	log.SetPrefix("[Logger] ")
	log.Println("Sleeping for midnight: ", midnight)
	log.Println("Time left: ", duration)
	time.Sleep(duration)
	// Closing the file after midnight
	<-time.After(duration)
	logFile.Close()
}

func Logger(){
	// Create a new log file every day
	go func() {
		for {
			CreateLogFile()
		}
	}()
}