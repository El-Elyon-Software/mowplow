package main

import (
	"log"
	"mowplow/api/v1/controllers"
	"net/http"
)

func main() {
	router := controllers.Routes()
	log.Fatal(http.ListenAndServe(":8080", router))
}
