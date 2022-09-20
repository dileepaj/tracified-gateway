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

func BuildReferredConstantManageData(element model.FormulaItemRequest) (txnbuild.ManageData, error) {
	fmt.Println("element      ", element)
	valueType := 3
	valueId := ""
	unit := 2
	referredConstantDataType := 2
	referredConstantDescription := ""
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
		logrus.Info(valueId + " is already recorded in the DB Map")
		valueMapData := valueMap.(model.ValueIDMap)
		valueId = fmt.Sprintf("%08d", valueMapData.MapID)
	} else {
		// if not add with incrementing id
		logrus.Info(valueId + " is not recorded in the DB Map")
		data, err := object.GetNextSequenceValue("VALUEID")
		if err != nil {
			logrus.Error("GetNextSequenceValue was failed" + err.Error())
			return txnbuild.ManageData{}, errors.New("GetNextSequenceValue of value map was failed")
		}
		valueIdMap := model.ValueIDMap{
			ValueId:   valueId,
			ValueType: "REFERREDCONSTANT",
			MapID:     data.SequenceValue,
		}
		err1 := object.InsertToValueIDMap(valueIdMap)
		if err1 != nil {
			logrus.Error("Insert Value map ID was failed" + err1.Error())
		}
		valueId = fmt.Sprintf("%08d", data.SequenceValue)
	}
	// check variable name is 20 character
	if len(element.Description) > 30 || element.Description == "" {
		logrus.Error("Description is greater than 30 character limit or Empty")
		return txnbuild.ManageData{}, errors.New("Description is greater than 30 character limit")
	} else {
		if len(element.Description) < 30 {
			// add 0s to the rest of the name
			remain := 30 - len(element.Description)
			setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain-1))
			referredConstantDescription = element.Description + `\` + setReaminder
		} else {
			referredConstantDescription = element.Description
		}
	}
	// check value is 20 character
	if len(referredConstantValue) > 8 {
		fmt.Println(referredConstantValue)
		logrus.Error("Value is greater than 8 character limit")
		return txnbuild.ManageData{}, errors.New("Value is greater than 20 character limit")
	} else {
		if len(referredConstantValue) < 8 {
			// add 0s to the rest of the name
			remain := 8 - len(referredConstantValue)
			setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
			referredConstantValue = setReaminder + referredConstantValue
		}
	}
	strFetureUsed := fmt.Sprintf("%014d", 0)
	strUnit, err := UnitToBinary(int64(unit))
	if err != nil {
		return txnbuild.ManageData{}, errors.New("Value is greater than 20 character limit " + err.Error())
	}
	srtValueType, err := StringToBinary(int64(valueType))
	if err != nil {
		return txnbuild.ManageData{}, errors.New("Value is greater than 20 character limit " + err.Error())
	}
	srtDataType, err := StringToBinary(int64(referredConstantDataType))
	if err != nil {
		return txnbuild.ManageData{}, errors.New("Value is greater than 20 character limit " + err.Error())
	}
	fmt.Println(strUnit+"    cnv             ", ConvertingBinaryToByteString(strUnit))
	fmt.Println(len(ConvertingBinaryToByteString(srtValueType)))
	fmt.Println(len(valueId))
	fmt.Println(len(ConvertingBinaryToByteString(srtDataType)))
	fmt.Println(len(referredConstantValue))
	fmt.Println(len(referredConstantDescription))
	fmt.Println(len(ConvertingBinaryToByteString(strUnit)))
	fmt.Println(len(strFetureUsed))
	// referred constant's manage data key and value
	keyString := ConvertingBinaryToByteString(srtValueType) + valueId + ConvertingBinaryToByteString(srtDataType) + referredConstantValue + referredConstantDescription + ConvertingBinaryToByteString(strUnit) + strFetureUsed
	valueString := element.MetricReference.Url
	logrus.Println("Semantic constant key:   ", keyString)
	logrus.Println("Semantic constant value:   ", valueString)

	// Building the manage data operation
	semanticConstManageData := txnbuild.ManageData{
		Name:  valueString,
		Value: []byte(keyString),
	}

	if len(keyString) != 64 {
		logrus.Error("Length ", len(keyString))
		return txnbuild.ManageData{}, errors.New("Referred contant  key length not equal to 64")
	}
	if len(valueString) > 64 {
		logrus.Error("Length ", len(valueString))
		return txnbuild.ManageData{}, errors.New("Referred contant value length should be less than or equal to 64")
	}
	return semanticConstManageData, nil
}
