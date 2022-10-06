package metricBinding

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	"github.com/sirupsen/logrus"
)

type MetricBinding struct{}

func (metric *MetricBinding) BuildMemo(mapMetricId uint64, metricName string, tenantId uint32, noOfFormula int32) (string, error) {
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
	memo := stellarprotocols.UInt64ToByteString(mapMetricId) + rebuildMetricName + stellarprotocols.UInt32ToByteString(tenantId) + stellarprotocols.UInt32ToByteString(uint32(noOfFormula))
	if len(memo) > 28 {
		return "", errors.New("Metric binding memo sholud be a 28 bytes")
	}
	return memo, nil
}
