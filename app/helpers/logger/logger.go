package logger

import (
	"filemanager/app/routes"
	"log"
	"os"

	"github.com/go-martini/martini"
)

var m *martini.ClassicMartini

func Logger(){
	m = routes.Martini
	f, err := os.OpenFile("app/logs/logs.1.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	m.Map(log.New(f, "[martini] ", log.LstdFlags))

	// add custom prints
	m.Use(func(c martini.Context) {
		log.SetPrefix("[martini] ")
		log.Println("Custom prints")
	})

	
	log.SetOutput(f)
}