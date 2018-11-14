package businessFacades

import (
	"encoding/json"
	"fmt"
	"main/dao"
	"strings"

	// "main/proofs/retriever/stellarRetriever"
	// "bytes"
	"net/http"

	// "github.com/stellar/go/build"
	// "github.com/stellar/go/clients/horizon"
	// "github.com/stellar/go/keypair"
	// "github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"

	"github.com/gorilla/mux"

	"main/api/apiModel"
	"main/model"
	"main/proofs/builder"
	// "main/proofs/builder"
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

	var TDP model.TransactionCollectionBody
	err := json.NewDecoder(r.Body).Decode(&TDP)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while Decoding the body")
		fmt.Println(err)
		return
	}

	var txe xdr.Transaction
	err = xdr.SafeUnmarshalBase64(TDP.XDR, &txe)
	if err != nil {
		fmt.Println(err)
	}

	var test xdr.TransactionEnvelope
	err = xdr.SafeUnmarshalBase64(TDP.XDR, &test)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(txe.SourceAccount.Address())
	TDP.PublicKey = txe.SourceAccount.Address()
	// TDP.TxnHash=txe.

	// fmt.Println(txe.Operations[1].Body.ManageDataOp.DataValue)

	fmt.Println(len(txe.Operations))
	TxnType := strings.TrimLeft(fmt.Sprintf("%s", txe.Operations[0].Body.ManageDataOp.DataValue), "&")
	TDP.TxnType = TxnType
	fmt.Println(TxnType)
	fmt.Println(txe.SeqNum)
	TDP.Status = "pending"

	// txe.SourceAccount.Address()

	// previous := model.TransactionCollectionBody{}
	object := dao.Connection{}
	// p := object.GetLastTransactionbyIdentifier(TDP.Identifier)
	// p.Then(func(data interface{}) interface{} {
	// 	var body = data.(map[string]string)
	// 	previous.TxnHash = body["TxnHash"]

	// 	return nil
	// }).Catch(func(error error) error {
	// 	return error
	// })
	// p.Await()
	// if previous != (model.TransactionCollectionBody{}) {
	// 	obj := stellarRetriever.ConcretePrevious{Count: 0}
	// 	data, er := obj.RetrievePrevious8Transactions(previous.TxnHash)
	// 	if er != nil {
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		json.NewEncoder(w).Encode(er)
	// 		fmt.Println(er)
	// 		return
	// 	}
	// 	if len(data.HashList) == 8 {

	/////////////////

	err1 := object.InsertTransaction(TDP)
	if err1 != nil {
		w.WriteHeader(http.StatusNotFound)
		result := apiModel.InsertCOCCollectionResponse{
			Message: "Failed"}
		json.NewEncoder(w).Encode(result)
		return
	}

	display := &builder.AbstractTDPInsert{XDR: TDP.XDR}
	response := display.TDPInsert()

	if response.Error.Code == 404 {
		w.WriteHeader(response.Error.Code)
		result := apiModel.SubmitXDRSuccess{
			Message:    response.Error.Message,
			TdpId:      TDP.TdpID,
			PublicKey:  TDP.PublicKey,
			Identifier: TDP.Identifier,
			Type:       TDP.TxnType}
		json.NewEncoder(w).Encode(result)
		return
	}

	upd := model.TransactionCollectionBody{TxnHash: response.TXNID, Status: "done"}
	err2 := object.UpdateTransaction(TDP, upd)

	if err2 != nil {
		w.WriteHeader(response.Error.Code)
		result := apiModel.SubmitXDRSuccess{
			Message:    response.Error.Message,
			TxNHash:    response.TXNID,
			TdpId:      TDP.TdpID,
			PublicKey:  TDP.PublicKey,
			Identifier: TDP.Identifier,
			Type:       TDP.TxnType,
			Status:     "pending"}
		json.NewEncoder(w).Encode(result)
		return
	}

	w.WriteHeader(response.Error.Code)
	result := apiModel.SubmitXDRSuccess{
		Message:    response.Error.Message,
		TxNHash:    response.TXNID,
		TdpId:      TDP.TdpID,
		PublicKey:  TDP.PublicKey,
		Identifier: TDP.Identifier,
		Type:       TDP.TxnType,
		Status:     "done"}
	json.NewEncoder(w).Encode(result)
	return

	//////////////////

	// 	} else {
	// 		// 	b := object.GetTransactionsbyIdentifier(TDP.Identifier)
	// 		// 	b.Then(func(data interface{}) interface{} {
	// 		// 		var body = data.(map[string]string)

	// 		// 		return nil
	// 		// 	}).Catch(func(error error) error {
	// 		// 		return error
	// 		// 	})
	// 		// 	b.Await()
	// 	}
	// }
	// // var TransactionBD model.TransactionCollectionBody
	// // TransactionBD := model.TransactionCollectionBody{XDR: TDP.XDR, Identifier: TDP.Identifier, PublicKey: TDP.PublicKey, TdpID: TDP.TdpId}
	// // object := dao.Connection{}
	// // err1 := object.InsertTransaction(TransactionBD)

	// // response := model.InsertDataResponse{}

	// // // display := &builder.AbstractTDPInsert{Hash: TObj.Data, InsertType: TType, PreviousTXNID: TObj.PreviousTXNID[0], ProfileId: TObj.ProfileID[0]}
	// // display := &stellarExecuter.ConcreteSubmitXDR{InsertTDP: TDP}
	// // response = display.SubmitXDR()

	// w.WriteHeader(http.StatusOK)
	// result := ""
	// json.NewEncoder(w).Encode(result)
}

func LastTxn(w http.ResponseWriter, r *http.Request) {
	//
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
		w.WriteHeader(http.StatusNotFound)
		response := model.Error{Message: "Not Found"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}
