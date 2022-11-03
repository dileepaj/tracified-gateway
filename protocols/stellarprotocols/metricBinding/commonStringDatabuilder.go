package metricBinding

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

/*
	des - build metric name manage data for metric binding endpoint
		Metric name - 127 bytes
		Length of the base64 metric name - 1 byte
*/

func (metric *MetricBinding) CommonStringBuilder(name string, typename string) (txnbuild.ManageData, error) {
	// base64 encode thetypename+ to overcome stellar UTE-8 issues
	metricName := base64.StdEncoding.EncodeToString([]byte(name))
	nameKey := ""
	nameValue := ""
	nameLength := 127

	if len(metricName) >= 127 {
		nameKey = metricName[0:64]
		nameValue = metricName[64:127]
	} else if len(metricName) > 64 && len(metricName) < 127 {
		nameKey = metricName[0:64]
		nameValue = metricName[64:]
		nameLength = len(metricName)
	} else if len(metricName) < 64 || len(metricName) == 64 {
		nameKey = metricName
		nameValue = strings.Repeat("0", 63)
		nameLength = len(metricName)
	}

	actualLength, errInLengthConvert := stellarprotocols.Int8ToByteString(uint8(nameLength))
	if errInLengthConvert != nil {
		logrus.Info("Error when converting length ", errInLengthConvert)
		return txnbuild.ManageData{}, errors.New("Error when converting length " + errInLengthConvert.Error())
	}
	// check if the key is greater than 64 characters

	// check the lengths and append 0s if needed
	if len(nameKey) < 64 {
		nameKey = nameKey + "/"
		if len(nameKey) < 64 {
			nameKey = fmt.Sprintf("%s%s", nameKey, strings.Repeat("0", 64-len(nameKey)))
		}
	}
	if len(nameValue) < 63 {
		nameValue = nameValue + "/"
		if len(nameValue) < 63 {
			nameValue = fmt.Sprintf("%s%s", nameValue, strings.Repeat("0", 63-len(nameValue)))
		}
	}

	nameValue = nameValue + actualLength

	logrus.Info(typename+" key : ", nameKey)
	logrus.Info(typename+"value : ", nameValue)

	// check the lengths
	if len(nameKey) > 64 || len(nameValue) > 64 {
		logrus.Error(typename+" Key string length : ", len(nameKey))
		logrus.Error(typename+" Value string length : ", len(nameValue))
		return txnbuild.ManageData{}, errors.New("length issue on key or value fields on " + typename)
	}

	metricNameBuilder := txnbuild.ManageData{
		Name:  nameKey,
		Value: []byte(nameValue),
	}

	return metricNameBuilder, nil
}
