package deprecatedStellarExecuter

import (
	"fmt"
	"net/http"

	// "github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
)

type ConcreteInsertData struct {
	XDR string
}

func (cd *ConcreteInsertData) InsertDataHash() model.SubmitXDRResponse {

	// publicKey := "GD3EEFYWEP2XLLHONN2TRTQV4H5GSXJGCSUXZJGXGNZT4EFACOXEVLDJ"
	// secretKey := "SA46OTS655ZDALIAODVCBWLWBXZWO6VUS6TU4U4GAIUVCKS2SYPDS7N4"
	var response model.SubmitXDRResponse
	// response.ProfileID = cd.InsertTDP.ProfileID
	// response.TxnType = cd.InsertTDP.Type

	// // save data
	// tx, err := build.Transaction(
	// 	commons.GetHorizonNetwork(),
	// 	build.SourceAccount{publicKey},
	// 	build.AutoSequence{commons.GetHorizonClient()},
	// 	build.SetData("Transaction Type", []byte(cd.InsertTDP.Type)),
	// 	build.SetData("PreviousTXNID", []byte(cd.InsertTDP.PreviousTXNID)),
	// 	build.SetData("ProfileID", []byte(cd.InsertTDP.ProfileID)),
	// 	build.SetData("Identifier", []byte(cd.InsertTDP.Identifier)),
	// 	build.SetData("TDPHash", []byte(cd.InsertTDP.DataHash)),
	// )

	// if err != nil {
	// 	// panic(err)
	// 	response.Error.Code = http.StatusNotFound
	// 	response.Error.Message = "The HTTP request failed for InsertDataHash "
	// 	return response
	// }

	// // Sign the transaction to prove you are actually the person sending it.
	// txe, err := tx.Sign(secretKey)
	// if err != nil {
	// 	// panic(err)
	// 	response.Error.Code = http.StatusNotFound
	// 	response.Error.Message = "signing request failed for the Transaction"
	// 	return response
	// }

	// txeB64, err := txe.Base64()
	// if err != nil {
	// 	// panic(err)
	// 	response.Error.Code = http.StatusNotFound
	// 	response.Error.Message = "Base64 conversion failed for the Transaction"
	// 	return response
	// }

	// And finally, send it off to Stellar!
	resp, err := commons.GetHorizonClient().SubmitTransactionXDR(cd.XDR)
	if err != nil {
		// panic(err)
		fmt.Println(err.Error())
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "Test net client crashed"
		return response
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)

	response.Error.Code = http.StatusOK
	response.Error.Message = "Transaction performed in the blockchain."
	response.TXNID = resp.Hash

	return response

}
