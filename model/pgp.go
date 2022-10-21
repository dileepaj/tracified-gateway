package model

//proof of POE
type PGPAccount struct {
	PGPPublickKey    string `json:"pgpPublickKey"`
	PGPPKsha256      string `json:"pgppksha256"`
	StellarPublicKey string `json:"stellarPublicKey"`
	UserName         string `json:"userName"`
}
