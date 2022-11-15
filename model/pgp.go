package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type RSAKeyPair struct {
	Sha256PublicKey string `json:"sha256pk"`
	Publickey       string `json:"publicKey"`
	PrivateKey      string `json:"secretKey"`
}

type RSAKeySaveSuccess struct {
	Message string `json:"message"`
}

type RSAPublickey struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	PgpPublickkey    string
	PgpPksha256      string
	StellarPublickey string
	Username         string
}
