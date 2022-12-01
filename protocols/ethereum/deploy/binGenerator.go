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
Generate the BIN file for the given smart contract
*/
func GenerateBIN(contractName string) (string, error) {
	//TODO check the request type Metric or Expert and then call the relevant contract and build location
	var out bytes.Buffer
	var stderr bytes.Buffer
	binString := ""
	cmdBINGen := exec.Command("cmd", "/C", "solcjs --bin "+contractName+".sol -o build")
	cmdBINGen.Dir = commons.GoDotEnvVariable("EXPERTCONTRACTLOCATION")
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
	location := commons.GoDotEnvVariable("EXPERTBUILDLOCATION") + "/" + fileName

	binInByte, errWhenReadingFile := os.ReadFile(location)
	if errWhenReadingFile != nil {
		logrus.Info("Error when reading the ABI file")
		return binString, errWhenReadingFile
	}

	binString = "0x" + string(binInByte)

	return binString, nil
}
