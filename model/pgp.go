package model

//proof of POE
type PGPAccount struct {
	PGPPublickKey    string `json:"pgppublickkey"`
	PGPPKsha256      string `json:"pgppksha256"`
	StellarPublicKey string `json:"stellarpublickey"`
	UserName         string `json:"username"`
}
