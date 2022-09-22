package validations

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/go-playground/validator/v10"
)

func ValidateFormulaBuilder(element model.FormulaBuildingRequest) error {
	validate := validator.New()
	//validate inner object array
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

func ValidateMetricBindingRequest(element model.MetricBindingRequest) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}
	//validate the inner object array
	for i := 0; i < len(element.Formula); i++ {
		errInValidateFormulasInMetricBinding := ValidateFormulaForMetricBuilding(element.Formula[i])
		if errInValidateFormulasInMetricBinding != nil {
			return errInValidateFormulasInMetricBinding
		}
	}

	return nil
}

func ValidateFormulaForMetricBuilding(element model.FormulaForMetricBinding) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}
	//validate the inner object array
	for i := 0; i < len(element.Variable); i++ {
		errInValidateVriablesInMetricBinding := ValidateVariablesForMetricBuilding(element.Variable[i])
		if errInValidateVriablesInMetricBinding != nil {
			return errInValidateVriablesInMetricBinding
		}
	}
	return nil
}

func ValidateVariablesForMetricBuilding(element model.VariableStructure) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}
	// errInValidateBindDataType := ValidateBindDataType(element.BindData.Master, element.BindData.Stage, element.BindingType)
	// if errInValidateBindDataType != nil {
	// 	return errInValidateBindDataType
	// }
	return nil
}

func ValidateBindDataType(element1 model.Master, element2 model.Stage, bindDataType int) error {
	if bindDataType == 0 {
		//validation of master data type
		errInMasterDataValidations := ValidateMasterData(element1)
		if errInMasterDataValidations != nil {
			return errInMasterDataValidations
		}
	} else if bindDataType == 1 {
		errInStageDataValidations := ValidateStageData(element2)
		if errInStageDataValidations != nil {
			return errInStageDataValidations
		}
	} else {
		return errors.New("Invalid data bind type")
	}
	return errors.New("Invalid data bind type")
}

func ValidateMasterData(element model.Master) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}
	return nil
}

func ValidateStageData(element model.Stage) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}
	return nil
}
