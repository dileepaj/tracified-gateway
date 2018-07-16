// _Interfaces_ are named collections of method
// signatures.

package main

import (
	"log"
	"main/api/routes"
	"net/http"
)

func main() {

	router := routes.NewRouter()

	log.Fatal(http.ListenAndServe(":8030", router))

	// computer := businessFacades.NewComputerFacade()
	// // computer.Start(0)
	// fmt.Println("----------")
	// computer.End()

	// businessFacades.TestPOE()

}
