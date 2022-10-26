package equationbuilding

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

/**
 * This is the template for the type 1 execution template manage data
 * Returns an array of manage data operations

 * key -> 64 bytes -> "Type 1 Execution Template" padded with 0s
 * value -> 1 Byte -> Type
			8 Bytes -> Start variable id
			4 Bytes -> Number of commands
			51 Bytes -> Future use
*/

func Type1TemplateBuilder(formulaId string, executionTemplate model.ExecutionTemplate) ([]txnbuild.Operation, error) {

	var startVariableID uint64											// mapped id in the db for Start variable
	numberOfCommands := uint32(len(executionTemplate.Lst_Commands))		// number of commands in the execution template
	var manageDataOpArray []txnbuild.Operation							// Array of manage data operations

	// get the mapped id for start variable from db using formula id and key
	object := dao.Connection{}
	valueMapDetails, errValueMapDetails := object.GetValueMapDetails(formulaId, "$"+executionTemplate.S_StartVarName).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errValueMapDetails != nil {
		logrus.Error("Error in getting value map from db(template1Builder.go) ", errValueMapDetails)
		return manageDataOpArray, errors.New("error in getting value map from db " + errValueMapDetails.Error())
	}
	startVariableID = valueMapDetails.(model.ValueIDMap).MapID

	// key field 
	keyString := "Type 1 Execution Template"
	if len(keyString) > 64 {
		logrus.Error("Length of the key is greater than 64(template1Builder.go)")
		return manageDataOpArray, errors.New("length of the key is greater than 64 ")
	} else if len(keyString) < 64 {
		keyString = keyString + "/"
		keyString = keyString + strings.Repeat("0", 64-len(keyString))
	}

	// value field
	// convert typeOfTemplate from int8 to string
	typeOfTemplate, errInConvertion := stellarprotocols.Int8ToByteString(uint8(1))
	if errInConvertion != nil {
		logrus.Info("Error when converting template type from int8 to string(template1Builder.go) ", errInConvertion)
		return manageDataOpArray, errors.New("error when converting template type to string " + errInConvertion.Error())
	}

	// futureUse
	decodedStrFutureUse, err := hex.DecodeString(fmt.Sprintf("%0102d", 0))
	if err != nil {
		return manageDataOpArray, errors.New("error in decoding future use string " + err.Error())
	}
	futureUse := string(decodedStrFutureUse)

	valueString := typeOfTemplate + stellarprotocols.UInt64ToByteString(startVariableID) + stellarprotocols.UInt32ToByteString(numberOfCommands) + futureUse

	logrus.Info("Template1Builder Key String ", keyString)
	logrus.Info("Template1Builder Value String ", valueString)

	// check if the length of key and value is 64
	if len(valueString) != 64 || len(keyString) != 64 {
		logrus.Error("Length of the key: ", len(keyString), " and value: ", len(valueString))
		logrus.Error("Length of the key or value is not 64(template1Builder.go)")
		return manageDataOpArray, errors.New("length issue on key or value fields in the template type 1 building")
	}

	// build the manage data for this template
	template1Builder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}
	// append the manage data for template 1 to the array
	manageDataOpArray = append(manageDataOpArray, &template1Builder)

	// loop through the lst_commands array and build the manage data for each command
	// build manage data related to a command will be returned as an array
	for _, command := range executionTemplate.Lst_Commands {
		manageDataOp, err := CommandBuilder(formulaId, command)	
		if err != nil {
			logrus.Error("Error in building the command(template1Builder.go) ", err)
			return manageDataOpArray, errors.New("error in building the command " + err.Error())
		}
		// append the manage data for command to the array
		manageDataOpArray = append(manageDataOpArray, manageDataOp...)
	}
	return manageDataOpArray, nil
}