package stellarexecuter

import (
	"fmt"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

func InsertDataHash(hash string, secret string, profileId string, rootHash string) string {

	// save data
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{secret},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.SetData(profileId, []byte(hash)),
	)

	if err != nil {
		panic(err)
	}

	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(secret)
	if err != nil {
		panic(err)
	}

	txeB64, err := txe.Base64()
	if err != nil {
		panic(err)
	}

	// And finally, send it off to Stellar!
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)

	return AppendToTree(resp.Hash, rootHash, secret)

}

func AppendToTree(current string, pre string, secret string) string {
	// save data
	tx, err := build.Transaction()
	if pre != "" {
		tx, err = build.Transaction(
			build.TestNetwork,
			build.SourceAccount{secret},
			build.AutoSequence{horizon.DefaultTestNetClient},
			build.SetData("previous", []byte(pre)),
			build.SetData("current", []byte(current)),
		)
	} else {
		tx, err = build.Transaction(
			build.TestNetwork,
			build.SourceAccount{secret},
			build.AutoSequence{horizon.DefaultTestNetClient},
			build.SetData("current", []byte(current)),
		)
	}
	if err != nil {
		panic(err)
	}

	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(secret)
	if err != nil {
		panic(err)
	}

	txeB64, err := txe.Base64()
	if err != nil {
		panic(err)
	}

	// And finally, send it off to Stellar!
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successful Transaction Tree Appended:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)
	return resp.Hash
}
