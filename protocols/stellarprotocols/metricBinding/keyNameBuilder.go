package metricBinding

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

/*
des - build the key name manage data
	128 byttes
*/

func (metric *MetricBinding) BuildKeyName(keyName string) (txnbuild.ManageData, error) {
	keyNameKey := ""
	keyNameValue := ""

	if len(keyName) > 128 {
		logrus.Error("key name is greater than 128 character limit")
		return txnbuild.ManageData{}, errors.New("Key name is greater than 128 character limit")
	} else {
		// check and divide the key name into key and value
		if len(keyName) > 64 {
			keyNameKey = keyName[0:64]
			keyNameValue = keyName[64:]
		} else if len(keyName) < 64 || len(keyName) == 64 {
			keyNameKey = keyName
			keyNameValue = fmt.Sprintf("%s", strings.Repeat("0", 64))
		}
	}

	// check the lengths and append 0s if needed
	if len(keyNameKey) < 64 {
		keyNameKey = keyNameKey + "/"
		if len(keyNameKey) < 64 {
			keyNameKey = fmt.Sprintf("%s%s", keyNameKey, strings.Repeat("0", 64-len(keyNameKey)))
		}
	}
	if len(keyNameValue) < 64 {
		keyNameValue = keyNameValue + "/"
		if len(keyNameValue) < 64 {
			keyNameValue = fmt.Sprintf("%s%s", keyNameValue, strings.Repeat("0", 64-len(keyNameValue)))
		}
	}

	logrus.Info("Key name key : ", keyNameKey)
	logrus.Info("Key name value : ", keyNameValue)

	//check the lengths
	if len(keyNameKey) > 64 || len(keyNameValue) > 64 {
		logrus.Error("Key string length : ", len(keyNameKey))
		logrus.Error("Value string length : ", len(keyNameValue))
		return txnbuild.ManageData{}, errors.New("Length issue on key or value fields on the key name building")
	}

	keyNameBuilder := txnbuild.ManageData{
		Name:  keyNameKey,
		Value: []byte(keyNameValue),
	}

	return keyNameBuilder, nil
}
