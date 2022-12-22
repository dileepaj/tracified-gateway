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

func ValueCodeGenerator(formulaJSON model.FormulaBuildingRequest) (string, []string, error) {
	valueList := formulaJSON.MetricExpertFormula.Formula // list of values in the formula
	var valueInitializations []string                    // list of initializations for the values
	var valueSetters []string                            // list of setters for the values(only for variables)
	var setterNames []string                             // list of setter names for the values(only for variables)
	var parts []string

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
				return "", setterNames, errors.New("Unable to connect to gateway datastore(valueCodeGenerator) Error: " + errInGettingNextSequence.Error() + "\n")
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
				return "", setterNames, errors.New("Unable to connect to gateway datastore(valueCodeGenerator) Error: " + err.Error() + "\n")
			}
		}

		keyNew := strings.ReplaceAll(selectedValue.Key, "$", "")

		if selectedValue.Type == "VARIABLE" {
			// adding comments
			comment := "\n\t" + `// value initialization for ` + selectedValue.Type + ` -> ` + keyNew + "\n"
			valueInitializations = append(valueInitializations, comment)

			// variable initialization
			valueInitializer := "\t" + `Variable ` + keyNew + ` = Variable(Value("`
			valueInitializer = valueInitializer +
				selectedValue.Type + `", "` +
				selectedValue.ID + `", "` +
				selectedValue.Name + `", 0, 0, "` +
				selectedValue.Description + `"), "` +
				selectedValue.MeasurementUnit + `", ` +
				strconv.Itoa(int(selectedValue.Precision)) + `);`
			valueInitializations = append(valueInitializations, valueInitializer)

			// variable setter
			commentForSetter := "\n\t" + `// value setter for ` + selectedValue.Type + ` ` + keyNew + "\n"
			valueSetter := "\t" + `function set` + keyNew + `(int256 _` + keyNew + `, int256 _EXPONENT) public {` + "\n\t"
			valueSetter = valueSetter + "\t" + keyNew + `.value.value = _` + keyNew + ";\n\t"
			valueSetter = valueSetter + "\t" + keyNew + `.value.exponent = _EXPONENT;` + "\n\t" + `}`
			valueSetters = append(valueSetters, commentForSetter)
			valueSetters = append(valueSetters, valueSetter)

			// add the setter name to the list
			setterNames = append(setterNames, "set"+keyNew)

		} else {
			// convert the value to string anf then to float
			valueAsString := fmt.Sprintf("%g", selectedValue.Value)
			exponentOfTheValueLen := 0
			// make the constant values according to the format (m, n) where (m x 10^n)
			if strings.Contains(valueAsString, ".") {
				exponentOfTheValueLen = len(valueAsString[strings.Index(valueAsString, ".")+1:]) * -1
				valueAsString = strings.Replace(valueAsString, ".", "", 1)
				valueAsString = strings.TrimLeft(valueAsString, "0")
				parts = strings.Split(valueAsString, "e")
				valueAsString = parts[0]
				logrus.Info("Value as a whole number: ", parts[0])
				logrus.Info("Exponent of the value Len: ", exponentOfTheValueLen)
			}

			if selectedValue.Type == "SEMANTICCONSTANT" {
				// adding comments
				comment := "\n\t" + `// value initialization for ` + selectedValue.Type + ` -> ` + keyNew + "\n"
				valueInitializations = append(valueInitializations, comment)

				// constant initialization
				valueInitializer := "\t" + `SemanticConstant ` + keyNew + ` = SemanticConstant(Value("` +
					selectedValue.Type + `", "` +
					selectedValue.ID + `", "` +
					selectedValue.Name + `", ` +
					valueAsString + `, ` +
					strconv.Itoa(exponentOfTheValueLen) + `, "` +
					selectedValue.Description + `"));`
				valueInitializations = append(valueInitializations, valueInitializer)
			} else if selectedValue.Type == "REFERREDCONSTANT" {
				// adding comments
				comment := "\n\t" + `// value initialization for ` + selectedValue.Type + ` -> ` + keyNew + "\n"
				valueInitializations = append(valueInitializations, comment)

				// constant initialization
				valueInitializer := "\t" + `ReferredConstant ` + keyNew + ` = ReferredConstant(Value("` +
					selectedValue.Type + `", "` +
					selectedValue.ID + `", "` +
					selectedValue.Name + `", ` +
					valueAsString + `, ` +
					strconv.Itoa(exponentOfTheValueLen) + `, "` +
					selectedValue.Description + `"), "` +
					selectedValue.MeasurementUnit + `", "` +
					selectedValue.MetricReference.Reference + `");`
				valueInitializations = append(valueInitializations, valueInitializer)
			}
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

	return valueCode, setterNames, nil
}
