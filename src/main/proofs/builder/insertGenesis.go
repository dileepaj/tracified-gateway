package builder

import (
	"main/api/apiModel"
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

type GenesisInsertInterface interface {
	InsertGenesis() model.InsertGenesisResponse
	InsertProfile() model.InsertProfileResponse
}

type AbstractGenesisInsert struct {
	InsertProfileStruct apiModel.InsertProfileStruct
	// Identifiers string
	// InsertType  string
	// PreviousTXNID string
}

func (AP *AbstractGenesisInsert) GenesisInsert() model.InsertGenesisResponse {

	object1 := stellarExecuter.ConcreteGenesis{InsertProfileStruct: AP.InsertProfileStruct}

	result := object1.InsertGenesis()

	if result.GenesisTxn == "" {
		return result
	}
	AP.InsertProfileStruct.Identifier = result.Identifiers
	AP.InsertProfileStruct.PreviousTXNID = result.GenesisTxn

	object2 := stellarExecuter.ConcreteProfile{InsertProfileStruct: AP.InsertProfileStruct}

	result2 := object2.InsertProfile()
	if result2.ProfileTxn == "" {
		result.Error.Message = result2.Error.Message
		result.Error.Code = result2.Error.Code
		return result
	}
	result.Error.Message = result2.Error.Message
	result.Error.Code = result2.Error.Code
	result.ProfileTxn = result2.ProfileTxn

	return result
}
