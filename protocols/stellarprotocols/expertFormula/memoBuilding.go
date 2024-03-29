package expertformula

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
)

type ExpertFormula struct{}

/*
des- build the memo according to the protocol
return the txnbuild.ManageData object
	manifest - 10 bytes
	formula id - 8 bytes
	no of values - 4
	future use 6
*/
// types = 0 - strating manifest
// types = 1 - managedata overflow sign
// tenant user (publisher) public key length - 2 bytes (uint16)
func (expertFormula ExpertFormula) BuildMemo(types uint8, variableCount uint32, mappedFormulaID uint64) (string, string, error) {
	manifest := ""
	if types == 0 {
		manifest = "0000000000AAAAAAAAAA"
	} else if types == 1 {
		manifest = "00000000AAAABBBBCCCC"
	}
	decodedStrFutureUse, err := hex.DecodeString(fmt.Sprintf("%012d", 0))
	if err != nil {
		return "", "", errors.New("error in decoding future use " + err.Error())
	}
	strFetureUsed := string(decodedStrFutureUse)
	decodedManifest, err := hex.DecodeString(manifest)
	if err != nil {
		return "", "", errors.New("error in decoding manifest " + err.Error())
	}
	strManifest := string(decodedManifest)

	memo := strManifest + stellarprotocols.UInt64ToByteString(mappedFormulaID) + stellarprotocols.UInt32ToByteString(variableCount) + strFetureUsed
	return memo, manifest, nil
}
