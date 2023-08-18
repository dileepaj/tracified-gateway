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
Generate the BIN file for the given smart contract
*/
func GenerateBIN(contractName string, reqType string) (string, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	var cmdBINGen *exec.Cmd
	var location string
	binString := ""
	runningOs := commons.GoDotEnvVariable("RUNNING_OS")
	if reqType == "EXPERT" {
		if runningOs == "windows" {
			cmdBINGen = exec.Command("cmd", "/C", "solcjs --bin "+contractName+".sol -o build")
		} else if runningOs == "ubuntu" {
			cmdBINGen = exec.Command("solcjs", "--bin", contractName+".sol", "-o", "build")
		} else if runningOs == "linux" {
			cmdBINGen = exec.Command("solcjs", "--bin", contractName+".sol", "-o", "build")
		}
		cmdBINGen.Dir = commons.GoDotEnvVariable("EXPERTCONTRACTLOCATION")
	} else if reqType == "METRIC" {
		if runningOs == "windows" {
			cmdBINGen = exec.Command("cmd", "/C", "solcjs --bin "+contractName+".sol -o metricbuild")
		} else if runningOs == "ubuntu" {
			cmdBINGen = exec.Command("solcjs", "--bin", contractName+".sol", "-o", "metricbuild")
		} else if runningOs == "linux" {
			cmdBINGen = exec.Command("solcjs", "--bin", contractName+".sol", "-o", "metricbuild")
		}
		cmdBINGen.Dir = commons.GoDotEnvVariable("METRICCONTRACTLOCATION")
	} else if reqType == "POLYGONEXPERT" {
		if runningOs == "windows" {
			cmdBINGen = exec.Command("cmd", "/C", "solcjs --bin "+contractName+".sol -o polygonformulabuild")
		} else if runningOs == "ubuntu" {
			cmdBINGen = exec.Command("solcjs", "--bin", contractName+".sol", "-o", "polygonformulabuild")
		} else if runningOs == "linux" {
			cmdBINGen = exec.Command("solcjs", "--bin", contractName+".sol", "-o", "polygonformulabuild")
		}
		cmdBINGen.Dir = "./protocols/polygon/polygonCodeGenerator/polygonExpertFormula/contracts"
	} else {
		logrus.Error("Invalid request type for BIN generator , TYPE : ", reqType)
		return binString, errors.New("Invalid request type for BIN generator , TYPE : " + reqType)
	}
	cmdBINGen.Stdout = &out
	cmdBINGen.Stderr = &stderr
	errWhenGettingBIN := cmdBINGen.Run()
	if errWhenGettingBIN != nil {
		logrus.Info("Error when getting the BIN file")
		logrus.Info(fmt.Sprint(errWhenGettingBIN) + ": " + stderr.String())
		return binString, errWhenGettingBIN
	}
	logrus.Info("BIN file generated" + out.String())

	//build file name
	fileName := contractName + "_sol_" + contractName + ".bin"
	if reqType == "EXPERT" {
		location = commons.GoDotEnvVariable("EXPERTBUILDLOCATION") + "/" + fileName
	} else if reqType == "METRIC" {
		location = commons.GoDotEnvVariable("METRICBUILDLOCATION") + "/" + fileName
	} else if reqType == "POLYGONEXPERT" {
		location = "./protocols/polygon/polygonCodeGenerator/polygonExpertFormula/contracts/polygonformulabuild" + "/" + fileName
	} else {
		logrus.Error("Invalid request type for BIN reader , TYPE : ", reqType)
		return binString, errors.New("Invalid request type for BIN reader , TYPE : " + reqType)
	}

	binInByte, errWhenReadingFile := os.ReadFile(location)
	if errWhenReadingFile != nil {
		logrus.Info("Error when reading the BIN file")
		return binString, errWhenReadingFile
	}

	binString = "0x" + string(binInByte)

	return binString, nil
}
