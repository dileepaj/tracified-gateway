package request

type TransactionCollectionBodyGenesis struct {
	Identifier    string `json:"Identifier" bson:"identifier" validate:"required"`
	SequenceNo    int64  `json:"SequenceNo" bson:"sequenceno" validate:"required"`
	ToIdentifier  string `json:"ToIdentifier" bson:"toidentifier" validate:"required"`
	MapIdentifier string `json:"MapIdentifier" bson:"mapidentifier" validate:"required"`
	XDR           string `json:"XDR" bson:"xdr" validate:"required"`
	Orphan        bool   `json:"Orphan" bson:"orphan" validate:"required"`
	MergeBlock    int    `json:"MergeBlock" bson:"mergeblock" validate:"required"`
	TenantID      string `json:"TenantID" bson:"tenantid" validate:"required"`
	StageID       string `json:"StageID" bson:"stageid" validate:"required"`
}
