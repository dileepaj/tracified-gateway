package model

type COCCollectionBody struct {
	Sender   string
	Receiver string
	Xdr      string
	Flag     string
}

type COCCollectionList struct {
	List []COCCollectionBody
}