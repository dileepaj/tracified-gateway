package metricBinding

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

func (metric *MetricBinding) BuildPrimaryKeyColumn(pkColName string) (txnbuild.ManageData, error) {
	pkColKey := ""
	pkColValue := ""

	if len(pkColName) > 128 {
		logrus.Error("Primary key column name is greater than 128 character limit")
		return txnbuild.ManageData{}, errors.New("Primary key column name is greater than 128 character limit")
	} else {
		if len(pkColName) > 64 {
			pkColKey = pkColName[0:64]
			pkColValue = pkColName[64:]
		} else if len(pkColName) < 64 || len(pkColName) == 64 {
			pkColKey = pkColName
			pkColValue = fmt.Sprintf("%s", strings.Repeat("0", 64))
		}
	}

	//append 0s if needed
	if len(pkColKey) < 64 {
		pkColKey = pkColKey + "/"
		if len(pkColKey) < 64 {
			pkColKey = fmt.Sprintf("%s%s", pkColKey, strings.Repeat("0", 64-len(pkColKey)))
		}
	}
	if len(pkColValue) < 64 {
		pkColValue = pkColValue + "/"
		if len(pkColValue) < 64 {
			pkColValue = fmt.Sprintf("%s%s", pkColValue, strings.Repeat("0", 64-len(pkColValue)))
		}
	}

	logrus.Info("Primary key column name key : ", pkColKey)
	logrus.Info("Primary key column name value : ", pkColValue)

	//check the lengths
	if len(pkColKey) > 64 || len(pkColValue) > 64 {
		logrus.Error("Key string length : ", len(pkColKey))
		logrus.Error("Value string length : ", len(pkColValue))
		return txnbuild.ManageData{}, errors.New("Length issue on key or value fields on the Primary key column name building")
	}

	pkColumnNameBuilder := txnbuild.ManageData{
		Name:  pkColKey,
		Value: []byte(pkColValue),
	}

	return pkColumnNameBuilder, nil
}
