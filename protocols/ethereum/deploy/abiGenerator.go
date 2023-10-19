package deploy

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/sirupsen/logrus"
)

/*
Generate the ABI file for the given smart contract
*/
func GenerateABI(contractName string, reqType string) (string, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	var cmdABIGen *exec.Cmd
	var location string
	abiString := ""
	if reqType == "EXPERT" {
		cmdABIGen = exec.Command("solcjs", "--abi", contractName+".sol", "-o", "build")
		cmdABIGen.Dir = commons.GoDotEnvVariable("EXPERTCONTRACTLOCATION")
	} else if reqType == "METRIC" {
		cmdABIGen = exec.Command("solcjs", "--abi", contractName+".sol", "-o", "metricbuild")
		cmdABIGen.Dir = commons.GoDotEnvVariable("METRICCONTRACTLOCATION")
	} else if reqType == "POLYGONEXPERT" {
		cmdABIGen = exec.Command("solcjs", "--abi", contractName+".sol", "-o", "polygonformulabuild")
		cmdABIGen.Dir = "./assets/contracts/polygon"
	} else {
		logrus.Error("Invalid request type for ABI generator , TYPE : ", reqType)
		return abiString, errors.New("Invalid request type for ABI generator , TYPE : " + reqType)
	}
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
	if reqType == "EXPERT" {
		location = commons.GoDotEnvVariable("EXPERTBUILDLOCATION") + "/" + fileName
	} else if reqType == "METRIC" {
		location = commons.GoDotEnvVariable("METRICBUILDLOCATION") + "/" + fileName
	} else if reqType == "POLYGONEXPERT" {
		location = "./assets/contracts/polygon/polygonformulabuild/" + fileName
	} else {
		logrus.Error("Invalid request type for ABI reader , TYPE : ", reqType)
		return abiString, errors.New("Invalid request type for ABI reader , TYPE : " + reqType)
	}

	abiInByte, errWhenReadingFile := os.ReadFile(location)
	if errWhenReadingFile != nil {
		logrus.Info("Error when reading the ABI file")
		return abiString, errWhenReadingFile
	}

	abiString = string(abiInByte)

	return abiString, nil
}
