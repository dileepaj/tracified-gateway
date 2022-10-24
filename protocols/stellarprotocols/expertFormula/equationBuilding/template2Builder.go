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
* This is the template for the type 2 execution template manage data
* key -> 64 bytes   -> constant value
* value ->  1 Byte   -> Type of the template
			1 Byte   -> IsVariable
* 			8 Bytes  -> Variable Id
* 			1 Bytes  -> Type of the constant
* 			53 Bytes -> Future use
*/

func Type2TemplateBuilder(formulaID string, executionTemplate model.ExecutionTemplate) (txnbuild.Operation, error) {
	var isVariable bool = false
	var variableID uint64
	var constantValue = ""
	constType := uint(executionTemplate.P_Entity.Ul_type)

	if executionTemplate.P_Entity.Ul_type != 0 {
		isVariable = true
	}

	if isVariable {
		object := dao.Connection{}
		valueMapDetails, errValueMapDetails := object.GetValueMapDetails(formulaID, "$"+executionTemplate.S_StartVarName).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if errValueMapDetails != nil {
			logrus.Error("Error in getting value map from key ", errValueMapDetails)
			return &txnbuild.ManageData{}, errors.New("Error in getting value map from key " + errValueMapDetails.Error())
		}
		variableID = valueMapDetails.(model.ValueIDMap).MapID
		constantValue = strings.Repeat("0", 64)

	} else {
		constantValue = fmt.Sprintf("%g", executionTemplate.P_Entity.Value.(float64))
		// check value is 64 character
		if len(constantValue) > 64 {
			logrus.Error("Value is greater than 8 character limit")
			return &txnbuild.ManageData{}, errors.New("value is greater than 64 character limit")
		} else {
			if len(constantValue) < 64 {
				// add 0s to the rest of the name
				remain := 64 - len(constantValue)
				setReaminder := strings.Repeat("0", remain)
				constantValue = setReaminder + constantValue
			}
		}
	}

	// value field
	// convert typeOfTemplate from uint8 to string
	typeOfTemplate, errInConvertion := stellarprotocols.Int8ToByteString(uint8(2))
	if errInConvertion != nil {
		logrus.Info("Error when converting type of template ", errInConvertion)
		return &txnbuild.ManageData{}, errors.New("Error when converting type of template " + errInConvertion.Error())
	}

	// convert isVariable from bool to string
	isVariableString := ""
	if isVariable {
		isVariable, errInConvertion := stellarprotocols.Int8ToByteString(uint8(1))
		if errInConvertion != nil {
			logrus.Info("Error when converting type of template ", errInConvertion)
		}

		isVariableString = isVariable
	} else {
		isVariable, errInConvertion := stellarprotocols.Int8ToByteString(uint8(0))
		if errInConvertion != nil {
			logrus.Info("Error when converting type of template ", errInConvertion)
		}

		isVariableString = isVariable
	}

	// convert variableID from uint8 to string
	constTypeString, errInConversion := stellarprotocols.Int8ToByteString(uint8(constType))
	if errInConversion != nil {
		logrus.Info("Error when converting type of template ", errInConversion)
		return &txnbuild.ManageData{}, errors.New("Error when converting type of template " + errInConversion.Error())
	}

	// for future use
	decodedStrFutureUsed, err := hex.DecodeString(fmt.Sprintf("%0106d", 0))
	if err != nil {
		return &txnbuild.ManageData{}, err
	}
	futureUse := string(decodedStrFutureUsed)

	keyString := constantValue
	valueString := typeOfTemplate + isVariableString + stellarprotocols.UInt64ToByteString(variableID) + constTypeString + futureUse

	logrus.Info("Key String ", keyString)
	logrus.Info("Value String ", valueString)

	// check if the key and value is 64 bytes
	if len(valueString) != 64 || len(keyString) != 64 {
		logrus.Error("Length of the key: ", len(keyString), " and value: ", len(valueString))
		logrus.Error("Length of the key or value is not 64")
		return &txnbuild.ManageData{}, errors.New("length of the key or value is not 64")
	}

	// build the manage data
	template2Builder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	return &template2Builder, nil

}
