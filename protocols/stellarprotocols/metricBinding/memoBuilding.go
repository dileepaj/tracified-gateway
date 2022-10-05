package metricBinding

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
)

func BuildMemo(mapMetricId uint64, metricName string, tenantId uint32, noOfFormula int32) (string, error) {
	if len(metricName) > 12 {
		return "", errors.New("metricName should be less than 12 chacter")
	}
	memo := stellarprotocols.UInt64ToByteString(mapMetricId) + metricName + stellarprotocols.UInt32ToByteString(tenantId) + stellarprotocols.UInt32ToByteString(uint32(noOfFormula))
	return memo, nil
}
