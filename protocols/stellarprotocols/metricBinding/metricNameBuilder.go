package metricBinding

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

func (metric *MetricBinding) BuildMetricNameManageData(name string, typename string) (txnbuild.ManageData, error) {
	// base64 encode thetypename+ to overcome stellar UTE-8 issues
	metricName := base64.StdEncoding.EncodeToString([]byte(name))
	nameKey := ""
	nameValue := ""

	if len(metricName) > 128 {
		logrus.Error(typename + " is greater than 128 character limit")
		return txnbuild.ManageData{}, errors.New(typename + "is greater than 128 character limit")
	} else {
		// check if the key is greater than 64 characters
		if len(metricName) > 64 {
			nameKey = metricName[0:64]
			nameValue = metricName[64:]
		} else if len(metricName) < 64 || len(metricName) == 64 {
			nameKey = metricName
			nameValue = fmt.Sprintf("%s", strings.Repeat("0", 64))
		}
	}

	// check the lengths and append 0s if needed
	if len(nameKey) < 64 {
		nameKey = nameKey + "/"
		if len(nameKey) < 64 {
			nameKey = fmt.Sprintf("%s%s", nameKey, strings.Repeat("0", 64-len(nameKey)))
		}
	}
	if len(nameValue) < 64 {
		nameValue = nameValue + "/"
		if len(nameValue) < 64 {
			nameValue = fmt.Sprintf("%s%s", nameValue, strings.Repeat("0", 64-len(nameValue)))
		}
	}

	logrus.Info(typename+" key : ", nameKey)
	logrus.Info(typename+"value : ", nameValue)

	// check the lengths
	if len(nameKey) > 64 || len(nameValue) > 64 {
		logrus.Error(typename+" Key string length : ", len(nameKey))
		logrus.Error(typename+" Value string length : ", len(nameValue))
		return txnbuild.ManageData{}, errors.New("Length issue on key or value fields on " + typename)
	}

	metricNameBuilder := txnbuild.ManageData{
		Name:  nameKey,
		Value: []byte(nameValue),
	}

	return metricNameBuilder, nil
}
