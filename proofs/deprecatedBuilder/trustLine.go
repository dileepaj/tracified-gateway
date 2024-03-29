package deprecatedBuilder

import (
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/deprecatedStellarExecuter"

)

type TrustlineInterface interface {
	CreateTrustline() model.InsertDataResponse
}

type AbstractTrustline struct {
	TrustlineStruct apiModel.TrustlineStruct
	// Code      string
	// Limit     string
	// Issuerkey string
	// Signerkey string
}

func (AP *AbstractTrustline) Trustline() string {

	object := deprecatedStellarExecuter.ConcreteTrustline{TrustlineStruct: AP.TrustlineStruct}
	// object := stellarExecuter.ConcreteTrustline{Code: AP.Code, Limit: AP.Limit, Issuerkey: AP.Issuerkey, Signerkey: AP.Signerkey}

	result := object.CreateTrustline()

	return result
}
