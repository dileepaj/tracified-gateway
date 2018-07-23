package builder

import (
	"main/model"
)

// type InsertProfile struct{}
type GenesisInsertInterface interface {
	InsertGenesis() model.InsertGenesisResponse
}

type AbstractGenesisInsert struct {
}

func (AP *AbstractGenesisInsert) GenesisInsert(GenesisInsertInterface GenesisInsertInterface) model.InsertGenesisResponse {
	result := GenesisInsertInterface.InsertGenesis()

	return result
}
