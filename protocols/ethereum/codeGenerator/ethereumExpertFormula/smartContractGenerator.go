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
	startOfTheExecutor = `function executeCalculation() public {`
	endOfTheExecutor   = "\n\t" + `}`
)

/*
Generate the smart contract for the solidity formula definitions
*/
func SmartContractGeneratorForFormula(w http.ResponseWriter, r *http.Request, formulaJSON model.FormulaBuildingRequest, fieldCount int) {
	object := dao.Connection{}
	var deployStatus string
	reqType := "EXPERT"

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
		w.WriteHeader(http.StatusBadRequest)
		response := model.SuccessResponseExpertFormula{
			Code:      http.StatusBadRequest,
			FormulaID: formulaJSON.MetricExpertFormula.ID,
			Message:   "Requested formula is in the blockchain :  Status : " + deployStatus,
		}
		json.NewEncoder(w).Encode(response)
		return
	} else if deployStatus == "QUEUE" {
		logrus.Info("Requested formula is in the queue, please try again")
		w.WriteHeader(http.StatusBadRequest)
		response := model.SuccessResponseExpertFormula{
			Code:      http.StatusBadRequest,
			FormulaID: formulaJSON.MetricExpertFormula.ID,
			Message:   "Requested formula is in the queue :  Status : " + deployStatus,
		}
		json.NewEncoder(w).Encode(response)
		return
	} else if deployStatus == "" || deployStatus == "FAILED" {
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
			TransactionUUID:     "",
			GOstring:            "",
			TransactionSender:   commons.GoDotEnvVariable("ETHEREUMPUBKEY"),
			Verify:              formulaJSON.Verify,
			ErrorMessage:        "",
			ActualStatus: 		 101, // SMART_CONTRACT_GENERATION_STARTED
		}

		if deployStatus == "" {
			//generate transaction UUID
			timeNow := time.Now().UTC()
			entropy := rand.New(rand.NewSource(timeNow.UnixNano()))
			id := ulid.MustNew(ulid.Timestamp(timeNow), entropy)
			logrus.Info("TXN UUID : ", id)
			ethFormulaObj.TransactionUUID = id.String()

			// add formula to the formula ID map
			// getting next sequence value
			data, errInGettingNextSequence := object.GetNextSequenceValue("ETHFORMULAID")
			if errInGettingNextSequence != nil {
				logrus.Info("Unable to connect to gateway datastore ", errInGettingNextSequence)
				commons.JSONErrorReturn(w, r, errInGettingNextSequence.Error(), http.StatusInternalServerError, "Error while getting next sequence value for ETHFORMULAID")
				return
			}
			formulaIDmap := model.EthFormulaIDMap{
				FormulaID: formulaJSON.MetricExpertFormula.ID,
				MapID:    data.SequenceValue,
			}
			errorWhenInsertingToFormulaIDMap := object.InsertEthFormulaIDMap(formulaIDmap)
			if errorWhenInsertingToFormulaIDMap != nil {
				logrus.Info("Unable to connect to gateway datastore ", errorWhenInsertingToFormulaIDMap)
				commons.JSONErrorReturn(w, r, errorWhenInsertingToFormulaIDMap.Error(), http.StatusInternalServerError, "Error while inserting to ETHFORMULAIDMAP")
				return
			}

		} else {
			ethFormulaObj.TransactionUUID = formulaDetails.(model.EthereumExpertFormula).TransactionUUID
		}

		//setting up the contract name and starting the contract
		contractName = cases.Title(language.English).String(formulaJSON.MetricExpertFormula.Name)
		contractName = strings.ReplaceAll(contractName, " ", "")
		contractName = contractName + "_" + formulaJSON.MetricExpertFormula.ID

		//call the general header writer
		generalValues, errWhenBuildingGeneralCodeSnippet := WriteGeneralCodeSnippets(formulaJSON, contractName)
		if errWhenBuildingGeneralCodeSnippet != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errWhenBuildingGeneralCodeSnippet.Error()
			ethFormulaObj.ActualStatus = 102	// SMART_CONTRACT_GENERATION_FAILED
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

		contractBody = generalValues.ResultVariable + generalValues.MetaDataStructure + generalValues.ValueDataStructure + generalValues.VariableStructure + generalValues.SemanticConstantStructure + generalValues.ReferredConstant + generalValues.MetadataDeclaration
		contractBody = contractBody + generalValues.ResultDeclaration + generalValues.CalculationObject 

		//call the value builder and get the string for the variable initialization and setter
		variableValues, setterNames, errInGeneratingValues := ValueCodeGenerator(formulaJSON)
		if errInGeneratingValues != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errInGeneratingValues.Error()
			ethFormulaObj.ActualStatus = 102	// SMART_CONTRACT_GENERATION_FAILED
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
		ethFormulaObj.SetterNames = setterNames

		//pass the query to the FCL and get the execution template
		executionTemplate, errInGettingExecutionTemplate := expertFormula.BuildExecutionTemplateByQuery(formulaJSON.MetricExpertFormula.FormulaAsQuery)
		if errInGettingExecutionTemplate != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errInGettingExecutionTemplate.Error()
			ethFormulaObj.ActualStatus = 102	// SMART_CONTRACT_GENERATION_FAILED
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
			ethFormulaObj.ActualStatus = 102	// SMART_CONTRACT_GENERATION_FAILED
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
		contractBody = contractBody + "\n\n\t" + commentForExecutor + "\n\t" + startOfTheExecutor + executorBody + endOfTheExecutor + "\n"

		//getter method
		commentForGetter := "\n\t" + `//get value and exponent`
		getterBody := "\n\t" + `function getValues() public view returns (int256, int256) {`
		getterBody = getterBody + "\n\t\t" + `return (result.value, result.exponent);`
		getterBody = getterBody + "\n\t" + `}` + "\n"
		contractBody = contractBody + commentForGetter + getterBody

		// metadata getter method
		contractBody += generalValues.MetadataGetter

		// create the contract
		template := generalValues.License + "\n\n" + generalValues.PragmaLine + "\n\n" + generalValues.ImportCalculationsSol + "\n\n" + generalValues.ContractStart + "\n\t" + contractBody + "\n" + generalValues.ContractEnd
		//convert the template to base64
		b64Template := base64.StdEncoding.EncodeToString([]byte(template))

		ethFormulaObj.TemplateString = b64Template
		ethFormulaObj.ActualStatus = 103	// SMART_CONTRACT_GENERATION_COMPLETED
		errorWhenUpdatingFormula1 := object.UpdateSelectedEthFormulaFields(ethFormulaObj.FormulaID, ethFormulaObj.TransactionUUID, ethFormulaObj)
		if errorWhenUpdatingFormula1 != nil {
			logrus.Error("Error while updating formula details to the DB after generating smart contract " + errorWhenUpdatingFormula1.Error())
			commons.JSONErrorReturn(w, r, errorWhenUpdatingFormula1.Error(), http.StatusInternalServerError, "Error while updating formula details to the DB after generating smart contract ")
			return
		}

		// write the contract to a solidity file
		fo, errInOutput := os.Create(commons.GoDotEnvVariable("EXPERTCONTRACTLOCATION") + "/" + contractName + `.sol`)
		if errInOutput != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errInOutput.Error()
			ethFormulaObj.ActualStatus = 104	// WRITING_CONTRACT_TO_FILE_FAILED
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
			ethFormulaObj.ActualStatus = 104	// WRITING_CONTRACT_TO_FILE_FAILED
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

		ethFormulaObj.ActualStatus = 105	// WRITING_CONTRACT_TO_FILE_COMPLETED
		errorWhenUpdatingFormula2 := object.UpdateSelectedEthFormulaFields(ethFormulaObj.FormulaID, ethFormulaObj.TransactionUUID, ethFormulaObj)
		if errorWhenUpdatingFormula2 != nil {
			logrus.Error("Error while updating formula details to the DB after writing the contract to file " + errorWhenUpdatingFormula2.Error())
			commons.JSONErrorReturn(w, r, errorWhenUpdatingFormula2.Error(), http.StatusInternalServerError, "Error while updating formula details to the DB after writing the contract to file ")
			return
		}

		//generate the ABI file
		abiString, errWhenGeneratingABI := deploy.GenerateABI(contractName, reqType)
		if errWhenGeneratingABI != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errWhenGeneratingABI.Error()
			ethFormulaObj.ActualStatus = 106	// GENERATING_ABI_FAILED
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
		ethFormulaObj.ActualStatus = 107	// GENERATING_ABI_COMPLETED
		errorWhenUpdatingFormula3 := object.UpdateSelectedEthFormulaFields(ethFormulaObj.FormulaID, ethFormulaObj.TransactionUUID, ethFormulaObj)
		if errorWhenUpdatingFormula3 != nil {
			logrus.Error("Error while updating formula details to the DB after getting abi " + errorWhenUpdatingFormula3.Error())
			commons.JSONErrorReturn(w, r, errorWhenUpdatingFormula3.Error(), http.StatusInternalServerError, "Error while updating formula details to the DB after getting abi ")
			return
		}
		
		//generate the BIN file
		binString, errWhenGeneratingBinFile := deploy.GenerateBIN(contractName, reqType)
		if errWhenGeneratingBinFile != nil {
			ethFormulaObj.Status = "FAILED"
			ethFormulaObj.ErrorMessage = errWhenGeneratingBinFile.Error()
			ethFormulaObj.ActualStatus = 108	// GENERATING_BIN_FAILED
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
		ethFormulaObj.ActualStatus = 109	// GENERATING_BIN_COMPLETED
		errorWhenUpdatingFormula4 := object.UpdateSelectedEthFormulaFields(ethFormulaObj.FormulaID, ethFormulaObj.TransactionUUID, ethFormulaObj)
		if errorWhenUpdatingFormula4 != nil {
			logrus.Error("Error while updating formula details to the DB after getting bin " + errorWhenUpdatingFormula4.Error())
			commons.JSONErrorReturn(w, r, errorWhenUpdatingFormula4.Error(), http.StatusInternalServerError, "Error while updating formula details to the DB after getting bin ")
			return
		}
		//generating go file by converting the code to bas64
		// goString, errWhenGeneratingGoCode := deploy.GenerateGoCode(contractName)
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
		commons.JSONErrorReturn(w, r, deployStatus, http.StatusInternalServerError, "Invalid formula status : ")
		return
	}

}
