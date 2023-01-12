package expertformula

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

/*
BuildVariableDefinitionManageData
des-Build the variable definition manage data
variable definition and byte used
	valueType - 1 byte defied by protocol -2 for sementic constant
	valueId  - 8 byte defied by protocol
	data type - 1 byte
	value name - 20 byte defied by protocol
	unit - 2 bytes
	precision - 1
	description - 40 bytes
	sementicConstantDataType - 1 byte defied by protocol -2 for flaot
	future use - 34 bytes

Manage data
	name 64 byte character - 	description - 64 byte defied by protocol + future use
	value 64 byte managedata - valueType + valueId +  + value name + fetureused + data type + unit + precision
*/
func (expertFormula ExpertFormula) BuildVariableDefinitionManageData(formulaID string, element model.FormulaItemRequest) (txnbuild.ManageData, model.ValueDefOutParmas, error) {
	VALUETYPE := 1
	DATATYPE := 2
	EMPTY := 0
	var valueId uint64
	var unitId uint16
	dataTypeString := ""
	valueTypeString := ""
	variableNameString := ""
	unitString := ""
	precisionString := ""
	descriptionString := ""
	errorRespObj := model.ValueDefOutParmas{
		ValueMapID: uint64(EMPTY),
		UnitMapID:  uint16(EMPTY),
	}
	// define the value type
	// this is a variable therefore the value type is 1
	tempValueType, errInValueTypeConvert := stellarprotocols.Int8ToByteString(uint8(VALUETYPE))
	if errInValueTypeConvert != nil {
		logrus.Info("Error when converting value type from int8 to string(variableBuilder) ", errInValueTypeConvert)
		return txnbuild.ManageData{}, errorRespObj, errors.New("error when converting value type from int8 to string ")
	}
	valueTypeString = tempValueType
	// DB validations for the variable id
	object := dao.Connection{}
	valueMap, errValueMap := object.GetValueMapID(element.ID, formulaID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errValueMap != nil {
		logrus.Info("Unable to connect to gateway datastore(variableBuilder) ", errValueMap)
		// return txnbuild.ManageData{}, errors.New("Unable to connect gateway datastore to get value map ID")
	}
	// check if the variable name for this formula is in the variale mapping
	if valueMap != nil {
		logrus.Info(element.Name + " is already recorded in the DB Map(variableBuilder) ")
		// add the value map part as the value id to the manage data key part string
		valueMapData := valueMap.(model.ValueIDMap)
		valueId = valueMapData.MapID
	} else {
		// if not add with incrementing id
		logrus.Info(element.Name + " is not recorded in the DB Map(variableBuilder) ")
		data, err := object.GetNextSequenceValue("VALUEID")
		if err != nil {
			logrus.Error("Retrieving value id from map was failed(variableBuilder) " + err.Error())
			return txnbuild.ManageData{}, errorRespObj, errors.New("retrieving value id from map was failed ")
		}
		valueIdMap := model.ValueIDMap{
			ValueId:   element.ID,
			ValueType: "VARIABLE",
			Key:       element.Key,
			FormulaID: formulaID,
			ValueName: element.Name,
			MapID:     data.SequenceValue,
		}
		err1 := object.InsertToValueIDMap(valueIdMap)
		if err1 != nil {
			logrus.Error("Inserting to Value map ID was failed(variableBuilder) " + err1.Error())
			return txnbuild.ManageData{}, errorRespObj, errors.New("inserting to Value map ID was failed ")
		}
		// add the data as the new value id to the manage data key part string
		valueId = data.SequenceValue
	}
	// check variable name is 20 character
	if len(element.Name) > 20 {
		logrus.Error("Variable name is greater than 20 character limit(variableBuilder) ")
		return txnbuild.ManageData{}, errorRespObj, errors.New("variable name is greater than 20 character limit ")
	} else {
		if len(element.Name) == 20 {
			variableNameString = element.Name
		} else if len(element.Name) < 20 {
			variableNameString = element.Name + "/"
		}
	}
	// finally check if the 20 bytes are filled if not append 0s at end
	if len(variableNameString) < 20 {
		remain := 20 - len(variableNameString)
		setRemainder := strings.Repeat("0", remain)
		variableNameString = variableNameString + setRemainder
	}
	// depending on the data type decide the integer to be assigned
	tempDataType, errInDataTypeConvert := stellarprotocols.Int8ToByteString(uint8(DATATYPE))
	if errInDataTypeConvert != nil {
		logrus.Info("Error when converting data type(variableBuilder) ", errInDataTypeConvert)
		return txnbuild.ManageData{}, errorRespObj, errors.New("error when converting data type " + errInDataTypeConvert.Error())
	}
	dataTypeString = tempDataType
	// depending on the unit type decide the integer to be assigned
	// convert unit type character -> byte -> bits
	unitMap, errInUnitIdMap := object.GetUnitMapID(element.MeasurementUnit).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInUnitIdMap != nil {
		logrus.Info("Unable to connect to gateway datastore(variableBuilder.go) ", errValueMap)
		// return txnbuild.ManageData{}, errors.New("Unable to connect gateway datastore to get value map ID")
	}
	// check if the unit is in the unit map
	if unitMap != nil {
		logrus.Info(element.MeasurementUnit + " is already recorded in the DB Map(variableBuilder.go) ")
		// add map id as the unit in the key string
		unitMapData := unitMap.(model.UnitIDMap)
		unitId = unitMapData.MapID
		unitString = stellarprotocols.UInt16ToByteString(uint16(unitMapData.MapID))
	} else {
		// if not add the incrementing id
		logrus.Info(element.MeasurementUnit + " is not recorded in the DB Map(variableBuilder.go) ")

		// get the current sequence for the units
		data, err := object.GetNextSequenceValue("UNITID")
		if err != nil {
			logrus.Error("GetNextSequenceValue was failed(variableBuilder.go) " + err.Error())
			return txnbuild.ManageData{}, errorRespObj, errors.New("getNextSequenceValue of unit map was failed " + err.Error())
		}
		unitIdMap := model.UnitIDMap{
			Unit:  element.MeasurementUnit,
			MapID: uint16(data.SequenceValue),
		}
		unitId = uint16(data.SequenceValue)
		err1 := object.InsertToUnitIDMap(unitIdMap)
		if err1 != nil {
			logrus.Error("Insert unit map ID was failed(variableBuilder.go) " + err1.Error())
			return txnbuild.ManageData{}, errorRespObj, errors.New("inserting to unit map ID was failed " + err1.Error())
		}
		unitString = stellarprotocols.UInt16ToByteString(uint16(data.SequenceValue))
	}
	// precision
	tempPrecision, errInPrecisionConvert := stellarprotocols.Int8ToByteString(uint8(element.Precision))
	if errInPrecisionConvert != nil {
		logrus.Info("Error when converting precision(variableBuilder.go) ", errInPrecisionConvert)
		return txnbuild.ManageData{}, errorRespObj, errors.New("error when converting precision " + errInPrecisionConvert.Error())
	}
	precisionString = tempPrecision
	// check if the description is 40 characters
	if len(element.Description) > 40 {
		logrus.Error("Description is greater than 40 character limit(variableBuilder.go) ")
		return txnbuild.ManageData{}, errorRespObj, errors.New("description is greater than 40 character limit ")
	} else {
		if len(element.Description) == 40 {
			descriptionString = element.Description
		} else if len(element.Description) < 40 {
			descriptionString = element.Description + "/"
		}
	}
	// finally check if the 40 bytes are filled if not append 0s to the end
	if len(descriptionString) < 40 {
		remain := 40 - len(descriptionString)
		setReaminder := strings.Repeat("0", remain)
		descriptionString = descriptionString + setReaminder
	}
	// define a 31 zeros string
	decodedStrFutureUse, err := hex.DecodeString(fmt.Sprintf("%062d", 0))
	if err != nil {
		return txnbuild.ManageData{}, errorRespObj, err
	}
	strFutureUse := string(decodedStrFutureUse)
	valueString := valueTypeString + stellarprotocols.UInt64ToByteString(valueId) + variableNameString + dataTypeString + unitString + precisionString + strFutureUse
	keyString := descriptionString + strings.Repeat("0", 24)
	logrus.Info("Building variable with Name string of   : ", keyString)
	logrus.Info("Building variable with value string of : ", valueString)
	variableDefinitionBuilder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}
	// check the lengths of the key and value
	if len(keyString) != 64 || len(valueString) != 64 {
		logrus.Error("Key string length : ", len(keyString))
		logrus.Error("Value string length : ", len(valueString))
		return txnbuild.ManageData{}, errorRespObj, errors.New("length issue on key or value fields on the variable building ")
	}
	respObj := model.ValueDefOutParmas{
		ValueMapID: valueId,
		UnitMapID:  unitId,
	}
	return variableDefinitionBuilder, respObj, nil
}
