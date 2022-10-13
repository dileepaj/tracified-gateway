package metricBinding

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

func (metric *MetricBinding) BuildFieldName(fieldName string) (txnbuild.ManageData, error) {
	fieldNameKey := ""
	fieldNameValue := ""

	if len(fieldName) > 128 {
		logrus.Error("Field name is greater than 128 character limit")
		return txnbuild.ManageData{}, errors.New("Field name is greater than 128 character limit")
	} else {
		if len(fieldName) > 64 {
			fieldNameKey = fieldName[0:64]
			fieldNameValue = fieldName[64:]
		} else if len(fieldName) < 64 || len(fieldName) == 64 {
			fieldNameKey = fieldName
			fieldNameValue = fmt.Sprintf("%s", strings.Repeat("0", 64))
		}
	}

	//append 0s if needed
	if len(fieldNameKey) < 64 {
		fieldNameKey = fieldNameKey + "/"
		if len(fieldNameKey) < 64 {
			fieldNameKey = fmt.Sprintf("%s%s", fieldNameKey, strings.Repeat("0", 64-len(fieldNameKey)))
		}
	}
	if len(fieldNameValue) < 64 {
		fieldNameValue = fieldNameValue + "/"
		if len(fieldNameValue) < 64 {
			fieldNameValue = fmt.Sprintf("%s%s", fieldNameValue, strings.Repeat("0", 64-len(fieldNameValue)))
		}
	}

	logrus.Info("Field name key : ", fieldNameKey)
	logrus.Info("Field name value : ", fieldNameValue)

	//check the lengths
	if len(fieldNameKey) > 64 || len(fieldNameValue) > 64 {
		logrus.Error("Key string length : ", len(fieldNameKey))
		logrus.Error("Value string length : ", len(fieldNameValue))
		return txnbuild.ManageData{}, errors.New("Length issue on key or value fields on the field name building")
	}

	fieldNameBuilder := txnbuild.ManageData{
		Name:  fieldNameKey,
		Value: []byte(fieldNameValue),
	}

	return fieldNameBuilder, nil
}
