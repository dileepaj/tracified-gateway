package metricBinding

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

func (metric *MetricBinding) BuildFieldKey(fieldKeyName string) (txnbuild.ManageData, error) {
	fieldKey := ""
	fieldValue := ""

	if len(fieldKeyName) > 128 {
		logrus.Error("Field key name is greater than 128 character limit")
		return txnbuild.ManageData{}, errors.New("Field key name is greater than 128 character limit")
	} else {
		if len(fieldKeyName) > 64 {
			fieldKey = fieldKeyName[0:64]
			fieldValue = fieldKeyName[64:]
		} else if len(fieldKeyName) < 64 || len(fieldKeyName) == 64 {
			fieldKey = fieldKeyName
			fieldValue = fmt.Sprintf("%s", strings.Repeat("0", 64))
		}
	}

	//append 0s if needed
	if len(fieldKey) < 64 {
		fieldKey = fieldKey + "/"
		if len(fieldKey) < 64 {
			fieldKey = fmt.Sprintf("%s%s", fieldKey, strings.Repeat("0", 64-len(fieldKey)))
		}
	}
	if len(fieldValue) < 64 {
		fieldValue = fieldValue + "/"
		if len(fieldValue) < 64 {
			fieldValue = fmt.Sprintf("%s%s", fieldValue, strings.Repeat("0", 64-len(fieldValue)))
		}
	}

	logrus.Info("Field key : ", fieldKey)
	logrus.Info("FIeld key value : ", fieldValue)

	//check the lengths
	if len(fieldKey) > 64 || len(fieldValue) > 64 {
		logrus.Error("Key string length : ", len(fieldKey))
		logrus.Error("Value string length : ", len(fieldValue))
		return txnbuild.ManageData{}, errors.New("Length issue on key or value fields on the field key name building")
	}

	fieldKeyNameBuilder := txnbuild.ManageData{
		Name:  fieldKey,
		Value: []byte(fieldValue),
	}

	return fieldKeyNameBuilder, nil
}
