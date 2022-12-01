package deploy

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/sirupsen/logrus"
)

/*
Generate the ABI file for the given smart contract
*/
func GenerateABI(contractName string) (string, error) {
	//TODO check the request type Metric or Expert and then call the relevant contract and build location
	var out bytes.Buffer
	var stderr bytes.Buffer
	abiString := ""
	cmdABIGen := exec.Command("cmd", "/C", "solcjs --abi "+contractName+".sol -o build")
	cmdABIGen.Dir = commons.GoDotEnvVariable("EXPERTCONTRACTLOCATION")
	cmdABIGen.Stdout = &out
	cmdABIGen.Stderr = &stderr
	errWhenGettingABI := cmdABIGen.Run()
	if errWhenGettingABI != nil {
		logrus.Info("Error when getting the ABI file")
		logrus.Info(fmt.Sprint(errWhenGettingABI) + ": " + stderr.String())
		return abiString, errWhenGettingABI
	}
	logrus.Info("ABI file generated" + out.String())

	//build the abi file name
	fileName := contractName + "_sol_" + contractName + ".abi"
	location := commons.GoDotEnvVariable("EXPERTBUILDLOCATION") + "/" + fileName

	abiInByte, errWhenReadingFile := os.ReadFile(location)
	if errWhenReadingFile != nil {
		logrus.Info("Error when reading the ABI file")
		return abiString, errWhenReadingFile
	}

	abiString = string(abiInByte)

	return abiString, nil
}
