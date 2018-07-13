package businessFacades

import (
	"fmt"
	"main/proofs/builder"
)

type ShapeMaker interface {
	builder.Geometry
}

// If a variable has an interface type, then we can call
// methods that are in the named interface. Here's a
// generic `measure` function taking advantage of this
// to work on any `geometry`.
func Measure(g ShapeMaker) {
	// var w g.Area
	// x := g.Geometry
	fmt.Println(g.Area())
	fmt.Println(g.Perim())
}
