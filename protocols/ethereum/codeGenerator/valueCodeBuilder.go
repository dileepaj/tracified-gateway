package codeGenerator

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum"
	"github.com/sirupsen/logrus"
)

func ValueCodeGenerator(formulaJSON model.FormulaBuildingRequest) (string, error) {
	valueList := formulaJSON.MetricExpertFormula.Formula // list of values in the formula
	var valueInitializations []string                    // list of initializations for the values
	var valueSetters []string                            // list of setters for the values(only for variables)

	/* loop through the values in the formula and generate the initializations and setters
	   for each variable -> initialize the variable and setter
	   for each constant -> initialize the constant */

	firstComment := "\n\t// Value initializations"
	valueInitializations = append(valueInitializations, firstComment)
	object := dao.Connection{}
	for _, selectedValue := range valueList {

		// check if the value is mapped in the DB and retrieve existing valueID or create a new valueID
		valueMapId, errGettingValueMap := object.GetValueMapID(selectedValue.ID, formulaJSON.MetricExpertFormula.ID).Then(func(data interface{}) interface{} {
			return data
		}).Await()
		if errGettingValueMap != nil {
			logrus.Info("Unable to connect to gateway datastore(valueCodeGenerator) ", errGettingValueMap)
		}
		if valueMapId != nil {
			logrus.Info(selectedValue.Name + " is already recorded in the DB Map(valueCodeGenerator) ")
		} else {
			// get the next sequence value for the value id
			data, errInGettingNextSequence := object.GetNextSequenceValue("VALUEID")
			if errInGettingNextSequence != nil {
				logrus.Info("Unable to connect to gateway datastore(valueCodeGenerator) ", errInGettingNextSequence)
				return "", errors.New("Unable to connect to gateway datastore(valueCodeGenerator) Error: " + errInGettingNextSequence.Error() + "\n")
			}
			valueIdMap := model.ValueIDMap{
				ValueId:   selectedValue.ID,
				ValueType: selectedValue.Type,
				Key:       selectedValue.Key,
				FormulaID: formulaJSON.MetricExpertFormula.ID,
				ValueName: selectedValue.Name,
				MapID:     data.SequenceValue,
			}

			// insert the valueID map into the DB
			err := object.InsertToValueIDMap(valueIdMap)
			if err != nil {
				logrus.Info("Unable to connect to gateway datastore(valueCodeGenerator) ", err)
				return "", errors.New("Unable to connect to gateway datastore(valueCodeGenerator) Error: " + err.Error() + "\n")
			}
		}

		keyNew := strings.ReplaceAll(selectedValue.Key, "$", "")

		if selectedValue.Type == "VARIABLE" {
			// adding comments
			comment := "\n\t" + `// value initialization for ` + selectedValue.Type + ` -> ` + keyNew + "\n"
			valueInitializations = append(valueInitializations, comment)

			// variable initialization
			valueInitializer := "\t" + `Variable ` + keyNew + ` = Variable(Value("` 
			valueInitializer= valueInitializer + selectedValue.Type + `", "` + selectedValue.ID + `", "` + selectedValue.Name + `", 0, "` + selectedValue.Description + `"), "` + selectedValue.MeasurementUnit + `", ` + strconv.Itoa(int(selectedValue.Precision)) + `);` 
			valueInitializations = append(valueInitializations, valueInitializer)

			// variable setter
			commentForSetter := "\n\t" + `// value setter for ` + selectedValue.Type + ` ` + keyNew + "\n"
			valueSetter := "\t" + `function set` + keyNew + `(int _` + keyNew + `) public {` + "\n\t"
			valueSetter = valueSetter + "\t" + keyNew + `.value.value = _` + keyNew + ";\n\t" + `}`
			valueSetters = append(valueSetters, commentForSetter)
			valueSetters = append(valueSetters, valueSetter)

		} else if selectedValue.Type == "SEMANTICCONSTANT" {
			// adding comments
			comment := "\n\t" + `// value initialization for ` + selectedValue.Type + ` -> ` + keyNew + "\n"
			valueInitializations = append(valueInitializations, comment)

			// constant initialization
			valueAsString := fmt.Sprintf("%f", selectedValue.Value)
			valueInitializer := "\t" + `SemanticConstant ` + keyNew + ` = SemanticConstant(Value("` + selectedValue.Type + `", "` + selectedValue.ID + `", "` + selectedValue.Name + `", ` + valueAsString + `, "` + selectedValue.Description + `"));`  
			valueInitializations = append(valueInitializations, valueInitializer)
		} else if selectedValue.Type == "REFERREDCONSTANT" {
			// adding comments
			comment := "\n\t" + `// value initialization for ` + selectedValue.Type + ` -> ` + keyNew + "\n"
			valueInitializations = append(valueInitializations, comment)

			// constant initialization
			valueAsString := fmt.Sprintf("%f", selectedValue.Value)
			valueInitializer := "\t" + `ReferredConstant ` + keyNew + ` = ReferredConstant(Value("` + selectedValue.Type + `", "` + selectedValue.ID + `", "` + selectedValue.Name + `", ` + valueAsString + `, "` + selectedValue.Description + `"), "` + selectedValue.MeasurementUnit + `", "` + selectedValue.MetricReference.Reference + `");` 
			valueInitializations = append(valueInitializations, valueInitializer)
		}
	}

	// remove the duplicates in the value initializations and setters
	valueInitializations = ethereum.RemoveDuplicatesInAnArray(valueInitializations)
	valueSetters = ethereum.RemoveDuplicatesInAnArray(valueSetters)

	// generate the string
	var valueCode string
	for _, value := range valueInitializations {
		valueCode = valueCode + value
	}
	valueCode = valueCode + "\n"
	for _, value := range valueSetters {
		valueCode = valueCode + value
	}

	return valueCode, nil
}
