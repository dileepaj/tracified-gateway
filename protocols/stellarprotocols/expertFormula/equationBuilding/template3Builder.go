package equationbuilding

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

/**
* This is the template for the type 3 execution template manage data
* Key -> 64 Bytes  -> "Type 3 Execution Template" padded with 0s
* Value -> 1 Byte   -> Type of the template
			4 Bytes  -> Special command id
* 			59 Bytes -> Future use
*/

func Type3TemplateBuilder(specialCommandId uint32) (txnbuild.ManageData, error) {

	// key field
	keyString := "Type 3 Execution Template"
	if len(keyString) > 64 {
		logrus.Error("Length of the key is greater than 64(template3Builder.go)")
		return txnbuild.ManageData{}, errors.New("length of the key is greater than 64")
	} else if len(keyString) < 64 {
		keyString = keyString + "/"
		keyString = keyString + strings.Repeat("0", 64-len(keyString))
	}

	// value field
	// convert typeOfTemplate from uint8 to string
	typeOfTemplate, errInConvertion := stellarprotocols.Int8ToByteString(uint8(2))
	if errInConvertion != nil {
		logrus.Info("Error when converting type of template from int8 to String(template3Builder.go) ", errInConvertion)
		return txnbuild.ManageData{}, errors.New("error when converting type of template to String " + errInConvertion.Error())
	}

	// futureUse in the value
	decodedStrFutureUse, err := hex.DecodeString(fmt.Sprintf("%0118d", 0))
	if err != nil {
		return txnbuild.ManageData{}, errors.New("error when decoding future use in the value " + err.Error())
	}
	futureUse := string(decodedStrFutureUse)

	valueString := typeOfTemplate + stellarprotocols.UInt32ToByteString(specialCommandId) + futureUse

	logrus.Info("Template3Builder Key String ", keyString)
	logrus.Info("Template3Builder Value String ", valueString)

	// check if the length of key and value is 64
	if len(valueString) != 64 || len(keyString) != 64 {
		logrus.Error("Length of the key: ", len(keyString), " and value: ", len(valueString))
		logrus.Error("Length of the key or value is not 64(template3Builder.go)")
		return txnbuild.ManageData{}, errors.New("length issue on key or value fields in the template type 3 building")
	}

	// build the manage data
	template3Builder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	return template3Builder, nil
}