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

func BuildSemanticConstantManageData(valueId string, dataType string, description string, Value string) (txnbuild.ManageData, error) {

	valueType := ""
	semanticConstValueID := ""
	semanticConstDataType := ""
	semanticConstDescription := ""
	semanticConstValue := ""

	//---------------------------------------------- valueType ----------------------------------------------------
	// this is a semantic constant therefore the value type is 2
	// convert value type -> string -> bits
	valueTypeInBits := stringToBin("2")
	valueTypeAsBitString := bitString(valueTypeInBits)

	// convert value type character bits -> bytes
	valueTypeInBytes := valueTypeAsBitString.AsByteSlice()
	valueType = string(valueTypeInBytes)
	fmt.Println(valueType)
	fmt.Println(len(valueType))

	//---------------------------------------------- valueId ----------------------------------------------------
	//DB validations for the valueID
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
		logrus.Info("Semantic constant is already recorded in the DB Map")

	} else {
		//if not add with incrementing id
		logrus.Info("Semantic constant is not recorded in the DB Map")
		data, err := object.GetNextSequenceValue("VALUEID")
		if err != nil {
			logrus.Error("GetNextSequenceValue was failed" + err.Error())
			return txnbuild.ManageData{}, errors.New("GetNextSequenceValue of value map was failed")
		}

		valueIdMap := model.ValueIDMap{
			ValueId:   valueId,
			ValueType: "SEMANTICCONSTANT",
			ValueName: "",
			MapID:     data.SequenceValue,
		}

		err1 := object.InsertToValueIDMap(valueIdMap)
		if err1 != nil {
			logrus.Error("Insert Value map ID was failed" + err1.Error())
		}

		// add data as the new value id to the manage data key part string
		tempValueID := int64(data.SequenceValue)
		semanticConstValueID = strconv.FormatInt(tempValueID, 10)
	}

	//---------------------------------------------- dataType ----------------------------------------------------
	// getting the value according to the data type
	dataTypeString := ""
	valueLength := 0
	switch dataType {
	case "Integer":
		dataTypeString = "1"
		valueLength = 4
	case "Float":
		dataTypeString = "2"
		valueLength = 8
	case "Char":
		dataTypeString = "3"
		valueLength = 1
	case "String":
		dataTypeString = "4"
		valueLength = 64
	case "Boolean":
		dataTypeString = "5"
		valueLength = 1
	}
	dataTypeInBits := stringToBin(dataTypeString)
	dataTypeAsBitString := bitString(dataTypeInBits)

	// convert value type character bits -> bytes
	dataTypeInBytes := dataTypeAsBitString.AsByteSlice()
	semanticConstDataType = string(dataTypeInBytes)
	fmt.Println(semanticConstDataType)
	fmt.Println(len(semanticConstDataType))

	//---------------------------------------------- description ----------------------------------------------------
	//check if the description have 40 characters
	if len(description) > 40 {
		logrus.Error("Description is greater than 40 character limit")
		return txnbuild.ManageData{}, errors.New("Description is greater than 40 character limit")

	} else if len(description) < 40 {

		semanticConstDescription = description + "/"

		// get the remaining number of characters to be filled with 0 and add it to the description
		remain := 40 - len(semanticConstDescription)
		setRemainder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		semanticConstDescription = semanticConstDescription + setRemainder

	} else {
		semanticConstDescription = description
	}

	//---------------------------------------------- semanticConstValue ----------------------------------------------------
	//check if the semantic constant value have 64 characters
	if len(semanticConstValue) > 64 {
		logrus.Error("Value of the Semantic Constant is greater than 64 character limit")
		return txnbuild.ManageData{}, errors.New("Value of the Semantic Constant is greater than 64 character limit")
	} else {
		// check whether the character limit is greater than the maximum character limit of the given data type
		if len(semanticConstValue) > valueLength {
			logrus.Error("Value of the Semantic Constant is greater than " + strconv.Itoa(valueLength) + " character limit")
			return txnbuild.ManageData{}, errors.New("Value of the Semantic Constant is greater than " + strconv.Itoa(valueLength) + " character limit")
		}
	}

	keyString := valueType + semanticConstValueID + semanticConstDataType + semanticConstDescription
	valueString := semanticConstValue

	logrus.Info("Semantic constant key: ", keyString)
	logrus.Info("Semantic constant value: ", valueString)

	// Building the manage data operation
	semanticConstManageData := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}

	return semanticConstManageData, nil
}
