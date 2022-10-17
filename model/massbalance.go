package model

//"github.com/stellar/go/txnbuild"

type Splits struct {
	Destination []Destination `json:"Destination" bson:"Destination" validate:"required"`
	Sender      string        `json:"Sender" bson:"Sender" validate:"required"`
	NFTName     string        `json:"NFTName" bson:"NFTName" validate:"required"`
	Status      string        `json:"Status" bson:"Status" validate:"required"`
	Issuer      string        `json:"Issuer" bson:"Issuer" validate:"required"`
	Limit       string        `json:"Limit" bson:"Limit" validate:"required"`
}

type Merges struct {
	Sender      []Destination `json:"Sender" bson:"Sender" validate:"required"`
	Destination string        `json:"Destination" bson:"Destination" validate:"required"`
	NFTName     string        `json:"NFTName" bson:"NFTName" validate:"required"`
	Status      string        `json:"Status" bson:"Status" validate:"required"`
	Issuer      string        `json:"Issuer" bson:"Issuer" validate:"required"`
	Limit       string        `json:"Limit" bson:"Limit" validate:"required"`
}

type Destination struct {
	Source string `json:"Source" bson:"Source" validate:"required"`
	Amount string `json:"Amount" bson:"Amount" validate:"required"`
}

type Batches struct {
	NFTName       string `json:"NFTName" bson:"NFTName" validate:"required"`
	TXNHash       string `json:"TXNHash" bson:"TXNHash" validate:"required"`
	CurrentOwner  string `json:"CurrentOwner" bson:"CurrentOwner" validate:"required"`
	PreviousOwner string `json:"PreviousOwner" bson:"PreviousOwner" validate:"required"`
}
