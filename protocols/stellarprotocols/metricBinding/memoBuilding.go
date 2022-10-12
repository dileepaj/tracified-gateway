package metricBinding

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	"github.com/sirupsen/logrus"
)

/*
des- build the memo according to the protocol
return the txnbuild.ManageData object

Fields
	1. Metric ID      - 8 bytes (uint64)  - mapped id stored in the DB for the metric
	2. Metric name    - 12 bytes (string) - metric name
	3. Tenant ID      - 4 bytes (uint32)  - mapped id stored in the DB for the tenant
	4. No of formulas - 2 bytes (uint16)  - no of formulas in the defined metric
	5. Future Use     - 2 bytes
*/

type MetricBinding struct{}

func (metric *MetricBinding) BuildMemo(mapMetricId uint64, metricName string, tenantId uint32, noOfFormula uint16) (string, error) {
	rebuildMetricName := ""
	if len(metricName) > 12 {
		logrus.Error("metric name is greater than 12 character limit")
		return "", errors.New("Metric name is greater than 12 character limit")
	} else {
		if len(metricName) == 12 {
			rebuildMetricName = metricName
		} else if len(metricName) < 12 {
			rebuildMetricName = metricName + "/"
		}
	}
	if len(rebuildMetricName) < 12 {
		remain := 12 - len(rebuildMetricName)
		setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		rebuildMetricName = rebuildMetricName + setReaminder
	}
	decodedStrFetureUsed, err := hex.DecodeString(fmt.Sprintf("%04d", 0))
	if err != nil {
		return "", errors.New("Feture used byte building issue in memo building")
	}
	memo := stellarprotocols.UInt64ToByteString(mapMetricId) + rebuildMetricName + stellarprotocols.UInt32ToByteString(tenantId) + stellarprotocols.UInt16ToByteString(noOfFormula) + string(decodedStrFetureUsed)
	if len(memo) > 28 {
		return "", errors.New("Metric binding memo sholud be a 28 bytes")
	}
	return memo, nil
}
