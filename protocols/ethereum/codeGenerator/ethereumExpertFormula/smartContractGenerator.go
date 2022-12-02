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
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/codeGenerator/ethereumExpertFormula/executionTemplates"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/deploy"
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
	startOfTheExecutor = `function executeCalculation() public returns (int256, int256) {`
	endOfTheExecutor   = "\n\t" + `}`
)

/*
Generate the smart contract for the solidity formula definitions
*/
func SmartContractGeneratorForFormula(w http.ResponseWriter, r *http.Request, formulaJSON model.FormulaBuildingRequest, fieldCount int) {
	object := dao.Connection{}
	var deployStatus string

	// check if the contract name is already present in the database with the status SUCCESS
	existingFormula, errInCheckingDuplicateNames := object.GetEthFormulaByName(formulaJSON.MetricExpertFormula.Name).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInCheckingDuplicateNames != nil {
		logrus.Error("An error occurred when checking duplicate formula names, ERROR: ", errInCheckingDuplicateNames)
	}
	if existingFormula != nil {
		logrus.Info("Contract for formula " + formulaJSON.MetricExpertFormula.Name + " already exists")
		commons.JSONErrorReturn(w, r, "", 400, "Requested contract name already exists")
		return
	}

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
	} else if deployStatus == "" || deployStatus == "FAILED" {
		// insert to map or get the mapped formula ID
		formulaMapID, errWhenMapping := MapFormulaID(formulaJSON.MetricExpertFormula.ID, deployStatus) 
		if errWhenMapping != nil {
			logrus.Error("An error occurred when mapping formula ID, ERROR : ", errWhenMapping)
		}
		logrus.Info("Formula map ID : ", formulaMapID)
		
		if deployStatus == "FAILED" {
			logrus.Info("Requested formula is in the failed status, trying to redeploy")
		} else {
			logrus.Info("New expert formula request, initiating new deployment")
		}
		ethFormulaObj := model.EthereumExpertFormula{
			FormulaID:           formulaJSON.MetricExpertFormula.ID,
			FormulaName:         formulaJSON.MetricExpertFormula.Name,
			MetricExpertFormula: formulaJSON.MetricExpertFormula,
			VariableCount:       int32(fieldCount),
			ContractAddress:     "",
			Timestamp:           time.Now().String(),
			TransactionHash:     "",
			TransactionCost:     "",
			TransactionTime:     "",
			TransactionUUID:     "",
			GOstring:            "",
			TransactionSender:   commons.GoDotEnvVariable("ETHEREUMPUBKEY"),
			User:                formulaJSON.User,
			ErrorMessage:        "",
		}

		if deployStatus == "" {
			//generate transaction UUID
			timeNow := time.Now().UTC()
			entropy := rand.New(rand.NewSource(timeNow.UnixNano()))
			id := ulid.MustNew(ulid.Timestamp(timeNow), entropy)
			logrus.Info("TXN UUID : ", id)
			ethFormulaObj.TransactionUUID = id.String()
		} else {
			ethFormulaObj.TransactionUUID = formulaDetails.(model.EthereumExpertFormula).TransactionUUID
		}

		//setting up the contract name and starting the contract
		contractName = cases.Title(language.English).String(formulaJSON.MetricExpertFormula.Name)
		contractName = strings.ReplaceAll(contractName, " ", "")

		//call the general header writer
		generalValues, errWhenBuildingGeneralCodeSnippet := WriteGeneralCodeSnippets(formulaJSON, contractName)
		if errWhenBuildingGeneralCodeSnippet != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errWhenBuildingGeneralCodeSnippet.Error()
			if deployStatus == "" {
				errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
				if errWhenInsertingFormulaToDB != nil {
					logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
					return
				}
			} else if deployStatus == "FAILED" {
				errWhenUpdatingFormulaToDB := object.UpdateEthereumFormulaStatus(ethFormulaObj.FormulaID, ethFormulaObj.TransactionUUID, ethFormulaObj)
				if errWhenUpdatingFormulaToDB != nil {
					logrus.Error("Error while updating formula details to the DB " + errWhenUpdatingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenUpdatingFormulaToDB.Error(), http.StatusInternalServerError, "Error while updating formula details to the DB ")
					return
				}
			}

			logrus.Error("Error when writing the general code snippet, ERROR : " + errWhenBuildingGeneralCodeSnippet.Error())
			commons.JSONErrorReturn(w, r, errWhenBuildingGeneralCodeSnippet.Error(), http.StatusInternalServerError, "Error when writing the general code snippet, ERROR : ")
			return
		}

		contractBody = contractBody + generalValues.ResultVariable + generalValues.MetaDataStructure + generalValues.ValueDataStructure + generalValues.VariableStructure + generalValues.SemanticConstantStructure + generalValues.ReferredConstant + generalValues.MetadataDeclaration
		contractBody = contractBody + generalValues.ResultDeclaration + generalValues.CalculationObject

		//call the value builder and get the string for the variable initialization and setter
		variableValues, errInGeneratingValues := ValueCodeGenerator(formulaJSON)
		if errInGeneratingValues != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errInGeneratingValues.Error()
			if deployStatus == "" {
				errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
				if errWhenInsertingFormulaToDB != nil {
					logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
					return
				}
			} else if deployStatus == "FAILED" {
				errWhenUpdatingFormulaToDB := object.UpdateEthereumFormulaStatus(ethFormulaObj.FormulaID, ethFormulaObj.TransactionUUID, ethFormulaObj)
				if errWhenUpdatingFormulaToDB != nil {
					logrus.Error("Error while updating formula details to the DB " + errWhenUpdatingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenUpdatingFormulaToDB.Error(), http.StatusInternalServerError, "Error while updating formula details to the DB ")
					return
				}
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
			if deployStatus == "" {
				errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
				if errWhenInsertingFormulaToDB != nil {
					logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
					return
				}
			} else if deployStatus == "FAILED" {
				errWhenUpdatingFormulaToDB := object.UpdateEthereumFormulaStatus(ethFormulaObj.FormulaID, ethFormulaObj.TransactionUUID, ethFormulaObj)
				if errWhenUpdatingFormulaToDB != nil {
					logrus.Error("Error while updating formula details to the DB " + errWhenUpdatingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenUpdatingFormulaToDB.Error(), http.StatusInternalServerError, "Error while updating formula details to the DB ")
					return
				}
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
			if deployStatus == "" {
				errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
				if errWhenInsertingFormulaToDB != nil {
					logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
					return
				}
			} else if deployStatus == "FAILED" {
				errWhenUpdatingFormulaToDB := object.UpdateEthereumFormulaStatus(ethFormulaObj.FormulaID, ethFormulaObj.TransactionUUID, ethFormulaObj)
				if errWhenUpdatingFormulaToDB != nil {
					logrus.Error("Error while updating formula details to the DB " + errWhenUpdatingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenUpdatingFormulaToDB.Error(), http.StatusInternalServerError, "Error while updating formula details to the DB ")
					return
				}
			}
			logrus.Error("Error in generating codes for values ", errInGeneratingValues.Error())
			commons.JSONErrorReturn(w, r, errInExecutionTemplateString.Error(), http.StatusInternalServerError, "Error in getting execution template from FCL ")
			return
		}

		// remove the substring from the last comma
		lenOfLastCommand := len(", calculations.GetExponent()")
		executionTemplateString = executionTemplateString[:len(executionTemplateString)-lenOfLastCommand]

		//setting up the executor (Result)
		commentForExecutor := `// method to get the result of the calculation`
		executorBody := "\n\t\t" + `result.value` + " = " + executionTemplateString + ";" + "\n\t\t"
		executorBody = executorBody + `result.exponent = calculations.GetExponent();` + "\n\t\t"
		executorBody = executorBody + "\n\t\t" + `return (result.value, result.exponent);`
		contractBody = contractBody + "\n\n\t" + commentForExecutor + "\n\t" + startOfTheExecutor + executorBody + endOfTheExecutor

		// create the contract
		template := generalValues.License + "\n\n" + generalValues.PragmaLine + "\n\n" + generalValues.ImportCalculationsSol + "\n\n" + generalValues.ContractStart + "\n\t" + contractBody + "\n" + generalValues.ContractEnd
		//convert the template to base64
		b64Template := base64.StdEncoding.EncodeToString([]byte(template))

		ethFormulaObj.TemplateString = b64Template

		// write the contract to a solidity file
		fo, errInOutput := os.Create(commons.GoDotEnvVariable("EXPERTCONTRACTLOCATION") + "/" + contractName + `.sol`)
		if errInOutput != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errInOutput.Error()
			//call the DB insert method
			if deployStatus == "" {
				errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
				if errWhenInsertingFormulaToDB != nil {
					logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
					return
				}
			} else if deployStatus == "FAILED" {
				errWhenUpdatingFormulaToDB := object.UpdateEthereumFormulaStatus(ethFormulaObj.FormulaID, ethFormulaObj.TransactionUUID, ethFormulaObj)
				if errWhenUpdatingFormulaToDB != nil {
					logrus.Error("Error while updating formula details to the DB " + errWhenUpdatingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenUpdatingFormulaToDB.Error(), http.StatusInternalServerError, "Error while updating formula details to the DB ")
					return
				}
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
			if deployStatus == "" {
				errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
				if errWhenInsertingFormulaToDB != nil {
					logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
					return
				}
			} else if deployStatus == "FAILED" {
				errWhenUpdatingFormulaToDB := object.UpdateEthereumFormulaStatus(ethFormulaObj.FormulaID, ethFormulaObj.TransactionUUID, ethFormulaObj)
				if errWhenUpdatingFormulaToDB != nil {
					logrus.Error("Error while updating formula details to the DB " + errWhenUpdatingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenUpdatingFormulaToDB.Error(), http.StatusInternalServerError, "Error while updating formula details to the DB ")
					return
				}
			}
			logrus.Error("Error in writing the output file " + errInWritingOutput.Error())
			commons.JSONErrorReturn(w, r, errInWritingOutput.Error(), http.StatusInternalServerError, "Error in writing the output file ")
			return
		}

		//generate the ABI file
		abiString, errWhenGeneratingABI := deploy.GenerateABI(contractName)
		if errWhenGeneratingABI != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errWhenGeneratingABI.Error()
			if deployStatus == "" {
				errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
				if errWhenInsertingFormulaToDB != nil {
					logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
					return
				}
			} else if deployStatus == "FAILED" {
				errWhenUpdatingFormulaToDB := object.UpdateEthereumFormulaStatus(ethFormulaObj.FormulaID, ethFormulaObj.TransactionUUID, ethFormulaObj)
				if errWhenUpdatingFormulaToDB != nil {
					logrus.Error("Error while updating formula details to the DB " + errWhenUpdatingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenUpdatingFormulaToDB.Error(), http.StatusInternalServerError, "Error while updating formula details to the DB ")
					return
				}
			}
			logrus.Info("Error when generating ABI file, ERROR : " + errWhenGeneratingABI.Error())
			commons.JSONErrorReturn(w, r, errWhenGeneratingABI.Error(), http.StatusInternalServerError, "Error when generating ABI file, ERROR : ")
			return
		}
		ethFormulaObj.ABIstring = abiString

		//generate the BIN file
		binString, errWhenGeneratingBinFile := deploy.GenerateBIN(contractName)
		if errWhenGeneratingBinFile != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errWhenGeneratingBinFile.Error()
			if deployStatus == "" {
				errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
				if errWhenInsertingFormulaToDB != nil {
					logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
					return
				}
			} else if deployStatus == "FAILED" {
				errWhenUpdatingFormulaToDB := object.UpdateEthereumFormulaStatus(ethFormulaObj.FormulaID, ethFormulaObj.TransactionUUID, ethFormulaObj)
				if errWhenUpdatingFormulaToDB != nil {
					logrus.Error("Error while updating formula details to the DB " + errWhenUpdatingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenUpdatingFormulaToDB.Error(), http.StatusInternalServerError, "Error while updating formula details to the DB ")
					return
				}
			}
			logrus.Info("Error when generating BIN file, ERROR : " + errWhenGeneratingBinFile.Error())
			commons.JSONErrorReturn(w, r, errWhenGeneratingBinFile.Error(), http.StatusInternalServerError, "Error when generating BIN file, ERROR : ")
			return
		}
		ethFormulaObj.BINstring = binString

		//generating go file by converting the code to bas64
		// goString, errWhenGeneratingGoCode := deploy.GenerateGoCode("Calculations")
		// if errWhenGeneratingGoCode != nil {
		// 	ethFormulaObj.Status = "FAILED"
		// 	ethFormulaObj.ErrorMessage = errWhenGeneratingGoCode.Error()
		// 	//call the DB insert method
		// 	errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
		// 	if errWhenInsertingFormulaToDB != nil {
		// 		logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
		// 		commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
		// 		return
		// 	}
		// 	logrus.Info("Error when generating Go file, ERROR : " + errWhenGeneratingGoCode.Error())
		// 	commons.JSONErrorReturn(w, r, errWhenGeneratingGoCode.Error(), http.StatusInternalServerError, "Error when generating Go file, ERROR : ")
		// 	return
		// }
		// ethFormulaObj.GOstring = goString

		ethFormulaObj.ContractName = contractName

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
			if deployStatus == "" {
				errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
				if errWhenInsertingFormulaToDB != nil {
					logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
					return
				}
			} else if deployStatus == "FAILED" {
				errWhenUpdatingFormulaToDB := object.UpdateEthereumFormulaStatus(ethFormulaObj.FormulaID, ethFormulaObj.TransactionUUID, ethFormulaObj)
				if errWhenUpdatingFormulaToDB != nil {
					logrus.Error("Error while updating formula details to the DB " + errWhenUpdatingFormulaToDB.Error())
					commons.JSONErrorReturn(w, r, errWhenUpdatingFormulaToDB.Error(), http.StatusInternalServerError, "Error while updating formula details to the DB ")
					return
				}
			}
			logrus.Error("Error when sending request to queue " + errWhenSendingToQueue.Error())
			commons.JSONErrorReturn(w, r, errWhenSendingToQueue.Error(), http.StatusInternalServerError, "Error when sending request to queue ")
			return
		}

		//call the DB insert method and send to queue
		logrus.Info("Expert formula is added to the queue")
		ethFormulaObj.Status = "QUEUE"
		if deployStatus == "" {
			errWhenInsertingFormulaToDB := object.InsertToEthFormulaDetails(ethFormulaObj)
			if errWhenInsertingFormulaToDB != nil {
				logrus.Error("Error while inserting formula details to the DB " + errWhenInsertingFormulaToDB.Error())
				commons.JSONErrorReturn(w, r, errWhenInsertingFormulaToDB.Error(), http.StatusInternalServerError, "Error while inserting formula details to the DB ")
				return
			}
		} else if deployStatus == "FAILED" {
			errWhenUpdatingFormulaToDB := object.UpdateEthereumFormulaStatus(ethFormulaObj.FormulaID, ethFormulaObj.TransactionUUID, ethFormulaObj)
			if errWhenUpdatingFormulaToDB != nil {
				logrus.Error("Error while updating formula details to the DB " + errWhenUpdatingFormulaToDB.Error())
				commons.JSONErrorReturn(w, r, errWhenUpdatingFormulaToDB.Error(), http.StatusInternalServerError, "Error while updating formula details to the DB ")
				return
			}
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
