package deploy

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/sirupsen/logrus"
)

/*
Generate the ABI file for the given smart contract
*/
func GenerateABI(contractName string) (string, error) {
	//TODO add the contract name as the sol name
	var out bytes.Buffer
	var stderr bytes.Buffer
	abiString := ""
	cmdABIGen := exec.Command("cmd", "/C", "solcjs --abi Calculations.sol -o build")
	cmdABIGen.Dir = commons.GoDotEnvVariable("CONTRACTLOCATION")
	cmdABIGen.Stdout = &out
	cmdABIGen.Stderr = &stderr
	errWhenGettingABI := cmdABIGen.Run()
	if errWhenGettingABI != nil {
		logrus.Info("Error when getting the ABI file")
		logrus.Info(fmt.Sprint(errWhenGettingABI) + ": " + stderr.String())
		return abiString, errWhenGettingABI
	}
	logrus.Info("ABI file generated" + out.String())

	//TODO read the abi file and pass the string to abistring

	return abiString, nil
}
