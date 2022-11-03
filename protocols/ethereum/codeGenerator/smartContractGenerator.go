package codeGenerator

import (
	"net/http"
	"os"
	"strings"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/codeGenerator/executionTemplates"
	expertFormula "github.com/dileepaj/tracified-gateway/protocols/stellarprotocols/expertFormula"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// initial keywords for the contract
var (
	contractName       = ``
	contractBody       = `string public name = ` + contractName + `;`
	startOfTheExecutor = `function Executor() public {`
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
	generalValues := GeneralCodeWriter(contractName)

	//setting up variable to store the results as the first line of the contract body
	contractBody = `uint public Result;`

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

	//meta variable definition
	metaDataVariables := WriteMetaData()
	contractBody = contractBody + metaDataVariables

	contractBody = contractBody + "\n"

	//meta data setter
	metaDataSetter := MetaDataSetter()
	contractBody = contractBody + metaDataSetter

	//setting up the executor (Result)
	executorBody := "\n\t\t" + `Result` + " = " + executionTemplateString + ";"
	contractBody = contractBody + "\n\n\t" + startOfTheExecutor + executorBody + endOfTheExecutor

	// getter for the result
	contractBody = contractBody + "\n\n\t" + `function getResult() public view returns (uint) {` + "\n\t\t" + `return Result;` + "\n\t" + `}`

	//formulaID getter getter
	formulaIDGetter := MetaDataFormulaIDGetter()
	contractBody = contractBody + formulaIDGetter

	//expert PK getter
	expertPKGetter := MetaDataExpertPKGetter()
	contractBody = contractBody + expertPKGetter

	//tenet ID getter
	tenetIDGetter := MetaDataTenantIDGetter()
	contractBody = contractBody + tenetIDGetter

	// create the contract
	template := generalValues.License + "\n\n" + generalValues.StartingCodeLine + "\n\n" + generalValues.ContractStart + "\n\t" + contractBody + "\n" + generalValues.ContractEnd

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
