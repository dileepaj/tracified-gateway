package builder

import (
	"main/api/apiModel"
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

type assetTransferInterface interface {
	SendAsset() model.SendAssetResponse
}

type AbstractAssetTransfer struct {
	// Code       string
	// Amount     string
	// Issuerkey  string
	// Reciverkey string
	// Signer     string
	SendAssest apiModel.SendAssest
}

func (AP *AbstractAssetTransfer) AssetTransfer() model.SendAssetResponse {

	object := stellarExecuter.ConcreteSendAssest{Assest: AP.SendAssest}

	result := object.SendAsset()

	return result
}
