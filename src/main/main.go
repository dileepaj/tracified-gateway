package main

import (
	"log"
	"main/api/routes"
	"main/proofs/proofExecuter/stellarexecuter"
	"net/http"
)

func main() {
	//create := stellarexecuter.CreateAccount()
	stellarexecuter.InsertDataHash("", "", "")
	router := routes.NewRouter()
	//fmt.Println(transfer)
	log.Fatal(http.ListenAndServe(":8080", router))
}
