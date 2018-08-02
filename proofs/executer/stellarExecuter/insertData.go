package stellarExecuter

import (
	"fmt"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"

	"github.com/Tracified-Gateway/model"
)

func InsertDataHash(hash string, secret string, profileId string, rootHash string) model.RootTree {
	result := model.RootTree{}
	// type rootTree model.rootTre
	// save data
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{secret},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.SetData(profileId, []byte(hash)),
	)

	if err != nil {
		// panic(err)
		result.Error.Code = 1
		result.Error.Message = "error"
		return result
	}

	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(secret)
	if err != nil {
		// panic(err)
		result.Error.Code = 1
		result.Error.Message = "error"
		return result
	}

	txeB64, err := txe.Base64()
	if err != nil {
		// panic(err)
		result.Error.Code = 1
		result.Error.Message = "error"
		return result
	}

	// And finally, send it off to Stellar!
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		// panic(err)
		result.Error.Code = 1
		result.Error.Message = "error"
		return result
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)

	return AppendToTree(resp.Hash, rootHash, secret)

}

func AppendToTree(current string, pre string, secret string) model.RootTree {
	// save data
	result := model.RootTree{}
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
		// panic(err)
		result.Error.Code = 1
		result.Error.Message = "error"
		return result
	}

	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(secret)
	if err != nil {
		// panic(err)
		result.Error.Code = 1
		result.Error.Message = "error"
		return result
	}

	txeB64, err := txe.Base64()
	if err != nil {
		// panic(err)
		result.Error.Code = 1
		result.Error.Message = "error"
		return result
	}

	// And finally, send it off to Stellar!
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		// panic(err)
		result.Error.Code = 1
		result.Error.Message = "error"
		return result
	}
	fmt.Println("Successful Transaction Tree Appended:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)

	result.Hash = resp.Hash
	return result
}
