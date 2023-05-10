package expertformula

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dileepaj/tracified-gateway/expertformulas/stellarprotocols"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

/*
des - build the reference manage data for referred constant
return txnbuild.manage data

	reference url - 127 bytes (key + value)
	length - 1 byte (in value)
*/

func (expertFormula ExpertFormula) BuildReference(refUrl string) (txnbuild.ManageData, error) {
	nameKey := ""
	nameValue := ""
	if len(refUrl) == 0 {
		refUrl = "URL Not Provided"
	}

	actualLength, errInLength := stellarprotocols.Int8ToByteString(uint8(len(refUrl)))
	if errInLength != nil {
		logrus.Info("Error when converting length(referenceBuilder.go) ", errInLength)
		return txnbuild.ManageData{}, errors.New("error when converting reference url length to String " + errInLength.Error())
	}

	if len(refUrl) > 127 {
		logrus.Error(refUrl + " is greater than 127 character limit(referenceBuilder.go)")
		return txnbuild.ManageData{}, errors.New(refUrl + "is greater than 127 character limit")
	} else {
		// check if the key is greater than 64 characters
		if len(refUrl) > 64 {
			nameKey = refUrl[0:64]
			nameValue = refUrl[64:]
		} else if len(refUrl) < 64 || len(refUrl) == 64 {
			nameKey = refUrl
			nameValue = strings.Repeat("0", 63)
		}
	}

	// check the lengths and append 0s if needed
	if len(nameKey) < 64 {
		if len(nameKey) < 64 {
			nameKey = fmt.Sprintf("%s%s", nameKey, strings.Repeat("0", 64-len(nameKey)))
		}
	}
	if len(nameValue) < 63 {
		if len(nameValue) < 63 {
			nameValue = fmt.Sprintf("%s%s", nameValue, strings.Repeat("0", 63-len(nameValue)))
		}
	}

	nameValue = nameValue + actualLength

	logrus.Info("Referece URL key : ", nameKey)
	logrus.Info("Reference URL value : ", nameValue)

	// check the lengths
	if len(nameKey) > 64 || len(nameValue) > 64 {
		logrus.Error("Referece URL key : ", len(nameKey))
		logrus.Error("Reference URL value : ", len(nameValue))
		return txnbuild.ManageData{}, errors.New("length issue on key or value fields on " + refUrl)
	}

	urlBuilder := txnbuild.ManageData{
		Name:  nameKey,
		Value: []byte(nameValue),
	}

	return urlBuilder, nil
}
