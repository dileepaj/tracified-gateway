package metricBinding

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

func (metric *MetricBinding) BuildMetaData(metadataName string) (txnbuild.ManageData, error) {
	metadataKey := ""
	metadataValue := ""

	if len(metadataName) > 128 {
		logrus.Error("Meta data name is greater than 128 character limit")
		return txnbuild.ManageData{}, errors.New("Meta data name is greater than 128 character limit")
	} else {
		if len(metadataName) > 64 {
			metadataKey = metadataName[0:64]
			metadataValue = metadataName[64:]
		} else if len(metadataName) < 64 || len(metadataName) == 64 {
			metadataKey = metadataName
			metadataValue = fmt.Sprintf("%s", strings.Repeat("0", 64))
		}
	}

	//append 0s if needed
	if len(metadataKey) < 64 {
		metadataKey = metadataKey + "/"
		if len(metadataKey) < 64 {
			metadataKey = fmt.Sprintf("%s%s", metadataKey, strings.Repeat("0", 64-len(metadataKey)))
		}
	}
	if len(metadataValue) < 64 {
		metadataValue = metadataValue + "/"
		if len(metadataValue) < 64 {
			metadataValue = fmt.Sprintf("%s%s", metadataValue, strings.Repeat("0", 64-len(metadataValue)))
		}
	}

	logrus.Info("Metadata name key : ", metadataKey)
	logrus.Info("Metadata name value : ", metadataValue)

	//check the lengths
	if len(metadataKey) > 64 || len(metadataValue) > 64 {
		logrus.Error("Key string length : ", len(metadataKey))
		logrus.Error("Value string length : ", len(metadataValue))
		return txnbuild.ManageData{}, errors.New("Length issue on key or value fields on the metadata name building")
	}

	metaDataBuilder := txnbuild.ManageData{
		Name:  metadataKey,
		Value: []byte(metadataValue),
	}

	return metaDataBuilder, nil
}
