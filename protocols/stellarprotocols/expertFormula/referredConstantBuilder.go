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
des- build the refreed constant manage data's name and value according to the protocol
return the txnbuild.ManageData object
variable definition and byte used
	valueType - 1 byte defieded by protocol -3 referred constant
	valueId  - 8 byte defieded by protocol
	unit  - 2 byte defieded by protocol  we maintan a map for each unit defineded by expert
	referredConstantDataType - 1 byte defieded by protocol -2 for flaot
	referredConstantDescription - 30 byte defieded by protocol
Manage data
	name 64 byte character - 64 byte refrence Short URL
	value 64 byte managedata - valueType + valueId + referredConstantDataType + referredConstantDescription + unit + fetureused
*/
func (expertFormula ExpertFormula) BuildReferredConstantManageData(formulaID string, element model.FormulaItemRequest) (txnbuild.ManageData, model.ValueDefOutParmas, error) {
	valueType := 3
	var valueId int64
	var unit int64
	referredConstantDataType := 2
	referredConstantDescription := ""
	EMPTY := 0
	errorRespObj := model.ValueDefOutParmas{
		ValueMapID: int64(EMPTY),
		UnitMapID:  int64(EMPTY),
	}

	referredConstantValue := fmt.Sprintf("%g", element.Value)
	// DB validations for the variable id
	object := dao.Connection{}
	valueMap, errValueMap := object.GetValueMapID(element.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errValueMap != nil {
		logrus.Info("Unable to connect gateway datastore ", errValueMap)
	}
	// check if the variable id for this formula is in the variale mapping
	if valueMap != nil {
		logrus.Info("Value ID is already recorded in the DB Map")
		valueMapData := valueMap.(model.ValueIDMap)
		valueId = valueMapData.MapID
	} else {
		// if not add with incrementing id
		logrus.Info("Value ID is already recorded in the DB Map")
		data, err := object.GetNextSequenceValue("VALUEID")
		if err != nil {
			logrus.Error("GetNextSequenceValue was failed" + err.Error())
			return txnbuild.ManageData{}, errorRespObj, errors.New("GetNextSequenceValue of value map was failed")
		}
		valueIdMap := model.ValueIDMap{
			ValueId:   element.ID,
			ValueType: "REFERREDCONSTANT",
			Key:       element.Key,
			FormulaID: formulaID,
			MapID:     data.SequenceValue,
		}
		err1 := object.InsertToValueIDMap(valueIdMap)
		if err1 != nil {
			logrus.Error("Insert Value map ID was failed" + err1.Error())
		}
		valueId = data.SequenceValue
	}
	// check variable name is 30 character
	if len(element.Description) > 30 || element.Description == "" {
		logrus.Error("Description is greater than 30 character limit or Empty")
		return txnbuild.ManageData{}, errorRespObj, errors.New("Description is greater than 30 character limit")
	} else {
		if len(element.Description) < 30 {
			// add 0s to the rest of the DESCRIPTION
			remain := 30 - len(element.Description)
			setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain-1))
			referredConstantDescription = element.Description + `/` + setReaminder
		} else {
			referredConstantDescription = element.Description
		}
	}
	// checked referred ConstantValue value is 8 character
	if len(referredConstantValue) > 8 {
		logrus.Error("Value is greater than 8 character limit")
		return txnbuild.ManageData{}, errorRespObj, errors.New("Value is greater than 8 character limit")
	} else {
		if len(referredConstantValue) < 8 {
			// add 0s to the rest of the name
			remain := 8 - len(referredConstantValue)
			setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
			referredConstantValue = setReaminder + referredConstantValue
		}
	}
	// define a 14 zeros string
	strFetureUsed := fmt.Sprintf("%014d", 0)
	// convert value type Int to binary string
	srtValueType, err := stellarprotocols.StringToBinary(int64(valueType))
	if err != nil {
		return txnbuild.ManageData{}, errorRespObj, errors.New("srtValueType " + err.Error())
	}
	// convert data type Int to binary string
	srtDataType, err := stellarprotocols.StringToBinary(int64(referredConstantDataType))
	if err != nil {
		return txnbuild.ManageData{}, errorRespObj, errors.New("srtValueType " + err.Error())
	}
	// unit building// convert value type Int to binary string
	unitMap, errInUnitIdMap := object.GetUnitMapID(element.MeasurementUnit).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInUnitIdMap != nil {
		logrus.Info("Unable to connect gateway datastore ", errValueMap)
		// return txnbuild.ManageData{}, errors.New("Unable to connect gateway datastore to get value map ID")
	}
	if unitMap != nil {
		logrus.Info(element.MeasurementUnit + " is already recorded in the DB Map")

		// add map id as the unit in the key string
		unitMapData := unitMap.(model.UnitIDMap)
		unit = unitMapData.MapID

	} else {
		// if not add the incrementing id
		logrus.Info(element.MeasurementUnit + " is not recorded in the DB Map")

		// get the current sequence for the units
		data, err := object.GetNextSequenceValue("UNITID")
		if err != nil {
			logrus.Error("GetNextSequenceValue was failed" + err.Error())
			return txnbuild.ManageData{}, errorRespObj, errors.New("GetNextSequenceValue of unit map was failed")
		}

		unitIdMap := model.UnitIDMap{
			Unit:  element.MeasurementUnit,
			MapID: data.SequenceValue,
		}

		err1 := object.InsertToUnitIDMap(unitIdMap)
		if err1 != nil {
			logrus.Error("Insert unit map ID was failed" + err1.Error())
			return txnbuild.ManageData{}, errorRespObj, errors.New("Insert unit map ID was failed")
		}
		unit = data.SequenceValue
	}
	// convert unit Int to binary string
	strUnit, err := stellarprotocols.UnitToBinary(unit)
	if err != nil {
		return txnbuild.ManageData{}, errorRespObj, errors.New("Error coverting unit to binary")
	}
	// referred constant's manage data key and value
	nameString := element.MetricReference.Url
	valueString := stellarprotocols.ConvertingBinaryToByteString(srtValueType) + stellarprotocols.UInt64ToByteString(valueId) + stellarprotocols.ConvertingBinaryToByteString(srtDataType) + referredConstantValue + referredConstantDescription + stellarprotocols.ConvertingBinaryToByteString(strUnit) + strFetureUsed

	fmt.Println("referred constant Name:   ", nameString)
	fmt.Println("referred constant value:   ", valueString)

	// Building the manage data operation
	semanticConstManageData := txnbuild.ManageData{
		Name:  nameString,
		Value: []byte(valueString),
	}

	if len(valueString) != 64 {
		logrus.Error("Length ", len(valueString))
		return txnbuild.ManageData{}, errorRespObj, errors.New("Referred contant  value length not equal to 64")
	}
	if len(nameString) > 64 || len(nameString) == 0 {
		logrus.Error("Length ", len(nameString))
		return txnbuild.ManageData{}, errorRespObj, errors.New("Referred contant name length should be less than or equal to 64")
	}

	respObj := model.ValueDefOutParmas{
		ValueMapID: valueId,
		UnitMapID:  unit,
	}

	return semanticConstManageData, respObj, nil
}
