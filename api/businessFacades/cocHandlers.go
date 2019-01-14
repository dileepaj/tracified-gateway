package businessFacades

import (
	// "io/ioutil"
	"github.com/dileepaj/tracified-gateway/dao"
	// "github.com/dileepaj/tracified-gateway/proofs/retriever/stellarRetriever"

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

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	// "github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/builder"
	// "github.com/dileepaj/tracified-gateway/proofs/interpreter"
)

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

			var TXNS []model.TransactionCollectionBody
			TXN:=model.TransactionCollectionBody{
				XDR:selection.AcceptXdr,
			}
			TXNS=append(TXNS,TXN)
			status,response:= builder.XDRSubmitter(TXNS)

			selection = data.(model.COCCollectionBody)
			// display := &builder.AbstractTDPInsert{XDR: GObj.AcceptXdr}
			// response := display.TDPInsert()

			if !status {
				w.WriteHeader(400)
				result = apiModel.InsertCOCCollectionResponse{
					Message: "Failed"}
				json.NewEncoder(w).Encode(result)
			}else{
				GObj.TxnHash=response.TXNID
				fmt.Println(response.TXNID)

				err1 := object.UpdateCOC(selection, GObj)
				if err1 != nil {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(400)
					result = apiModel.InsertCOCCollectionResponse{
						Message: "Failed"}
					json.NewEncoder(w).Encode(result)
					
				} else {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusOK)
					body:=GObj
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
			w.WriteHeader(400)
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
		
			if response.Error.Code == 404 {
				w.WriteHeader(400)
				result = apiModel.InsertCOCCollectionResponse{
					Message: "Failed"}
				json.NewEncoder(w).Encode(result)
			}else{
				GObj.TxnHash=response.TXNID
				fmt.Println(response.TXNID)
				err1 := object.UpdateCOC(selection, GObj)
				if err1 != nil {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(400)
					result = apiModel.InsertCOCCollectionResponse{
						Message: "Failed"}
					json.NewEncoder(w).Encode(result)
					
				} else {
					w.Header().Set("Content-Type", "application/json; charset=UTF-8")
					w.WriteHeader(http.StatusOK)
					body:=GObj
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
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(error)
			return error
		})
		p.Await()
		break

	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(400)
		result = apiModel.InsertCOCCollectionResponse{
			Message: "Failed, Status invalid"}
		json.NewEncoder(w).Encode(result)
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

func CheckAccountsStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var GObj apiModel.GetSubAccountStatus
	var result []apiModel.GetSubAccountStatusResponse

	err := json.NewDecoder(r.Body).Decode(&GObj)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while Decoding the body")
		fmt.Println(err)
		return
	}
	fmt.Println(GObj)
	object := dao.Connection{}
	for i:=0;i<len(GObj.SubAccounts);i++ {
		
		p := object.GetLastCOCbySubAccount(GObj.SubAccounts[i])
		p.Then(func(data interface{}) interface{} {
			result=append(result,data.(apiModel.GetSubAccountStatusResponse))
			return data
		}).Catch(func(error error) error {
			result=append(result,apiModel.GetSubAccountStatusResponse{SubAccount:GObj.SubAccounts[i],Available:true})

			return error
		})
		p.Await()

	
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
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
