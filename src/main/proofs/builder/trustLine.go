package builder

import (
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

type TrustlineInterface interface {
	CreateTrustline() model.InsertDataResponse
}

type AbstractTrustline struct {
	Code      string
	Limit     string
	Issuerkey string
	Signerkey string
}

func (AP *AbstractTrustline) Trustline() string {

	object := stellarExecuter.ConcreteTrustline{Code: AP.Code, Limit: AP.Limit, Issuerkey: AP.Issuerkey, Signerkey: AP.Signerkey}

	result := object.CreateTrustline()

	return result
}
