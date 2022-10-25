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
 * This is the template for the manage data for commands in type 1 execution template
 * Key -> 64 Bytes -> "Command" padded with 0s
 * Value -> 4 Byte -> Command Type
 * 			1 Bytes -> Has argument
 * 			59 Bytes -> Future use
 */

func CommandBuilder(commandtype uint32, hasArgument bool) (txnbuild.ManageData, error) {

	// key field
	keyString := "Command"
	if len(keyString) > 64 {
		logrus.Error("Length of the key is greater than 64")
		return txnbuild.ManageData{}, errors.New("length of the key is greater than 64")
	} else if len(keyString) < 64 {
		keyString = keyString + "/"
		keyString = keyString + strings.Repeat("0", 64-len(keyString))
	}

	// value field
	// convert hasArgument from bool to string
	hasArgumentInt := 0
	if hasArgument {
		hasArgumentInt = 1
	}
	hasArgumentString, errInConvertion := stellarprotocols.Int8ToByteString(uint8(hasArgumentInt))
	if errInConvertion != nil {
		logrus.Info("Error when converting has argument ", errInConvertion)
		return txnbuild.ManageData{}, errors.New("Error when converting has argument " + errInConvertion.Error())
	}

	// futureUse
	decodedStrFutureUsed, err := hex.DecodeString(fmt.Sprintf("%0118d", 0))
	if err != nil {
		return txnbuild.ManageData{}, err
	}
	futureUse := string(decodedStrFutureUsed)

	valueString := stellarprotocols.UInt32ToByteString(commandtype) + hasArgumentString + futureUse

	logrus.Info("Key String ", keyString)
	logrus.Info("Value String ", valueString)

	// check if the length of key and value is 64
	if len(valueString) != 64 || len(keyString) != 64 {
		logrus.Error("Length of the key: ", len(keyString), " and value: ", len(valueString))
		logrus.Error("Length of the key or value is not 64")
		return txnbuild.ManageData{}, errors.New("length of the key or value is not 64")
	}

	// build the manage data
	commandBuilder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	return commandBuilder, nil

}