package model

type RSAKeyPair struct {
	Sha256PublicKey string `json:"sha256pk"`
	Publickey       string `json:"publicKey"`
	PrivateKey      string `json:"secretKey"`
}

type RSAKeySaveSuccess struct {
	Message string `json:"message"`
}

type RSAPublickey struct {
	Publickey string `json:"publicKey"`
}
