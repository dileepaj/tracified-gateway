package builder

import (
	"main/api/apiModel"
	"main/model"
	"main/proofs/executer/stellarExecuter"
)

type cocLinkInterface interface {
	CreateTrustline() model.InsertDataResponse
}

type AbstractcocLink struct {
	// Code       string
	// Amount     string
	// IssuerKey  string
	// Reciverkey string
	// Sender     string
	ChangeOfCustodyLink apiModel.ChangeOfCustodyLink
}

func (AP *AbstractcocLink) CoCLink() string {

	object2 := stellarExecuter.ConcreteProfile{InsertProfileStruct: AP.ChangeOfCustodyLink.InsertProfileStruct}

	result2 := object2.InsertProfile()

	object := stellarExecuter.ConcreteCoCLinkage{ChangeOfCustodyLink: AP.ChangeOfCustodyLink, ProfileId: result2.ProfileTxn}

	result := object.CoCLinkage()

	return result
}
