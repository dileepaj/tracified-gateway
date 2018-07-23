package builder

import (
	"main/model"
)

// type InsertData struct{}
type TDPInsertInterface interface {
	InsertDataHash() model.InsertDataResponse
}

type AbstractTDPInsert struct {
}

func (AP *AbstractTDPInsert) TDPInsert(TDPInsertInterface TDPInsertInterface) model.InsertDataResponse {
	result := TDPInsertInterface.InsertDataHash()

	return result
}
