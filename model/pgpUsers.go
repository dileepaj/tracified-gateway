package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type TrustNetWorkUser struct {
	ID                 primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name               string             `json:"name"`
	Company            string             `json:"company"`
	Email              string             `json:"email"`
	Password           string             `json:"password"`
	PasswordResetCode  []byte             `json:"pwdresetcode"`
	Contact            string             `json:"contact"`
	Industry           string             `json:"industry"`
	StellerPK          string             `json:"stellarPK"`
	PGPPK              string             `json:"pgppk"`
	PGPPKHash          string             `json:"pgppkhash"`
	PgpSecret          string             `json:"pgpseceret"`
	StellarSecret      string             `json:"stellarsecret"`
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
	PasswordResetCode  []byte             `json:"pwdresetcode"`
	Contact            string             `json:"contact"`
	Industry           string             `json:"industry"`
	StellerPK          string             `json:"stellarPK"`
	PGPPK              string             `json:"pgppk"`
	PGPPKHash          string             `json:"pgppkhash"`
	PgpSecret          string             `json:"pgpseceret"`
	StellarSecret      string             `json:"stellarsecret"`
	DigitalSignature   string             `json:"signature"`
	Signaturehash      string             `json:"signaturehash"`
	Date               string             `json:"regDate"`
	Endorsments        []Endorsments      `json:"endorsements"`
	TXNOrgRegistration string             `json:"txnOrgRegistration"`
}

type AcceptUserEndorsment struct {
	EndorseePKHash string      `json:"pgppkhash"`
	EndorsmentData Endorsments `json:"endorsmentData"`
}

type Endorsments struct {
	UserID              primitive.ObjectID `json:"userid"`
	EndorserPK          string             `json:"endorserpk"`
	EndorsmentsStatus   string             `json:"endorsmentsstatus"`
	TxnVerificationHash string             `json:"txnverificationhash"`
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

type UpdateTrustNetworkUserPassword struct {
	Email             string `json:"email"`
	PasswordResetCode string `json:"pwdresetcode"`
	Password          string `json:"password"`
}

type ResetPsswordRequest struct {
	Email string `json:"email"`
}
