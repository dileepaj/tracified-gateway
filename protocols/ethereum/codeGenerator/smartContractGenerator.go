package codeGenerator

import (
	"net/http"
	"os"
	"strings"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/codeGenerator/executionTemplates"
	expertFormula "github.com/dileepaj/tracified-gateway/protocols/stellarprotocols/expertFormula"
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
func SmartContractGeneratorForFormula(w http.ResponseWriter, r *http.Request, formulaJSON model.FormulaBuildingRequest) {

	//--------------------------------------------TODO----------------------------------------------------------------
	//Metadata DB validations + Common code writer

	//Variable builder
	//DB Check, Initialization, Setter

	//Semantic Constant
	//DB Check , Initialization

	//Ref Constant
	//DB check, Initialization

	//run the execution template
	//build the equation
	//execution writer

	//-------------------------------------------------------------------------------------------------------------

	//setting up the contract name and starting the contract
	contractName = cases.Title(language.English).String(formulaJSON.MetricExpertFormula.Name)
	contractName = strings.ReplaceAll(contractName, " ", "")

	//call the general header writer
	generalValues, errWhenBuildingGeneralCodeSnippet := WriteGeneralCodeSnippets(formulaJSON, contractName)
	if errWhenBuildingGeneralCodeSnippet != nil {
		logrus.Error("Error when writing the general code snippet, ERROR : " + errWhenBuildingGeneralCodeSnippet.Error())
		commons.JSONErrorReturn(w, r, errWhenBuildingGeneralCodeSnippet.Error(), http.StatusInternalServerError, "Error when writing the general code snippet, ERROR : ")
		return
	}

	contractBody = contractBody + generalValues.ResultVariable
	contractBody = contractBody + generalValues.MetaDataStructure
	contractBody = contractBody + generalValues.ValueDataStructure
	contractBody = contractBody + generalValues.VariableStructure
	contractBody = contractBody + generalValues.SemanticConstantStructure
	contractBody = contractBody + generalValues.ReferredConstant
	contractBody = contractBody + generalValues.MetadataDeclaration

	//call the value builder and get the string for the variable initialization and setter
	variableValues, errInGeneratingValues := ValueCodeGenerator(formulaJSON)
	if errInGeneratingValues != nil {
		commons.JSONErrorReturn(w, r, errInGeneratingValues.Error(), http.StatusInternalServerError, "Error in getting codes for values ")
	}
	contractBody = contractBody + variableValues

	//pass the query to the FCL and get the execution template
	executionTemplate, errInGettingExecutionTemplate := expertFormula.BuildExecutionTemplateByQuery(formulaJSON.MetricExpertFormula.FormulaAsQuery)
	if errInGettingExecutionTemplate != nil {
		commons.JSONErrorReturn(w, r, errInGettingExecutionTemplate.Error(), http.StatusInternalServerError, "Error in getting execution template from FCL ")
		return
	}

	//loop through the execution template and getting the built equation
	executionTemplateString, errInExecutionTemplateString := executionTemplates.ExecutionTemplateDivider(executionTemplate)
	if errInExecutionTemplateString != nil {
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

	// write the contract to a solidity file
	fo, errInOutput := os.Create(`protocols/ethereum/contracts/` + contractName + `.sol`)
	if errInOutput != nil {
		commons.JSONErrorReturn(w, r, errInOutput.Error(), http.StatusInternalServerError, "Error in creating the output file ")
		return
	}
	defer fo.Close()
	_, errInWritingOutput := fo.Write([]byte(template))
	if errInWritingOutput != nil {
		commons.JSONErrorReturn(w, r, errInWritingOutput.Error(), http.StatusInternalServerError, "Error in writing the output file ")
		return
	}
}
