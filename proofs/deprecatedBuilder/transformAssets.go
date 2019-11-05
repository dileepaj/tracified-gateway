package deprecatedBuilder

import (
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/deprecatedStellarExecuter"

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

	object2 := deprecatedStellarExecuter.ConcreteProfile{InsertProfileStruct: temp}

	result2 := object2.InsertProfile()

	object := deprecatedStellarExecuter.ConcreteTransform{
		AssetTransfer: AP.AssetTransfer, ProfileID: result2.ProfileTxn}

	result := object.TransformMerge()

	return result
}
