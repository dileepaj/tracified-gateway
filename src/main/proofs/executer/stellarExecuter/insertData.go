package stellarExecuter

import (
	"fmt"
	"net/http"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"

	"main/model"
	"main/proofs/builder"
)

type ConcreteInsertData struct {
	*builder.AbstractTDPInsert
	Hash          string
	InsertType    string
	PreviousTDPID string
	ProfileId     string
}

func (cd *ConcreteInsertData) InsertDataHash() model.InsertDataResponse {

	// publicKey := "GALRYOKBUJFH33OKLIQQWR2W466256XQRXK53YXFUF377T4MY7AZ6AXL"
	secretKey := "SCXW23KAWO444ZKKW4UJDM43PEPPOOFWMTN75VUPAAYKOZFJMQ3T6J5D"
	var response model.InsertDataResponse
	response.ProfileID = cd.ProfileId
	response.TxnType = cd.InsertType

	// save data
	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{secretKey},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.SetData("PreviousTDPID", []byte(cd.PreviousTDPID)),
		build.SetData("TDPHash", []byte(cd.Hash)),
		build.SetData("ProfileID", []byte(cd.ProfileId)),
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
	response.Txn = resp.Hash

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
