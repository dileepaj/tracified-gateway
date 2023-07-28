package polygonexpertformula

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	codeGenerator "github.com/dileepaj/tracified-gateway/protocols/ethereum/codeGenerator/ethereumExpertFormula"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/codeGenerator/ethereumExpertFormula/executionTemplates"
	deletecontract "github.com/dileepaj/tracified-gateway/protocols/polygon/polygonCodeGenerator/polygonExpertFormula/deleteContract"
	experthelpers "github.com/dileepaj/tracified-gateway/protocols/polygon/polygonCodeGenerator/polygonExpertFormula/expertHelpers"
	expertformula "github.com/dileepaj/tracified-gateway/protocols/stellarprotocols/expertFormula"
	"github.com/dileepaj/tracified-gateway/utilities"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	contractName = ``
	contractBody = ``
)

func PolygonExpertFormulaContractGenerator(w http.ResponseWriter, r *http.Request, formulaJSON model.FormulaBuildingRequest, fieldCount int) {
	object := dao.Connection{}
	var deployStatus int
	reqType := "POLYGONEXPERT"
	logger := utilities.NewCustomLogger()
	formulaDetails, errWhenGettingFormulaDetails := object.GetPolygonFormulaStatus(formulaJSON.MetricExpertFormula.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errWhenGettingFormulaDetails != nil {
		logger.LogWriter("An error occurred when getting formula status, ERROR : "+errWhenGettingFormulaDetails.Error(), constants.ERROR)
	}
	if formulaDetails == nil {
		deployStatus = 0
	}
	if formulaDetails != nil {
		deployStatus = formulaDetails.(model.EthereumExpertFormula).Status
		logger.LogWriter("Polygon formula contract deploy status : "+strconv.FormatInt(int64(deployStatus), 10), constants.INFO)
	}
	if deployStatus != 0 || deployStatus != 119 {
		//handle Queue, Success, invalid status
		experthelpers.SuccessOrQueueResponse(w, r, formulaJSON, deployStatus)
	} else {
		if deployStatus == 119 {
			logger.LogWriter("Requested formula is in the failed status, trying to redeploy", constants.INFO)
		} else {
			logger.LogWriter("New expert formula request, initiating new deployment", constants.INFO)
		}
		//create expert formula
		formulaObj := experthelpers.BuildExpertObject(formulaJSON.MetricExpertFormula.ID, formulaJSON.MetricExpertFormula.Name, formulaJSON.MetricExpertFormula, fieldCount, formulaJSON.Verify)
		if deployStatus == 0 {
			transactionUuid := experthelpers.GenerateTransactionUUID()
			formulaObj.TransactionUUID = transactionUuid
			errWhenInsertingToFormulaIdMap := experthelpers.InsertToFormulaIdMap(formulaJSON.MetricExpertFormula.ID, 2)
			if errWhenInsertingToFormulaIdMap != nil {
				logger.LogWriter("Error when inserting the new formula ID to the polygon formula ID map : "+errWhenInsertingToFormulaIdMap.Error(), constants.ERROR)
				commons.JSONErrorReturn(w, r, errWhenInsertingToFormulaIdMap.Error(), http.StatusInternalServerError, "Error when inserting the new formula ID to the polygon formula ID map : "+errWhenInsertingToFormulaIdMap.Error())
				return
			}
		} else {
			formulaObj.TransactionUUID = formulaDetails.(model.EthereumExpertFormula).TransactionUUID
		}
		//setting up the contract name and starting the contract
		contractName = cases.Title(language.English).String(formulaJSON.MetricExpertFormula.Name)
		contractName = strings.ReplaceAll(contractName, " ", "")
		contractName = contractName + "_" + formulaJSON.MetricExpertFormula.ID
		//call the general header writer
		generalValues, errWhenBuildingGeneralCodeSnippets := codeGenerator.WriteGeneralCodeSnippets(formulaJSON, contractName)
		if errWhenBuildingGeneralCodeSnippets != nil {
			errWhenUpdatingOrInsertingFormulaDetails := experthelpers.InsertAndUpdateExpertFormulaDetailsToPolygonCollections(deployStatus, 119, errWhenBuildingGeneralCodeSnippets.Error(), 102, formulaObj, formulaObj.FormulaID, formulaObj.TransactionUUID)
			if errWhenUpdatingOrInsertingFormulaDetails != nil {
				logger.LogWriter("Error when updating/inserting to polygon collections : "+errWhenUpdatingOrInsertingFormulaDetails.Error(), constants.INFO)
				commons.JSONErrorReturn(w, r, errWhenUpdatingOrInsertingFormulaDetails.Error(), http.StatusInternalServerError, "Error when updating/inserting to polygon collections")
				return
			}
			logger.LogWriter("Error when writing the general code snippet, ERROR : "+errWhenBuildingGeneralCodeSnippets.Error(), constants.ERROR)
			commons.JSONErrorReturn(w, r, errWhenBuildingGeneralCodeSnippets.Error(), http.StatusInternalServerError, "Error when writing the general code snippet, ERROR : ")
			return
		}
		contractBody = generalValues.ResultVariable + generalValues.MetaDataStructure + generalValues.ValueDataStructure + generalValues.VariableStructure + generalValues.SemanticConstantStructure + generalValues.ReferredConstant + generalValues.MetadataDeclaration
		contractBody = contractBody + generalValues.ResultDeclaration + generalValues.CalculationObject
		variableValue, setterName, errInGeneratingValues := codeGenerator.ValueCodeGenerator(formulaJSON)
		if errInGeneratingValues != nil {
			errWhenUpdatingOrInsertingFormulaDetails := experthelpers.InsertAndUpdateExpertFormulaDetailsToPolygonCollections(deployStatus, 119, errInGeneratingValues.Error(), 102, formulaObj, formulaObj.FormulaID, formulaObj.TransactionUUID)
			if errWhenUpdatingOrInsertingFormulaDetails != nil {
				logger.LogWriter("Error when updating/inserting to polygon collections : "+errWhenUpdatingOrInsertingFormulaDetails.Error(), constants.INFO)
				commons.JSONErrorReturn(w, r, errWhenUpdatingOrInsertingFormulaDetails.Error(), http.StatusInternalServerError, "Error when updating/inserting to polygon collections")
				return
			}
			logger.LogWriter("Error in generating codes for values "+errInGeneratingValues.Error(), constants.ERROR)
			commons.JSONErrorReturn(w, r, errInGeneratingValues.Error(), http.StatusInternalServerError, "Error in getting codes for values ")
			return
		}
		contractBody = contractBody + variableValue
		formulaObj.SetterNames = setterName
		executionTemplate, errInGettingExecutionTemplate := expertformula.BuildExecutionTemplateByQuery(formulaObj.MetricExpertFormula.FormulaAsQuery)
		if errInGettingExecutionTemplate != nil {
			errWhenUpdatingOrInsertingFormulaDetails := experthelpers.InsertAndUpdateExpertFormulaDetailsToPolygonCollections(deployStatus, 119, errInGettingExecutionTemplate.Error(), 102, formulaObj, formulaObj.FormulaID, formulaObj.TransactionUUID)
			if errWhenUpdatingOrInsertingFormulaDetails != nil {
				logger.LogWriter("Error when updating/inserting to polygon collections : "+errWhenUpdatingOrInsertingFormulaDetails.Error(), constants.INFO)
				commons.JSONErrorReturn(w, r, errWhenUpdatingOrInsertingFormulaDetails.Error(), http.StatusInternalServerError, "Error when updating/inserting to polygon collections")
				return
			}
			logger.LogWriter("Error in getting execution template "+errInGettingExecutionTemplate.Error(), constants.ERROR)
			commons.JSONErrorReturn(w, r, errInGettingExecutionTemplate.Error(), http.StatusInternalServerError, "Error in getting execution template from FCL ")
			return
		}
		formulaObj.ExecutionTemplate = executionTemplate
		executionTemplateString, errInGettingExecutionTemplateString := executionTemplates.ExecutionTemplateDivider(executionTemplate)
		if errInGettingExecutionTemplateString != nil {
			errWhenUpdatingOrInsertingFormulaDetails := experthelpers.InsertAndUpdateExpertFormulaDetailsToPolygonCollections(deployStatus, 119, errInGettingExecutionTemplateString.Error(), 102, formulaObj, formulaObj.FormulaID, formulaObj.TransactionUUID)
			if errWhenUpdatingOrInsertingFormulaDetails != nil {
				logger.LogWriter("Error when updating/inserting to polygon collections : "+errWhenUpdatingOrInsertingFormulaDetails.Error(), constants.INFO)
				commons.JSONErrorReturn(w, r, errWhenUpdatingOrInsertingFormulaDetails.Error(), http.StatusInternalServerError, "Error when updating/inserting to polygon collections")
				return
			}
			logger.LogWriter("Error in getting execution template string "+errInGettingExecutionTemplateString.Error(), constants.ERROR)
			commons.JSONErrorReturn(w, r, errInGettingExecutionTemplateString.Error(), http.StatusInternalServerError, "Error in getting execution template string ")
			return
		}
		// remove the substring from the last comma
		lenOfLastCommand := len(", calculations.GetExponent()")
		executionTemplateString = executionTemplateString[:len(executionTemplateString)-lenOfLastCommand]
		contractBody = contractBody + experthelpers.WriteCalculationGetterCode(executionTemplateString) + experthelpers.WriteGetterMethods(generalValues.MetadataGetter)
		template, errWhenGeneratingTemplate := experthelpers.ContractTemplateBuilder(formulaObj, generalValues.License, generalValues.PragmaLine, generalValues.ImportCalculationsSol, generalValues.ContractStart, contractBody, generalValues.ContractEnd)
		if errWhenGeneratingTemplate != nil {
			errWhenUpdatingOrInsertingFormulaDetails := experthelpers.InsertAndUpdateExpertFormulaDetailsToPolygonCollections(deployStatus, 119, errWhenGeneratingTemplate.Error(), 102, formulaObj, formulaObj.FormulaID, formulaObj.TransactionUUID)
			if errWhenUpdatingOrInsertingFormulaDetails != nil {
				logger.LogWriter("Error when updating/inserting to polygon collections : "+errWhenUpdatingOrInsertingFormulaDetails.Error(), constants.INFO)
				commons.JSONErrorReturn(w, r, errWhenUpdatingOrInsertingFormulaDetails.Error(), http.StatusInternalServerError, "Error when updating/inserting to polygon collections")
				return
			}
			logger.LogWriter("Error when generating the contract template string : "+errWhenGeneratingTemplate.Error(), constants.ERROR)
			commons.JSONErrorReturn(w, r, errWhenGeneratingTemplate.Error(), http.StatusInternalServerError, "Error when generating the contract template string")
			return
		}
		errWhenWritingSolidityFile := experthelpers.WriteFormulaContractToFile(contractName, template, formulaObj.FormulaID, formulaObj.TransactionUUID, formulaObj)
		if errWhenWritingSolidityFile != nil {
			errWhenUpdatingOrInsertingFormulaDetails := experthelpers.InsertAndUpdateExpertFormulaDetailsToPolygonCollections(deployStatus, 119, errWhenWritingSolidityFile.Error(), 104, formulaObj, formulaObj.FormulaID, formulaObj.TransactionUUID)
			if errWhenUpdatingOrInsertingFormulaDetails != nil {
				logger.LogWriter("Error when updating/inserting to polygon collections : "+errWhenUpdatingOrInsertingFormulaDetails.Error(), constants.INFO)
				commons.JSONErrorReturn(w, r, errWhenUpdatingOrInsertingFormulaDetails.Error(), http.StatusInternalServerError, "Error when updating/inserting to polygon collections")
				return
			}
			logger.LogWriter("Error when generating the contract template string : "+errWhenWritingSolidityFile.Error(), constants.ERROR)
			commons.JSONErrorReturn(w, r, errWhenWritingSolidityFile.Error(), http.StatusInternalServerError, "Error when generating the contract template string")
			return
		}
		//ABI and Bin generator
		errWhenGeneratingAbiAndBin := experthelpers.AbiAndBinGenerator(contractName, reqType, formulaObj.FormulaID, formulaObj.TransactionUUID, formulaObj)
		if errWhenGeneratingAbiAndBin != nil {
			//108 - ABI AND BIN generation failed
			errWhenUpdatingOrInsertingFormulaDetails := experthelpers.InsertAndUpdateExpertFormulaDetailsToPolygonCollections(deployStatus, 119, errWhenGeneratingAbiAndBin.Error(), 108, formulaObj, formulaObj.FormulaID, formulaObj.TransactionUUID)
			if errWhenUpdatingOrInsertingFormulaDetails != nil {
				logger.LogWriter("Error when updating/inserting to polygon collections : "+errWhenUpdatingOrInsertingFormulaDetails.Error(), constants.INFO)
				commons.JSONErrorReturn(w, r, errWhenUpdatingOrInsertingFormulaDetails.Error(), http.StatusInternalServerError, "Error when updating/inserting to polygon collections")
				return
			}
			logger.LogWriter("Error when generating ABI and BIN for contract : "+errWhenGeneratingAbiAndBin.Error(), constants.ERROR)
			commons.JSONErrorReturn(w, r, errWhenGeneratingAbiAndBin.Error(), http.StatusInternalServerError, "Error when generating ABI and BIN for contract : ")
			return
		}
		formulaObj.ContractName = contractName
		errWhenSendingToQueue := experthelpers.BuildExpertQueueObjectAndSendToQueue(formulaObj, "POLYGONEXPERTFORMULA", "QUEUE")
		if errWhenSendingToQueue != nil {
			errWhenUpdatingOrInsertingFormulaDetails := experthelpers.InsertAndUpdateExpertFormulaDetailsToPolygonCollections(deployStatus, 119, errWhenSendingToQueue.Error(), formulaObj.ActualStatus, formulaObj, formulaObj.FormulaID, formulaObj.TransactionUUID)
			if errWhenUpdatingOrInsertingFormulaDetails != nil {
				logger.LogWriter("Error when updating/inserting to polygon collections : "+errWhenUpdatingOrInsertingFormulaDetails.Error(), constants.INFO)
				commons.JSONErrorReturn(w, r, errWhenUpdatingOrInsertingFormulaDetails.Error(), http.StatusInternalServerError, "Error when updating/inserting to polygon collections")
				return
			}
			logger.LogWriter("Error when sending request to Queue : "+errWhenSendingToQueue.Error(), constants.ERROR)
			commons.JSONErrorReturn(w, r, errWhenGeneratingAbiAndBin.Error(), http.StatusInternalServerError, "Error when sending to Queue : ")
			return
		}
		logger.LogWriter("Expert formula is added to Queue", constants.INFO)
		errWhenUpdatingOrInsertingAfterQueue := experthelpers.InsertAndUpdateExpertFormulaDetailsToPolygonCollections(deployStatus, 116, "", formulaObj.ActualStatus, formulaObj, formulaObj.FormulaID, formulaObj.TransactionUUID)
		if errWhenUpdatingOrInsertingAfterQueue != nil {
			logger.LogWriter("Error when updating/inserting to polygon collections : "+errWhenUpdatingOrInsertingAfterQueue.Error(), constants.INFO)
			commons.JSONErrorReturn(w, r, errWhenUpdatingOrInsertingAfterQueue.Error(), http.StatusInternalServerError, "Error when updating/inserting to polygon collections")
			return
		}

		//contract, ABI and BIN file deletion
		errWhenDeleting := deletecontract.DeleteExpertContract(contractName)
		if errWhenDeleting != nil {
			logger.LogWriter("Error when deleting the source files : "+errWhenDeleting.Error(), constants.ERROR)
		}

		w.WriteHeader(http.StatusOK)
		response := model.SuccessResponseExpertFormula{
			Code:      http.StatusOK,
			FormulaID: formulaJSON.MetricExpertFormula.ID,
			Message:   "Expert formula request sent to queue",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

}
