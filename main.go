package main

import (
	"log"
	"github.com/Tracified-Gateway/api/routes"
	"net/http"
)

func main() {

	router := routes.NewRouter()

	log.Fatal(http.ListenAndServe(":8030", router))

}
