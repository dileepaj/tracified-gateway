package deletecontract

import (
	"errors"
	"os"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/utilities"
)

func DeleteExpertContract(contractName string) error {
	logger := utilities.NewCustomLogger()
	//delete the solidity file
	contractFilePath := commons.GoDotEnvVariable("POLYGONEXPERTLOCATION") + "/" + contractName + `.sol`
	_, err := os.Stat(contractFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.LogWriter("Solidity file "+contractName+" is not found", constants.ERROR)
			return errors.New("Solidity file is not found")
		}
		return err
	}
	// Attempt to delete the file
	err = os.Remove(contractFilePath)
	if err != nil {
		return err
	}

	//delete the ABI file
	abiFilePath := commons.GoDotEnvVariable("POLYGONEXPERTBUILDLOCATION") + "/" + contractName + "_sol_" + contractName + ".abi"
	_, errAbi := os.Stat(abiFilePath)
	if errAbi != nil {
		if os.IsNotExist(errAbi) {
			logger.LogWriter("ABI file "+contractName+" is not found", constants.ERROR)
			return errors.New("ABI file is not found")
		}
		return errAbi
	}
	// Attempt to delete the file
	errAbi = os.Remove(contractFilePath)
	if errAbi != nil {
		return errAbi
	}

	//delete the BIN file
	binFilePath := commons.GoDotEnvVariable("POLYGONEXPERTBUILDLOCATION") + "/" + contractName + "_sol_" + contractName + ".bin"
	_, errBin := os.Stat(binFilePath)
	if errBin != nil {
		if os.IsNotExist(errBin) {
			logger.LogWriter("BIN file "+contractName+" is not found", constants.ERROR)
			return errors.New("BIN file is not found")
		}
		return errBin
	}
	// Attempt to delete the file
	errBin = os.Remove(contractFilePath)
	if errBin != nil {
		return errBin
	}

	return nil
}
