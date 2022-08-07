package validations

import (
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/go-playground/validator/v10"
)

func ValidateBatchCoinConvert(e model.CoinConvertBody) error {
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

func ValidateCreatePool(e model.CreatePoolBody) error {
	validate := validator.New()
	err := validate.Struct(e)
	if err != nil {
		return err
	}
	return nil
}
