package builder

import (
	"main/model"
)

// type InsertProfile struct{}
type ProflieInsertInterface interface {
	InsertProfile() model.InsertProfileResponse
}

type AbstractProfileInsert struct {
}

func (AP *AbstractProfileInsert) ProfileInsert(ProflieInsertInterface ProflieInsertInterface) model.InsertProfileResponse {
	result := ProflieInsertInterface.InsertProfile()

	return result
}
