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

	//this is a variable therefore the value type is 1
	//convert value type character -> byte -> bits

	//check if the variable name for this formula is in the variale mapping
	//if not add with incrementing id
	//if so get that id to be added to the manage data
	object := dao.Connection{}
	data, err := object.GetNextSequenceValue("VALUEID")
	if err != nil {
		logrus.Error("GetNextSequenceValu was failed" + err.Error())
	}

	valueIdMap := model.ValueIDMap{
		ValueId:   valueId,
		ValueType: "Variable",
		ValueName: variableName,
		MapID:     data.SequenceValue,
	}

	err1 := object.InsertToValueIDMap(valueIdMap)
	if err1 != nil {
		logrus.Error("Insert ExpertIDMap was failed" + err1.Error())
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

	keyString := valueTypeString + variableNameString + descriptionString

	variableDefinitionBuilder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(""),
	}

	return variableDefinitionBuilder, nil
}
