package deploy

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/sirupsen/logrus"
)

/*
Generate the BIN file for the given smart contract
*/
func GenerateBIN(contractName string) (string, error) {
	//TODO pass the contract name as the sol file
	var out bytes.Buffer
	var stderr bytes.Buffer
	binString := ""
	cmdBINGen := exec.Command("cmd", "/C", "solcjs --bin Calculations.sol -o build")
	cmdBINGen.Dir = commons.GoDotEnvVariable("CONTRACTLOCATION")
	cmdBINGen.Stdout = &out
	cmdBINGen.Stderr = &stderr
	errWhenGettingBIN := cmdBINGen.Run()
	if errWhenGettingBIN != nil {
		logrus.Info("Error when getting the BIN file")
		logrus.Info(fmt.Sprint(errWhenGettingBIN) + ": " + stderr.String())
		return binString, errWhenGettingBIN
	}
	logrus.Info("BIN file generated" + out.String())

	//TODO read the bin file and assign to binString

	return binString, nil
}
