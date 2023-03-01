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
	var out bytes.Buffer
	var stderr bytes.Buffer
	goString := ""
	cmdGoGen := exec.Command("powershell", "/C", "./abigen --bin="+contractName+"_sol_"+contractName+".bin --abi="+contractName+"_sol_"+contractName+".abi --pkg="+"build"+" --out="+contractName+".go")
	cmdGoGen.Dir = commons.GoDotEnvVariable("EXPERTBUILDLOCATION")
	cmdGoGen.Stdout = &out
	cmdGoGen.Stderr = &stderr
	errWhenGettingGo := cmdGoGen.Run()
	if errWhenGettingGo != nil {
		logrus.Info("Error when getting the ABI file")
		logrus.Info(fmt.Sprint(errWhenGettingGo) + ": " + stderr.String())
		return goString, errWhenGettingGo
	}
	logrus.Info("Go file generated" + out.String())

	return goString, nil
}
