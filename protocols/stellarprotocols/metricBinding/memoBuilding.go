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
	1. Metric ID      - 8 bytes (uint64)  - mapped id stored in the DB for the metric
	2. Tenant ID      - 4 bytes (uint32)  - mapped id stored in the DB for the tenant
	3. No of formulas - 2 bytes (uint16)  - no of formulas in the defined metric
	4. tenant user (publisher) public key length - 2 bytes (uint16)
	4. Future Use     - 12 bytes
*/

type MetricBinding struct{}

func (metric *MetricBinding) BuildMemo(mapMetricId uint64, tenantId uint32, noOfFormula, publickeyLength uint16) (string, error) {
	decodedStrFetureUsed, err := hex.DecodeString(fmt.Sprintf("%024d", 0))
	if err != nil {
		logrus.Error("Future used byte building issue in memo building")
		return "", errors.New("Future used byte building issue in memo building")
	}
	memo := stellarprotocols.UInt64ToByteString(mapMetricId) + stellarprotocols.UInt32ToByteString(tenantId) + stellarprotocols.UInt16ToByteString(noOfFormula) + stellarprotocols.UInt16ToByteString(noOfFormula) + string(decodedStrFetureUsed)
	if len(memo) > 28 {
		return "", errors.New("Metric binding memo sholud be a 28 bytes")
	}
	logrus.Info("Builded memo : ", memo)
	return memo, nil
}
