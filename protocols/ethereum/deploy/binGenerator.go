package deploy

import (
	"bytes"
	"fmt"
	"os/exec"
)

/*
Generate the BIN file for the given smart contract
*/
func GenerateBIN() {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmdBINGen := exec.Command("cmd", "/C", "solcjs --bin UrbanWaterUsage.sol -o build")
	cmdBINGen.Stdout = &out
	cmdBINGen.Stderr = &stderr
	errWhenGettingBIN := cmdBINGen.Run()
	if errWhenGettingBIN != nil {
		fmt.Println("Error when getting the BIN file")
		fmt.Println(fmt.Sprint(errWhenGettingBIN) + ": " + stderr.String())
		return
	}
	fmt.Println("Result : " + out.String())
}
