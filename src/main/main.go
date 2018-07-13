// _Interfaces_ are named collections of method
// signatures.

package main

import (
	"fmt"
	"main/api/businessFacades"
)

func main() {

	computer := businessFacades.NewComputerFacade()
	// computer.Start(0)
	fmt.Println("----------")
	computer.End()

}
