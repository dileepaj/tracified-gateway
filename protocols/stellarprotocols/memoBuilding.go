package stellarprotocols

import (
	"encoding/hex"
	"fmt"
)

// types = 0 - strating manifest
// types = 1 - managedata overflow sign
func BuildMemo(types, variableCount int32, sequeceValue int64) (string, error) {
	manifest := ""
	if types == 0 {
		manifest = "0000000000AAAAAAAAAA"
	} else if types == 1 {
		manifest = "00000000AAAABBBBCCCC"
	}
	strSequenceValue := fmt.Sprintf("%08d", sequeceValue)
	strVariableCount := fmt.Sprintf("%04d", variableCount)
	strFetureUsed := fmt.Sprintf("%06d", 0)
	decodedManifest, err := hex.DecodeString(manifest)
	strManifest := string(decodedManifest)
	if err != nil {
		return "", err
	}
	memo := strManifest + strSequenceValue + strVariableCount + strFetureUsed
	return memo, nil
}
