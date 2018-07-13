package builder

import (
	"fmt"

	"main/model"
	"main/proofs/executer/stellarExecuter"
)

type InsertData struct{}

func InsertTDP(hash string, secret string, profileId string, rootHash string) model.RootTree {
	result := stellarExecuter.InsertDataHash(hash, secret, profileId, rootHash)

	if result.Hash == "" {
		fmt.Println("Error in Stellar Executer!")
	}

	return result
}

func (I *InsertData) TDPInsert(hash string, secret string, profileId string, rootHash string) model.RootTree {
	result := stellarExecuter.InsertDataHash(hash, secret, profileId, rootHash)

	if result.Hash == "" {
		fmt.Println("Error in Stellar Executer!")
	}

	return result
}

// // For our example we'll implement this interface on
// // `rect` and `circle` types.
// type Rect struct {
// 	Width, Height float64
// }
// type Circle struct {
// 	Radius float64
// }

// // To implement an interface in Go, we just need to
// // implement all the methods in the interface. Here we
// // implement `geometry` on `rect`s.
// func (r Rect) Area() float64 {
// 	return r.Width * r.Height
// }
// func (r Rect) Perim() float64 {
// 	return 2*r.Width + 2*r.Height
// }

// // The implementation for `circle`s.
// func (c Circle) Area() float64 {
// 	return math.Pi * c.Radius * c.Radius
// }
// func (c Circle) Perim() float64 {
// 	return 2 * math.Pi * c.Radius
// }

// type CPU struct{}

// func (c *CPU) Freeze() {
// 	fmt.Println("CPU.Freeze()")
// }

// func (c *CPU) Jump(position int) {
// 	fmt.Println("CPU.Jump()")
// }

// func (c *CPU) Execute() {
// 	fmt.Println("CPU.Execute()")
// }
