package businessFacades

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/builder"
	"github.com/dileepaj/tracified-gateway/proofs/deprecatedBuilder"
	"github.com/gorilla/mux"
	"github.com/stellar/go/build"
	"github.com/stellar/go/xdr"
)

/*Transaction - Deprecated
@author - Azeem Ashraf, Jajeththanan Sabapathipillai
@params - ResponseWriter,Request
*/
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

			display := &deprecatedBuilder.AbstractGenesisInsert{InsertGenesisStruct: GObj}
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

			display := &deprecatedBuilder.AbstractProfileInsert{InsertProfileStruct: PObj}
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
			display := &deprecatedBuilder.AbstractTDPInsert{XDR: TDP.XDR}
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

			display := &deprecatedBuilder.AbstractSplitProfile{SplitProfileStruct: SplitObj}
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

			display := &deprecatedBuilder.AbstractMergeProfile{MergeProfileStruct: MergeObj}
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

			display := &deprecatedBuilder.AbstractPOAInsert{InsertPOAStruct: POA}
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

			display := &deprecatedBuilder.AbstractPOCertInsert{InsertPOCertStruct: Cert}
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


/*SubmitGenesis @desc Handles an incoming request and calls the genesisBuilder
@author - Azeem Ashraf
@params - ResponseWriter,Request
*/
func SubmitGenesis(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var TDP []model.TransactionCollectionBody

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

	err := json.NewDecoder(r.Body).Decode(&TDP)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		fmt.Println(err)

		return
	}
	fmt.Println(TDP)

	display := &builder.AbstractXDRSubmiter{TxnBody: TDP}
	display.SubmitSpecial(w, r)
	// 	display.SubmitGenesis(w,r)

	return
}

/*SubmitData - @desc Handles an incoming request and calls the dataBuilder
@author - Azeem Ashraf
@params - ResponseWriter,Request
*/
func SubmitData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var TDP []model.TransactionCollectionBody

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

	err := json.NewDecoder(r.Body).Decode(&TDP)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		fmt.Println(err)
		return
	}
	fmt.Println(TDP)

	display := &builder.AbstractXDRSubmiter{TxnBody: TDP}
	// display.SubmitData(w,r,true)
	display.SubmitSpecial(w, r)

	return
}

/*SubmitSplit - @desc Handles an incoming request and calls the splitBuilder
@author - Azeem Ashraf
@params - ResponseWriter,Request
*/
func SubmitSplit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var TDP []model.TransactionCollectionBody

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

	err := json.NewDecoder(r.Body).Decode(&TDP)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		fmt.Println(err)
		return
	}
	fmt.Println(TDP)

	display := &builder.AbstractXDRSubmiter{TxnBody: TDP}
	display.SubmitSplit(w, r)

	return
}

/*SubmitMerge - @desc Handles an incoming request and calls the mergeBuilder
@author - Azeem Ashraf
@params - ResponseWriter,Request
*/
func SubmitMerge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var TDP []model.TransactionCollectionBody

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

	err := json.NewDecoder(r.Body).Decode(&TDP)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		fmt.Println(err)
		return

	}
	fmt.Println(TDP)

	display := &builder.AbstractXDRSubmiter{TxnBody: TDP}
	display.SubmitMerge(w, r)

	return
}

/*SubmitTransformation - Needs to be Refurbished @desc Handles an incoming request and calls the TransformationBuilder
@author - Azeem Ashraf
@params - ResponseWriter,Request
*/
func SubmitTransformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var TDP []model.TransactionCollectionBody

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

	err := json.NewDecoder(r.Body).Decode(&TDP)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		fmt.Println(err)
		return
	}
	fmt.Println(TDP)

	display := &builder.AbstractXDRSubmiter{TxnBody: TDP}
	display.SubmitMerge(w, r)

	return
}

/*SubmitTransfer - Needs to be Refurbished @desc Handles an incoming request and calls the TransferBuilder
@author - Azeem Ashraf
@params - ResponseWriter,Request
*/
func SubmitTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var TDP []model.TransactionCollectionBody

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

	err := json.NewDecoder(r.Body).Decode(&TDP)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		fmt.Println(err)
		return
	}
	fmt.Println(TDP)

	display := &builder.AbstractXDRSubmiter{TxnBody: TDP}
	display.SubmitSpecialTransfer(w, r)

	// status, _ := builder.XDRSubmitter(TDP)
	// if status {
	// 	w.WriteHeader(http.StatusOK)
	// 	result := apiModel.SubmitXDRSuccess{
	// 		Status: "Success",
	// 	}
	// 	json.NewEncoder(w).Encode(result)
	// } else {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	result := apiModel.SubmitXDRSuccess{
	// 		Status: "Failed",
	// 	}
	// 	json.NewEncoder(w).Encode(result)
	// }
	return
}


/*SubmitCertificateInsert - @desc Handles an incoming request and calls the CertificateInsertBuilder
@author - Azeem Ashraf
@params - ResponseWriter,Request
*/
func SubmitCertificateInsert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var TDP model.CertificateCollectionBody

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

	err := json.NewDecoder(r.Body).Decode(&TDP)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		fmt.Println(err)
		return
	}
	fmt.Println(TDP)

	var temp []model.CertificateCollectionBody
	temp = append(temp, TDP)
	display := &builder.AbstractCertificateSubmiter{TxnBody: temp}
	display.SubmitInsertCertificate(w, r)
	return
}

/*SubmitCertificateRenewal - @desc Handles an incoming request and calls the CertificateRevewalBuilder
@author - Azeem Ashraf
@params - ResponseWriter,Request
*/
func SubmitCertificateRenewal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var TDP model.CertificateCollectionBody

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

	err := json.NewDecoder(r.Body).Decode(&TDP)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		fmt.Println(err)
		return
	}
	fmt.Println(TDP)

	var temp []model.CertificateCollectionBody
	temp = append(temp, TDP)
	display := &builder.AbstractCertificateSubmiter{TxnBody: temp}
	display.SubmitRenewCertificate(w, r)
	return
}

/*SubmitCertificateRevoke - @desc Handles an incoming request and calls the CertificateRevokeBuilder
@author - Azeem Ashraf
@params - ResponseWriter,Request
*/
func SubmitCertificateRevoke(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var TDP model.CertificateCollectionBody

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

	err := json.NewDecoder(r.Body).Decode(&TDP)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		fmt.Println(err)
		return
	}
	fmt.Println(TDP)

	var temp []model.CertificateCollectionBody
	temp = append(temp, TDP)
	display := &builder.AbstractCertificateSubmiter{TxnBody: temp}
	display.SubmitRevokeCertificate(w, r)
	return
}

/*LastTxn - @desc Handles an incoming request and Returns the Last TXN for the Identifier in the Params
@author - Azeem Ashraf
@params - ResponseWriter,Request
*/
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

type TranXDR struct {
	XDR string `json:"XDR"`
}

/*ConvertXDRToTXN - Test Endpoint @desc Handles an incoming request and Returns the TXN Hash for teh XDR Provided
@author - Azeem Ashraf
@params - ResponseWriter,Request
*/
func ConvertXDRToTXN(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var Trans xdr.Transaction
	// var lol string

	var TDP TranXDR
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
		fmt.Println(err)
		return
	}
	fmt.Println(TDP)

	err1 := xdr.SafeUnmarshalBase64(TDP.XDR, &Trans)
	if err1 != nil {
		fmt.Println(err1)
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

/*TDPForTXN - Test Endpoint @desc Handles an incoming request and Returns the TDP ID for the TXN Provided.
@author - Azeem Ashraf
@params - ResponseWriter,Request
*/
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

/*TXNForTDP - Test Endpoint @desc Handles an incoming request and Returns the TXN ID for the TDP ID Provided.
@author - Azeem Ashraf
@params - ResponseWriter,Request
*/
func TXNForTDP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)

	object := dao.Connection{}
	p := object.GetTransactionForTdpId(vars["Txn"])
	p.Then(func(data interface{}) interface{} {

		result := data.(model.TransactionCollectionBody)

		// res := TDP{TdpId: result.TdpId}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "TdpId Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	})
	p.Await()

}

func ArtifactTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	fmt.Println("lol")
	var Artifacts model.ArtifactTransaction
	fmt.Println("lol")
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
	err := json.NewDecoder(r.Body).Decode(&Artifacts)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		fmt.Println(err)
		return
	}
	fmt.Println(Artifacts)
	// fmt.Println(TDPs)
	object := dao.Connection{}
	err2 := object.InsertArtifact(Artifacts)
	if err2 != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Failed",
		}
		json.NewEncoder(w).Encode(result)
		return

	} else {
		w.WriteHeader(http.StatusOK)
		result := apiModel.SubmitXDRSuccess{
			Status: "Success",
		}
		json.NewEncoder(w).Encode(result)
		return
	}

}
