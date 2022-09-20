package stellarprotocols

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

func BuildVariableDefinitionManageData(valueId string, variableName string, dataType string, unit string, precision string, description string) (txnbuild.ManageData, error) {

	VALUETYPE := "1"
	valueIdString := ""
	valueTypeString := ""
	variableNameString := ""
	unitString := ""

	descriptionString := ""

	fmt.Println("------------------Called -------------------------------------")

	//define the value type
	//this is a variable therefore the value type is 1
	valueTypeInBits := stringToBin(VALUETYPE)
	x := bitString(valueTypeInBits)
	valueInBytes := x.AsByteSlice()
	valueType := string(valueInBytes)
	valueTypeString = valueType

	//DB validations for the variable id
	object := dao.Connection{}
	valueMap, errValueMap := object.GetValueMapID(valueId).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errValueMap != nil {
		logrus.Info("Unable to connect gateway datastore ", errValueMap)
		return txnbuild.ManageData{}, errors.New("Unable to connect gateway datastore to get value map ID")
	}
	//check if the variable name for this formula is in the variale mapping
	if valueMap != nil {
		logrus.Info(variableName + " is already recorded in the DB Map")

		//add the value map part as the value id to the manage data key part string
		valueMapData := valueMap.(model.ValueIDMap)
		tempValueID := int64(valueMapData.MapID)
		valueIdString = strconv.FormatInt(tempValueID, 10)

	} else {
		//if not add with incrementing id
		logrus.Info(variableName + " is not recorded in the DB Map")
		data, err := object.GetNextSequenceValue("VALUEID")
		if err != nil {
			logrus.Error("GetNextSequenceValue was failed" + err.Error())
			return txnbuild.ManageData{}, errors.New("GetNextSequenceValue of value map was failed")
		}

		valueIdMap := model.ValueIDMap{
			ValueId:   valueId,
			ValueType: "VARIABLE",
			ValueName: variableName,
			MapID:     data.SequenceValue,
		}

		err1 := object.InsertToValueIDMap(valueIdMap)
		if err1 != nil {
			logrus.Error("Insert Value map ID was failed" + err1.Error())
			return txnbuild.ManageData{}, errors.New("Insert Value map ID was failed")
		}

		//add the data as the new value id to the manage data key part string
		tempValueID := int64(data.SequenceValue)
		valueIdString = strconv.FormatInt(tempValueID, 10)
	}

	//check weather the valueIdString has 8 characters if not append 0s to the left side
	if len(valueIdString) < 8 {
		remain := 8 - len(valueIdString)
		setRemainder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		valueIdString = setRemainder + valueIdString
	}

	//check variable name is 20 character
	if len(variableName) > 20 {
		logrus.Error("Variable name is greater than 20 character limit")
		return txnbuild.ManageData{}, errors.New("Variable name is greater than 20 character limit")
	} else {
		if len(variableName) == 20 {
			variableNameString = variableName
		} else if len(variableName) < 20 {
			variableNameString = variableName + "/"
		}
	}

	//finally check if the 20 bytes are filled if not append 0s at end
	if len(variableNameString) < 20 {
		remain := 20 - len(variableNameString)
		setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		variableNameString = variableNameString + setReaminder
	}

	//TODO: depending on the data type decide the integer to be assigned

	//TODO: depending on the unit type decide the integer to be asigned
	//convert unit type character -> byte -> bits
	unitMap, errInUnitIdMap := object.GetUnitMapID(unit).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errInUnitIdMap != nil {
		logrus.Info("Unable to connect gateway datastore ", errValueMap)
		//return txnbuild.ManageData{}, errors.New("Unable to connect gateway datastore to get value map ID")
	}
	//check if the unit is in the unit map
	if unitMap != nil {
		logrus.Info(unit + " is already recorded in the DB Map")

		//TODO: add map id as the unit in the key string

	} else {
		//if not add the incrementing id
		logrus.Info(unit + " is not recorded in the DB Map")

		//get the current sequence for the units
		data, err := object.GetNextSequenceValue("UNITID")
		if err != nil {
			logrus.Error("GetNextSequenceValue was failed" + err.Error())
			return txnbuild.ManageData{}, errors.New("GetNextSequenceValue of unit map was failed")
		}

		unitIdMap := model.UnitIDMap{
			Unit:  unit,
			MapID: data.SequenceValue,
		}

		err1 := object.InsertToUnitIDMap(unitIdMap)
		if err1 != nil {
			logrus.Error("Insert unit map ID was failed" + err1.Error())
			return txnbuild.ManageData{}, errors.New("Insert unit map ID was failed")
		}

		//TODO: add the unit id to the key string
	}

	//TODO:
	//precision

	//check if the description is 40 characters
	if len(description) > 40 {
		logrus.Error("Description is greater than 40 character limit")
		return txnbuild.ManageData{}, errors.New("Description is greater than 40 character limit")
	} else {
		if len(description) == 40 {
			descriptionString = description
		} else if len(description) < 40 {
			descriptionString = description + "/"
		}
	}

	//finally check if the 40 bytes are filled if not append 0s to the end
	if len(descriptionString) < 40 {
		remain := 40 - len(descriptionString)
		setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		descriptionString = descriptionString + setReaminder
	}

	keyString := valueTypeString + valueIdString + variableNameString + unitString
	valueString := descriptionString

	logrus.Info("Building variable with key string of   : ", keyString)
	logrus.Info("Building variable with value string of : ", valueString)

	variableDefinitionBuilder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	return variableDefinitionBuilder, nil
}
