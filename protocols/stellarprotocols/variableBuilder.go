package stellarprotocols

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/txnbuild"
)

func BuildVariableDefinitionManageData(valueId string, variableName string, dataType string, unit string, precision string, description string) (txnbuild.ManageData, error) {

	valueTypeString := ""
	variableNameString := ""
	descriptionString := ""

	fmt.Println("------------------Called -------------------------------------")
	//this is a variable therefore the value type is 1
	//convert value type character -> byte -> bits

	//DB validations for the variable id
	object := dao.Connection{}
	valueMap, errValueMap := object.GetValueMapID(valueId).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errValueMap != nil {
		logrus.Info("Unable to connect gateway datastore ", errValueMap)
		//return txnbuild.ManageData{}, errors.New("Unable to connect gateway datastore to get value map ID")
	}
	//check if the variable name for this formula is in the variale mapping
	if valueMap != nil {
		logrus.Info(variableName + " is already recorded in the DB Map")

		//TODO:add thE value map part as the value id to the manage data key part string

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
			ValueType: "Variable",
			ValueName: variableName,
			MapID:     data.SequenceValue,
		}

		err1 := object.InsertToValueIDMap(valueIdMap)
		if err1 != nil {
			logrus.Error("Insert Value map ID was failed" + err1.Error())
		}

		//TODO: add the data as the new value id to the manage data key part string
	}

	//check variable name is 20 character
	if len(variableName) > 20 {
		logrus.Error("Variable name is greater than 20 character limit")
		return txnbuild.ManageData{}, errors.New("Variable name is greater than 20 character limit")
	} else {
		if len(variableName) < 20 {
			//add 0s to the rest of the name
			remain := 20 - len(variableName)
			setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
			variableNameString = variableName + setReaminder
		} else if len(variableName) == 20 {
			variableNameString = variableName
		}
	}

	//depending on the data type decide the integer to be asigned
	//convert data type character -> byte -> bits

	//depending on the unit type decide the integer to be asigned
	//convert unit type character -> byte -> bits

	//precision

	//check if the description is 40 characters
	if len(description) > 40 {
		logrus.Error("Description is greater than 40 character limit")
		return txnbuild.ManageData{}, errors.New("Description is greater than 40 character limit")
	} else {
		if len(description) < 40 {
			//add 0s to the rest of the name
			remain := 40 - len(description)
			setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
			descriptionString = description + setReaminder
		} else if len(description) == 40 {
			descriptionString = description
		}
	}

	keyString := valueTypeString + variableNameString
	valueString := descriptionString

	logrus.Info("Building variable with key string of   : ", keyString)
	logrus.Info("Building variable with value string of : ", valueString)

	variableDefinitionBuilder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	return variableDefinitionBuilder, nil
}
