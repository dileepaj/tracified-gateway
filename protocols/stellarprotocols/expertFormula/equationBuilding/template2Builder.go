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
 * Returns one manage data operation

 * key   ->  6 bytes  -> "Type2/"
 			 58 bytes -> constant value
 * value ->  1 Byte   -> Type of the template
			 1 Byte   -> IsVariable if the query has a #value, then it is not a variable, else it is a variable(all the variables, semantic and referred constants considered as variables)
 * 			 8 Bytes  -> Variable Id
 * 			 1 Bytes  -> Type of the constant
 * 			 53 Bytes -> Future use
*/

func Type2TemplateBuilder(formulaID string, executionTemplate model.ExecutionTemplate) (txnbuild.Operation, error) {
	var isVariable = 0			//   0 -> not a variable, 1 -> variable
	var variableID uint64		// mapped id in the db for the variable if isvariable is 1
	var constantValue = ""		// value of the constant if isvariable is 0
	constType := uint(executionTemplate.P_Entity.Ul_type)	// type of the constant if the isvariable is 0

	// if the ul_Type in the p_Entity is 0, then it is a variable
	if executionTemplate.P_Entity.Ul_type == 0 {
		isVariable = 1
	}

	// if isvariable is 1, then get the mapped variable id from the db
	if isVariable == 1 {
		object := dao.Connection{}
		valueMapDetails, errValueMapDetails := object.GetValueMapDetails(formulaID, "$"+executionTemplate.S_StartVarName).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if errValueMapDetails != nil {
			logrus.Error("Error in getting value map from db(template2Builder) ", errValueMapDetails)
			return &txnbuild.ManageData{}, errors.New("error in getting value map from db(template2Builder) " + errValueMapDetails.Error())
		}
		variableID = valueMapDetails.(model.ValueIDMap).MapID
		constantValue = strings.Repeat("0", 58)

	} else {
		constantValue = fmt.Sprintf("%g", executionTemplate.P_Entity.Value.(float64))
		// check value is 58 character
		if len(constantValue) > 58 {
			logrus.Error("Constant value is greater than 58 character limit(temlate2Builder)")
			return &txnbuild.ManageData{}, errors.New("constant value is greater than 58 character limit(template2Builder)")
		} else {
			if len(constantValue) < 58 {
				// add 0s to the rest of the name
				remain := 58 - len(constantValue)
				setReaminder := strings.Repeat("0", remain)
				constantValue = setReaminder + constantValue
			}
		}
	}

	// value field
	// convert typeOfTemplate from uint8 to string
	typeOfTemplate, errInConvertion := stellarprotocols.Int8ToByteString(uint8(2))
	if errInConvertion != nil {
		logrus.Info("Error when converting type of template from int8 to String(template2Builder) ", errInConvertion)
		return &txnbuild.ManageData{}, errors.New("error when converting type of template from int8 to String(template2Builder) " + errInConvertion.Error())
	}

	// convert isVariable from int8 to string
	isVariableString := ""
	if isVariable == 1 {
		isVariable, errInConvertion := stellarprotocols.Int8ToByteString(uint8(isVariable))
		if errInConvertion != nil {
			logrus.Info("Error when converting isVariable from int8 to String(template2Builder) ", errInConvertion)
		}
		isVariableString = isVariable
	} 

	// convert variableID from uint8 to string
	constTypeString, errInConversion := stellarprotocols.Int8ToByteString(uint8(constType))
	if errInConversion != nil {
		logrus.Info("Error when converting constant type from int8 to String(template2Builder) ", errInConversion)
		return &txnbuild.ManageData{}, errors.New("error when converting constant type from int8 to String(template2Builder) " + errInConversion.Error())
	}

	// for future use in the value
	decodedStrFutureUse, err := hex.DecodeString(fmt.Sprintf("%0106d", 0))
	if err != nil {
		return &txnbuild.ManageData{}, errors.New("error in decoding future use (template2Builder) " + err.Error())
	}
	futureUse := string(decodedStrFutureUse)

	keyString := "Type2/" + constantValue
	valueString := typeOfTemplate + isVariableString + stellarprotocols.UInt64ToByteString(variableID) + constTypeString + futureUse

	logrus.Info("Template2Builder Key String ", keyString)
	logrus.Info("Template2Builder Value String ", valueString)

	// check if the key and value is 64 bytes
	if len(valueString) != 64 || len(keyString) != 64 {
		logrus.Error("Length of the key: ", len(keyString), " and value: ", len(valueString))
		logrus.Error("Length of the key or value is not 64(template2Builder)")
		return &txnbuild.ManageData{}, errors.New("length of the key or value is not 64(template2Builder)")
	}

	// build the manage data
	template2Builder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	return &template2Builder, nil

}
