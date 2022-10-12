package metricBinding

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

func (metric *MetricBinding) BuildStageName(stageName string) (txnbuild.ManageData, error) {
	stageNameKey := ""
	stageNameValue := ""

	if len(stageName) > 128 {
		logrus.Error("stage name is greater than 128 character limit")
		return txnbuild.ManageData{}, errors.New("Stage name is greater than 128 character limit")
	} else {
		// check and divide the stage name into key and value
		if len(stageName) > 64 {
			stageNameKey = stageName[0:64]
			stageNameValue = stageName[64:]
		} else if len(stageName) < 64 || len(stageName) == 64 {
			stageNameKey = stageName
			stageNameValue = fmt.Sprintf("%s", strings.Repeat("0", 64))
		}
	}

	// check the lengths and append 0s if needed
	if len(stageNameKey) < 64 {
		stageNameKey = stageNameKey + "/"
		if len(stageNameKey) < 64 {
			stageNameKey = fmt.Sprintf("%s%s", stageNameKey, strings.Repeat("0", 64-len(stageNameKey)))
		}
	}
	if len(stageNameValue) < 64 {
		stageNameValue = stageNameValue + "/"
		if len(stageNameValue) < 64 {
			stageNameValue = fmt.Sprintf("%s%s", stageNameValue, strings.Repeat("0", 64-len(stageNameValue)))
		}
	}
	
	logrus.Info("Stage name key : ", stageNameKey)
	logrus.Info("Stage name value : ", stageNameValue)

	//check the lengths
	if len(stageNameKey) > 64 || len(stageNameValue) > 64 {
		logrus.Error("Key string length : ", len(stageNameKey))
		logrus.Error("Value string length : ", len(stageNameValue))
		return txnbuild.ManageData{}, errors.New("Length issue on key or value fields on the stage name building")
	}

	stageNameBuilder := txnbuild.ManageData{
		Name:  stageNameKey,
		Value: []byte(stageNameValue),
	}

	return stageNameBuilder, nil
}