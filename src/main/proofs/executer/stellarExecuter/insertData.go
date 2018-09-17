package stellarExecuter

import (
	"fmt"
	"net/http"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"

	"main/api/apiModel"
	"main/model"
)

type ConcreteInsertData struct {
	InsertTDP apiModel.InsertTDP
	// Hash          string
	// InsertType    string
	// PreviousTXNID string
	// ProfileId     string
}

func (cd *ConcreteInsertData) InsertDataHash() model.InsertDataResponse {

	publicKey := "GD3EEFYWEP2XLLHONN2TRTQV4H5GSXJGCSUXZJGXGNZT4EFACOXEVLDJ"
	secretKey := "SA46OTS655ZDALIAODVCBWLWBXZWO6VUS6TU4U4GAIUVCKS2SYPDS7N4"
	var response model.InsertDataResponse
	response.ProfileID = cd.InsertTDP.ProfileID
	response.TxnType = cd.InsertTDP.Type

	// save data
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{publicKey},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.SetData("Transaction Type", []byte(cd.InsertTDP.Type)),
		build.SetData("PreviousTXNID", []byte(cd.InsertTDP.PreviousTXNID)),
		build.SetData("ProfileID", []byte(cd.InsertTDP.ProfileID)),
		build.SetData("TDPHash", []byte(cd.InsertTDP.DataHash)),
	)

	if err != nil {
		// panic(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "The HTTP request failed for InsertDataHash "
		return response
	}

	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(secretKey)
	if err != nil {
		// panic(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "signing request failed for the Transaction"
		return response
	}

	txeB64, err := txe.Base64()
	if err != nil {
		// panic(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "Base64 conversion failed for the Transaction"
		return response
	}

	// And finally, send it off to Stellar!
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		// panic(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "Test net client crashed"
		return response
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)

	response.Error.Code = http.StatusOK
	response.Error.Message = "Transaction performed in the blockchain."
	response.TDPID = resp.Hash

	return response

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
