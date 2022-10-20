package expertformula

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
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
	formulaArray := formulaJSON.Formula             // formula array sent by the backend
	var manageDataOpArray []txnbuild.Operation      // manageDataOpArray all manage data append to to this array
	var transactionArray []model.FormulaTransaction // transaction array
	expertIDMap := model.ExpertIDMap{}
	var ValueDefinitionManageDataArray []model.ValueDefinition // value definition array
	var status string
	var startTransactionTime time.Time
	var endTransactionTime time.Time
	var expertMapID uint64

	object := dao.Connection{}
	// checked whether given formulaID already in the database or not
	formulaMap, err := object.GetExpertFormulaCount(formulaJSON.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		logrus.Info(err)
	}
	// if formulA already in Database, not allowed to  build expert formula to that ID
	if formulaMap.(int64) != 0 {
		commons.JSONErrorReturn(w, r, "Formula Id is in gateway datastore", http.StatusBadRequest, "")
		return
	}
	// if not,  retrived the current latest sequence number for formulaID
	dataFormulaID, err := object.GetNextSequenceValue("FORMULAID")
	if err != nil {
		commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "GetNextSequenceValu was failed")
		return
	}
	expertFormula := ExpertFormula{}
	// build memo
	memo, manifest, err := expertFormula.BuildMemo(0, uint32(fieldCount), dataFormulaID.SequenceValue)
	if err != nil {
		commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Memo hex converting issue")
		return
	}
	if len(memo) != 28 {
		commons.JSONErrorReturn(w, r, "Memo length error ", http.StatusInternalServerError, memo)
		return
	}
	// checked whether given ExpertID already in the database or not
	expertMapdata, err := object.GetExpertMapID(formulaJSON.User.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Unable to connect to gateway datastore ")
	}
	// if not,  retrived the current latest sequence number for expertID , map the expertID with incrementing interger
	if expertMapdata == nil {
		data, err := object.GetNextSequenceValue("EXPERTID")
		if err != nil {
			commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Mapping expert ID failed ")
			return
		}
		expertIDMap = model.ExpertIDMap{
			ExpertID:  formulaJSON.User.ID,
			ExpertPK:  formulaJSON.User.Publickey,
			MapID:     data.SequenceValue,
			FormulaID: formulaJSON.ID,
		}
		err1 := object.InsertExpertIDMap(expertIDMap)
		if err1 != nil {
			commons.JSONErrorReturn(w, r, err1.Error(), http.StatusInternalServerError, "Insert to ExpertIDMap was failed")
			return
		}
		expertMapID = data.SequenceValue
	} else {
		expertMap := expertMapdata.(model.ExpertIDMap)
		expertMapID = expertMap.MapID
	}
	// formula identity operation
	formulaIdentityBuilder, errInFormulaIdentity := expertFormula.BuildFormulaIdentity(expertMapID, formulaJSON.Name, formulaJSON.Name)
	if errInFormulaIdentity != nil {
		commons.JSONErrorReturn(w, r, errInFormulaIdentity.Error(), http.StatusInternalServerError, "An error occured when building formula identity ")
		return
	}
	// build formula object to be inserted
	formulaManageDataObj := model.FormulaIdentity{
		ManageDataOrder: 1,
		ManageDataName:  "FORMULA IDENTITY",
		FormulaMapID:    dataFormulaID.SequenceValue,
		ManageDataKey:   formulaIdentityBuilder.Name,
		ManageDataValue: formulaIdentityBuilder.Value,
	}
	// append to the manage data array
	manageDataOpArray = append(manageDataOpArray, &formulaIdentityBuilder)
	// author details opreation
	authorDetailsBuilder, errInAuthorBuilder := expertFormula.BuildPublisherManageData(formulaJSON.User.Publickey)
	if errInAuthorBuilder != nil {
		commons.JSONErrorReturn(w, r, errInAuthorBuilder.Error(), http.StatusInternalServerError, "An error occured when building author identity ")
		return
	}
	// build author identity manage data object(Expert is the author)
	// author-->Expert
	authorObj1 := model.AuthorIdentity{
		ManageDataOrder: 2,
		ManageDataName:  "AUTHOR IDENTITY",
		Expert:          expertIDMap,
		ManageDataKey:   authorDetailsBuilder.Name,
		ManageDataValue: authorDetailsBuilder.Value,
	}
	// append to the manage data array
	manageDataOpArray = append(manageDataOpArray, &authorDetailsBuilder)
	// loop through the formulaArray to see build the field definitions
	c := 2
	for i := 0; i < len(formulaArray); i++ {
		if formulaArray[i].Type == "VARIABLE" {
			// excute the variable builder
			variableBuilder, respObj, err := expertFormula.BuildVariableDefinitionManageData(formulaJSON.ID, formulaArray[i])
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "VARIABLE ")
				return
			}
			c++
			variableDefMGOObj := model.ValueDefinition{
				ManageDataOrder:   c,
				ValueType:         "VARIABLE",
				ValueMapID:        respObj.ValueMapID,
				UnitMapID:         respObj.UnitMapID,
				Precision:         formulaArray[i].Precision,
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
			referredConstant, respObj, err := expertFormula.BuildReferredConstantManageData(formulaJSON.ID, formulaArray[i])
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Referred Constant ")
				return
			}
			//url builder
			urlBuilder, err := expertFormula.BuildReference(formulaArray[i].MetricReference.Url)
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Referred URL ")
				return
			}
			c++
			// Reffered constant object
			referredConstObj := model.ValueDefinition{
				ManageDataOrder:   c,
				ValueType:         "REFERREDCONSTANT",
				ValueMapID:        respObj.ValueMapID,
				UnitMapID:         respObj.UnitMapID,
				Precision:         formulaArray[i].Precision,
				Value:             formulaArray[i].Value.(float64),
				MetricReferenceID: formulaArray[i].MetricReferenceId,
				ManageDataKey:     referredConstant.Name,
				ManageDataValue:   referredConstant.Value,
			}
			// append to the value definition array
			ValueDefinitionManageDataArray = append(ValueDefinitionManageDataArray, referredConstObj)
			// append to the manage data array
			manageDataOpArray = append(manageDataOpArray, &referredConstant)
			manageDataOpArray = append(manageDataOpArray, &urlBuilder)
		} else if formulaArray[i].Type == "SEMANTICCONSTANT" {
			// execute the semantic constant builder
			sematicConstant, respObj, err := expertFormula.BuildSemanticConstantManageData(formulaJSON.ID, formulaArray[i])
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "SEMANTI CCONSTANT ")
				return
			}
			//value builder
			valueBuilder, err := expertFormula.BuildSemanticValue(formulaArray[i].Value)
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "SEMANTI CCONSTANT ")
				return
			}
			c++
			// semantic constant object
			semanticConstObj := model.ValueDefinition{
				ValueType:         "SEMANTICCONSTANT",
				ValueMapID:        respObj.ValueMapID,
				UnitMapID:         respObj.UnitMapID,
				Precision:         formulaArray[i].Precision,
				Value:             formulaArray[i].Value.(float64),
				MetricReferenceID: formulaArray[i].MetricReferenceId,
				ManageDataKey:     sematicConstant.Name,
				ManageDataValue:   sematicConstant.Value,
			}
			// append to the value definition array
			ValueDefinitionManageDataArray = append(ValueDefinitionManageDataArray, semanticConstObj)
			// append to the manage data array
			manageDataOpArray = append(manageDataOpArray, &sematicConstant)
			manageDataOpArray = append(manageDataOpArray, &valueBuilder)
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
			commons.JSONErrorReturn(w, r, err.Error(), errCode, "Error when submitting transaction to blockchain ")
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
			logrus.Error("Insert formula to the formula map was failed" + err1.Error())
		}
		// build ManageData Array
		manageDataObj := model.FormulaManageData{
			FormulaIdentity:  formulaManageDataObj,
			AuthorIdentity:   []model.AuthorIdentity{authorObj1},
			ValueDefinitions: ValueDefinitionManageDataArray,
		}
		transactionCost := float64(int64(len(manageDataOpArray))) * 0.00001
		// build transaction
		// memo put to DB as a []byte to overcome invalid UTF-8 basonformate
		transactionObj := model.FormulaTransaction{
			TransactionHash:   string(hash),
			TransactionStatus: status,
			Memo:              []byte(memo),
			Manifest:          manifest,
			FormulaMapID:      dataFormulaID.SequenceValue,
			NoOfVariables:     fieldCount,
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
			ExpertID:               formulaJSON.User.ID,
			ExpertPK:               formulaJSON.User.Publickey,
			VariableCount:          len(formulaArray),
			FormulaJsonRequestBody: formulaJSON,
			Transactions:           transactionArray,
			OverflowAmount:         len(transactionArray),
			Status:                 status,
			CreatedAt:              time.Now().String(),
			CiperText:              formulaJSON.CiperText,
		}
		Id, errResult := object.InsertExpertFormula(expertFormulaBuilder)
		if errResult != nil {
			logrus.Error("Error while inserting the expert formula into DB: ", errResult)
		} else {
			w.WriteHeader(http.StatusOK)
			response := model.SuccessResponseExpertFormula{
				Code:              http.StatusOK,
				ID:                Id,
				FormulaID:         formulaJSON.ID,
				TransactionHashes: []string{hash},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}
