package deploy

import (
	"bytes"
	"fmt"
	"os/exec"
)

/*
Generate the ABI file for the given smart contract
*/
func GenerateABI() {
	cmd := exec.Command("cmd", "/C", "cd protocols"+`\`+"ethereum"+`\`+"contracts")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error when changing the location")
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}

	cmdABIGen := exec.Command("cmd", "/C", "solcjs --abi UrbanWaterUsage.sol")
	cmdABIGen.Stdout = &out
	cmdABIGen.Stderr = &stderr
	errWhenGettingABI := cmdABIGen.Run()
	if errWhenGettingABI != nil {
		fmt.Println("Error when getting the ABI file")
		fmt.Println(fmt.Sprint(errWhenGettingABI) + ": " + stderr.String())
		return
	}
	fmt.Println("Result : " + out.String())
}
