package main

import (
	"net/http"
	"filemanager/app/routes"
)

func main(){
	http.ListenAndServe(":8080", routes.Martini)
}