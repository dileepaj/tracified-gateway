package expertformula

import (
	"encoding/hex"
	"fmt"

	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
)

type ExpertFormula struct{}
/*
des- build the memo according to the protocol
return the txnbuild.ManageData object
*/
// types = 0 - strating manifest
// types = 1 - managedata overflow sign
func (expertFormula ExpertFormula)BuildMemo(types, variableCount int32, mappedFormulaID int64) (string, error) {
	manifest := ""
	if types == 0 {
		manifest = "0000000000AAAAAAAAAA"
	} else if types == 1 {
		manifest = "00000000AAAABBBBCCCC"
	}
	strVariableCount := fmt.Sprintf("%04d", variableCount)
	strFetureUsed := fmt.Sprintf("%06d", 0)
	decodedManifest, err := hex.DecodeString(manifest)
	strManifest := string(decodedManifest)
	if err != nil {
		return "", err
	}
	memo := strManifest +  stellarprotocols.UInt64ToByteString(mappedFormulaID) + strVariableCount + strFetureUsed
	return memo, nil
}
