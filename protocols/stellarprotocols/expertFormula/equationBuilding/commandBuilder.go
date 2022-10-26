package equationbuilding

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

/**
 * This is the template for the manage data for commands in type 1 execution template
 * Returns an array of manage data operations

 * Key -> 64 Bytes -> "Command" padded with 0s
 * Value -> 4 Byte -> Command Type
 * 			1 Bytes -> Has argument
 * 			59 Bytes -> Future use
 */

func CommandBuilder(formulaId string, command model.Command) ([]txnbuild.Operation, error) {
	hasArgument := 0								// command has argument 0 - false and 1 - true
	commandType := command.Ul_CommandType			// command type (Values given in FCL for Add, Subtract, Multiply, Divide, etc)
	var manageDataOpArray []txnbuild.Operation		// Array of manage data operations

	// Check if the command has an argument by checking the s_StartVariableName in the p_Arg
	if command.P_Arg.S_StartVarName != "" {
		hasArgument = 1
	}

	// key field
	keyString := "Command"
	if len(keyString) > 64 {
		logrus.Error("Length of the key is greater than 64(commandBuilder)")
		return manageDataOpArray, errors.New("length of the key is greater than 64(commandBuilder)")
	} else if len(keyString) < 64 {
		keyString = keyString + "/"
		keyString = keyString + strings.Repeat("0", 64-len(keyString))
	}

	// value field
	// convert hasArgument from int8 to string
	hasArgumentString, errInConvertion := stellarprotocols.Int8ToByteString(uint8(hasArgument))
	if errInConvertion != nil {
		logrus.Info("Error when converting hasArgument from int8 to String(commandBuilder) ", errInConvertion)
		return manageDataOpArray, errors.New("error when converting hasArgument from int8 to String(commandBuilder) " + errInConvertion.Error())
	}

	// futureUse
	decodedStrFutureUse, err := hex.DecodeString(fmt.Sprintf("%0118d", 0))
	if err != nil {
		return manageDataOpArray, errors.New("error when decoding future use(commandBuilder) " + err.Error())
	}
	futureUse := string(decodedStrFutureUse)

	valueString := stellarprotocols.UInt32ToByteString(commandType) + hasArgumentString + futureUse

	logrus.Info("CommandBuilder Key String ", keyString)
	logrus.Info("CommandBuilder Value String ", valueString)

	// check if the length of key and value is 64
	if len(valueString) != 64 || len(keyString) != 64 {
		logrus.Error("Length of the key: ", len(keyString), " and value: ", len(valueString))
		logrus.Error("Length of the key or value is not 64(commandBuilder)")
		return manageDataOpArray, errors.New("length of the key or value is not 64(commandBuilder)")
	}

	// build the manage data
	commandBuilder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	manageDataOpArray = append(manageDataOpArray, &commandBuilder)

	/* if the command has argument then the template type in the arugment should be identified
			if the lst_commands in P_Arg is not empty	-> call the template 1 builder
			if the lst_commands in P_Arg is empty		-> call the template 2 builder
	*/
	if hasArgument == 1 {
		// build manage data for argument
		if command.P_Arg.Lst_Commands != nil {
			manageDataOp, err := Type1TemplateBuilder(formulaId, command.P_Arg)
			if err != nil {
				logrus.Error("Error when building the manage data for the argument type 1 execution template(commandBuilder) ", err)
				return manageDataOpArray, errors.New("error when building the manage data for the argument type 1 execution template(commandBuilder) " + err.Error())
			}
			// append the manage data to the array
			manageDataOpArray = append(manageDataOpArray, manageDataOp...)
		} else {
			manageDataOp, err := Type2TemplateBuilder(formulaId, command.P_Arg)
			if err != nil {
				logrus.Error("Error when building the manage data for the argument type 2 execution template(commandBuilder) ", err)
				return manageDataOpArray, errors.New("error when building the manage data for the argument type 2 execution template(commandBuilder) " + err.Error())
			}
			// append the manage data to the array
			manageDataOpArray = append(manageDataOpArray, manageDataOp)
		}
	}

	return manageDataOpArray, nil

}