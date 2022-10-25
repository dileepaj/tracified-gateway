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
des- build the refreed constant manage data's name and value according to the protocol
return the txnbuild.ManageData object
variable definition and byte used
	valueType - 1 byte defieded by protocol -3 referred constant
	valueId  - 8 byte defieded by protocol
	unit  - 2 byte defieded by protocol  we maintan a map for each unit defineded by expert
	referredConstantDataType - 1 byte defieded by protocol -2 for flaot
	variable name - 30 byte defieded by protocol
	description 40 bytes
Manage data
	name 64 byte character - decription + future use
	value 64 byte managedata - valueType + valueId + referredConstantDataType + variable name + unit + fetureused
*/
func (expertFormula ExpertFormula) BuildReferredConstantManageData(formulaID string, element model.FormulaItemRequest) (txnbuild.ManageData, model.ValueDefOutParmas, error) {
	valueType := 3
	var valueId uint64
	var unit uint16
	referredConstantDataType := 2
	referredConstantDescription := ""
	variableName := ""
	EMPTY := 0
	errorRespObj := model.ValueDefOutParmas{
		ValueMapID: uint64(EMPTY),
		UnitMapID:  uint16(EMPTY),
	}

	referredConstantValue := fmt.Sprintf("%g", element.Value.(float64))
	// DB validations for the variable id
	object := dao.Connection{}
	valueMap, errValueMap := object.GetValueMapID(element.ID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errValueMap != nil {
		logrus.Info("Unable to connect to gateway datastore ", errValueMap)
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
			logrus.Error("Get Next Sequence Value for value id was failed" + err.Error())
			return txnbuild.ManageData{}, errorRespObj, errors.New("Get Next Sequence Value for value id was failed" + err.Error())
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
			logrus.Error("Inserting Value map ID was failed" + err1.Error())
		}
		valueId = data.SequenceValue
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
	decodedStrFetureUsed, err := hex.DecodeString(fmt.Sprintf("%028d", 0))
	if err != nil {
		return txnbuild.ManageData{}, errorRespObj, err
	}
	strFetureUsed := string(decodedStrFetureUsed)
	// convert value type Int to binary string
	srtValueType, err := stellarprotocols.Int8ToByteString(uint8(valueType))
	if err != nil {
		return txnbuild.ManageData{}, errorRespObj, errors.New("Error when converting value type to byte  " + err.Error())
	}
	// convert data type Int to byte string
	srtDataType, err := stellarprotocols.Int8ToByteString(uint8(referredConstantDataType))
	if err != nil {
		return txnbuild.ManageData{}, errorRespObj, errors.New("Error when converting data type to byte  " + err.Error())
	}
	// unit building// convert value type Int to byte string
	unitMap, errInUnitIdMap := object.GetUnitMapID(element.MeasurementUnit).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInUnitIdMap != nil {
		logrus.Info("Unable to connect to gateway datastore ", errValueMap)
		// return txnbuild.ManageData{}, errors.New("Unable to connect gateway datastore to get value map ID")
	}
	if unitMap != nil {
		logrus.Info(element.MeasurementUnit + " is already recorded in the DB Map")
		// add map id as the unit in the key string
		unitMapData := unitMap.(model.UnitIDMap)
		unit = unitMapData.MapID
	} else {
		logrus.Info(element.MeasurementUnit + " is not recorded in the DB Map")
		// get the current sequence for the units
		data, err := object.GetNextSequenceValue("UNITID")
		if err != nil {
			logrus.Error("Get next sequence value ID was failed " + err.Error())
			return txnbuild.ManageData{}, errorRespObj, errors.New("Get next sequence value ID was failed " + err.Error())
		}
		unitIdMap := model.UnitIDMap{
			Unit:  element.MeasurementUnit,
			MapID: uint16(data.SequenceValue),
		}
		err1 := object.InsertToUnitIDMap(unitIdMap)
		if err1 != nil {
			logrus.Error("Inserting unit map ID was failed" + err1.Error())
			return txnbuild.ManageData{}, errorRespObj, errors.New("Inserting unit map ID was failed " + err1.Error())
		}
		unit = uint16(data.SequenceValue)
	}

	//variable builder
	if len(element.Name) > 30 || element.Name == "" {
		logrus.Error("Variable name is greater than 30 character limit or Empty")
		return txnbuild.ManageData{}, errorRespObj, errors.New("Variable name is greater than 30 character limit")
	} else {
		if len(element.Name) < 30 {
			// add 0s to the rest of the DESCRIPTION
			remain := 30 - len(element.Name)
			setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain-1))
			variableName = element.Name + `/` + setReaminder
		} else {
			variableName = element.Name
		}
	}

	//build description for 40 bytes
	if len(element.MetricReference.Description) > 40 {
		logrus.Error("Description is greater than 30 character limit")
		return txnbuild.ManageData{}, errorRespObj, errors.New("Description is greater than 30 character limit")
	} else {
		if len(element.MetricReference.Description) < 40 {
			// add 0s to the rest of the DESCRIPTION
			remain := 40 - len(element.MetricReference.Description)
			setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain-1))
			referredConstantDescription = element.MetricReference.Description + `/` + setReaminder
		} else {
			referredConstantDescription = element.MetricReference.Description
		}
	}

	keyFutureUse := fmt.Sprintf("%s", strings.Repeat("0", 24))

	// referred constant's manage data key and value
	nameString := referredConstantDescription + keyFutureUse
	valueString := srtValueType + stellarprotocols.UInt64ToByteString(valueId) + srtDataType + stellarprotocols.Float64ToByteString(element.Value.(float64)) + variableName + stellarprotocols.UInt16ToByteString(uint16(unit)) + strFetureUsed

	logrus.Println("referred constant Name:   ", nameString)
	logrus.Println("referred constant value:   ", valueString)
	// Building the manage data operation
	semanticConstManageData := txnbuild.ManageData{
		Name:  nameString,
		Value: []byte(valueString),
	}
	if len(valueString) != 64 {
		logrus.Error("Length ", len(valueString))
		return txnbuild.ManageData{}, errorRespObj, errors.New("Referred constant value length not equal to 64")
	}
	if len(nameString) > 64 || len(nameString) == 0 {
		logrus.Error("Length ", len(nameString))
		return txnbuild.ManageData{}, errorRespObj, errors.New("Referred constant name length should be less than or equal to 64")
	}
	respObj := model.ValueDefOutParmas{
		ValueMapID: valueId,
		UnitMapID:  unit,
	}
	return semanticConstManageData, respObj, nil
}
