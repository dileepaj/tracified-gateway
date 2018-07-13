// _Interfaces_ are named collections of method
// signatures.

package main

import (
	"fmt"
	"main/api/businessFacades"
)

func main() {
	// r := builder.Rect{3, 4}
	// c := builder.Circle{5}

	// // The `circle` and `rect` struct types both
	// // implement the `geometry` interface so we can use
	// // instances of
	// // these structs as arguments to `measure`.
	// businessFacades.Measure(r)
	// businessFacades.Measure(c)

	computer := businessFacades.NewComputerFacade()
	// computer.Start(0)
	fmt.Println("----------")
	computer.End()

}
