package model

type RSAKeyPair struct {
	Sha256PyublicKye string `json:"rsapk"`
	Publickey        string `json:"publicKey"`
	PrivateKey       string `json:"secretKey"`
}

type RSAKeySaveSuccess struct {
	Message string `json:"message"`
}
