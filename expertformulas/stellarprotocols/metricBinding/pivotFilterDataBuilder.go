package metricBinding

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/dileepaj/tracified-gateway/expertformulas/stellarprotocols"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

/*
 * des- build the manage data for pivot definition (not compulsory for each formula)
 * return the txnbuild.ManageData object
 * variable definitions and byte used
		1. artifactMapId  - 8 bytes (uint64)  - mapped id stored in the DB for the formula
		2. artifactMapId- 8 bytes (uint64)  - number of dynamic variables in the formula
		3. condition    - 2 bytes (uint64)  - mapped type stored in the DB for the pivot filter
		4. future use    - 46 bytes
 * Manage data
		name 64 bytes  - the string "PIVOT DEFINITION/"
		value 64 bytes -  condition + artifactMapId + artifactMapId + future use
*/

func (metric *MetricBinding) PivotFieldDefinition(artifactTemplateId, artifactId uint64, condition string) (txnbuild.ManageData,error) {
	decodedStrFutureUsed, err := hex.DecodeString(fmt.Sprintf("%092d", 0))
	conditionNumber := 0
	if err != nil {
		return txnbuild.ManageData{}, errors.New("future use byte building issue in building formula definition")
	}
	strFutureUsed := string(decodedStrFutureUsed)

	// String for key
	strForKey := "PIVOT DEFINITION/"
	if len(strForKey) < 64 {
		remain := 64 - len(strForKey)
		setRemainder := strings.Repeat("0", remain)
		strForKey = strForKey + setRemainder
	}
	// can add more condition
	if condition == "EQUAL" {
		conditionNumber = 1
	}

	valueString := stellarprotocols.UInt16ToByteString(uint16(conditionNumber)) + stellarprotocols.UInt64ToByteString(artifactTemplateId) + stellarprotocols.UInt64ToByteString(artifactId) + strFutureUsed
	keyString := strForKey

	pivotDefinitionBuilder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	logrus.Info("pivot filter definition - generated key : ", keyString)
	logrus.Info("pivot filter definition - generated value : ", valueString)

	if len(valueString) != 64 {
		logrus.Error("Length ", len(valueString))
		return txnbuild.ManageData{}, errors.New("pivot filter definition builder value length not equal to 64")
	}
	if len(keyString) > 64 || len(keyString) == 0 {
		logrus.Error("Length ", len(keyString))
		return txnbuild.ManageData{}, errors.New("pivot filter definition builder value length not equal to 64")
	}

	return pivotDefinitionBuilder, nil
}
