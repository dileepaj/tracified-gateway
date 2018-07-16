package stellarRetriever

import (
	"fmt"

	"main/proofs/interpreter"
)

type ConcretePOE struct {
	*interpreter.AbstractPOE
	Txn       string
	ProfileID string
	Hash      string
}

func (db *ConcretePOE) RetrievePOETest() (string, string, string) {
	fmt.Println(db.Hash)
	fmt.Println(db.ProfileID)
	fmt.Println(db.Txn)
	bcHash := "cf68e34967e10837d629b941bb8ec85d0ef016bc324340bd54e0ccae08a30b7a"

	// https://horizon-testnet.stellar.org/transactions/e903f5ef813002295e97c0f08cf26d1fd411615e18384890395f6b0943ed83b5/operations
	return db.Txn, bcHash, db.Hash
}
