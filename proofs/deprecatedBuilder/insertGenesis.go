package deprecatedBuilder

import (
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/deprecatedStellarExecuter"

)

type GenesisInsertInterface interface {
	InsertGenesis() model.InsertGenesisResponse
	InsertProfile() model.InsertProfileResponse
}

type AbstractGenesisInsert struct {
	InsertGenesisStruct apiModel.InsertGenesisStruct
	// Identifiers string
	// InsertType  string
	// PreviousTXNID string
}

func (AP *AbstractGenesisInsert) GenesisInsert() model.InsertGenesisResponse {

	object1 := deprecatedStellarExecuter.ConcreteGenesis{InsertGenesisStruct: AP.InsertGenesisStruct}

	result := object1.InsertGenesis()

	if result.GenesisTxn == "" {
		return result
	}
	temp := apiModel.InsertProfileStruct{Type: "1",
		PreviousProfileID: "",
		PreviousTXNID:     result.GenesisTxn,
		Identifier:        AP.InsertGenesisStruct.Identifier}

	object2 := deprecatedStellarExecuter.ConcreteProfile{InsertProfileStruct: temp}

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
