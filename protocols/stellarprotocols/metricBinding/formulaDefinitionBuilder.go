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

func (metric *MetricBinding) BuildFormulaDefinition(formulaMapID, activityMapID, stageID uint64, variableCount uint16, stageName string) (txnbuild.ManageData, string, []byte, error) {
	rebuildStargeName := ""
	// covert ulint to byte array anf then to string
	decodedStrFetureUsed, err := hex.DecodeString(fmt.Sprintf("%036d", 0))
	if err != nil {
		return txnbuild.ManageData{}, "", []byte{}, errors.New("Feture used byte building issue in formula definition")
	}
	strFetureUsed := string(decodedStrFetureUsed)

	// String for key
	strForKey := "FORMULA METADATA"
	if len(strForKey) < 64 {
		remain := 64 - len(strForKey)
		setRemainder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		strForKey = strForKey + setRemainder
	}
	if len(stageName) > 20 {
		logrus.Error("Stage name is greater than 20 character limit")
		return txnbuild.ManageData{}, "", []byte{}, errors.New("Strage name is greater than 20 character limit")
	} else {
		if len(stageName) == 20 {
			rebuildStargeName = stageName
		} else if len(stageName) < 20 {
			rebuildStargeName = stageName + "/"
		}
	}

	if len(rebuildStargeName) < 20 {
		remain := 20 - len(rebuildStargeName)
		setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		rebuildStargeName = rebuildStargeName + setReaminder
	}

	valueString := stellarprotocols.UInt64ToByteString(formulaMapID) + stellarprotocols.UInt16ToByteString(variableCount) + stellarprotocols.UInt64ToByteString(activityMapID) + stellarprotocols.UInt64ToByteString(stageID) +
		rebuildStargeName + strFetureUsed
	keyString := strForKey

	formulaDefinitionBuilder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	if len(valueString) != 64 {
		logrus.Error("Length ", len(valueString))
		return txnbuild.ManageData{}, "", []byte{}, errors.New("Formula Definition Builder value length not equal to 64")
	}
	if len(keyString) > 64 || len(keyString) == 0 {
		logrus.Error("Length ", len(keyString))
		return txnbuild.ManageData{}, "", []byte{}, errors.New("Formula Definition Builder name length should be less than or equal to 64")
	}

	return formulaDefinitionBuilder, keyString, []byte(valueString), nil
}
