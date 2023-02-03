package main

import (
	"filemanager/app/routes"
	"filemanager/app/helpers/logger"
)

func main(){
	logger.Logger()

	routes.Martini.RunOnAddr(":8080")
}