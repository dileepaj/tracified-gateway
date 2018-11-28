package businessFacades

import (
	// "io/ioutil"
	"github.com/tracified-gateway/dao"
	"github.com/tracified-gateway/proofs/retriever/stellarRetriever"

	// "gopkg.in/mgo.v2/bson"
	// "github.com/fanliao/go-promise"

	// "gopkg.in/mgo.v2"
	"github.com/stellar/go/build"
	"github.com/stellar/go/xdr"

	"encoding/json"
	"fmt"
	// "gopkg.in/mgo.v2"

	"net/http"

	// "github.com/fanliao/go-promise"
	"github.com/gorilla/mux"

	"github.com/tracified-gateway/api/apiModel"
	// "github.com/tracified-gateway/dao"
	"github.com/tracified-gateway/model"
	"github.com/tracified-gateway/proofs/builder"
	// "github.com/tracified-gateway/proofs/interpreter"
)

func CreateTrust(w http.ResponseWriter, r *http.Request) {

	// var response model.POE
	var TObj apiModel.TrustlineStruct
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Please send a request body")
		return
	} else {
		err := json.NewDecoder(r.Body).Decode(&TObj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Error while Decoding the body")
			fmt.Println(err)
			return
		}
		display := &builder.AbstractTrustline{TrustlineStruct: TObj}
		result := display.Trustline()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result2 := apiModel.PoeSuccess{Message: "TrustLine Created", TxNHash: result}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return
}
func SendAssests(w http.ResponseWriter, r *http.Request) {

	// var response model.POE
	var TObj apiModel.SendAssest
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Please send a request body")
		return
	} else {
		err := json.NewDecoder(r.Body).Decode(&TObj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Error while Decoding the body")
			fmt.Println(err)
			return
		}
		display := &builder.AbstractAssetTransfer{SendAssest: TObj}
		response := display.AssetTransfer()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(response.Error.Code)
		result := apiModel.SendAssetRes{Message: response.Error.Message, PreviousTXNID: response.PreviousTXNID, PreviousProfileID: response.PreviousProfileID, Code: response.Code, Amount: response.Amount, Txn: response.Txn, To: response.To, From: response.From}
		json.NewEncoder(w).Encode(result)
		return
	}
	return
}

func MultisigAccount(w http.ResponseWriter, r *http.Request) {

	var TObj apiModel.RegistrarAccount
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Please send a request body")
		return
	} else {
		err := json.NewDecoder(r.Body).Decode(&TObj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Error while Decoding the body")
			fmt.Println(err)
			return
		}

		// var response model.POE

		display := &builder.AbstractCreateRegistrar{RegistrarAccount: TObj}
		result := display.CreateRegistrarAcc()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result2 := apiModel.PoeSuccess{Message: "Success", TxNHash: result}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return

}

func AppointRegistrar(w http.ResponseWriter, r *http.Request) {

	var TObj apiModel.AppointRegistrar
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Please send a request body")
		return
	} else {
		err := json.NewDecoder(r.Body).Decode(&TObj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Error while Decoding the body")
			fmt.Println(err)
			return
		}

		display := &builder.AbstractAppointRegistrar{AppointRegistrar: TObj}
		result := display.AppointReg()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result2 := apiModel.RegSuccess{Message: "Success", Xdr: result}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return
}
func TransformV2(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)

	// var response model.POE
	var TObj apiModel.AssetTransfer
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Please send a request body")
		return
	} else {
		err := json.NewDecoder(r.Body).Decode(&TObj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Error while Decoding the body")
			fmt.Println(err)
			return
		}
		display := &builder.AbstractTransformAssets{AssetTransfer: TObj}
		result := display.TransformAssets()
		// display := &builder.AbstractTransformAssets{Code1: vars["code1"], Limit1: vars["limit1"], Code2: vars["code2"], Limit2: vars["limit2"], Code3: vars["code3"], Limit3: vars["limit3"], Code4: vars["code4"], Limit4: vars["limit4"]}
		// display := &builder.AbstractTransformAssets{Code1: TObj.Asset[0].Code, Limit1: TObj.Asset[0].Limit, Code2: TObj.Asset[1].Code, Limit2: TObj.Asset[1].Limit, Code3: TObj.Asset[2].Code, Limit3: TObj.Asset[2].Limit, Code4: TObj.Asset[3].Code, Limit4: TObj.Asset[3].Limit}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result2 := apiModel.RegSuccess{Message: "Success", Xdr: result}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return

}

func COC(w http.ResponseWriter, r *http.Request) {

	// var response model.POE
	var TObj apiModel.ChangeOfCustody
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Please send a request body")
		return
	} else {
		err := json.NewDecoder(r.Body).Decode(&TObj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Error while Decoding the body")
			fmt.Println(err)
			return
		}
		display := &builder.AbstractCoCTransaction{ChangeOfCustody: TObj}
		response := display.CoCTransaction()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(response.Error.Code)
		result2 := apiModel.COCRes{Message: response.Error.Message, PreviousTXNID: response.PreviousTXNID, PreviousProfileID: response.PreviousProfileID, Code: response.Code, Amount: response.Amount, To: response.To, From: response.From, TxnXDR: response.TxnXDR}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return
}

func COCLink(w http.ResponseWriter, r *http.Request) {

	// var response model.POE
	var TObj apiModel.ChangeOfCustodyLink
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Please send a request body")
		return
	} else {
		err := json.NewDecoder(r.Body).Decode(&TObj)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("Error while Decoding the body")
			fmt.Println(err)
			return
		}
		display := &builder.AbstractcocLink{ChangeOfCustodyLink: TObj}
		result := display.CoCLink()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result2 := apiModel.PoeSuccess{Message: "Success", TxNHash: result}
		json.NewEncoder(w).Encode(result2)
		return
	}
	return
}

func DeveloperRetriever(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var response model.POC

	pocStructObj := apiModel.POCStruct{Txn: vars["Txn"]}
	display := &stellarRetriever.ConcretePOC{POCStruct: pocStructObj}
	// display := &stellarRetriever.ConcretePOC{Txn: vars["Txn"]}
	response.RetrievePOC = display.RetrieveFullPOC()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(response.RetrievePOC.Error.Code)
	// w.WriteHeader(http.StatusBadRequest)

	// result := apiModel.PoeSuccess{Message: "response.RetrievePOC.Error.Message", TxNHash: "response.RetrievePOC.Txn"}
	result := apiModel.PocSuccess{
		Chain: response.RetrievePOC.BCHash}
	json.NewEncoder(w).Encode(result)

	return

}

func GetCocBySender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	object := dao.Connection{}
	p := object.GetCOCbySender(vars["Sender"])
	p.Then(func(data interface{}) interface{} {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		// result := apiModel.GetMultiCOCCollection{
		// 	Collection: data}
		json.NewEncoder(w).Encode(data)
		return data
	}).Catch(func(error error) error {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		// result := model.Error{Code: http.StatusNotFound,
		// 	Message: "No Results Found"}
		json.NewEncoder(w).Encode(error)
		return error
	})
	p.Await()

}

func GetCocByReceiver(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	object := dao.Connection{}
	p := object.GetCOCbyReceiver(vars["Receiver"])
	p.Then(func(data interface{}) interface{} {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		// result := apiModel.GetMultiCOCCollection{
		// 	Collection: data}
		json.NewEncoder(w).Encode(data)
		return data
	}).Catch(func(error error) error {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		// result := model.Error{Code: http.StatusNotFound,
		// 	Message: "No Results Found"}
		json.NewEncoder(w).Encode(error)
		return error
	})
	p.Await()

}

func InsertCocCollection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var GObj model.COCCollectionBody
	err := json.NewDecoder(r.Body).Decode(&GObj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while Decoding the body")
		fmt.Println(err)
		return
	}
	var accept xdr.Transaction
	var reject xdr.Transaction
	err = xdr.SafeUnmarshalBase64(GObj.AcceptXdr, &accept)
	if err != nil {
		fmt.Println(err)
	}

	brr := build.TransactionBuilder{TX: &accept, NetworkPassphrase: build.TestNetwork.Passphrase}
	fmt.Println(build.TestNetwork.Passphrase)
	// fmt.Println(brr.Hash())
	t, _ := brr.Hash()
	test := fmt.Sprintf("%x", t)

	err = xdr.SafeUnmarshalBase64(GObj.RejectXdr, &reject)
	if err != nil {
		fmt.Println(err)
	}

	brr1 := build.TransactionBuilder{TX: &reject, NetworkPassphrase: build.TestNetwork.Passphrase}
	fmt.Println(build.TestNetwork.Passphrase)
	// fmt.Println(brr.Hash())
	t1, _ := brr1.Hash()
	test1 := fmt.Sprintf("%x", t1)

	GObj.AcceptTxn = test
	GObj.RejectTxn = test1
	fmt.Println(GObj)
	object := dao.Connection{}
	err1 := object.InsertCoc(GObj)

	if err1 != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		result := apiModel.InsertCOCCollectionResponse{
			Message: "Failed"}
		json.NewEncoder(w).Encode(result)
		return
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		result := apiModel.InsertCOCCollectionResponse{
			Message: "Success", Body: GObj}
		json.NewEncoder(w).Encode(result)
		return
	}
}

func UpdateCocCollection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var GObj model.COCCollectionBody
	var selection model.COCCollectionBody
	var result apiModel.InsertCOCCollectionResponse

	err := json.NewDecoder(r.Body).Decode(&GObj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while Decoding the body")
		fmt.Println(err)
		return
	}
	fmt.Println(GObj)
	object := dao.Connection{}
	switch GObj.Status {
	case "accepted":
		p := object.GetCOCbyAcceptTxn(GObj.AcceptTxn)
		p.Then(func(data interface{}) interface{} {
			selection = data.(model.COCCollectionBody)
			// fmt.Println("OLD BODY")

			// fmt.Println(selection)
			display := &builder.AbstractTDPInsert{XDR: GObj.AcceptXdr}
			response := display.TDPInsert()
		
			if response.TXNID == "" {
				w.WriteHeader(response.Error.Code)
				result = apiModel.InsertCOCCollectionResponse{
					Message: "Failed"}
				json.NewEncoder(w).Encode(result)
			}else{
				GObj.TxnHash=response.TXNID
				err1 := object.UpdateCOC(selection, GObj)
				if err1 != nil {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusNotFound)
					result = apiModel.InsertCOCCollectionResponse{
						Message: "Pending"}
					json.NewEncoder(w).Encode(result)
					
				} else {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusOK)
					body:=selection
					body.AcceptTxn=GObj.AcceptTxn
					body.AcceptXdr=GObj.AcceptXdr
					body.Status=GObj.Status
					result = apiModel.InsertCOCCollectionResponse{
						Message: "Success", Body: body}
					json.NewEncoder(w).Encode(result)
					
				}
			}
			

			
			return data
		}).Catch(func(error error) error {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(error)
			return error
		})
		p.Await()
		break
	case "rejected":
		p := object.GetCOCbyRejectTxn(GObj.RejectTxn)
		p.Then(func(data interface{}) interface{} {
			selection = data.(model.COCCollectionBody)
			display := &builder.AbstractTDPInsert{XDR: GObj.RejectXdr}
			response := display.TDPInsert()
		
			if response.TXNID == "" {
				w.WriteHeader(response.Error.Code)
				result = apiModel.InsertCOCCollectionResponse{
					Message: "Pending"}
				json.NewEncoder(w).Encode(result)
			}else{
				GObj.TxnHash=response.TXNID
				err1 := object.UpdateCOC(selection, GObj)
				if err1 != nil {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusNotFound)
					result = apiModel.InsertCOCCollectionResponse{
						Message: "Failed"}
					json.NewEncoder(w).Encode(result)
					
				} else {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusOK)
					body:=selection
					body.RejectTxn=GObj.RejectTxn
					body.RejectXdr=GObj.RejectXdr
					body.Status=GObj.Status
					result = apiModel.InsertCOCCollectionResponse{
						Message: "Success", Body: body}
					json.NewEncoder(w).Encode(result)
					
				}
			}
			

			
			return data
		}).Catch(func(error error) error {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(error)
			return error
		})
		p.Await()
		break
	}
	
	

	// err1 := object.UpdateCOC(selection, GObj)
	// if err1 != nil {
	// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// 	w.WriteHeader(http.StatusNotFound)
	// 	result := apiModel.InsertCOCCollectionResponse{
	// 		Message: "Failed"}
	// 	json.NewEncoder(w).Encode(result)
	// 	return
	// } else {
	// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// 	w.WriteHeader(http.StatusOK)
	// 	result := apiModel.InsertCOCCollectionResponse{
	// 		Message: "Success", Body: GObj}
	// 	json.NewEncoder(w).Encode(result)
	// 	return
	// }
	return
}

// func InsertTransactionCollection(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	var GObj model.TransactionCollectionBody
// 	err := json.NewDecoder(r.Body).Decode(&GObj)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode("Error while Decoding the body")
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Println(GObj)
// 	object := dao.Connection{}
// 	err1 := object.InsertTransaction(GObj)

// 	if err1 != nil {
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		w.WriteHeader(http.StatusNotFound)
// 		result := apiModel.InsertTransactionCollectionResponse{
// 			Message: "Failed"}
// 		json.NewEncoder(w).Encode(result)
// 		return
// 	} else {
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		w.WriteHeader(http.StatusOK)
// 		result := apiModel.InsertTransactionCollectionResponse{
// 			Message: "Success", Body: GObj}
// 		json.NewEncoder(w).Encode(result)
// 		return
// 	}
// }
// func UpdateTransactionCollection(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	var GObj model.TransactionUpdate
// 	err := json.NewDecoder(r.Body).Decode(&GObj)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode("Error while Decoding the body")
// 		fmt.Println(err)
// 		return
// 	}

// 	fmt.Println(GObj)
// 	object := dao.Connection{}
// 	err1 := object.UpdateTransaction(GObj.Selector,GObj.Update)

// 	if err1 != nil {
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		w.WriteHeader(http.StatusNotFound)
// 		result := apiModel.InsertTransactionCollectionResponse{
// 			Message: "Failed"}
// 		json.NewEncoder(w).Encode(result)
// 		return
// 	} else {
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 		w.WriteHeader(http.StatusOK)
// 		result := apiModel.InsertTransactionCollectionResponse{
// 			Message: "Success", Body: GObj.Update}
// 		json.NewEncoder(w).Encode(result)
// 		return
// 	}
// }
