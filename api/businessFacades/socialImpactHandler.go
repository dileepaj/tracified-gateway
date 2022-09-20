package businessFacades

import (
	"encoding/json"
	"net/http"

	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

type bitString string

// BuildSocialImpactFormula --> this method take the farmulaJSOn from backend and map the equationID to 4 byte value and store it DB,
// generate  executionTemplete object via passing FCL query to FCL library,
// base on the executionTemplete objcet store the formula in stellar blckchin
func BuildSocialImpactFormula(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var formulaJSON model.FormulaBuildingRequest

	stellarprotocols.BuildVariableDefinitionManageData("12", "Testing", "Variable", "unit1", "5", "Hello there")

	err := json.NewDecoder(r.Body).Decode(&formulaJSON)
	if err != nil {
		logrus.Error(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Error while decoding the body")
		return
	}

	formulaArray := formulaJSON.Formula
	fieldCount:=0;
	for i,element := range formulaJSON.Formula{
		if element.Type =="DATA"{
			formulaArray[i].Type="VARIABLE"
		}else if element.Type =="CONSTANT" && element.MetricReferenceId!=""{
			formulaArray[i].Type="REFERREDCONSTANT"
		}else if element.Type =="CONSTANT" && element.MetricReferenceId==""{
			formulaArray[i].Type="SEMANTICCONSTANT"
		}
		if element.Type!="OPERATOR"{
			fieldCount++
		}
	}

	formulaJSON.Formula=formulaArray

	object := dao.Connection{}
	formulaMap, err5 := object.GetFormulaMapID(formulaJSON.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err5 != nil {
		logrus.Info("Unable to connect gateway datastore ", err5)
		w.WriteHeader(http.StatusNotFound)
	}
	if formulaMap != nil {
		logrus.Error("Formula Id is in gateway datastore")
		w.WriteHeader(http.StatusNoContent)
		response := model.Error{Code: http.StatusNoContent, Message: "Formula Id is in gateway datastore"}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := object.GetNextSequenceValue("FORMULAID")
	if err != nil {
		logrus.Error("GetNextSequenceValu was failed" + err.Error())
	}

	// build memo
	// types = 0 - strating manifest
	// types = 1 - managedata overflow sign
	memo, err := stellarprotocols.BuildMemo(0, 4, data.SequenceValue)
	if err != nil {
		logrus.Error("Memo ", err)
		w.WriteHeader(http.StatusInternalServerError)
		response := model.Error{Code: http.StatusInternalServerError, Message: "Memo hex converting issue  " + memo}
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(memo) != 28 {
		logrus.Error("Memo length error ", memo)
		w.WriteHeader(http.StatusInternalServerError)
		response := model.Error{Code: http.StatusInternalServerError, Message: "Memo length error  " + memo}
		json.NewEncoder(w).Encode(response)
		return
	}

	expertMapID := 0
	expertMapdata, err := object.GetExpertMapID(formulaJSON.Expert.ExpertID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		logrus.Info("Unable to connect gateway datastore ", err5)
		w.WriteHeader(http.StatusNotFound)
	}
	if expertMapdata == nil {
		data, err := object.GetNextSequenceValue("EXPERTID")
		if err != nil {
			logrus.Error("GetNextSequenceValu was failed" + err.Error())
		}
		expertIDMap := model.ExpertIDMap{
			ExpertID:  formulaJSON.Expert.ExpertID,
			ExpertPK:  formulaJSON.Expert.ExpertPK,
			MapID:     data.SequenceValue,
			FormulaID: formulaJSON.ID,
		}
		err1 := object.InsertExpertIDMap(expertIDMap)
		if err1 != nil {
			logrus.Error("Insert ExpertIDMap was failed" + err1.Error())
		}
		expertMapID = int(data.SequenceValue)
	} else {
		expertMap := expertMapdata.(model.ExpertIDMap)
		expertMapID = int(expertMap.MapID)
	}

	// formula identity operation
	formulaIdentityBuilder, errInFormulaIdentity := stellarprotocols.BuildFormulaIdentity(expertMapID, formulaJSON.Name, formulaJSON.Name)
	if errInFormulaIdentity != nil {
		logrus.Error("Building formula identity manage data failed : Error : " + errInFormulaIdentity.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "An error occured when building formula identity",
		}
		json.NewEncoder(w).Encode(result)
		return
	}

	// author details opreation
	authorDetailsBuilder, errInAuthorBuilder := stellarprotocols.BuildAuthorManageData(formulaJSON.Expert.ExpertPK)
	if errInAuthorBuilder != nil {
		logrus.Error("Building author details manage data failed : Error : " + errInAuthorBuilder.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		result := apiModel.SubmitXDRSuccess{
			Status: "An error occured when building author identity",
		}
		json.NewEncoder(w).Encode(result)
		return
	}

	// load account
	publicKey := constants.PublicKey
	secretKey := constants.SecretKey
	tracifiedAccount, err := keypair.ParseFull(secretKey)
	client := commons.GetHorizonClient()
	pubaccountRequest := horizonclient.AccountRequest{AccountID: publicKey}
	pubaccount, err := client.AccountDetail(pubaccountRequest)

	// build manage date with expert information
	// experInforBuilder := txnbuild.ManageData{
	// 	Name:  expertManageDataKey,
	// 	Value: []byte(expertManageDataValue),
	// }

	// if len(expertManageDataKey) > 64 || len(expertManageDataValue) > 64 {
	// 	logrus.Error("expert mange data length issue ", memo)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	response := model.Error{Code: http.StatusInternalServerError, Message: "Expert mange data length issue  " + " Value  " + expertManageDataValue + " key  " + expertManageDataKey}
	// 	json.NewEncoder(w).Encode(response)
	// 	return
	// }

	// check if any builder has failed
	// TODO: should add other manage data operations error handling here as well
	if errInFormulaIdentity == nil && errInAuthorBuilder == nil {
		// BUILD THE GATEWAY XDR
		tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
			SourceAccount:        &pubaccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&formulaIdentityBuilder, &authorDetailsBuilder},
			BaseFee:              txnbuild.MinBaseFee,
			Memo:                 txnbuild.MemoText(memo),
			Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		})
		if err != nil {
			logrus.Println("Error while buliding XDR " + err.Error())
		}
		// SIGN THE GATEWAY BUILT XDR WITH GATEWAYS PRIVATE KEY
		GatewayTXE, err := tx.Sign(commons.GetStellarNetwork(), tracifiedAccount)
		if err != nil {
			logrus.Error("Error while signing the XDR by secretKey  ", err)
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Error{Code: http.StatusInternalServerError, Message: "Error while signing the XDR by secretKey    " + err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		// CONVERT THE SIGNED XDR TO BASE64 to SUBMIT TO STELLAR
		resp, err := client.SubmitTransaction(GatewayTXE)
		if err != nil {
			logrus.Error("XDR submitting issue  ", err)
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Error{Code: http.StatusInternalServerError, Message: "XDR submitting issue  " + err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		logrus.Info("----------- Txn Hash -------- ", resp.Hash)

		formulaIDMap := model.FormulaIDMap{
			FormulaID: formulaJSON.ID,
			MapID:     data.SequenceValue,
		}

		err1 := object.InsertFormulaIDMap(formulaIDMap)
		if err1 != nil {
			logrus.Error("Insert FormulaIDMap was failed" + err1.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(formulaJSON)
		return
	}

	//! acessing the manage data in blockchain and convert byte to bit string
	// url := "https://horizon-testnet.stellar.org/transactions/" + resp.Hash + "/operations"
	// result, err := http.Get(url)
	// if err != nil {
	// 	logrus.Error("Unable to reach Stellar network", url)
	// }
	// if result.StatusCode != 200 {
	// 	logrus.Error(result)
	// }
	// defer result.Body.Close()
	// assertInfo, err := ioutil.ReadAll(result.Body)
	// if err != nil {
	// 	logrus.Error(err)
	// }

	// var raw map[string]interface{}
	// var raw1 map[string]interface{}

	// json.Unmarshal(assertInfo, &raw)

	// out, err := json.Marshal(raw["_embedded"])
	// if err != nil {
	// 	logrus.Error("Unable to marshal embedded")
	// }
	// err = json.Unmarshal(out, &raw1)
	// if err != nil {
	// 	logrus.Error("Unable to unmarshal  json.Unmarshal(out, &raw1)")
	// }
	// out1, _ := json.Marshal(raw1["records"])
	// json.Unmarshal(out1, &raw1)

	// keysBody := out1
	// keys := make([]PublicKeyPOCOC, 0)
	// err = json.Unmarshal(keysBody, &keys)
	// if err != nil {
	// 	logrus.Error("Unable to unmarshal keys data")
	// }
	// for i := range keys {
	// 	//fmt.Println(keys[i])
	// 	acceptTxn_byteData, err := base64.StdEncoding.DecodeString(keys[i].Value)
	// 	if err != nil {
	// 		logrus.Error("Unable to base64 decoding")
	// 	}
	// 	for _, nq := range acceptTxn_byteData {
	// 		// fmt.Printf("% 09b", nq) // prints 00000000 11111101
	// 		fmt.Println(nq)
	// 	}
	// }
}
