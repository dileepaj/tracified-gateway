package businessFacades

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/configs"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	fosponsoring "github.com/dileepaj/tracified-gateway/nft/stellar/FOSponsoring"
	"github.com/dileepaj/tracified-gateway/proofs/builder"
	"github.com/dileepaj/tracified-gateway/proofs/deprecatedBuilder"
	"github.com/dileepaj/tracified-gateway/services"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
)

/*
Transaction - Deprecated
@author - Azeem Ashraf, Jajeththanan Sabapathipillai
@params - ResponseWriter,Request
*/

var accountStatus string

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
				Type:        GObj.Type,
			}
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
				Type:              PObj.Type,
			}
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
				Type:      "TDP.Type",
			}
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
				Type:             TType,
			}
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
				MergeTXNs:          response.MergeTXNs,
			}
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
				Type:      POA.Type,
			}
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
				Type:      Cert.Type,
			}
			json.NewEncoder(w).Encode(result)

		default:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Please send a valid Transaction Type")
			return
		}
	}
	return
}

/*
SubmitData - @desc Handles an incoming request and calls the dataBuilder
@author - Azeem Ashraf
@params - ResponseWriter,Request
*/
func SubmitData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	requestId := r.Header.Get("Custom-Request-Tag-Id")
	var TDPs []model.TransactionCollectionBody
	utilities.BenchmarkLog(configs.BenchmarkLogsTag.TDP_REQUEST, configs.BenchmarkLogsAction.RECEIVED_FORM_BACKEND, requestId, configs.BenchmarkLogsStatus.OK)
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

	err := json.NewDecoder(r.Body).Decode(&TDPs)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "Error while Decoding the body",
		}
		json.NewEncoder(w).Encode(result)
		return
	}
	url := r.URL.String()
	var response []apiModel.TDPOperationRequest
	for i, TxnBody := range TDPs {
		TDPs[i].Status = "pending"
		TDPs[i].RequestId = requestId
		// Check if the URL contains "/merge"
		if strings.Contains(url, "/merge") {
			// "/merge" is part of the URL
			TDPs[i].MergeBlock = i
		}
		// Convert the struct to a JSON string using encoding/json
		jsonStr, err := json.Marshal(TDPs[i])
		if err != nil {
			response = append(response, apiModel.TDPOperationRequest{i, TxnBody.MapIdentifier, TxnBody.TdpId, TxnBody.XDR, "Error"})
			log.Error("Error in convert the struct to a JSON string using encoding/json:", err, " TxnBody: ", TxnBody)
			continue
		}
		services.PublishToQueue(configs.QueueTransaction.Prefix+TxnBody.UserID, string(jsonStr), configs.QueueTransaction.Method)
		utilities.BenchmarkLog("tdp-request", configs.BenchmarkLogsAction.PUBLISH_TO+configs.QueueTransaction.Prefix+TxnBody.UserID, requestId, configs.BenchmarkLogsStatus.OK)
		response = append(response, apiModel.TDPOperationRequest{i, TxnBody.MapIdentifier, TxnBody.TdpId, TxnBody.XDR, "Success"})
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	return
}

/*
SubmitCertificateInsert - @desc Handles an incoming request and calls the CertificateInsertBuilder
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

/*
SubmitCertificateRenewal - @desc Handles an incoming request and calls the CertificateRevewalBuilder
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

/*
SubmitCertificateRevoke - @desc Handles an incoming request and calls the CertificateRevokeBuilder
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

/*
LastTxn - @desc Handles an incoming request and Returns the Last TXN for the Identifier in the Params
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

/*
ConvertXDRToTXN - Test Endpoint @desc Handles an incoming request and Returns the TXN Hash for teh XDR Provided
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

	brr, _ := txnbuild.TransactionFromXDR(TDP.XDR)
	t, _ := brr.Hash(network.TestNetworkPassphrase)
	test := fmt.Sprintf("%x", t)

	w.WriteHeader(http.StatusOK)
	response := Transuc{TXN: test}
	json.NewEncoder(w).Encode(response)
	return
}

type TDP struct {
	TdpId string `json:"tdpId"`
}

/*
TDPForTXN - Test Endpoint @desc Handles an incoming request and Returns the TDP ID for the TXN Provided.
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

/*
TXNForTDP - Test Endpoint @desc Handles an incoming request and Returns the TXN ID for the TDP ID Provided.
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

func TxnForIdentifier(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var result []model.TransactionHashWithIdentifier
	vars := mux.Vars(r)
	object := dao.Connection{}
	p := object.GetRealIdentifierByMapValue(vars["identifier"])
	p.Then(func(data interface{}) interface{} {
		dbResult := data.([]model.TransactionCollectionBody)
		for _, TxnBody := range dbResult {
			result1, err := http.Get(commons.GetHorizonClient().HorizonURL + "transactions/" + TxnBody.TxnHash)
			if err != nil {
				logrus.Error("Unable to reach Stellar network in result1")
			}
			if result1.StatusCode != 200 {
				logrus.Error("Transaction could not be retrieved from Stellar Network in result1")
			}
			data, _ := ioutil.ReadAll(result1.Body)
			var raw map[string]interface{}
			json.Unmarshal(data, &raw)
			createdAt := fmt.Sprintf("%s", raw["created_at"])

			temp := model.TransactionHashWithIdentifier{
				Status:          TxnBody.Status,
				Txnhash:         TxnBody.TxnHash,
				TxnType:         GetTransactiontype(TxnBody.TxnType),
				Identifier:      TxnBody.Identifier,
				FromIdentifier1: TxnBody.FromIdentifier1,
				FromIdentifier2: TxnBody.FromIdentifier2,
				ToIdentifier:    TxnBody.ToIdentifier,
				AvailableProof:  GetProofName(TxnBody.TxnType),
				ProductName:     TxnBody.ProductName,
				ProductID:       TxnBody.ProductID,
				Timestamp:       createdAt,
			}
			result = append(result, temp)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return nil
	}).Catch(func(error error) error {
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Identifier Not Found in Gateway DataStore"}
		json.NewEncoder(w).Encode(response)
		return error
	}).Await()
}

func TxnForArtifact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var result []model.TransactionHashWithIdentifier
	vars := mux.Vars(r)
	object := dao.Connection{}

	// call backend to get identider by artifactId
	url := constants.TracifiedBackend + "/api/v2/identifiers/artifact/" + vars["artifactid"]
	bearer := "Bearer " + constants.BackendToken
	// Create a new request using http
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error("Error while create new request using http " + err.Error())
	}
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resq, err1 := client.Do(req)
	if err1 != nil {
		logrus.Error("Error while getting response " + err1.Error())
		w.WriteHeader(http.StatusBadRequest)
		response := model.Error{Message: "Connection to the Traceability DataStore was interupted " + err1.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	body, err3 := ioutil.ReadAll(resq.Body)
	if err3 != nil {
		logrus.Error("Error while ioutil.ReadAll(resq.Body) " + err3.Error())
	}
	var identifiers []string
	json.Unmarshal(body, &identifiers)
	if resq.StatusCode == 200 || resq.StatusCode == 204 {
		if len(identifiers) > 0 {
			p := object.GetRealIdentifiersByArtifactId(identifiers)
			p.Then(func(data interface{}) interface{} {
				dbResult := data.([]model.TransactionCollectionBody)
				for _, TxnBody := range dbResult {
					temp := model.TransactionHashWithIdentifier{
						Status:          TxnBody.Status,
						Txnhash:         TxnBody.TxnHash,
						Identifier:      TxnBody.Identifier,
						FromIdentifier1: TxnBody.FromIdentifier1,
						FromIdentifier2: TxnBody.FromIdentifier2,
						ToIdentifier:    TxnBody.ToIdentifier,
						TxnType:         GetTransactiontype(TxnBody.TxnType),
						AvailableProof:  GetProofName(TxnBody.TxnType),
						ProductID:       TxnBody.ProductID,
						ProductName:     TxnBody.ProductName,
					}
					result = append(result, temp)
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(result)
				return nil
			}).Catch(func(error error) error {
				w.WriteHeader(http.StatusBadRequest)
				response := model.Error{Message: error.Error()}
				json.NewEncoder(w).Encode(response)
				return error
			}).Await()
		} else {
			w.WriteHeader(http.StatusNoContent)
			response := model.Error{Message: "Can not find the identires for artifactid"}
			json.NewEncoder(w).Encode(response)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadGateway)
		response := model.Error{Message: "Connection to the Traceability DataStore was interupted "}
		json.NewEncoder(w).Encode(response)
		return
	}
}

func SubmitFOData(w http.ResponseWriter, r *http.Request) {
	if commons.GoDotEnvVariable("FONEW_FLAG") == "TRUE" {
		logger := utilities.NewCustomLogger()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var Response model.TransactionData
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&Response)
		if err != nil {
			logger.LogWriter("Error submitting data to the blockchain : "+err.Error(), constants.ERROR)
			w.WriteHeader(http.StatusBadRequest)
			response := model.Error{Message: "Error submitting data to the blockchain"}
			json.NewEncoder(w).Encode(response)
			return
		}
		if Response.XDR != "" && Response.FOUser != "" && Response.AccountIssuer != "" {
			resp, err := http.Get(commons.GetHorizonClient().HorizonURL + "accounts/" + Response.FOUser)
			if err != nil {
				logrus.Error("Error making HTTP request:", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusNotFound {
				accountStatus = "0"
			} else {
				accountStatus = "1"
			}
			if accountStatus == "0" {
				logrus.Error("Account of FO User is inactive")
			}

			if accountStatus == "1" {
				object := dao.Connection{}
				p := object.GetIssuerAccountByFOUser(Response.FOUser)
				rst, err := p.Await()
				if err != nil {
					logrus.Error("There was an error cant get issuer!")
				} else {
					data := rst.(model.TransactionDataKeys)
					if rst != nil || rst != "" {
						result1, err := http.Get(commons.GetHorizonClient().HorizonURL + "accounts/" + data.AccountIssuerPK)
						body, err := ioutil.ReadAll(result1.Body)
						if err != nil {
							log.Error("Error while read response " + err.Error())
						}
						var balances model.BalanceResponse
						err = json.Unmarshal(body, &balances)
						if err != nil {
							log.Error("Error while json.Unmarshal(body, &balance) " + err.Error())
						}

						balance := balances.Balances[0].Balance

						if balance < "10" {
							hash, err := fosponsoring.FundAccount(data.AccountIssuerPK)
							if err != nil {
								log.Error("Error while funding issuer " + err.Error())
							}
							logrus.Info("funded and hash is : ", hash)
						}

						TransactionPayload := model.TransactionData{
							FOUser:        Response.FOUser,
							AccountIssuer: data.AccountIssuerPK,
							XDR:           Response.XDR,
						}

						hash, err := fosponsoring.BuildSignedSponsoredXDR(TransactionPayload)
						if err != nil {
							log.Error(err)
						} else {
							w.Header().Set("Content-Type", "application/json;")
							w.WriteHeader(http.StatusOK)
							result := model.Hash{
								Hash: hash,
							}
							logrus.Info("Hash been passed to frontend : ", result)
							json.NewEncoder(w).Encode(result)
						}

					}
				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				response := model.Error{Message: "Can not create XDR and submit"}
				json.NewEncoder(w).Encode(response)
			}
		}
	}
}

func CreateSponsorer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if commons.GoDotEnvVariable("FONEW_FLAG") == "TRUE" {
		object := dao.Connection{}
		vars := mux.Vars(r)
		p := object.GetIssuerAccountByFOUser((vars["foUser"]))
		rst, err := p.Await()
		if err != nil {
			IssuerPK, EncodedIssuerSK, encSK, err := fosponsoring.CreateIssuerAccountForFOUser()
			if err != nil && IssuerPK == "" && EncodedIssuerSK == "" {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(err)
			} else {
				Keys := model.TransactionDataKeys{
					FOUser:          (vars["foUser"]),
					AccountIssuerPK: IssuerPK,
					AccountIssuerSK: encSK,
				}
				// adding the credentials to the DB
				object := dao.Connection{}
				err := object.InsertIssuingAccountKeys(Keys)
				if err != nil {
					logrus.Error(err)
				}

				// send the response
				result := model.NFTIssuerAccount{
					NFTIssuerPK: IssuerPK,
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(result)
			}
		} else {
			data := rst.(model.TransactionDataKeys)
			if rst != nil || rst != "" {
				logrus.Info("Issuer is:", data.AccountIssuerPK)
			}
			result := model.NFTIssuerAccount{
				NFTIssuerPK: data.AccountIssuerPK,
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(result)

		}
	}
}

func ActivateFOUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if commons.GoDotEnvVariable("FONEW_FLAG") == "TRUE" {
		var Response model.ActivateXDR
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&Response)
		tx, account, err := fosponsoring.SubmittingXDRs(Response.XDR, 4)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Error{Message: "Can not activate user " + account}
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusOK)
		result := model.Hash{
			Hash: tx,
		}
		json.NewEncoder(w).Encode(result)
	}
}
