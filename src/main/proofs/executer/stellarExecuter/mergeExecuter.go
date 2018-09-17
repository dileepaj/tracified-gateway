package stellarExecuter

import (
	// "encoding/base64"
	// "encoding/json"
	"fmt"
	"main/api/apiModel"
	"main/model"
	"net/http"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

type ConcreteMerge struct {
	MergeProfileStruct apiModel.MergeProfileStruct
	// MergingTXNs        []string
	// PreviousTXNID      string
	// PreviousProfileID  string
	// Identifiers        string
	// MergingIdentifiers []string
	// InsertType         string
	// ProfileID          string
	// Assets             string
	// Code               string
}

func (cd *ConcreteMerge) InsertMerge() model.MergeProfileResponse {

	// publicKey := "GAEO4AVTWOD6YRC3WFYYXFR6EYYRD2MYKLBB6XTHC3YDUPIEXEIKD5C3"
	secretKey := "SBSEIZJJXYL6SIC5Y2RDYEQYSBBSRTPSAPGBQPKXGLHC5TZZBC3TSYLC"
	var response model.MergeProfileResponse

	var MergeTXN []string

	if len(cd.MergeProfileStruct.MergingIdentifiers) >= 1 {

		for i := 0; i < len(cd.MergeProfileStruct.MergingIdentifiers); i++ {

			tx, err := build.Transaction(
				build.TestNetwork,
				build.SourceAccount{secretKey},
				build.AutoSequence{horizon.DefaultTestNetClient},
				build.SetData("TransactionType", []byte(cd.MergeProfileStruct.Type)),
				build.SetData("PreviousTXNID", []byte(cd.MergeProfileStruct.PreviousTXNID)),
				build.SetData("ProfileID", []byte(cd.MergeProfileStruct.ProfileID)),
				build.SetData("MergingTXN", []byte(cd.MergeProfileStruct.MergingTXNs[i])),
				build.SetData("Assets", []byte(cd.MergeProfileStruct.Assets)),
				build.SetData("Code", []byte(cd.MergeProfileStruct.Code)),
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
			response.PreviousTXNID = cd.MergeProfileStruct.PreviousTXNID

			cd.MergeProfileStruct.PreviousTXNID = resp.Hash
			MergeTXN = append(MergeTXN, resp.Hash)

		}

	}

	// data, _ := json.Marshal(&cd.MergeProfiles)
	// sEnc := base64.StdEncoding.EncodeToString(data)
	// fmt.Println(string(sEnc))

	// save data
	// response.PreviousProfileID=cd.PreviousProfileID

	response.MergeTXNs = MergeTXN
	response.PreviousIdentifiers = cd.MergeProfileStruct.MergingIdentifiers
	response.ProfileID = cd.MergeProfileStruct.ProfileID
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
