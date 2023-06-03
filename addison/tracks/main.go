package main

import (
	"log"
	"net/http"
	"tracks/repository"
	"tracks/resources"
)

func main() {
	repository.Init()
	repository.Create()
	repository.Clear()
	log.Fatal(http.ListenAndServe(":3000", resources.Router())) 
}
