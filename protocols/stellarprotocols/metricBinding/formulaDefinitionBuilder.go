package metricBinding

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

func BuildFormulaDefinition(formulaID uint64, variableCount uint32) (txnbuild.ManageData, error) {

	// covert ulint to byte array anf then to string
	formulaIDByteArray := make([]byte, 8)
	binary.LittleEndian.PutUint64(formulaIDByteArray, formulaID)
	strFormulaID := string(formulaIDByteArray)

	// convert unit to byte array and then to string
	variableCountByteArray := make([]byte, 4)
	binary.LittleEndian.PutUint32(variableCountByteArray, variableCount)
	strVariableCount := string(variableCountByteArray)

	strFutureUse := ""
	if len(strFormulaID+strVariableCount) < 64 {
		remain := 64 - len(strFormulaID+strVariableCount)
		setRemainder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		strFutureUse = setRemainder
	}

	// String for key
	strForKey := "FORMULA METADATA"
	if len(strForKey) < 64 {
		remain := 64 - len(strForKey)
		setRemainder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		strForKey = strForKey + setRemainder
	}

	valueString := strFormulaID + strVariableCount + strFutureUse
	keyString := strForKey

	formulaDefinitionBuilder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	if len(valueString) != 64 {
		logrus.Error("Length ", len(valueString))
		return txnbuild.ManageData{}, errors.New("Formula Definition Builder value length not equal to 64")
	}
	if len(keyString) > 64 || len(keyString) == 0 {
		logrus.Error("Length ", len(keyString))
		return txnbuild.ManageData{}, errors.New("Formula Definition Builder name length should be less than or equal to 64")
	}

	return formulaDefinitionBuilder, nil
}
