package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type TrustNetWorkUser struct {
	ID                 primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name               string             `json:"name"`
	Company            string             `json:"company"`
	Email              string             `json:"email"`
	Password           string             `json:"password"`
	Contact            string             `json:"contact"`
	Industry           string             `json:"industry"`
	StellerPK          string             `json:"stellarPK"`
	PGPPK              string             `json:"pgppk"`
	PGPPKHash          string             `json:"pgppkhash"`
	DigitalSignature   string             `json:"signature"`
	Signaturehash      string             `json:"signaturehash"`
	Date               string             `json:"regDate"`
	Endorsments        []Endorsments      `json:"endorsements"`
	TXNOrgRegistration string             `json:"txnOrgRegistration"`
}

type LoggedInTrustNetworkUser struct {
	ID                 primitive.ObjectID `json:"_id" bson:"_id"`
	Name               string             `json:"name"`
	Company            string             `json:"company"`
	Email              string             `json:"email"`
	Password           string             `json:"password"`
	Contact            string             `json:"contact"`
	Industry           string             `json:"industry"`
	StellerPK          string             `json:"stellarPK"`
	PGPPK              string             `json:"pgppk"`
	PGPPKHash          string             `json:"pgppkhash"`
	DigitalSignature   string             `json:"signature"`
	Signaturehash      string             `json:"signaturehash"`
	Date               string             `json:"regDate"`
	Endorsments        []Endorsments      `json:"endorsements"`
	TXNOrgRegistration string             `json:"txnOrgRegistration"`
}

type AcceptUserEndorsment struct {
	EndorseePKHash string      `json:"endoserpkhash"`
	EndorsmentData Endorsments `json:"endorsmentData"`
}

type Endorsments struct {
	UserID              primitive.ObjectID `json:"userID"`
	EndorserPK          string             `json:"endorserPK"`
	EndorsmentsStatus   string             `json:"status"`
	TxnVerificationHash string             `json:"txnSigVerify"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EndorsmentUpdateSuccess struct {
	Message string `json:"message"`
}

type TrustNetworkUserEndorsmentCount struct {
	Totalendorsements   int `json:"totalendorsements"`
	FullEndorsements    int `json:"fullendorsement"`
	PartialEndorsements int `json:"partialendorsement"`
}
