package codeGenerator

import (
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/codeGenerator/executionTemplates"
	expertFormula "github.com/dileepaj/tracified-gateway/protocols/stellarprotocols/expertFormula"
	"github.com/dileepaj/tracified-gateway/services"
	"github.com/oklog/ulid"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// initial keywords for the contract
var (
	contractName       = ``
	contractBody       = ``
	startOfTheExecutor = `function executeCalculation() public returns (int) {`
	endOfTheExecutor   = "\n\t" + `}`
)

/*
Generate the smart contract for the solidity formula definitions
*/
func SmartContractGeneratorForFormula(w http.ResponseWriter, r *http.Request, formulaJSON model.FormulaBuildingRequest, fieldCount int) {
	object := dao.Connection{}
	var deployStatus string

	//TODO check the DB with the formula ID to see wether its a duplicate with the deploy and verify status
	formulaDetails, errWhenGettingFormulaDetailsFromDB := object.GetEthFormulaStatus(formulaJSON.MetricExpertFormula.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errWhenGettingFormulaDetailsFromDB != nil {
		logrus.Error("An error occurred when getting formula status, ERROR : ", errWhenGettingFormulaDetailsFromDB)
	}
	if formulaDetails == nil {
		deployStatus = ""
	}
	if formulaDetails != nil {
		deployStatus = formulaDetails.(model.EthereumExpertFormula).Status
		logrus.Info("Deploy status : ", deployStatus)
	}

	if deployStatus == "SUCCESS" {
		logrus.Info("Contract for formula " + formulaJSON.MetricExpertFormula.Name + " has been added to the blockchain")
		commons.JSONErrorReturn(w, r, "Status : "+deployStatus, 400, "Requested formula is in the blockchain")
		return
	} else if deployStatus == "QUEUE" {
		logrus.Info("Requested formula is in the queue, please try again")
		commons.JSONErrorReturn(w, r, "Status : "+deployStatus, 400, "Requested formula is in the queue, please try again")
		return
	} else if deployStatus == "FAILED" {
		//deploy the contract another time
		logrus.Info("Requested formula is in the failed status")
		commons.JSONErrorReturn(w, r, "Status : "+deployStatus, 400, "Requested formula is in the failed status, redeploy")
		return
	} else if deployStatus == "" {
		ethFormulaObj := model.EthereumExpertFormula{
			FormulaID:           formulaJSON.MetricExpertFormula.ID,
			FormulaName:         formulaJSON.MetricExpertFormula.Name,
			MetricExpertFormula: formulaJSON.MetricExpertFormula,
			VariableCount:       int32(fieldCount),
			BINstring:           "",
			ABIstring:           "",
			ContractAddress:     "",
			Timestamp:           time.Now().String(),
			TransactionHash:     "",
			TransactionCost:     "",
			TransactionTime:     "",
			TransactionUUID:     "",
			TransactionSender:   commons.GoDotEnvVariable("ETHEREUMPUBKEY"),
			User:                formulaJSON.User,
			ErrorMessage:        "",
		}
		//generate transaction UUID
		timeNow := time.Now().UTC()
		entropy := rand.New(rand.NewSource(timeNow.UnixNano()))
		id := ulid.MustNew(ulid.Timestamp(timeNow), entropy)
		logrus.Info("TXN UUID : ", id)
		ethFormulaObj.TransactionUUID = id.String()

		//setting up the contract name and starting the contract
		contractName = cases.Title(language.English).String(formulaJSON.MetricExpertFormula.Name)
		contractName = strings.ReplaceAll(contractName, " ", "")

		//call the general header writer
		generalValues, errWhenBuildingGeneralCodeSnippet := WriteGeneralCodeSnippets(formulaJSON, contractName)
		if errWhenBuildingGeneralCodeSnippet != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errWhenBuildingGeneralCodeSnippet.Error()
			//call the DB insert method
			errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
			if errWhenInsertingFormulaToDB != nil {
				logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
				commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
				return
			}
			logrus.Error("Error when writing the general code snippet, ERROR : " + errWhenBuildingGeneralCodeSnippet.Error())
			commons.JSONErrorReturn(w, r, errWhenBuildingGeneralCodeSnippet.Error(), http.StatusInternalServerError, "Error when writing the general code snippet, ERROR : ")
			return
		}

		contractBody = contractBody + generalValues.ResultVariable + generalValues.MetaDataStructure + generalValues.ValueDataStructure + generalValues.VariableStructure + generalValues.SemanticConstantStructure + generalValues.ReferredConstant + generalValues.MetadataDeclaration

		//call the value builder and get the string for the variable initialization and setter
		variableValues, errInGeneratingValues := ValueCodeGenerator(formulaJSON)
		if errInGeneratingValues != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errInGeneratingValues.Error()
			//call the DB insert method
			errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
			if errWhenInsertingFormulaToDB != nil {
				logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
				commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
				return
			}
			logrus.Error("Error in generating codes for values ", errInGeneratingValues.Error())
			commons.JSONErrorReturn(w, r, errInGeneratingValues.Error(), http.StatusInternalServerError, "Error in getting codes for values ")
			return
		}
		contractBody = contractBody + variableValues

		//pass the query to the FCL and get the execution template
		executionTemplate, errInGettingExecutionTemplate := expertFormula.BuildExecutionTemplateByQuery(formulaJSON.MetricExpertFormula.FormulaAsQuery)
		if errInGettingExecutionTemplate != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errInGettingExecutionTemplate.Error()
			//call the DB insert method
			errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
			if errWhenInsertingFormulaToDB != nil {
				logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
				commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
				return
			}
			logrus.Error("Error in generating codes for values ", errInGeneratingValues.Error())
			commons.JSONErrorReturn(w, r, errInGettingExecutionTemplate.Error(), http.StatusInternalServerError, "Error in getting execution template from FCL ")
			return
		}
		ethFormulaObj.ExecutionTemplate = executionTemplate

		//loop through the execution template and getting the built equation
		executionTemplateString, errInExecutionTemplateString := executionTemplates.ExecutionTemplateDivider(executionTemplate)
		if errInExecutionTemplateString != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errInExecutionTemplateString.Error()
			//call the DB insert method
			errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
			if errWhenInsertingFormulaToDB != nil {
				logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
				commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
				return
			}
			logrus.Error("Error in generating codes for values ", errInGeneratingValues.Error())
			commons.JSONErrorReturn(w, r, errInExecutionTemplateString.Error(), http.StatusInternalServerError, "Error in getting execution template from FCL ")
			return
		}

		//setting up the executor (Result)
		commentForExecutor := `// method to get the result of the calculation`
		executorBody := "\t\n\t\t" + `if (result == -9999) {`
		executorBody = executorBody + "\n\t\t\t" + `result` + " = " + executionTemplateString + ";" + "\n\t\t" + `}`
		executorBody = executorBody + "\n\t\t" + `return result;`
		contractBody = contractBody + "\n\n\t" + commentForExecutor + "\n\t" + startOfTheExecutor + executorBody + endOfTheExecutor

		// create the contract
		template := generalValues.License + "\n\n" + generalValues.PragmaLine + "\n\n" + generalValues.ContractStart + "\n\t" + contractBody + "\n" + generalValues.ContractEnd
		//convert the template to base64
		b64Template := base64.StdEncoding.EncodeToString([]byte(template))

		ethFormulaObj.TemplateString = b64Template

		// write the contract to a solidity file
		fo, errInOutput := os.Create(`protocols/ethereum/contracts/` + contractName + `.sol`)
		if errInOutput != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errInOutput.Error()
			//call the DB insert method
			errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
			if errWhenInsertingFormulaToDB != nil {
				logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
				commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
				return
			}
			logrus.Error("Error in creating the output file " + errInOutput.Error())
			commons.JSONErrorReturn(w, r, errInOutput.Error(), http.StatusInternalServerError, "Error in creating the output file ")
			return
		}
		defer fo.Close()
		_, errInWritingOutput := fo.Write([]byte(template))
		if errInWritingOutput != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errInWritingOutput.Error()
			//call the DB insert method
			errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
			if errWhenInsertingFormulaToDB != nil {
				logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
				commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
				return
			}
			logrus.Error("Error in writing the output file " + errInWritingOutput.Error())
			commons.JSONErrorReturn(w, r, errInWritingOutput.Error(), http.StatusInternalServerError, "Error in writing the output file ")
			return
		}

		buildQueueObj := model.SendToQueue{
			EthereumExpertFormula: ethFormulaObj,
			Type:                  "ETHEXPERTFORMULA",
			User:                  formulaJSON.User,
			Status:                "QUEUE",
		}

		//add to queue
		errWhenSendingToQueue := services.SendToQueue(buildQueueObj)
		if errWhenSendingToQueue != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errWhenSendingToQueue.Error()
			//call the DB insert method
			errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
			if errWhenInsertingFormulaToDB != nil {
				logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
				commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
				return
			}
			logrus.Error("Error when sending request to queue " + errWhenSendingToQueue.Error())
			commons.JSONErrorReturn(w, r, errWhenSendingToQueue.Error(), http.StatusInternalServerError, "Error when sending request to queue ")
			return
		}

		//call the DB insert method and send to queue
		logrus.Info("Expert formula is added to the queue")
		ethFormulaObj.Status = "QUEUE"
		errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
		if errWhenInsertingFormulaToDB != nil {
			logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
			commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
			return
		}

		//success response
		w.WriteHeader(http.StatusOK)
		response := model.SuccessResponseExpertFormula{
			Code:      http.StatusOK,
			FormulaID: formulaJSON.MetricExpertFormula.ID,
			Message:   "Expert formula request sent to queue",
		}
		json.NewEncoder(w).Encode(response)
		return

	} else {
		logrus.Info("Invalid formula status " + deployStatus)
		commons.JSONErrorReturn(w, r, deployStatus, http.StatusInternalServerError, "Invalid formula status ")
		return
	}

}
