package builder

import (
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/stellarExecuter"
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
