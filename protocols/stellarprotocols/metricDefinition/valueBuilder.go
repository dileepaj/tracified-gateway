package metricdefinition

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/stellarprotocols"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

//Used to define the general metadat of a single variable
//Key - “VALUE METADATA”
/* VALUE
*	value id - 8 bytes (unsigned long integer)
*	variable name - 20 bytes (string)
* 	binding type - 1 byte (byte)
*	future use - 35 bytes
 */
func generalValueDefinitionBuilder(element model.GeneralValueDefBuildRequest) (txnbuild.ManageData, error) {
	var valueID int64
	valueTypeString := ""
	variableNameString := ""
	futureUse := ""
	//get value id from the map for the type, key and formula ID
	object := dao.Connection{}
	variableDefMap, errWhenRetrievingVariableInfo := object.GetVariableDetails(element.FormulaID, element.Key, element.VariableName).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errWhenRetrievingVariableInfo != nil {
		logrus.Error("Unable to connect gateway datastore ", errWhenRetrievingVariableInfo)
		return txnbuild.ManageData{}, errors.New("Unable to connect gateway datastore " + errWhenRetrievingVariableInfo.Error())
	}
	if variableDefMap == nil {
		//requested varible does not exists in the DB
		logrus.Error("Requested variable " + element.VariableName + " does not exists in the gateway DB")
		return txnbuild.ManageData{}, errors.New("Requested variable " + element.VariableName + " does not exists in the gateway DB")
	} else {
		valueMapData := variableDefMap.(model.ValueIDMap)
		valueID = valueMapData.MapID
	}
	//convert value id to binart
	strValueId, errWhenConvertingValueIDToBinary := stellarprotocols.IDToBinary(valueID)
	if errWhenConvertingValueIDToBinary != nil {
		return txnbuild.ManageData{}, errors.New("Error when converting value ID to the binary format " + errWhenConvertingValueIDToBinary.Error())
	}

	//check the variable name length for 20 bytes
	if len(element.VariableName) > 20 {
		logrus.Error("Variable name us greater than 20 character limit")
		return txnbuild.ManageData{}, errors.New("Variable name us greater than 20 character limit")
	} else {
		if len(element.VariableName) == 20 {
			variableNameString = element.VariableName
		} else if len(element.VariableName) < 20 {
			variableNameString = element.VariableName + "/"
		}
	}
	//check whether the new variable name string saturated the 20 byte limit, if not add 0s to rest
	if len(variableNameString) < 20 {
		remain := 20 - len(variableNameString)
		setRemainder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		variableNameString = variableNameString + setRemainder
	}

	//check if the binding type is 0 or 1
	if element.BindingType != 0 || element.BindingType != 1 {
		logrus.Error("Invalid binding type, should be 1 or 0")
		return txnbuild.ManageData{}, errors.New("Invalid binding type, should be 1 or 0")
	} else {
		tempValueType, errInValueTypeConvert := stellarprotocols.StringToBinary(int64(element.BindingType))
		if errInValueTypeConvert != nil {
			logrus.Error("Error when converting value type ", errInValueTypeConvert)
			return txnbuild.ManageData{}, errors.New("Error when converting value type " + errInValueTypeConvert.Error())
		}

		valueTypeString = stellarprotocols.ConvertingBinaryToByteString(tempValueType)
	}

	//future use
	futureUse = fmt.Sprintf("%s", strings.Repeat("0", 35))

	//build key and value strings
	keyString := "VALUE METADATA/" + fmt.Sprintf("%s", strings.Repeat("0", (64-len("VALUE METADATA/"))))
	valueString := stellarprotocols.ConvertingBinaryToByteString(strValueId) + variableNameString + valueTypeString + futureUse

	//check the key value string length for 64 byte limit
	if len(keyString) > 61 {
		logrus.Error("Key string exceeding the given 64 byte limit in variable manage data builder. Length : ", len(keyString))
		return txnbuild.ManageData{}, errors.New("Key string exceeding the given 64 byte limit in variable manage data builder")
	}

	if len(valueString) > 64 {
		logrus.Error("Value string exceeding the given 64 byte limit in variable manage data builder. Length : ", len(keyString))
		return txnbuild.ManageData{}, errors.New("Value string exceeding the given 64 byte limit in variable manage data builder")
	}

	generalValueDefinitionBuilder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	return generalValueDefinitionBuilder, nil
}
