package expertformula

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

/*
BuildVariableDefinitionManageData
des-Build the variable definition manage data
*/
func (expertFormula ExpertFormula) BuildVariableDefinitionManageData(element model.FormulaItemRequest) (txnbuild.ManageData, model.ValueDefOutParmas, error) {
	VALUETYPE := 1
	DATATYPE := 2
	EMPTY := 0
	var valueId int64
	var unitId int64
	dataTypeString := ""
	valueTypeString := ""
	variableNameString := ""
	unitString := ""
	precisionString := ""

	descriptionString := ""

	errorRespObj := model.ValueDefOutParmas{
		ValueMapID: int64(EMPTY),
		UnitMapID:  int64(EMPTY),
	}

	//define the value type
	//this is a variable therefore the value type is 1
	tempValueType, errInValueTypeConvert := stellarprotocols.StringToBinary(int64(VALUETYPE))
	if errInValueTypeConvert != nil {
		logrus.Info("Error when converting value type ", errInValueTypeConvert)
		return txnbuild.ManageData{}, errorRespObj, errors.New("Error when converting value type")
	}
	valueTypeString = stellarprotocols.ConvertingBinaryToByteString(tempValueType)

	//DB validations for the variable id
	object := dao.Connection{}
	valueMap, errValueMap := object.GetValueMapID(element.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errValueMap != nil {
		logrus.Info("Unable to connect gateway datastore ", errValueMap)
		//return txnbuild.ManageData{}, errors.New("Unable to connect gateway datastore to get value map ID")
	}
	//check if the variable name for this formula is in the variale mapping
	if valueMap != nil {
		logrus.Info(element.Name + " is already recorded in the DB Map")

		//add the value map part as the value id to the manage data key part string
		valueMapData := valueMap.(model.ValueIDMap)
		valueId = valueMapData.MapID

	} else {
		//if not add with incrementing id
		logrus.Info(element.Name + " is not recorded in the DB Map")
		data, err := object.GetNextSequenceValue("VALUEID")
		if err != nil {
			logrus.Error("GetNextSequenceValue was failed" + err.Error())
			return txnbuild.ManageData{}, errorRespObj, errors.New("GetNextSequenceValue of value map was failed")
		}

		valueIdMap := model.ValueIDMap{
			ValueId:   element.ID,
			ValueType: "VARIABLE",
			ValueName: element.Name,
			MapID:     data.SequenceValue,
		}

		err1 := object.InsertToValueIDMap(valueIdMap)
		if err1 != nil {
			logrus.Error("Insert Value map ID was failed" + err1.Error())
			return txnbuild.ManageData{}, errorRespObj, errors.New("Insert Value map ID was failed")
		}

		//add the data as the new value id to the manage data key part string
		valueId = data.SequenceValue
	}

	//check variable name is 20 character
	if len(element.Name) > 20 {
		logrus.Error("Variable name is greater than 20 character limit")
		return txnbuild.ManageData{}, errorRespObj, errors.New("Variable name is greater than 20 character limit")
	} else {
		if len(element.Name) == 20 {
			variableNameString = element.Name
		} else if len(element.Name) < 20 {
			variableNameString = element.Name + "/"
		}
	}

	//finally check if the 20 bytes are filled if not append 0s at end
	if len(variableNameString) < 20 {
		remain := 20 - len(variableNameString)
		setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		variableNameString = variableNameString + setReaminder
	}

	//depending on the data type decide the integer to be assigned
	tempDataType, errInDataTypeConvert := stellarprotocols.StringToBinary(int64(DATATYPE))
	if errInDataTypeConvert != nil {
		logrus.Info("Error when converting data type ", errInDataTypeConvert)
		return txnbuild.ManageData{}, errorRespObj, errors.New("Error when converting data type")
	}
	dataTypeString = stellarprotocols.ConvertingBinaryToByteString(tempDataType)

	//depending on the unit type decide the integer to be asigned
	//convert unit type character -> byte -> bits
	unitMap, errInUnitIdMap := object.GetUnitMapID(element.MeasurementUnit).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInUnitIdMap != nil {
		logrus.Info("Unable to connect gateway datastore ", errValueMap)
		//return txnbuild.ManageData{}, errors.New("Unable to connect gateway datastore to get value map ID")
	}
	//check if the unit is in the unit map
	if unitMap != nil {
		logrus.Info(element.MeasurementUnit + " is already recorded in the DB Map")

		//add map id as the unit in the key string
		unitMapData := unitMap.(model.UnitIDMap)
		unitId = unitMapData.MapID
		strUnit, err := stellarprotocols.UnitToBinary(unitMapData.MapID)
		if err != nil {
			return txnbuild.ManageData{}, errorRespObj, errors.New("Error coverting unit to binary")
		}

		unitString = stellarprotocols.ConvertingBinaryToByteString(strUnit)

	} else {
		//if not add the incrementing id
		logrus.Info(element.MeasurementUnit + " is not recorded in the DB Map")

		//get the current sequence for the units
		data, err := object.GetNextSequenceValue("UNITID")
		if err != nil {
			logrus.Error("GetNextSequenceValue was failed" + err.Error())
			return txnbuild.ManageData{}, errorRespObj, errors.New("GetNextSequenceValue of unit map was failed")
		}

		unitIdMap := model.UnitIDMap{
			Unit:  element.MeasurementUnit,
			MapID: data.SequenceValue,
		}
		unitId = data.SequenceValue

		err1 := object.InsertToUnitIDMap(unitIdMap)
		if err1 != nil {
			logrus.Error("Insert unit map ID was failed" + err1.Error())
			return txnbuild.ManageData{}, errorRespObj, errors.New("Insert unit map ID was failed")
		}
		strUnit, err := stellarprotocols.UnitToBinary(data.SequenceValue)
		if err != nil {
			return txnbuild.ManageData{}, errorRespObj, errors.New("Error coverting unit to binary")
		}

		unitString = stellarprotocols.ConvertingBinaryToByteString(strUnit)
	}

	strValueID, err := stellarprotocols.IDToBinary(valueId)
	if err != nil {
		return txnbuild.ManageData{}, errorRespObj, errors.New("Error coverting unit to binary")
	}

	//precision
	tempPrecision, errInPrecisionConvert := stellarprotocols.StringToBinary(int64(element.Precision))
	if errInPrecisionConvert != nil {
		logrus.Info("Error when converting precision ", errInPrecisionConvert)
		return txnbuild.ManageData{}, errorRespObj, errors.New("Error when converting precision")
	}
	precisionString = stellarprotocols.ConvertingBinaryToByteString(tempPrecision)

	//check if the description is 40 characters
	if len(element.Description) > 40 {
		logrus.Error("Description is greater than 40 character limit")
		return txnbuild.ManageData{}, errorRespObj, errors.New("Description is greater than 40 character limit")
	} else {
		if len(element.Description) == 40 {
			descriptionString = element.Description
		} else if len(element.Description) < 40 {
			descriptionString = element.Description + "/"
		}
	}

	//finally check if the 40 bytes are filled if not append 0s to the end
	if len(descriptionString) < 40 {
		remain := 40 - len(descriptionString)
		setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		descriptionString = descriptionString + setReaminder
	}

	keyString := valueTypeString + stellarprotocols.ConvertingBinaryToByteString(strValueID) + variableNameString + dataTypeString + unitString + precisionString
	valueString := descriptionString

	logrus.Info("Building variable with key string of   : ", keyString)
	logrus.Info("Building variable with value string of : ", valueString)

	variableDefinitionBuilder := txnbuild.ManageData{
		Name:  valueString,
		Value: []byte(keyString),
	}

	//check the lengths of the key and value
	if len(keyString) > 64 || len(valueString) > 64 {
		logrus.Error("Key string length : ", len(keyString))
		logrus.Error("Value string length : ", len(valueString))
		return txnbuild.ManageData{}, errorRespObj, errors.New("Length issue on key or value fields on the variable building")
	}

	respObj := model.ValueDefOutParmas{
		ValueMapID: valueId,
		UnitMapID:  unitId,
	}

	return variableDefinitionBuilder, respObj, nil
}
