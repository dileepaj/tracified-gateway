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
	Sign   string `json:"Sign" bson:"Sign" validate:"required"`
	Amount string `json:"Amount" bson:"Amount" validate:"required"`
}

type Batches struct {
	NFTName         string `json:"NFTName" bson:"NFTName" validate:"required"`
	TXNHashTrust    string `json:"TXNHashTrust" bson:"TXNHashTrust" validate:"required"`
	TXNHashTransfer string `json:"TXNHashTransfer" bson:"TXNHashTransfer" validate:"required"`
	CurrentOwner    string `json:"CurrentOwner" bson:"CurrentOwner" validate:"required"`
	PreviousOwner   string `json:"PreviousOwner" bson:"PreviousOwner" validate:"required"`
}

type Conversions struct {
	Sender      string `json:"Sender" bson:"Sender" validate:"required"`
	Amount      string `json:"Amount" bson:"Amount" validate:"required"`
	SellAsset   string `json:"SellAsset" bson:"SellAsset" validate:"required"`
	BuyAsset    string `json:"BuyAsset" bson:"BuyAsset" validate:"required"`
	SellIssuer  string `json:"SellIssuer" bson:"SellIssuer" validate:"required"`
	BuyIssuer   string `json:"BuyIssuer" bson:"BuyIssuer" validate:"required"`
	Numerator   int    `json:"Numerator" bson:"Numerator" validate:"required"`
	Denominator int    `json:"Denominator" bson:"Denominator" validate:"required"`
	TXNHash     string `json:"TXNHash" bson:"TXNHash" validate:"required"`
}

type TokenCoversion struct {
	SellerSourceAccount string      `json:"sellerSource" bson:"sellerSource"`
	BuyerSourceAccount  string      `json:"buyerSource" bson:"buyerSource"`
	ManageSellOffer     ManageOffer `json:"manageSellOffer" bson:"manageSellOffer"`
	ManageBuyOffer      ManageOffer `json:"manageBuyOffer" bson:"manageBuyOffer"`
}

type ManageOffer struct {
	TokenName string `json:"tokename" bson:"tokenname" validate:"required"`
	//* OffertType : Sell-Offer || Buy-Offer
	OfferType              string `json:"offerType" bson:"OfferType" validate:"required"`
	TokenIssuerAccount     string `json:"tokenIssuer" bson:"TokenIssuer" validate:"required"`
	NoTokenForSale         int    `json:"no_token_for_sale" bson:"notkenOnsell" validate:"required"`
	UnitPriceofBuyingToken int    `json:"unit_price_buying_token" bson:"unitPriceofBuyingToken" validate:"required"`
}
