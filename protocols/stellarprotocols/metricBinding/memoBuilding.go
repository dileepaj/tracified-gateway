package metricbinding

import (
	"errors"
	"fmt"

	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
)

func BuildMemo(mapMetricId int64, metricName string, tenantId, noOfFormula int32) (string, error) {
	if len(metricName) > 12 {
		return "", errors.New("metricName should be less than 12 chacter")
	}
	strNoOfFormula := fmt.Sprintf("%04d", noOfFormula)
	if len(strNoOfFormula) > 4 {
		return "", errors.New("numer of formula count should be less than 4 chacter")
	}
	strTenatID, err := stellarprotocols.TenantIDToBinary(int64(tenantId))
	if err != nil {
		return "", errors.New("BuildMemo issue (faormula ID convert to type) " + err.Error())
	}
	memo := stellarprotocols.UInt64ToByteString(mapMetricId) + metricName + stellarprotocols.ConvertingBinaryToByteString(strTenatID) + strNoOfFormula
	return memo, nil
}
