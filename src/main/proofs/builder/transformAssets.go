package builder

import (
	"main/api/apiModel"
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

type transformAssetsInterface interface {
	CreateTrustline() model.InsertDataResponse
}

type AbstractTransformAssets struct {
	AssestTransfer apiModel.AssestTransfer
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
	object2 := stellarExecuter.ConcreteProfile{InsertProfileStruct: AP.AssestTransfer.InsertProfileStruct}

	result2 := object2.InsertProfile()

	object := stellarExecuter.ConcreteTransform{AssestTransfer: AP.AssestTransfer, ProfileID: result2.ProfileTxn}

	result := object.TransformMerge()

	return result
}
