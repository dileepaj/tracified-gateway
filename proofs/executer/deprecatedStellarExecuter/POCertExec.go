package deprecatedStellarExecuter

import (
	"fmt"
	"net/http"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
)

type ConcreteInsertPOCert struct {
	InsertPOCertStruct apiModel.InsertPOCertStruct
}

func (cd *ConcreteInsertPOCert) InsertPOCertHash() model.InsertDataResponse {

	publicKey := "GD3EEFYWEP2XLLHONN2TRTQV4H5GSXJGCSUXZJGXGNZT4EFACOXEVLDJ"
	secretKey := "SA46OTS655ZDALIAODVCBWLWBXZWO6VUS6TU4U4GAIUVCKS2SYPDS7N4"
	var response model.InsertDataResponse
	// response.ProfileID = cd.InsertPOCertStruct.ProfileID
	response.TxnType = cd.InsertPOCertStruct.Type

	// netClient := commons.GetHorizonClient()
	// accountRequest := horizonclient.AccountRequest{AccountID: publicKey}
	// account, err := netClient.AccountDetail(accountRequest)
	kp,_ := keypair.Parse(publicKey)
	client := horizonclient.DefaultTestNetClient
	accountRequest := horizonclient.AccountRequest{AccountID: kp.Address()}
	account, err := client.AccountDetail(accountRequest)
	if err != nil {
		// log.Fatal(err)
	}

	typeTXNBuilder := txnbuild.ManageData{Name: "Transaction Type", Value: []byte(cd.InsertPOCertStruct.Type)}
	CertTypeTXNBuilder := txnbuild.ManageData{Name: "CertType", Value: []byte(cd.InsertPOCertStruct.CertType)}
	CertBodyTXNBuilder := txnbuild.ManageData{Name: "CertBody", Value:  []byte(cd.InsertPOCertStruct.CertType)}
	ValidityTXNBuilder := txnbuild.ManageData{Name: "Validity", Value:[]byte(cd.InsertPOCertStruct.Validity)}
	IssuedTXNBuilder := txnbuild.ManageData{Name: "Issued", Value: []byte(cd.InsertPOCertStruct.Issued)}
	ExpiredTXNBuilder := txnbuild.ManageData{Name: "Expired", Value: []byte(cd.InsertPOCertStruct.Issued)}
	// BUILD THE GATEWAY XDR
	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount:        &account,
		IncrementSequenceNum: true,
		Operations:           []txnbuild.Operation{&typeTXNBuilder, &CertTypeTXNBuilder, &CertBodyTXNBuilder, &ValidityTXNBuilder,&IssuedTXNBuilder,&ExpiredTXNBuilder},
		BaseFee:              txnbuild.MinBaseFee,
		Memo:                 nil,
		Preconditions:        txnbuild.Preconditions{},
	})
	
	// // save data
	// tx, err := build.Transaction(
	// 	commons.GetHorizonNetwork(),
	// 	build.SourceAccount{publicKey},
	// 	build.AutoSequence{commons.GetHorizonClient()},
	// 	build.SetData("Transaction Type", []byte(cd.InsertPOCertStruct.Type)),
	// 	build.SetData("CertType", []byte(cd.InsertPOCertStruct.CertType)),
	// 	build.SetData("CertBody", []byte(cd.InsertPOCertStruct.CertBody)),
	// 	build.SetData("Validity", []byte(cd.InsertPOCertStruct.Validity)),
	// 	build.SetData("Issued", []byte(cd.InsertPOCertStruct.Issued)),
	// 	build.SetData("Expired", []byte(cd.InsertPOCertStruct.Expired)),
	// )

	if err != nil {
		// panic(err)
		response.Error.Code = http.StatusNotFound
		response.Error.Message = "The HTTP request failed for InsertPOAHash "
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
	resp, err := commons.GetHorizonClient().SubmitTransactionXDR(txeB64)
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
