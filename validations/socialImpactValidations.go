package validations

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/model"
	"github.com/go-playground/validator/v10"
)

func ValidateFormulaBuilder(element model.FormulaBuildingRequest) error {
	validate := validator.New()
	// validate inner object array
	for i := 0; i < len(element.MetricExpertFormula.Formula); i++ {
		errInValidateFormulaItem := ValidateFormulaItem(element.MetricExpertFormula.Formula[i])
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
	// checking the required fields in each of the field type
	if element.Type == "DATA" {
		// check the variable type validations
		if element.Name == "" || element.Key == "" || element.MeasurementUnit == "" {
			return errors.New("Incorrect data type fields")
		}
	} else if element.Type == "CONSTANT" && element.MetricReferenceId == "" {
		// check the semantic constant type validations
		if element.Name == "" || element.Key == "" {
			return errors.New("Incorrect semantic constant type fields")
		}
	} else if element.Type == "CONSTANT" && element.MetricReferenceId != "" {
		// check the referred contant type validations
		if element.Name == "" || element.MetricReference.Description == "" || element.Key == "" || element.MetricReferenceId == "" || element.MeasurementUnit == "" || element.MetricReference.Name == "" || element.MetricReference.MeasurementUnit == "" {
			return errors.New("Incorrect referred constant type fields")
		}
	}
	return nil
}

func ValidateFormulaForMetricBuilding(element model.FormulaForMetricBinding) error {
	// check if the required fields are empty
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}
	// validate the inner object array
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
	// check binding time and validate master data and stage data
	if element.BindingType == 0 {
		// validation of master data
		if element.BindData.Master.KeyDataType == "" || element.BindData.Master.KeyValue == "" || element.BindData.Master.MetaDataName == "" || element.BindData.Master.PrimaryKeyName == "" || element.BindData.Master.ValueColumnName == "" || element.BindData.Master.ValueDataType == "" {
			return errors.New("Body of the master data is incorrect")
		}
	} else if element.BindingType == 1 {
		// validation of stage data
		if element.BindData.Stage.StageId == "" || element.BindData.Stage.WorkflowId == "" || element.BindData.Stage.StageName == "" || element.BindData.Stage.FieldName == "" || element.BindData.Stage.FieldDataType == "" || element.BindData.Stage.FieldId == "" {
			return errors.New("Body of the stage data is incorrect")
		}
	} else {
		return errors.New("Bind data type is invalid")
	}

	return nil
}

func ValidateBindDataType(element1 model.Master, element2 model.Stage, bindDataType int) error {
	if bindDataType == 0 {
		// validation of master data type
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

func ValidateMetricDataBindingRequest(element model.MetricDataBindingRequest) error {
	// validate Metric Object
	errWhenValidatingMetricReq := ValidateMetricObject(element.Metric)
	if errWhenValidatingMetricReq != nil {
		return errWhenValidatingMetricReq
	}

	// validate User object
	errWhenValidatingUserDetails := ValidateUser(element.User)
	if errWhenValidatingUserDetails != nil {
		return errWhenValidatingUserDetails
	}

	return nil
}

func ValidateMetricObject(element model.MetricReq) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}

	for i := 0; i < len(element.Activities); i++ {
		// Validate activity array
		errWHenValidatingActivity := ValidateActivityArray(element.Activities[i])
		if errWHenValidatingActivity != nil {
			return errWHenValidatingActivity
		}
	}

	return nil
}

func ValidateActivityArray(element model.MetricDataBindActivityRequest) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}

	// validate stage array
	errWhenValidatingStageBlock := ValidateStageReq(element.Stage)
	if errWhenValidatingStageBlock != nil {
		return errWhenValidatingStageBlock
	}

	// validate metric formula
	errWhenValidatingMetricFormula := ValidateMetricFormulaReq(element.MetricFormula)
	if errWhenValidatingMetricFormula != nil {
		return errWhenValidatingMetricFormula
	}

	// validate ActivityFormulaDefinitionManageData
	errWhenValidatingActivityFormulaDefinitionManageData := ValidateActivityFormulaDefinitionManageData(element.ActivityFormulaDefinitionManageData)
	if errWhenValidatingActivityFormulaDefinitionManageData != nil {
		return errWhenValidatingActivityFormulaDefinitionManageData
	}
	return nil
}

func ValidateActivityFormulaDefinitionManageData(element model.ActivityFormulaDefinitionManageData) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}
	return nil
}

func ValidateMetricFormulaReq(element model.MetricFormulaReq) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}

	// validate formula
	for i := 0; i < len(element.Formula); i++ {
		errWhenValidatingFormula := ValidateFormula(element.Formula[i])
		if errWhenValidatingFormula != nil {
			return errWhenValidatingFormula
		}
	}

	// validate metric expert formula
	errWhenValidatingMetricExpertFormula := ValidateMetricExpertFormula(element.MetricExpertFormula)
	if errWhenValidatingMetricExpertFormula != nil {
		return errWhenValidatingMetricExpertFormula
	}

	return nil
}

func ValidateFormula(element model.FormulaDetails) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}

	// validate master data
	if element.ArtifactTemplateID != "" {
		// validate artifact template
		errWhenValidatingArtifactTemplate := ValidateArtifactTemplate(element.ArtifactTemplate)
		if errWhenValidatingArtifactTemplate != nil {
			return errWhenValidatingArtifactTemplate
		}

	}

	return nil
}

func ValidateArtifactTemplate(element model.ArtifactTemplate) error {
	if element.ID == "" || element.Name == "" || element.FieldName == "" {
		return errors.New("Artifact template validation failed")
	}
	return nil
}

func ValidateStageReq(element model.StageReq) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}
	return nil
}

func ValidateUser(element model.User) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}
	return nil
}

func ValidateMetricDataBindArtifactRequest(element model.MetricDataBindArtifactRequest) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}
	// validate metric formula
	errWhenValidatingMetricFormula := ValidateMetricFormula(element.MetricFormula)
	if errWhenValidatingMetricFormula != nil {
		return errWhenValidatingMetricFormula
	}

	return nil
}

func ValidateMetricFormula(element model.MetricFormula) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}

	// validate formula details
	for i := 0; i < len(element.Formula); i++ {
		errWhenValidateingFormulaDetails := ValidateFormulaDetails(element.Formula[i])
		if errWhenValidateingFormulaDetails != nil {
			return err
		}
	}
	// validate MetricExpertFormula
	errWHenValidatingMetricExpertFormula := ValidateMetricExpertFormula(element.MetricExpertFormula)
	if errWHenValidatingMetricExpertFormula != nil {
		return errWHenValidatingMetricExpertFormula
	}
	// validate pivot field
	for i := 0; i < len(element.PivotField); i++ {
		errWhenValidatingPivotFields := ValidatePivotField(element.PivotField[i])
		if errWhenValidatingPivotFields != nil {
			return errWhenValidatingPivotFields
		}
	}

	return nil
}

func ValidateFormulaDetails(element model.FormulaDetails) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}
	return nil
}

func ValidateMetricExpertFormula(element model.MetricExpertFormula) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}

	// validate full formula
	for i := 0; i < len(element.Formula); i++ {
		errWhenValidatingFullFormula := ValidateFullFormula(element.Formula[i])
		if errWhenValidatingFullFormula != nil {
			return errWhenValidatingFullFormula
		}
	}

	return nil
}

func ValidatePivotField(element model.PivotField) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}
	return nil
}

func ValidateFullFormula(element model.FullFormula) error {
	validate := validator.New()
	err := validate.Struct(element)
	if err != nil {
		return err
	}
	// checking the required fields in each of the field type
	if element.Type == "DATA" {
		// check the variable type validations
		if element.Name == "" || element.Key == "" || element.ID == "" {
			return errors.New("incorrect data type fields")
		}
	} else if element.Type == "CONSTANT" {
		// check the constant type validations
		if element.Name == "" || element.ID == "" || element.Key == "" || element.Value == "" {
			return errors.New("incorrect constant type fields")
		}
	} else if element.Type == "OPERATOR" {
		// check the operator contant type validations
		if element.ID == "" {
			return errors.New("incorrect operator type fields")
		}
	}
	return nil
}
