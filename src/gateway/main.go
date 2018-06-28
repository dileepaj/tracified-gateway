package main

import (
	"gateway/api/routes"
	"log"
	"net/http"
)

func main() {

	router := routes.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
