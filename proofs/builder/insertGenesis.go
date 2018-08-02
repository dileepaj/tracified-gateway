package builder

import (
	"github.com/Tracified-Gateway/model"
)

// type InsertProfile struct{}
type GenesisInsertInterface interface {
	InsertGenesis() model.InsertGenesisResponse
	InsertProfile() model.InsertGenesisResponse
}

type AbstractGenesisInsert struct {
}

func (AP *AbstractGenesisInsert) GenesisInsert(GenesisInsertInterface GenesisInsertInterface) model.InsertGenesisResponse {

	result := GenesisInsertInterface.InsertGenesis()

	if result.GenesisTxn == "" {
		return result
	}

	result2 := GenesisInsertInterface.InsertProfile()
	if result2.Txn == "" {
		return result2
	}

	return result2
}
