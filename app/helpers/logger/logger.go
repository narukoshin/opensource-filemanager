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

func Logger(){
	m = routes.Martini

	go func(){
		for {
			time.Sleep(time.Until(time.Now().Truncate(24 * time.Hour)))
			// It will create a new file every day
			file_name := fmt.Sprintf("logs_%s.log", time.Now().Format("2006_01_02"))
			f, _ := os.OpenFile(filepath.Join(config.LogsFolder, file_name), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

			m.Map(log.New(f, "[martini] ", log.LstdFlags))

			// add custom prints
			m.Use(func(c martini.Context) {
				log.SetPrefix("[martini] ")
			})

			log.SetOutput(f)
		}
	}()
}