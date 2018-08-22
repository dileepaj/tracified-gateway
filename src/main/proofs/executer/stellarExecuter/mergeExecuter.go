package stellarExecuter

import (
	// "encoding/base64"
	// "encoding/json"
	"fmt"
	"main/model"
	"net/http"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

type ConcreteMerge struct {
	// *builder.AbstractMergeProfile
	MergingTXNs        []string
	PreviousTXNID      string
	PreviousProfileID  string
	Identifiers        string
	MergingIdentifiers []string
	InsertType         string
	ProfileID          string
	Assets             string
	Code               string
}

func (cd *ConcreteMerge) InsertMerge() model.MergeProfileResponse {

	// publicKey := "GAEO4AVTWOD6YRC3WFYYXFR6EYYRD2MYKLBB6XTHC3YDUPIEXEIKD5C3"
	secretKey := "SBSEIZJJXYL6SIC5Y2RDYEQYSBBSRTPSAPGBQPKXGLHC5TZZBC3TSYLC"
	var response model.MergeProfileResponse

	var MergeTXN []string

	if len(cd.MergingIdentifiers) >= 1 {

		for i := 0; i < len(cd.MergingIdentifiers); i++ {

			tx, err := build.Transaction(
				build.TestNetwork,
				build.SourceAccount{secretKey},
				build.AutoSequence{horizon.DefaultTestNetClient},
				build.SetData("TransactionType", []byte(cd.InsertType)),
				build.SetData("PreviousTXNID", []byte(cd.PreviousTXNID)),
				build.SetData("ProfileID", []byte(cd.ProfileID)),
				build.SetData("MergingTXN", []byte(cd.MergingTXNs[i])),
				build.SetData("Assets", []byte(cd.Assets)),
				build.SetData("Code", []byte(cd.Code)),
			)

			if err != nil {
				// panic(err)
				response.Error.Code = http.StatusNotFound
				response.Error.Message = "The HTTP request failed for MergeProfile "
				fmt.Println(err)
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
			response.PreviousTXNID = cd.PreviousTXNID

			cd.PreviousTXNID = resp.Hash
			MergeTXN=append(MergeTXN,resp.Hash)

		}

	}

	// data, _ := json.Marshal(&cd.MergeProfiles)
	// sEnc := base64.StdEncoding.EncodeToString(data)
	// fmt.Println(string(sEnc))

	// save data
	// response.PreviousProfileID=cd.PreviousProfileID

	response.MergeTXNs = MergeTXN
	response.PreviousIdentifiers = cd.MergingIdentifiers
	response.ProfileID = cd.ProfileID
	return response

}

// func (cd *ConcreteMerge) InsertProfile() model.MergeProfileResponse {

// 	publicKey := "GAEO4AVTWOD6YRC3WFYYXFR6EYYRD2MYKLBB6XTHC3YDUPIEXEIKD5C3"
// 	// secretKey := "SBSEIZJJXYL6SIC5Y2RDYEQYSBBSRTPSAPGBQPKXGLHC5TZZBC3TSYLC"
// 	var response model.MergeProfileResponse
// 	response.PreviousTXNID = cd.PreviousTXNID
// 	response.PreviousProfileID = cd.ProfileID
// 	response.Identifiers = cd.Identifiers
// 	response.TxnType = cd.InsertType

// 	// save data
// 	tx, err := build.Transaction(
// 		build.TestNetwork,
// 		build.SourceAccount{publicKey},
// 		build.AutoSequence{horizon.DefaultTestNetClient},
// 		build.SetData("PreviousTXNID", []byte(cd.PreviousTXNID)),
// 		build.SetData("ProfileID", []byte(cd.ProfileID)),
// 		build.SetData("Identifiers", []byte(cd.Identifiers)),
// 	)

// 	if err != nil {
// 		// panic(err)
// 		response.Error.Code = http.StatusNotFound
// 		response.Error.Message = "The HTTP request failed for InsertProfile "
// 		return response
// 	}

// 	// Sign the transaction to prove you are actually the person sending it.
// 	txe, err := tx.Sign(publicKey)
// 	if err != nil {
// 		// panic(err)
// 		response.Error.Code = http.StatusNotFound
// 		response.Error.Message = "signing request failed for the Transaction"
// 		return response
// 	}

// 	txeB64, err := txe.Base64()
// 	if err != nil {
// 		// panic(err)
// 		response.Error.Code = http.StatusNotFound
// 		response.Error.Message = "Base64 conversion failed for the Transaction"
// 		return response
// 	}

// 	// And finally, send it off to Stellar!
// 	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
// 	if err != nil {
// 		// panic(err)
// 		response.Error.Code = http.StatusNotFound
// 		response.Error.Message = "Test net client crashed"
// 		return response
// 	}

// 	fmt.Println("Successful Transaction:")
// 	fmt.Println("Ledger:", resp.Ledger)
// 	fmt.Println("Hash:", resp.Hash)

// 	response.Error.Code = http.StatusOK
// 	response.Error.Message = "Transaction performed in the blockchain."
// 	response.Txn = resp.Hash

// 	return response

// }
