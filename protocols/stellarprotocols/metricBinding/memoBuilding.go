package metricbinding

import (
	"errors"
	"fmt"

	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
)

func BuildMemo(mapMetricId uint64, metricName string, tenantId uint32, noOfFormula int32) (string, error) {
	if len(metricName) > 12 {
		return "", errors.New("metricName should be less than 12 chacter")
	}
	strNoOfFormula := fmt.Sprintf("%04d", noOfFormula)
	if len(strNoOfFormula) > 4 {
		return "", errors.New("numer of formula count should be less than 4 chacter")
	}
	memo := stellarprotocols.UInt64ToByteString(mapMetricId) + metricName + stellarprotocols.UInt32ToByteString(tenantId) + strNoOfFormula
	return memo, nil
}
