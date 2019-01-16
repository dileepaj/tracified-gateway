package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"

	// "strings"
	"github.com/gorilla/mux"
	"github.com/stellar/go/build"
	"github.com/stellar/go/xdr"
	// "github.com/stellar/go/xdr"
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/builder"
)

func Transaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	TType := (vars["TType"])
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Please send a request body")
		return
	} else {
		switch TType {
		case "0":
			var GObj apiModel.InsertGenesisStruct
			err := json.NewDecoder(r.Body).Decode(&GObj)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("Error while Decoding the body")
				fmt.Println(err)
				return
			}
			fmt.Println(GObj)
			result := model.InsertGenesisResponse{}

			display := &builder.AbstractGenesisInsert{InsertGenesisStruct: GObj}
			result = display.GenesisInsert()

			w.WriteHeader(result.Error.Code)
			result2 := apiModel.GenesisSuccess{
				Message:     result.Error.Message,
				ProfileTxn:  result.ProfileTxn,
				GenesisTxn:  result.GenesisTxn,
				Identifiers: GObj.Identifier,
				Type:        GObj.Type}
			json.NewEncoder(w).Encode(result2)

		case "1":
			var PObj apiModel.InsertProfileStruct
			err := json.NewDecoder(r.Body).Decode(&PObj)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("Error while Decoding the body")
				fmt.Println(err)
				return
			}
			fmt.Println(PObj)
			response := model.InsertProfileResponse{}

			display := &builder.AbstractProfileInsert{InsertProfileStruct: PObj}
			response = display.ProfileInsert()

			w.WriteHeader(response.Error.Code)
			result := apiModel.ProfileSuccess{
				Message:           response.Error.Message,
				ProfileTxn:        response.ProfileTxn,
				PreviousTXNID:     response.PreviousTXNID,
				PreviousProfileID: response.PreviousProfileID,
				Identifiers:       PObj.Identifier,
				Type:              PObj.Type}
			json.NewEncoder(w).Encode(result)
		case "2":
			var TDP apiModel.TestTDP
			err := json.NewDecoder(r.Body).Decode(&TDP)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("Error while Decoding the body")
				fmt.Println(err)
				return
			}
			fmt.Println(TDP)
			response := model.SubmitXDRResponse{}

			// display := &builder.AbstractTDPInsert{Hash: TObj.Data, InsertType: TType, PreviousTXNID: TObj.PreviousTXNID[0], ProfileId: TObj.ProfileID[0]}
			display := &builder.AbstractTDPInsert{XDR: TDP.XDR}
			response = display.TDPInsert()

			w.WriteHeader(response.Error.Code)
			result := apiModel.InsertSuccess{
				Message:   response.Error.Message,
				TxNHash:   response.TDPID,
				ProfileID: "response.ProfileID",
				Type:      "TDP.Type"}
			json.NewEncoder(w).Encode(result)

		case "5":
			var SplitObj apiModel.SplitProfileStruct
			err := json.NewDecoder(r.Body).Decode(&SplitObj)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("Error while Decoding the body")
				fmt.Println(err)
				return
			}
			fmt.Println(SplitObj)
			response := model.SplitProfileResponse{}

			// for i := 0; i < len(TObj.Identifiers); i++ {

			display := &builder.AbstractSplitProfile{SplitProfileStruct: SplitObj}
			response = display.ProfileSplit()
			// 	SplitProfiles = append(SplitProfiles, response.Txn)
			// }

			w.WriteHeader(response.Error.Code)
			result := apiModel.SplitSuccess{
				Message:          response.Error.Message,
				TxnHash:          response.Txn,
				PreviousTXNID:    response.PreviousTXNID,
				SplitProfiles:    response.SplitProfiles,
				SplitTXN:         response.SplitTXN,
				Identifier:       SplitObj.Identifier,
				SplitIdentifiers: SplitObj.SplitIdentifiers,
				Type:             TType}
			json.NewEncoder(w).Encode(result)
		case "6":
			var MergeObj apiModel.MergeProfileStruct
			err := json.NewDecoder(r.Body).Decode(&MergeObj)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("Error while Decoding the body")
				fmt.Println(err)
				return
			}
			fmt.Println(MergeObj)
			response := model.MergeProfileResponse{}

			display := &builder.AbstractMergeProfile{MergeProfileStruct: MergeObj}
			response = display.ProfileMerge()

			w.WriteHeader(response.Error.Code)
			result := apiModel.MergeSuccess{
				Message:            response.Error.Message,
				TxnHash:            response.Txn,
				PreviousTXNID:      response.PreviousTXNID,
				ProfileID:          response.ProfileID,
				Identifier:         MergeObj.Identifier,
				Type:               TType,
				MergingIdentifiers: response.PreviousIdentifiers,
				MergeTXNs:          response.MergeTXNs}
			json.NewEncoder(w).Encode(result)

		case "10":
			var POA apiModel.InsertPOAStruct
			err := json.NewDecoder(r.Body).Decode(&POA)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("Error while Decoding the body")
				fmt.Println(err)
				return
			}
			fmt.Println(POA)
			response := model.InsertDataResponse{}

			display := &builder.AbstractPOAInsert{InsertPOAStruct: POA}
			response = display.POAInsert()

			w.WriteHeader(response.Error.Code)
			result := apiModel.InsertSuccess{
				Message:   response.Error.Message,
				TxNHash:   response.TDPID,
				ProfileID: response.ProfileID,
				Type:      POA.Type}
			json.NewEncoder(w).Encode(result)

		case "11":
			var Cert apiModel.InsertPOCertStruct
			err := json.NewDecoder(r.Body).Decode(&Cert)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode("Error while Decoding the body")
				fmt.Println(err)
				return
			}
			fmt.Println(Cert)
			response := model.InsertDataResponse{}

			display := &builder.AbstractPOCertInsert{InsertPOCertStruct: Cert}
			response = display.POCertInsert()

			w.WriteHeader(response.Error.Code)
			result := apiModel.InsertSuccess{
				Message:   response.Error.Message,
				TxNHash:   response.TDPID,
				ProfileID: response.ProfileID,
				Type:      Cert.Type}
			json.NewEncoder(w).Encode(result)

		default:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Please send a valid Transaction Type")
			return
		}

	}
	return

}

func SubmitXDR(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var TDP []model.TransactionCollectionBody
	// object := dao.Connection{}
	// var copy model.TransactionCollectionBody

	if r.Header == nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "No Header present!",
		}
		json.NewEncoder(w).Encode(result)

		return
	}

	if r.Header.Get("Content-Type") == "" {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "No Content-Type present!",
		}
		json.NewEncoder(w).Encode(result)

		return
	}

	// fmt.Println(TDP)
	err := json.NewDecoder(r.Body).Decode(&TDP)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		// fmt.Println(err)
		return
	}
	fmt.Println(TDP)

	status, _ := builder.XDRSubmitter(TDP)
	if status {
		w.WriteHeader(http.StatusOK)
		result := apiModel.SubmitXDRSuccess{
			Status: "Success",
		}
		json.NewEncoder(w).Encode(result)

	}

	return
}

func LastTxn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	object := dao.Connection{}
	p := object.GetLastTransactionbyIdentifier(vars["Identifier"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.TransactionCollectionBody)
		res := model.LastTxnResponse{LastTxn: result.TxnHash}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Identifier Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}

type Transuc struct {
	TXN string `json:"txn"`
}

func ConvertXDRToTXN(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	var Trans xdr.Transaction
	// var lol string

	err := xdr.SafeUnmarshalBase64(vars["XDR"], &Trans)
	if err != nil {
		fmt.Println(err)
	}

	brr := build.TransactionBuilder{TX: &Trans, NetworkPassphrase: build.TestNetwork.Passphrase}
	fmt.Println(build.TestNetwork.Passphrase)
	// fmt.Println(brr.Hash())
	t, _ := brr.Hash()
	test := fmt.Sprintf("%x", t)

	w.WriteHeader(http.StatusOK)
	response := Transuc{TXN: test}
	json.NewEncoder(w).Encode(response)
	return 

}

type TDP struct {
	TdpId string `json:"tdpId"`
}

func TDPForTXN(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	object := dao.Connection{}
	p := object.GetTdpIdForTransaction(vars["Txn"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.TransactionCollectionBody)

		res := TDP{TdpId: result.TdpId}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "TdpId Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}

// func SplitXDR(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	var TDP []model.TransactionCollectionBody
// 	object := dao.Connection{}
// 	var copy model.TransactionCollectionBody

// 	if r.Header == nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode("No Header present!")
// 		return
// 	}

// 	if r.Header.Get("Content-Type") == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode("No Content-Type present!")
// 		return
// 	}

// 	err := json.NewDecoder(r.Body).Decode(&TDP)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode("Error while Decoding the body")
// 		// fmt.Println(err)
// 		return
// 	}
// 	for i := 0; i < len(TDP); i++ {
// 		TDP[i].Status = "Pending"
// 		var txe xdr.Transaction
// 		err = xdr.SafeUnmarshalBase64(TDP[i].XDR, &txe)
// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 		TDP[i].PublicKey = txe.SourceAccount.Address()
// 		TxnType := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
// 		TDP[i].TxnType = TxnType
// 		TDP[i].Status = "pending"

// 		copy=TDP[i]
// 		err1 := object.InsertTransaction(TDP[i])
// 		if err1 != nil {
// 			TDP[i].Status = "failed"
// 		}

// 	}
// 	for i := 0; i < len(TDP); i++ {
// 		display := &builder.AbstractTDPInsert{XDR: TDP[i].XDR}
// 		response := display.TDPInsert()
// 		if response.Error.Code == 503 {
// 			TDP[i].Status = "pending"
// 		} else {
// 			TDP[i].TxnHash = response.TXNID

// 			upd := model.TransactionCollectionBody{TxnHash: response.TXNID, Status: "done"}
// 			err2 := object.UpdateTransaction(copy, upd)
// 			if err2 != nil {
// 				TDP[i].Status = "pending"
// 			} else {
// 				TDP[i].Status = "done"
// 			}
// 		}
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	result := apiModel.SubmitXDRSuccess{
// 		Message: "Success, Please check each transaction status below",
// 		Txns:    TDP,
// 	}
// 	json.NewEncoder(w).Encode(result)
// 	return
// }
