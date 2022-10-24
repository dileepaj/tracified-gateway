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
des- build the sematic constant manage data's name and value according to the protocol
return the txnbuild.ManageData object

variable definition and byte used
	valueType - 1 byte defieded by protocol -2 for sementic constant
	valueId  - 8 byte defieded by protocol
	value name - 40 byte defieded by protocol
	description - 40 bytes
	sementicConstantDataType - 1 byte defieded by protocol -2 for flaot

Manage data
	name 64 byte character - 	description - 64 byte defieded by protocol
	value 64 byte managedata - valueType + valueId +  + value name + fetureused
*/
func (expertFormula ExpertFormula) BuildSemanticConstantManageData(formulaID string, element model.FormulaItemRequest) (txnbuild.ManageData, model.ValueDefOutParmas, error) {
	valueType := 2
	var valueId uint64
	sementicConstantDataType := 2
	semanticConstantDescription := ""
	variableName := ""
	EMPTY := 0
	errorRespObj := model.ValueDefOutParmas{
		ValueMapID: uint64(EMPTY),
		UnitMapID:  uint16(EMPTY),
	}
	semanticConstantValue := fmt.Sprintf("%g", element.Value.(float64))
	// DB validations for the variable id
	object := dao.Connection{}
	valueMap, errValueMap := object.GetValueMapID(element.ID,formulaID).Then(func(data interface{}) interface{} {
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
		logrus.Info("Value ID is not recorded in the DB Map")
		data, err := object.GetNextSequenceValue("VALUEID")
		if err != nil {
			logrus.Error("GetNextSequenceValue was failed" + err.Error())
			return txnbuild.ManageData{}, errorRespObj, errors.New("Get next sequence value of value map was failed")
		}
		valueIdMap := model.ValueIDMap{
			ValueId:   element.ID,
			ValueType: "semanticCONSTANT",
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
	// check variable description is 20 character
	if len(element.Description) > 40 {
		logrus.Error("Description is greater than 40 character limit")
		return txnbuild.ManageData{}, errorRespObj, errors.New("Description is greater than 40 character limit")
	} else {
		if len(element.Description) < 40 {
			// add 0s to the rest of the name
			remain := 40 - len(element.Description)
			setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain-1))
			semanticConstantDescription = element.Description + `/` + setReaminder
		} else {
			semanticConstantDescription = element.Description
		}
	}

	//future use of 24 bytes
	keyFutureUse := fmt.Sprintf("%s", strings.Repeat("0", 24))

	//Variable name - 40 bytes
	if len(element.Name) > 40 || element.Name == "" {
		logrus.Error("Value name is greater than 40 character limit or Empty")
		return txnbuild.ManageData{}, errorRespObj, errors.New("Value name is greater than 40 character limit")
	} else {
		if len(element.Name) < 40 {
			// add 0s to the rest of the name
			remain := 40 - len(element.Name)
			setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain-1))
			variableName = element.Name + `/` + setReaminder
		} else {
			variableName = element.Name
		}
	}

	// check value is 20 character
	if len(semanticConstantValue) > 64 {
		logrus.Error("Value is greater than 8 character limit")
		return txnbuild.ManageData{}, errorRespObj, errors.New("Value is greater than 64 character limit")
	} else {
		if len(semanticConstantValue) < 64 {
			// add 0s to the rest of the name
			remain := 64 - len(semanticConstantValue)
			setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
			semanticConstantValue = setReaminder + semanticConstantValue
		}
	}
	// define a 14 zeros string
	decodedStrFetureUsed, err := hex.DecodeString(fmt.Sprintf("%028d", 0))
	if err != nil {
		return txnbuild.ManageData{}, errorRespObj, err
	}
	strFetureUsed := string(decodedStrFetureUsed)
	srtValueType, err := stellarprotocols.Int8ToByteString(uint8(valueType))
	if err != nil {
		return txnbuild.ManageData{}, errorRespObj, errors.New("Error when converting value type to biinary " + err.Error())
	}
	srtDataType, err := stellarprotocols.Int8ToByteString(uint8(sementicConstantDataType))
	if err != nil {
		return txnbuild.ManageData{}, errorRespObj, errors.New("Error when converting data type to binary " + err.Error())
	}
	// semantic constant's manage data key and value
	nameString := semanticConstantDescription + keyFutureUse
	valueString := srtValueType + stellarprotocols.UInt64ToByteString(valueId) + srtDataType + variableName + strFetureUsed
	logrus.Println("Semantic constant Name:   ", nameString)
	logrus.Println("Semantic constant value:   ", valueString)
	// Building the manage data operation
	semanticConstManageData := txnbuild.ManageData{
		Name:  nameString,
		Value: []byte(valueString),
	}
	if len(valueString) != 64 {
		logrus.Error("Length ", len(nameString))
		return txnbuild.ManageData{}, errorRespObj, errors.New("semantic constant name length not equal to 64")
	}
	if len(nameString) > 64 {
		logrus.Error("Length ", len(valueString))
		return txnbuild.ManageData{}, errorRespObj, errors.New("semantic constant value length should be less than or equal to 64")
	}
	respObj := model.ValueDefOutParmas{
		ValueMapID: valueId,
		UnitMapID:  uint16(EMPTY),
	}
	return semanticConstManageData, respObj, nil
}
