package metricBinding

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

/*
 * des- build the manage data for formula definition (compulsory for each formula)
 * return the txnbuild.ManageData object
 * variable definitions and byte used
		1. formulaMapID  - 8 bytes (uint64)  - mapped id stored in the DB for the formula
		2. variableCount - 2 bytes (uint16)  - number of dynamic variables in the formula
		3. activityID    - 8 bytes (uint64)  - mapped id stored in the DB for the activity
		4. future use    - 46 bytes
 * Manage data
		name 64 bytes  - the string "FORMULA METADATA"
		value 64 bytes - formulaMapID + variableCount + activityMapID + future use
*/

func (metric *MetricBinding) BuildFormulaDefinition(formulaMapID uint64, activityMapID uint64, variableCount uint16) (txnbuild.ManageData, string, []byte, error) {
	// covert ulint to byte array anf then to string
	decodedStrFetureUsed, err := hex.DecodeString(fmt.Sprintf("%092d", 0))
	if err != nil {
		return txnbuild.ManageData{}, "", []byte{}, errors.New("future use byte building issue in building formula definition")
	}
	strFetureUsed := string(decodedStrFetureUsed)

	// String for key
	strForKey := "FORMULA METADATA/"
	if len(strForKey) < 64 {
		remain := 64 - len(strForKey)
		setRemainder := strings.Repeat("0", remain)
		strForKey = strForKey + setRemainder
	}

	valueString := stellarprotocols.UInt64ToByteString(formulaMapID) + stellarprotocols.UInt16ToByteString(variableCount) + stellarprotocols.UInt64ToByteString(activityMapID) + strFetureUsed
	keyString := strForKey

	formulaDefinitionBuilder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	logrus.Info("Formula definition - generated key : ", keyString)
	logrus.Info("Formula definition - generated value : ", valueString)

	if len(valueString) != 64 {
		logrus.Error("Length ", len(valueString))
		return txnbuild.ManageData{}, "", []byte{}, errors.New("formula definition builder value length not equal to 64")
	}
	if len(keyString) > 64 || len(keyString) == 0 {
		logrus.Error("Length ", len(keyString))
		return txnbuild.ManageData{}, "", []byte{}, errors.New("formula definition builder value length not equal to 64")
	}

	return formulaDefinitionBuilder, keyString, []byte(valueString), nil
}
