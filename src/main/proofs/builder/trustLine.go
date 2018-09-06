package builder

import (
	"main/api/apiModel"
	"main/model"
	"main/proofs/executer/stellarExecuter"
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

	object := stellarExecuter.ConcreteTrustline{TrustlineStruct: AP.TrustlineStruct}
	// object := stellarExecuter.ConcreteTrustline{Code: AP.Code, Limit: AP.Limit, Issuerkey: AP.Issuerkey, Signerkey: AP.Signerkey}

	result := object.CreateTrustline()

	return result
}
