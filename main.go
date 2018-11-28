package main

import (
	// "fmt"

	"fmt"
	// "log"
	"github.com/tracified-gateway/api/routes"
	"net/http"

	"github.com/gorilla/handlers"

)

func main() {
	
	headersOk := handlers.AllowedHeaders([]string{"Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router := routes.NewRouter()	
	fmt.Println("Gateway Started @port :8030...!")
	http.ListenAndServe(":8030", handlers.CORS(originsOk, headersOk, methodsOk)(router))

}
