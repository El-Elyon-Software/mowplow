package main

import (
	"log"
	"net/http"
	"mowplow/v1/controllers"
)

func main() {
	router := controllers.Routes()
	log.Fatal(http.ListenAndServe(":8080", router))
}