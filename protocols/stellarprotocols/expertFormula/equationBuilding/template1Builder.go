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
 * This is the template for the type 1 execution template manage data
 * key -> 64 bytes -> "Type 1 Execution Template" padded with 0s
 * value -> 1 Byte -> Type
			8 Bytes -> Start variable id
			4 Bytes -> Number of commands
			51 Bytes -> Future use
*/

func Type1TemplateBuilder(startvariableId uint64, numberOfCommands uint32) (txnbuild.ManageData, error) {

	// key field 
	keyString := "Type 1 Execution Template"
	if len(keyString) > 64 {
		logrus.Error("Length of the key is greater than 64")
		return txnbuild.ManageData{}, errors.New("length of the key is greater than 64")
	} else if len(keyString) < 64 {
		keyString = keyString + "/"
		keyString = keyString + strings.Repeat("0", 64-len(keyString))
	}

	// value field
	// convert typeOfTemplate from uint8 to string
	typeOfTemplate, errInConvertion := stellarprotocols.Int8ToByteString(uint8(1))
	if errInConvertion != nil {
		logrus.Info("Error when converting type of template ", errInConvertion)
		return txnbuild.ManageData{}, errors.New("Error when converting type of template " + errInConvertion.Error())
	}

	// futureUse
	decodedStrFutureUsed, err := hex.DecodeString(fmt.Sprintf("%0102d", 0))
	if err != nil {
		return txnbuild.ManageData{}, err
	}
	futureUse := string(decodedStrFutureUsed)

	valueString := typeOfTemplate + stellarprotocols.UInt64ToByteString(startvariableId) + stellarprotocols.UInt32ToByteString(numberOfCommands) + futureUse

	logrus.Info("Key String ", keyString)
	logrus.Info("Value String ", valueString)

	// check if the length of key and value is 64
	if len(valueString) != 64 || len(keyString) != 64 {
		logrus.Error("Length of the key: ", len(keyString), " and value: ", len(valueString), " is not 64")
		logrus.Error("Length of the key or value is not 64")
		return txnbuild.ManageData{}, errors.New("length of the key or value is not 64")
	}

	// build the manage data
	template1Builder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	return template1Builder, nil
}