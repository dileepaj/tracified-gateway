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

/**
 * des : This function is used to build the general master data info manage data operation (only for the master data)
 * Key - 64 Bytes - "MASTER DATA DEFINITION"
 * Value - 8 bytes - Artifact ID
 *  	 - 1 byte - Traceability data type
 *  	 - 55 byte - Future use
 */

func (metric *MetricBinding) BuildGeneralMasterDataInfo(artifactID uint64, traceabilityDataType uint) (txnbuild.ManageData, error) {

	key := "MASTER DATA DEFINITION"
	// check the lengths and append 0s if needed
	if len(key) < 64 {
		key = key + "/"
		if len(key) < 64 {
			key = key + strings.Repeat("0", 64-len(key))
		}
	}

	// convert the traceability data type to string
	tdType, errInTDPTypeConvert := stellarprotocols.Int8ToByteString(uint8(traceabilityDataType))
	if errInTDPTypeConvert != nil {
		logrus.Error("Error when converting traceability data type(generalMasterDataInfoBuilding.go) " + errInTDPTypeConvert.Error())
		return txnbuild.ManageData{}, errors.New("error when converting traceability data type " + errInTDPTypeConvert.Error())
	}

	decodedStrFutureUsed, err := hex.DecodeString(fmt.Sprintf("%0110d", 0))
	if err != nil {
		logrus.Error("Error in decoding the future use string(generalMasterDataInfoBuilding.go)")
		return txnbuild.ManageData{}, errors.New("error in decoding the future use string")
	}
	futureUse := string(decodedStrFutureUsed)

	keyString := key
	valueString := stellarprotocols.UInt64ToByteString(artifactID) + tdType + futureUse

	logrus.Info("General master data info - Key : ", keyString)
	logrus.Info("General master data info - Value : ", valueString)
	
	generalInfoBuilder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	// check the lengths
	if len(keyString) > 64 {
		logrus.Error("Key string length : ", len(keyString))
		return txnbuild.ManageData{}, errors.New("length issue on key field on the general master data info building")
	}
	if len(valueString) > 64 {
		logrus.Error("Value string length : ", len(valueString))
		return txnbuild.ManageData{}, errors.New("length issue on value field on the general master data info building")
	}

	return generalInfoBuilder, nil
}