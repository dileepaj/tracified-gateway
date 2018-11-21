package builder

import (
	"github.com/tracified-gateway/api/apiModel"
	"github.com/tracified-gateway/model"
	"github.com/tracified-gateway/proofs/executer/stellarExecuter"
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

	temp := apiModel.InsertProfileStruct{
		Type:              "1",
		PreviousProfileID: AP.ChangeOfCustodyLink.PreviousProfileID,
		PreviousTXNID:     AP.ChangeOfCustodyLink.PreviousTXNID,
		Identifier:        AP.ChangeOfCustodyLink.Identifier}

	object2 := stellarExecuter.ConcreteProfile{
		InsertProfileStruct: temp}

	result2 := object2.InsertProfile()

	object := stellarExecuter.ConcreteCoCLinkage{ChangeOfCustodyLink: AP.ChangeOfCustodyLink, ProfileId: result2.ProfileTxn}

	result := object.CoCLinkage()

	return result
}
