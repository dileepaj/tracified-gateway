// _Interfaces_ are named collections of method
// signatures.

package main

import (
	"log"
	"github.com/tracified-gateway/api/routes"
	"net/http"
)

func main() {

	router := routes.NewRouter()

	log.Fatal(http.ListenAndServe(":8030", router))

}
