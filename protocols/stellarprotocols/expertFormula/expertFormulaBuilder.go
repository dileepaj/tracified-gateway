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
	equationbuilding "github.com/dileepaj/tracified-gateway/protocols/stellarprotocols/expertFormula/equationBuilding"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

/*
StellarExpertFormulBuilder
des- This method build stellar trasactions for expert formula

	 steps
		* map the formulaId and retrive the mapped id
		* build memo for the trasacions
		* map the experId and retive mapped id
		* build formula identity manageData opration
		* build auther identity manage data opration
		* loop through the formulaArray to see build the field definitions and build relevenat manage data oprations
		* get the execution template from fcld and build relevenat manage data oprations
		* load stellar account,build and sing the XDR
		* put XDR to stellar blockchain
*/
func StellarExpertFormulBuilder(w http.ResponseWriter, r *http.Request, formulaJSON model.FormulaBuildingRequest, fieldCount int) {
	w.Header().Set("Content-Type", "application/json")

	formulaArray := formulaJSON.MetricExpertFormula.Formula // formula array sent by the backend
	var hashArray []string                                  // formula array sent by the backend
	var manageDataOpArray []txnbuild.Operation              // manageDataOpArray all manage data append to to this array
	var transactionArray []model.FormulaTransaction         // transaction array
	expertIDMap := model.ExpertIDMap{}
	var ValueDefinitionManageDataArray []model.ValueDefinition // value definition array
	var status string
	var startTransactionTime time.Time
	var endTransactionTime time.Time
	var expertMapID uint64
	var memo0, memo1, manifest string
	object := dao.Connection{}
	// checked whether given formulaID already in the database or not
	formulaMap, err := object.GetExpertFormulaCount(formulaJSON.MetricExpertFormula.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		logrus.Info(err)
	}
	// if formulA already in Database, not allowed to  build expert formula to that ID
	if formulaMap.(int64) != 0 {
		commons.JSONErrorReturn(w, r, "Formula Id is in gateway datastore", http.StatusBadRequest, "Duplicate formula IDs not allowed(expertFormulaBuilder) ")
		return
	}
	// if not,  retrived the current latest sequence number for formulaID
	dataFormulaID, err := object.GetNextSequenceValue("FORMULAID")
	if err != nil {
		commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "GetNextSequenceValue for formula Id was failed(expertFormulaBuilder) ")
		return
	}
	expertFormula := ExpertFormula{}

	// checked whether given ExpertID already in the database or not
	expertMapdata, err := object.GetExpertMapID(formulaJSON.User.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Unable to connect to gateway datastore(expertFormulaBuilder) ")
	}
	// if not,  retrived the current latest sequence number for expertID , map the expertID with incrementing interger
	if expertMapdata == nil {
		data, err := object.GetNextSequenceValue("EXPERTID")
		if err != nil {
			commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Mapping expert ID failed(expertFormulaBuilder) ")
			return
		}
		expertIDMap = model.ExpertIDMap{
			ExpertID:  formulaJSON.User.ID,
			ExpertPK:  formulaJSON.User.Publickey,
			MapID:     data.SequenceValue,
			FormulaID: formulaJSON.MetricExpertFormula.ID,
		}
		err1 := object.InsertExpertIDMap(expertIDMap)
		if err1 != nil {
			commons.JSONErrorReturn(w, r, err1.Error(), http.StatusInternalServerError, "Insert to ExpertIDMap was failed(expertFormulaBuilder) ")
			return
		}
		expertMapID = data.SequenceValue
	} else {
		expertMap := expertMapdata.(model.ExpertIDMap)
		expertMapID = expertMap.MapID
	}
	// formula identity operation
	formulaIdentityBuilder, errInFormulaIdentity := expertFormula.BuildFormulaIdentity(expertMapID, formulaJSON.MetricExpertFormula.Name)
	if errInFormulaIdentity != nil {
		commons.JSONErrorReturn(w, r, errInFormulaIdentity.Error(), http.StatusInternalServerError, "An error occured when building formula identity(expertFormulaBuilder) ")
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
	authorDetailsBuilder, errInAuthorBuilder := expertFormula.BuildPublicManageData(formulaJSON.User.Publickey)
	if errInAuthorBuilder != nil {
		commons.JSONErrorReturn(w, r, errInAuthorBuilder.Error(), http.StatusInternalServerError, "An error occured when building author identity(expertFormulaBuilder) ")
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
			variableBuilder, respObj, err := expertFormula.BuildVariableDefinitionManageData(formulaJSON.MetricExpertFormula.ID, formulaArray[i])
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Building Variable was failed(expertFormulaBuilder) ")
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
			referredConstant, respObj, err := expertFormula.BuildReferredConstantManageData(formulaJSON.MetricExpertFormula.ID, formulaArray[i])
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Building Referred Constant was failed(expertFormulaBuilder) ")
				return
			}
			// url builder
			urlBuilder, err := expertFormula.BuildReference(formulaArray[i].MetricReference.Reference)
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Building Referred Constant Reference was failed(expertFormulaBuilder) ")
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
			sematicConstant, respObj, err := expertFormula.BuildSemanticConstantManageData(formulaJSON.MetricExpertFormula.ID, formulaArray[i])
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Building Semantic Constant was failed(expertFormulaBuilder) ")
				return
			}
			// value builder
			valueBuilder, err := expertFormula.BuildSemanticValue(formulaArray[i].Value)
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Building Semantic Constant Value was failed(expertFormulaBuilder) ")
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

	/* logic section of the expert formula builder 

		* BuildExecutionTemplateByQuery() method will return the execution template that returns from the FCL 
		* if the lst_commands in the returned execution template is not empty 
			-> Type 1 execution template(Start variable followed by a list of commands) - returns an array of manage data operations
		  else 
			-> Type 2 execution template(Entity) - returns a single manage data operation
	*/
	executionTemplate, errInGettingExecutionTemplate := BuildExecutionTemplateByQuery(formulaJSON.MetricExpertFormula.FormulaAsQuery)
	if errInGettingExecutionTemplate != nil {
		commons.JSONErrorReturn(w, r, errInGettingExecutionTemplate.Error(), http.StatusInternalServerError, "Error in getting execution template from FCL(expertFormulaBuilder) ")
		return
	}
	if executionTemplate.Lst_Commands != nil {
		manageDataOp, errTemplate1Builder := equationbuilding.Type1TemplateBuilder(formulaJSON.MetricExpertFormula.ID, executionTemplate)
		if errTemplate1Builder != nil {
			commons.JSONErrorReturn(w, r, errTemplate1Builder.Error(), http.StatusInternalServerError, "Error in building execution template type 1 failed(expertFormulaBuilder) ")
			return
		}
		// append to the manage data array
		manageDataOpArray = append(manageDataOpArray, manageDataOp...)
	} else {
		template1Builder, errInTemplate1Builder := equationbuilding.Type2TemplateBuilder(formulaJSON.MetricExpertFormula.ID, executionTemplate)
		if errInTemplate1Builder != nil {
			commons.JSONErrorReturn(w, r, errInTemplate1Builder.Error(), http.StatusInternalServerError, "Error in building execution template type 2 failed(expertFormulaBuilder) ")
			return
		}
		// append to the manage data array
		manageDataOpArray = append(manageDataOpArray, template1Builder)
	}

	// split the manage data array into two parts
	manageData2dArray := commons.ChunkSlice(manageDataOpArray, 25)
	for i, manadataOperationArray := range manageData2dArray {
		if i == 0 {
			// build memo0 send the transaction
			memo0, manifest, err = expertFormula.BuildMemo(0, uint32(fieldCount), dataFormulaID.SequenceValue)
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Hex conversion issue in building memo(expertFormulaBuilder) ")
				return
			}
			if len(memo0) != 28 {
				commons.JSONErrorReturn(w, r, "Memo length error(expertFormulaBuilder) ", http.StatusInternalServerError, memo0)
				return
			}

		}
		stellarProtocol := stellarprotocols.StellarTrasaction{
			PublicKey:  constants.PublicKey,
			Operations: manadataOperationArray,
			Memo:       memo0,
		}
		transaction := model.TransactionHash{
			Order: i,
			Memo:  []byte(memo0),
		}
		if i != 0 {
			// here for insted of no of values we pass the current index of the manadataOperationArray array
			memo1, manifest, err = expertFormula.BuildMemo(1, uint32(i), dataFormulaID.SequenceValue)
			if err != nil {
				commons.JSONErrorReturn(w, r, err.Error(), http.StatusInternalServerError, "Hex conversion issue in building memo(expertFormulaBuilder)")
				return
			}
			if len(memo1) != 28 {
				commons.JSONErrorReturn(w, r, "Memo length error(expertFormulaBuilder) ", http.StatusInternalServerError, memo1)
				return
			}
			stellarProtocol.Memo = memo1
			transaction.Memo = []byte(memo1)

		}
		startTransactionTime = time.Now()
		err, errCode, hash, _ := stellarProtocol.SubmitToStellerBlockchain()
		hashArray = append(hashArray, hash)
		endTransactionTime = time.Now()
		if err != nil {
			status = "Failed"
			commons.JSONErrorReturn(w, r, err.Error(), errCode, "Error when submitting transaction to blockchain(expertFormulaBuilder) ")
			return
		}
		status = "Success"
		logrus.Info("Transaction Hash of ", i+1, " transaction for the formula building : ", hash)
		timeForTransaction := endTransactionTime.Sub(startTransactionTime)
		formulaIDMap := model.FormulaIDMap{
			FormulaID: formulaJSON.MetricExpertFormula.ID,
			MapID:     dataFormulaID.SequenceValue,
		}
		// map the formulaID with incrementing Integer put those object to blockchain
		err1 := object.InsertFormulaIDMap(formulaIDMap)
		if err1 != nil {
			logrus.Error("Inserting formula to the formula map was failed(expertFormulaBuilder) " + err1.Error())
		}
		// build ManageData Array
		manageDataObj := model.FormulaManageData{
			FormulaIdentity:  formulaManageDataObj,
			AuthorIdentity:   []model.AuthorIdentity{authorObj1},
			ValueDefinitions: ValueDefinitionManageDataArray,
		}

		transactionCost := float64(int64(len(manadataOperationArray))) * 0.00001
		// build transaction
		// memo put to DB as a []byte to overcome invalid UTF-8 basonformate
		transactionObj := model.FormulaTransaction{
			TransactionHash:   string(hash),
			TransactionStatus: status,
			Manifest:          manifest,
			Memo:              []byte(transaction.Memo),
			FormulaMapID:      dataFormulaID.SequenceValue,
			NoOfVariables:     fieldCount,
			ManageData:        manageDataObj,
			TransactionTime:   timeForTransaction.String(),
			Cost:              fmt.Sprintf("%f", transactionCost),
		}
		// append this transaction to the transaction array
		transactionArray = append(transactionArray, transactionObj)
	}
	// save expert formula in the database
	expertFormulaBuilder := model.FormulaStore{
		Blockchain:              formulaJSON.MetricExpertFormula.Blockchain,
		FormulaID:               formulaJSON.MetricExpertFormula.ID,
		ExpertID:                formulaJSON.User.ID,
		ExpertPK:                formulaJSON.User.Publickey,
		VariableCount:           len(formulaArray),
		FormulaJsonRequestBody:  formulaJSON,
		Transactions:            transactionArray,
		OverflowAmount:          len(transactionArray),
		Status:                  status,
		CreatedAt:               time.Now().String(),
		CiperText:               formulaJSON.MetricExpertFormula.CiperText,
		ExecutionTemplate:       executionTemplate,
		TotalNumberOfManageData: len(manageDataOpArray),
	}
	Id, errResult := object.InsertExpertFormula(expertFormulaBuilder)
	if errResult != nil {
		logrus.Error("Error while inserting the expert formula into DB(expertFormulaBuilder): ", errResult)
	}
	w.WriteHeader(http.StatusOK)
	response := model.SuccessResponseExpertFormula{
		Code:              http.StatusOK,
		ID:                Id,
		FormulaID:         formulaJSON.MetricExpertFormula.ID,
		TransactionHashes: hashArray,
	}
	json.NewEncoder(w).Encode(response)
	return
}
