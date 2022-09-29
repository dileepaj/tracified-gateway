package expertformula

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

/*
StellarExpertFormulBuilder
des- This method build stellar trasactiond for expert formula

	 steps
		* map the formulaId and retrive the mapped id
		* build memo for the trasacions
		* map the experId and retive mapped id
		* build formula identity manageData opration
		* build auther identity manage data opration
		* loop through the formulaArray to see build the field definitions and build relevenat manage data oprations
		* load stellar account,build and sing the XDR
		* put XDR to stellar blockchain
*/
func StellarExpertFormulBuilder(w http.ResponseWriter, r *http.Request, formulaJSON model.FormulaBuildingRequest, fieldCount int) {
	w.Header().Set("Content-Type", "application/json")
	// formula array sent by the backend
	formulaArray := formulaJSON.Formula
	// manageDataOpArray all manage data append to to this array
	var manageDataOpArray []txnbuild.Operation
	// transaction array
	var transactionArray []model.FormulaTransaction
	var authorMapId int64
	// value definition array
	var ValueDefinitionManageDataArray []model.ValueDefinition
	var status string
	var startTransactionTime time.Time
	var endTransactionTime time.Time

	object := dao.Connection{}
	// checked whether given formulaID already in the database or not
	formulaMap, err5 := object.GetFormulaMapID(formulaJSON.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err5 != nil {
		logrus.Info(err5)
	}
	// if formulA already in Database, not allowed to  build expert formula to that ID
	if formulaMap != nil {
		logrus.Error("Formula Id is in gateway datastore")
		w.WriteHeader(http.StatusNoContent)
		response := model.Error{Code: http.StatusNoContent, Message: "Formula Id is in gateway datastore"}
		json.NewEncoder(w).Encode(response)
		return
	}
	// if not,  retrived the current latest sequence number for formulaID
	dataFormulaID, err := object.GetNextSequenceValue("FORMULAID")
	if err != nil {
		logrus.Error("GetNextSequenceValu was failed" + err.Error())
	}

	expertFormula := ExpertFormula{}
	// build memo

	// types = 0 - strating manifest
	// types = 1 - managedata overflow sign
	memo, err := expertFormula.BuildMemo(0, int32(fieldCount), dataFormulaID.SequenceValue)
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
	// checked whether given ExpertID already in the database or not
	expertMapdata, err := object.GetExpertMapID(formulaJSON.Expert.ExpertID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		logrus.Info("Unable to connect gateway datastore ", err5)
		w.WriteHeader(http.StatusNotFound)
	}
	// if not,  retrived the current latest sequence number for expertID , map the expertID with incrementing interger
	if expertMapdata == nil {
		data, err := object.GetNextSequenceValue("EXPERTID")
		if err != nil {
			logrus.Error("GetNextSequenceValu was failed" + err.Error())
		}
		authorMapId = data.SequenceValue
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
	formulaIdentityBuilder, errInFormulaIdentity := expertFormula.BuildFormulaIdentity(expertMapID, formulaJSON.Name, formulaJSON.Name)
	if errInFormulaIdentity != nil {
		logrus.Error("Building formula identity manage data failed : Error : " + errInFormulaIdentity.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("An error occured when building formula identity")
		return
	}
	// build formula object to be inserted
	formulaManageDataObj := model.FormulaIdentity{
		ManageDataName:  "FORMULA IDENTITY",
		FormulaMapID:    dataFormulaID.SequenceValue,
		ManageDataKey:   formulaIdentityBuilder.Name,
		ManageDataValue: formulaIdentityBuilder.Value,
	}
	// append to the manage data array
	manageDataOpArray = append(manageDataOpArray, &formulaIdentityBuilder)
	// author details opreation
	authorDetailsBuilder, errInAuthorBuilder := expertFormula.BuildAuthorManageData(formulaJSON.Expert.ExpertPK)
	if errInAuthorBuilder != nil {
		logrus.Error("Building author details manage data failed : Error : " + errInAuthorBuilder.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("An error occured when building author identity")
		return
	}

	// build author identity manage data object
	authorObj := model.AuthorIdentity{
		ManageDataName:  "AUTHOR IDENTITY",
		AuthorMapID:     authorMapId,
		ManageDataKey:   authorDetailsBuilder.Name,
		ManageDataValue: authorDetailsBuilder.Value,
	}
	// append to the manage data array
	manageDataOpArray = append(manageDataOpArray, &authorDetailsBuilder)
	// loop through the formulaArray to see build the field definitions
	for i := 0; i < len(formulaArray); i++ {
		if formulaArray[i].Type == "VARIABLE" {
			// excute the variable builder
			variableBuilder, respObj, err := expertFormula.BuildVariableDefinitionManageData(formulaArray[i])
			if err != nil {
				logrus.Error("Variable  ", err.Error())
				w.WriteHeader(http.StatusNoContent)
				response := model.Error{Code: http.StatusNoContent, Message: err.Error()}
				json.NewEncoder(w).Encode(response)
				return
			}
			variableDefMGOObj := model.ValueDefinition{
				ValueType:         "VARIABLE",
				ValueMapID:        respObj.ValueMapID,
				UnitMapID:         respObj.UnitMapID,
				Precision:         formulaArray[i].Precision,
				Value:             formulaArray[i].Value,
				MetricReferenceID: formulaArray[i].MetricReferenceId,
				ManageDataKey:     variableBuilder.Name,
				ManageDataValue:   variableBuilder.Value,
			}
			// append to value definition array to be inserted in to the DB
			ValueDefinitionManageDataArray = append(ValueDefinitionManageDataArray, variableDefMGOObj)

			// append to the manage data array
			manageDataOpArray = append(manageDataOpArray, &variableBuilder)
		} else if formulaArray[i].Type == "REFERREDCONSTANT" {
			// execute the referred constant builder

			referredConstant, respObj, err := expertFormula.BuildReferredConstantManageData(formulaArray[i])
			if err != nil {
				logrus.Error("referred Constant   ", err.Error())
				w.WriteHeader(http.StatusNoContent)
				response := model.Error{Code: http.StatusNoContent, Message: err.Error()}
				json.NewEncoder(w).Encode(response)
				return
			}
			// Reffered constant object
			referredConstObj := model.ValueDefinition{
				ValueType:         "REFERREDCONSTANT",
				ValueMapID:        respObj.ValueMapID,
				UnitMapID:         respObj.UnitMapID,
				Precision:         formulaArray[i].Precision,
				Value:             formulaArray[i].Value,
				MetricReferenceID: formulaArray[i].MetricReferenceId,
				ManageDataKey:     referredConstant.Name,
				ManageDataValue:   referredConstant.Value,
			}
			// append to the value definition array
			ValueDefinitionManageDataArray = append(ValueDefinitionManageDataArray, referredConstObj)

			// append to the manage data array
			manageDataOpArray = append(manageDataOpArray, &referredConstant)
		} else if formulaArray[i].Type == "SEMANTICCONSTANT" {
			// execute the semantic constant builder
			sematicConstant, respObj, err := expertFormula.BuildSemanticConstantManageData(formulaArray[i])
			if err != nil {
				logrus.Error("sementic Constant   ", err.Error())
				w.WriteHeader(http.StatusNoContent)
				response := model.Error{Code: http.StatusNoContent, Message: err.Error()}
				json.NewEncoder(w).Encode(response)
				return
			}
			// semantic constant object
			semanticConstObj := model.ValueDefinition{
				ValueType:         "SEMANTICCONSTANT",
				ValueMapID:        respObj.ValueMapID,
				UnitMapID:         respObj.UnitMapID,
				Precision:         formulaArray[i].Precision,
				Value:             formulaArray[i].Value,
				MetricReferenceID: formulaArray[i].MetricReferenceId,
				ManageDataKey:     sematicConstant.Name,
				ManageDataValue:   sematicConstant.Value,
			}
			// append to the value definition array
			ValueDefinitionManageDataArray = append(ValueDefinitionManageDataArray, semanticConstObj)

			// append to the manage data array
			manageDataOpArray = append(manageDataOpArray, &sematicConstant)
		}
	}
	if errInFormulaIdentity == nil && errInAuthorBuilder == nil {
		stellarProtocol := stellarprotocols.StellarTrasaction{
			PublicKey:  constants.PublicKey,
			SecretKey:  constants.SecretKey,
			Operations: manageDataOpArray,
			Memo:       memo,
		}
		startTransactionTime = time.Now()
		err, errCode, hash := stellarProtocol.SubmitToStellerBlockchain()
		endTransactionTime = time.Now()
		if err != nil {
			status = "Failed"
			logrus.Error("Error when transaction Submitting to blockchain  ", err)
			w.WriteHeader(errCode)
			response := model.Error{Code: errCode, Message: "Error when transaction Submitting to blockchain  " + err.Error()}
			json.NewEncoder(w).Encode(response)
			return
		}
		status = "Success"
		logrus.Info("Transaction Hash for the formula building : ", hash)

		timeForTransaction := endTransactionTime.Sub(startTransactionTime)

		formulaIDMap := model.FormulaIDMap{
			FormulaID: formulaJSON.ID,
			MapID:     dataFormulaID.SequenceValue,
		}
		// map the formulaID with incremting Integer put those object to blockchain
		err1 := object.InsertFormulaIDMap(formulaIDMap)
		if err1 != nil {
			logrus.Error("Insert FormulaIDMap was failed" + err1.Error())
		}

		// build ManageData Array
		manageDataObj := model.FormulaManageData{
			FormulaIdentity:  formulaManageDataObj,
			AuthorIdentity:   authorObj,
			ValueDefinitions: ValueDefinitionManageDataArray,
		}

		transactionCost := float64(int64(len(manageDataOpArray))) * 0.00001
		// build transaction
		// memo put to DB as a []byte to overcome invalid UTF-8 basonformate
		transactionObj := model.FormulaTransaction{
			TransactionHash:   string(hash),
			TransactionStatus: status,
			Memo:              []byte(memo),
			ManageData:        manageDataObj,
			TransactionTime:   timeForTransaction.String(),
			Cost:              fmt.Sprintf("%f", transactionCost),
		}

		// append this transaction to the transaction array
		transactionArray = append(transactionArray, transactionObj)

		// save expert formula in the database
		// Todo: Transactions, overflowAmount, status should be changed to actual values
		expertFormulaBuilder := model.FormulaStore{
			Blockchain:             formulaJSON.Blockchain,
			FormulaID:              formulaJSON.ID,
			ExpertPK:               formulaJSON.Expert.ExpertPK,
			VariableCount:          len(formulaArray),
			FormulaJsonRequestBody: formulaJSON,
			Transactions:           transactionArray,
			OverflowAmount:         len(transactionArray),
			Status:                 status,
			CreatedAt:              time.Now().String(),
		}

		object := dao.Connection{}
		errResult := object.InsertExpertFormula(expertFormulaBuilder)
		if errResult != nil {
			logrus.Error("Error while inserting the expert formula into DB : ", errResult)
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(formulaJSON)
			return
		}
	}
}
