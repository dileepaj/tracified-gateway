package builder

import (
	"main/api/apiModel"
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

type assetTransferInterface interface {
	CreateTrustline() model.InsertDataResponse
}

type AbstractAssetTransfer struct {
	// Code       string
	// Amount     string
	// Issuerkey  string
	// Reciverkey string
	// Signer     string
	SendAssest apiModel.SendAssest
}

func (AP *AbstractAssetTransfer) AssetTransfer() string {

	object := stellarExecuter.ConcreteSendAssest{Assest: AP.SendAssest}

	result := object.SendAsset()

	return result
}
