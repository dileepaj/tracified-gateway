package validations

import (
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/go-playground/validator/v10"
)

func ValidateBatchCoinConvert(e model.BatchCoinConvert) error {
	validate := validator.New()
	err := validate.Struct(e)
	if err != nil {
		return err
	}
	return nil
}

func ValidateArtifactCoinConvert(e model.ArtifactCoinConvert) error {
	validate := validator.New()
	err := validate.Struct(e)
	if err != nil {
		return err
	}
	return nil
}

func ValidateCreatePool(e model.CreatePool) error {
	validate := validator.New()
	err := validate.Struct(e)
	if err != nil {
		return err
	}
	return nil
}

func ValidateCreatePoolForArtifact(e model.CreatePoolForArtifact) error {
	validate := validator.New()
	err := validate.Struct(e)
	if err != nil {
		return err
	}
	return nil
}
