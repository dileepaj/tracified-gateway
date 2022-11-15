package deploy

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/sirupsen/logrus"
)

/*
Generate Go code from BIN and ABI files
*/
func GenerateGoCode(contractName string) (string, error) {
	//TODO add the contract name as the sol name
	var out bytes.Buffer
	var stderr bytes.Buffer
	goString := ""
	cmdGoGen := exec.Command("powershell", "/C", "./abigen --bin=Calculations_sol_Calculations.bin --abi=Calculations_sol_Calculations.abi --pkg=Calculations --out=Calculations.go")
	cmdGoGen.Dir = commons.GoDotEnvVariable("BUILDLOCATION")
	cmdGoGen.Stdout = &out
	cmdGoGen.Stderr = &stderr
	errWhenGettingGo := cmdGoGen.Run()
	if errWhenGettingGo != nil {
		logrus.Info("Error when getting the ABI file")
		logrus.Info(fmt.Sprint(errWhenGettingGo) + ": " + stderr.String())
		return goString, errWhenGettingGo
	}
	logrus.Info("Go file generated" + out.String())

	//TODO read the abi file and pass the string to abistring

	return goString, nil
}
