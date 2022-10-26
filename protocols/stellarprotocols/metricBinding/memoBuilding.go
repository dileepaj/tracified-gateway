package metricBinding

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	"github.com/sirupsen/logrus"
)

/*
des- build the memo according to the protocol
return the txnbuild.ManageData object

Fields
	1. Manifest      - 10 bytes (uint64)  - used to identify the transacions
	2. Metric ID    - 8 bytes (uint32)  - mapped id stored in the DB for the tenant
	4. Tenant ID     - 4 bytes - tracified tanant id
	3. No of formulas - 2 bytes (uint16)  - no of formulas in the defined metric
	5. No of manage data in the current transaction - 1 byte
	6. Future use - 3 bytes

types = 0 - strating manifest
types = 1 - managedata overflow sign
*/

type MetricBinding struct{}

func (metric *MetricBinding) BuildMemo(types uint8, mapMetricId uint64, tenantId uint32, noOfFormula uint16, managedatalenth uint8) (string, error) {
	manifest := ""
	if types == 0 {
		manifest = "0000000011AAAAAAAAAA"
	} else if types == 1 {
		manifest = "00000011AAAABBBBCCCC"
	}
	decodedManifest, err := hex.DecodeString(manifest)
	if err != nil {
		return "", err
	}
	strManifest := string(decodedManifest)
	decodedStrFetureUsed, err := hex.DecodeString(fmt.Sprintf("%06d", 0))
	if err != nil {
		logrus.Error("Future used byte building issue in memo building")
		return "", errors.New("Future used byte building issue in memo building")
	}
	// convert data type Int to byte string
	srtManageDataLength, err := stellarprotocols.Int8ToByteString(managedatalenth) // TODO limite the byte, if the user put 8-byte number this should give an error
	if err != nil {
		return "", errors.New("Error when converting data types to byte in memo " + err.Error())
	}
	memo := strManifest + stellarprotocols.UInt64ToByteString(mapMetricId) + stellarprotocols.UInt32ToByteString(tenantId) + stellarprotocols.UInt16ToByteString(noOfFormula) + srtManageDataLength + string(decodedStrFetureUsed)
	if len(memo) > 28 {
		return "", errors.New("Metric binding memo sholud be a 28 bytes")
	}
	logrus.Info("Builded memo : ", memo)
	return memo, nil
}
