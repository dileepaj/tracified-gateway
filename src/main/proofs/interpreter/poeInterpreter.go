package interpreter

import (
	"fmt"
)


type POEInterface interface {
	RetrievePOETest() (string, string, string)
}

type AbstractPOE struct {
}

func (AP *AbstractPOE) InterpretPOE(POEInterface POEInterface) (string, string) {
	Txn, bcHash, dbHash := POEInterface.RetrievePOETest()
	result := MatchingHash(bcHash, dbHash)
	if result == true {
		return "success", Txn
	} else {
		return "Error", Txn
	}
}

func MatchingHash(bcHash string, dbHash string) bool {
	isTrue := false
	if bcHash == dbHash {
		fmt.Println("Hash match!")
		isTrue = true

	} else {
		fmt.Println("Hash Din't match!")

	}
	return isTrue
}
