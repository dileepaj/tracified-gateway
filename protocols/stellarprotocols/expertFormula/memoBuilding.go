package expertformula

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
)

type AbstracMomoBuilder struct{
	Types int32
	VariableCount int32
	MappedFormulaID int
}
/*
des- build the memo according to the protocol
return the txnbuild.ManageData object
*/
// types = 0 - strating manifest
// types = 1 - managedata overflow sign
func BuildMemo(types, variableCount int32, mappedFormulaID int64) (string, error) {
	manifest := ""
	if types == 0 {
		manifest = "0000000000AAAAAAAAAA"
	} else if types == 1 {
		manifest = "00000000AAAABBBBCCCC"
	}
	srtValueID, err := stellarprotocols.IDToBinary(int64(mappedFormulaID))
	if err != nil {
		return "", errors.New("BuildMemo issue (faormula ID convert to type) " + err.Error())
	}
	strVariableCount := fmt.Sprintf("%04d", variableCount)
	strFetureUsed := fmt.Sprintf("%06d", 0)
	decodedManifest, err := hex.DecodeString(manifest)
	strManifest := string(decodedManifest)
	if err != nil {
		return "", err
	}
	memo := strManifest +  stellarprotocols.ConvertingBinaryToByteString(srtValueID) + strVariableCount + strFetureUsed
	return memo, nil
}
