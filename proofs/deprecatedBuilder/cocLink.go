package deprecatedBuilder

import (
	"github.com/dileepaj/tracified-gateway/api/apiModel"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/proofs/executer/deprecatedStellarExecuter"

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

	object2 := deprecatedStellarExecuter.ConcreteProfile{
		InsertProfileStruct: temp}

	result2 := object2.InsertProfile()

	object := deprecatedStellarExecuter.ConcreteCoCLinkage{ChangeOfCustodyLink: AP.ChangeOfCustodyLink, ProfileId: result2.ProfileTxn}

	result := object.CoCLinkage()

	return result
}
