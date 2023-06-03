package main

import (
	"log"
	"net/http"
	"cooltown/resources"
)

func main() {
	log.Fatal(http.ListenAndServe(":3002", resources.Router())) 
}

