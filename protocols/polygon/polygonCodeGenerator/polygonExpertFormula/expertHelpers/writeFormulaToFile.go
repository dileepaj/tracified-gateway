package experthelpers

import (
	"errors"
	"os"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/utilities"
)

func WriteFormulaContractToFile(contractName string, template string, formulaID string, transactionUUID string, formulaObj model.EthereumExpertFormula) error {
	logger := utilities.NewCustomLogger()
	object := dao.Connection{}
	fo, errInOutput := os.Create(commons.GoDotEnvVariable("") + "/" + contractName + `.sol`)
	if errInOutput != nil {
		logger.LogWriter("Error when creating output file : "+errInOutput.Error(), constants.ERROR)
		return errors.New("Error when creating output files : " + errInOutput.Error())
	}
	defer fo.Close()

	_, errWhenWritingOutput := fo.Write([]byte(template))
	if errWhenWritingOutput != nil {
		logger.LogWriter("Error when writing into the solidity file :"+errWhenWritingOutput.Error(), constants.ERROR)
		return errors.New("Error when writing into the solidity file :" + errWhenWritingOutput.Error())
	}

	//update in the collection
	formulaObj.ActualStatus = 105 // WRITING_CONTRACT_TO_FILE_COMPLETED
	errWhenUpdatingSelectedFields := object.UpdateSelectedPolygonFormulaFields(formulaID, transactionUUID, formulaObj)
	if errWhenUpdatingSelectedFields != nil {
		logger.LogWriter("Error when updating Polygon collections after the file write : "+errWhenUpdatingSelectedFields.Error(), constants.ERROR)
		return errors.New("Error when updating Polygon collections after the file write : " + errWhenUpdatingSelectedFields.Error())
	}
	return nil
}
