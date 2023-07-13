package polygonexpertformula

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	codeGenerator "github.com/dileepaj/tracified-gateway/protocols/ethereum/codeGenerator/ethereumExpertFormula"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/codeGenerator/ethereumExpertFormula/executionTemplates"
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
	// reqType := "POLYGONEXPERT"
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
			//generate transaction UUID
			transactionUuid := experthelpers.GenerateTransactionUUID()
			formulaObj.TransactionUUID = transactionUuid

			//add new formula to formula ID map
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
		}

		// remove the substring from the last comma
		lenOfLastCommand := len(", calculations.GetExponent()")
		executionTemplateString = executionTemplateString[:len(executionTemplateString)-lenOfLastCommand]
	}

}
