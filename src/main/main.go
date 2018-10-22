package main

import (
	"fmt"

	"log"
	"main/api/routes"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {

	headersOk := handlers.AllowedHeaders([]string{"Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router := routes.NewRouter()

	log.Fatal(http.ListenAndServe(":8030", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
	fmt.Println("Gateway Started @port :8030...!")
}
