package main

import (
	"main/proofs/interpreter"
	// "main/proofs/retriever/stellarRetriever"
	// "main/proofs/interpreter"
	"fmt"
	"log"
	"net/http"

	"main/api/routes"
)

func main() {
	//create := executer.CreateAccount()
	// newRootHash := executer.InsertDataHash("E3FC18CB4776193F8AD15A947406DBYE", "SDL26B3CQN4AQHPV3MDRMUB5BXNMCQLHY3HVAD7ZOP4QACX2OL7V2IOW", "001", "cda5c9845b218fdd8e8f04d2db3a82db6b59d1d78785a214cbce3ee13311ee6e")

	router := routes.NewRouter()

	// fmt.Println("Latest Root hash")
	// fmt.Println(newRootHash)

	// test
	fmt.Println("Scenario 1 - Success")
	interpreter.InterpretPOC("", "")
	fmt.Println("Scenario 1 - Sequence Error ")
	interpreter.InterpretPOCError("", "")
	fmt.Println("Scenario 1 - Fake Tree Error ")
	interpreter.InterpretPOCFakeTree("", "")
	fmt.Println("Scenario 2 - Success")
	interpreter.InterpretPOC2("", "")
	fmt.Println("Scenario 2 - Sequence Error")
	interpreter.InterpretPOC2Error("", "")
	fmt.Println("Scenario 2 - Fake Tree Error")
	interpreter.InterpretPOC2FakeTree("", "")
	// lol:=
	// fmt.Println()

	log.Fatal(http.ListenAndServe(":8030", router))

}
