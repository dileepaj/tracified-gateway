package experthelpers

import (
	"errors"
	"os"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/utilities"
)

func WriteFormulaContractToFile(contractName string, template string) error {
	logger := utilities.NewCustomLogger()
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
	return nil
}
