package validations

import (
	"github.com/dileepaj/tracified-gateway/apiDemo/model/dtos/request"
	"github.com/go-playground/validator/v10"
)

func ValidateGenesisTDPRequest(tdps []request.TransactionCollectionBodyGenesis) error {
    // Create a new validator instance
    v := validator.New()

    // Validate each TDP struct in the array
    for _, tdp := range tdps {
        if err := v.Struct(tdp); err != nil {
            return err
        }
    }

    // Return nil if all TDP structs in the array pass validation
    return nil
}