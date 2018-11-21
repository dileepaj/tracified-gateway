package builder

import (
	"github.com/tracified-gateway/api/apiModel"
	"github.com/tracified-gateway/model"
	"github.com/tracified-gateway/proofs/executer/stellarExecuter"
)

type transformAssetsInterface interface {
	CreateTrustline() model.InsertDataResponse
}

type AbstractTransformAssets struct {
	AssetTransfer apiModel.AssetTransfer
	// Code1  string
	// Limit1 string
	// Code2  string
	// Limit2 string
	// Code3  string
	// Limit3 string
	// Code4  string
	// Limit4 string

	// Reciver1 string
	// Reciver2 string
}

func (AP *AbstractTransformAssets) TransformAssets() string {

	temp := apiModel.InsertProfileStruct{
		Type:              "1",
		PreviousProfileID: AP.AssetTransfer.PreviousProfileID,
		PreviousTXNID:     AP.AssetTransfer.PreviousTXNID,
		Identifier:        AP.AssetTransfer.Identifier}

	object2 := stellarExecuter.ConcreteProfile{InsertProfileStruct: temp}

	result2 := object2.InsertProfile()

	object := stellarExecuter.ConcreteTransform{
		AssetTransfer: AP.AssetTransfer, ProfileID: result2.ProfileTxn}

	result := object.TransformMerge()

	return result
}
