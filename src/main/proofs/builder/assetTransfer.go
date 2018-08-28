package builder

import (
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

type assetTransferInterface interface {
	CreateTrustline() model.InsertDataResponse
}

type AbstractAssetTransfer struct {
	Code       string
	Amount     string
	Issuerkey  string
	Reciverkey string
	Signer     string
}

func (AP *AbstractAssetTransfer) AssetTransfer() string {

	object := stellarExecuter.ConcreteSendAssest{Code: AP.Code, Amount: AP.Amount, Issuerkey: AP.Issuerkey, Reciverkey: AP.Reciverkey, Signer: AP.Signer}

	result := object.SendAssest()

	return result
}
