package validations

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/go-playground/validator/v10"
)

func ValidateFormulaBuilder(element model.FormulaBuildingRequest) error {
	validate := validator.New()
	for i := 0; i < len(element.Formula); i++ {
		errInValidateFormulaItem := ValidateFormulaItem(element.Formula[i])
		if errInValidateFormulaItem != nil {
			return errInValidateFormulaItem
		}
	}
	err := validate.Struct(element)
	if err != nil {
		return err
	}
	return nil
}

func ValidateFormulaItem(element model.FormulaItemRequest) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}
	errInFieldValidation := ValidateFields(element)
	if errInFieldValidation != nil {
		return errInFieldValidation
	}
	return nil
}

func ValidateFields(element model.FormulaItemRequest) error {
	//checking the required fields in each of the field type
	if element.Type == "DATA" {
		//check the variable type validations
		if element.Name == "" || element.Key == "" || element.Description == "" || element.MeasurementUnit == "" {
			return errors.New("Incorrect data type fields")
		}
	} else if element.Type == "CONSTANT" && element.MetricReferenceId == "" {
		//check the semantic constant type validations
		if element.Description == "" || element.Name == "" || element.Key == "" {
			return errors.New("Incorrect semantic constant type fields")
		}
	} else if element.Type == "CONSTANT" && element.MetricReferenceId != "" {
		//check the referred contant type validations
		if element.Name == "" || element.Description == "" || element.Key == "" || element.MetricReferenceId == "" || element.MeasurementUnit == "" || element.MetricReference.Name == "" || element.MetricReference.MeasurementUnit == "" || element.MetricReference.Url == "" {
			return errors.New("Incorrect referred constant type fields")
		}
	}

	return nil
}
