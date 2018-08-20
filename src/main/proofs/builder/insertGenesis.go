package builder

import (
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

type GenesisInsertInterface interface {
	InsertGenesis() model.InsertGenesisResponse
	InsertProfile() model.InsertProfileResponse
}

type AbstractGenesisInsert struct {
	Identifiers string
	InsertType  string
	// PreviousTXNID string
}

func (AP *AbstractGenesisInsert) GenesisInsert() model.InsertGenesisResponse {

	object1 := stellarExecuter.ConcreteGenesis{Identifiers: AP.Identifiers, InsertType: AP.InsertType}

	result := object1.InsertGenesis()

	if result.GenesisTxn == "" {
		return result
	}

	object2 := stellarExecuter.ConcreteProfile{Identifiers: result.Identifiers, InsertType: result.TxnType, PreviousTXNID: result.GenesisTxn, PreviousProfileID: ""}

	result2 := object2.InsertProfile()
	if result2.Txn == "" {
		result.Error.Message = result2.Error.Message
		result.Error.Code = result2.Error.Code
		return result
	}
	result.Error.Message = result2.Error.Message
	result.Error.Code = result2.Error.Code
	result.Txn = result2.Txn

	return result
}
