package metricBinding

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

//Used to define the general metadata of a single variable
/* KEY
*	Field name or Stage name - 25 bytes (string) - Done
*	Key name - 20 bytes - Done
*	Future use - 19 bytes - Done
/
/* VALUE
*	Value id - 8 bytes (unsigned long integer) - Done
*	Variable name - 20 bytes (string) - Done
* 	Binding type - 1 byte (byte) - Done
*	Stage or Reference ID - 8 bytes (unsigned long integer)
*	future use - 27 bytes - Done
*/
func (metric *MetricBinding) ValueDefinitionBuilder(element model.GeneralValueDefBuildRequest) (txnbuild.ManageData, string, []byte, error) {
	// key string components
	resourceNameString := ""
	keyNameString := ""
	futureUseInKey := ""
	// value string components
	var valueID uint64
	var resourceID uint64
	valueTypeString := ""
	variableNameString := ""
	// Validate resource name 30 byte limit
	if len(element.ResourceName) > 30 {
		logrus.Error("Resource name(Key or filed) is greater than 30 character limit")
		return txnbuild.ManageData{}, "", []byte{}, errors.New("Resource name(Key or filed) is greater than 30 character limit")
	} else {
		if len(element.ResourceName) == 30 {
			resourceNameString = element.ResourceName
		} else if len(element.ResourceName) < 30 {
			resourceNameString = element.ResourceName + "/"
		}
	}
	if len(resourceNameString) < 30 {
		remain := 30 - len(resourceNameString)
		setRemainder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		resourceNameString = resourceNameString + setRemainder
	}
	// Validate key name 30 limit
	if len(element.Key) > 30 {
		logrus.Error("Key is greater than 30 character limit")
		return txnbuild.ManageData{}, "", []byte{}, errors.New("Key is greater than 30 character limit")
	} else {
		if len(element.Key) == 30 {
			keyNameString = element.Key
		} else if len(element.Key) < 30 {
			keyNameString = element.Key + "/"
		}
	}
	if len(keyNameString) < 30 {
		remain := 30 - len(keyNameString)
		setRemainder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		keyNameString = keyNameString + setRemainder
	}
	// Build future use in key
	futureUseInKey = fmt.Sprintf("%s", strings.Repeat("0", 4))

	// get value id from the map for the type, key and formula ID
	object := dao.Connection{}
	variableDefMap, errWhenRetrievingVariableInfo := object.GetValueMapID(element.VariableUUID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errWhenRetrievingVariableInfo != nil {
		logrus.Error("Unable to connect gateway datastore ", errWhenRetrievingVariableInfo)
		// return txnbuild.ManageData{}, errors.New("Unable to connect gateway datastore " + errWhenRetrievingVariableInfo.Error())
	}
	if variableDefMap == nil {
		// requested varible does not exists in the DB
		logrus.Error("Requested variable " + element.VariableName + " does not exists in the gateway DB")
		return txnbuild.ManageData{}, "", []byte{}, errors.New("Requested variable " + element.VariableName + " does not exists in the gateway DB")
	} else {
		valueMapData := variableDefMap.(model.ValueIDMap)
		valueID = valueMapData.MapID
	}
	// check the variable name length for 20 bytes
	if len(element.VariableName) > 20 {
		logrus.Error("Variable name is greater than 20 character limit")
		return txnbuild.ManageData{}, "", []byte{}, errors.New("Variable name is greater than 20 character limit")
	} else {
		if len(element.VariableName) == 20 {
			variableNameString = element.VariableName
		} else if len(element.VariableName) < 20 {
			variableNameString = element.VariableName + "/"
		}
	}
	// check whether the new variable name string saturated the 20 byte limit, if not add 0s to rest
	if len(variableNameString) < 20 {
		remain := 20 - len(variableNameString)
		setRemainder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		variableNameString = variableNameString + setRemainder
	}
	// check if the binding type is 0 or 1
	if element.BindingType == 0 || element.BindingType == 1 {
		tempValueType, errInValueTypeConvert := stellarprotocols.Int8ToByteString(uint8(element.BindingType))
		if errInValueTypeConvert != nil {
			logrus.Error("Error when converting value type ", errInValueTypeConvert)
			return txnbuild.ManageData{}, "", []byte{}, errors.New("Error when converting value type " + errInValueTypeConvert.Error())
		}
		valueTypeString = tempValueType
	} else {
		logrus.Error("Invalid binding type, should be 1 or 0")
		return txnbuild.ManageData{}, "", []byte{}, errors.New("Invalid binding type, should be 1 or 0")
	}
	// Stage/Ref id mapping and adding to string
	resourceIdMap, errResourceMap := object.GetResourceMapID(element.ResourceID).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errResourceMap != nil {
		logrus.Info(errResourceMap)
	}
	if resourceIdMap != nil {
		logrus.Info(element.ResourceID + " is already recorded in the DB Map")
		resourceMapData := resourceIdMap.(model.ResourceIdMap)
		resourceID = resourceMapData.MapID
	} else {
		logrus.Info(element.ResourceID + " is not recorded in the DB Map")
		data, err := object.GetNextSequenceValue("RESOURCEID")
		if err != nil {
			logrus.Error("Retrieving resource id from map was failed " + err.Error())
			return txnbuild.ManageData{}, "", []byte{}, errors.New("Retrieving resource id from map was failed")
		}
		insertResourceMap := model.ResourceIdMap{
			ResourceType: element.ResourceType,
			ResourceID:   element.ResourceID,
			MapID:        data.SequenceValue,
		}
		errWhenInsertingToResourceMap := object.InsertToResourceIDMap(insertResourceMap)
		if errWhenInsertingToResourceMap != nil {
			logrus.Error("Inserting to resource map ID was failed" + errWhenInsertingToResourceMap.Error())
			return txnbuild.ManageData{}, "", []byte{}, errors.New("Inserting to resource map ID was failed")
		}
		resourceID = data.SequenceValue
	}
	// future use in value
	decodedStrFetureUsed, err := hex.DecodeString(fmt.Sprintf("%054d", 0))
	if err != nil {
		return txnbuild.ManageData{}, "", []byte{}, errors.New("Feture used byte building issue in formula definition")
	}
	futureUseInValue := string(decodedStrFetureUsed)
	// build key and value string
	keyString := resourceNameString + keyNameString + futureUseInKey
	valueString := stellarprotocols.UInt64ToByteString(valueID) + variableNameString + valueTypeString + stellarprotocols.UInt64ToByteString(resourceID) + futureUseInValue
	// check the key value string length for 64 byte limit
	if len(keyString) > 64 {
		logrus.Error("Key string exceeding the given 64 byte limit in variable manage data builder. Length : ", len(keyString))
		return txnbuild.ManageData{}, "", []byte{}, errors.New("Key string exceeding the given 64 byte limit in variable manage data builder")
	}
	if len(valueString) > 64 {
		logrus.Error("Value string exceeding the given 64 byte limit in variable manage data builder. Length : ", len(valueString))
		return txnbuild.ManageData{}, "", []byte{}, errors.New("Value string exceeding the given 64 byte limit in variable manage data builder")
	}
	generalValueDefinitionBuilder := txnbuild.ManageData{
		Name:  keyString,
		Value: []byte(valueString),
	}
	return generalValueDefinitionBuilder, keyString, []byte(valueString), nil
}
