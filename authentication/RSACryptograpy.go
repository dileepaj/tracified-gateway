package authentication

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/sirupsen/logrus"
)

func CreateSignature(cipherText string, rsaPrivateKey rsa.PrivateKey) []byte {
	rng := rand.Reader
	message := []byte(cipherText)
	hashed := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rng, &rsaPrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		logrus.Error("Error from signing ", err.Error())
	}
	return signature
}

func VerifySignature(messeges string, signature []byte, pubKey rsa.PublicKey) bool {
	message := []byte(messeges)
	hashed := sha256.Sum256(message)

	err := rsa.VerifyPKCS1v15(&pubKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		logrus.Error("Sign verification issue ", err.Error())
		return false
	}
	return true
}

//Exporting keys to pem stringds
func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) (string, error) {
	privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	return string(privkey_pem), nil
}

func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	pubkey_bytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return "", err
	}
	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkey_bytes,
		},
	)

	return string(pubkey_pem), nil
}

func ParseRsaPrivateKeyFromPemStr(privPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func ParseRsaPublicKeyFromPemStr(pubPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("Key type is not RSA")
}

//Importing keys from pem string

// func RSA_OAEP_Encrypt(secretMessage string, key rsa.PublicKey) string {
// 	label := []byte("OAEP Encrypted")
// 	rng := rand.Reader
// 	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &key, []byte(secretMessage), label)
// 	CheckError(err)
// 	return base64.StdEncoding.EncodeToString(ciphertext)
// }

// func CheckError(e error) {
// 	if e != nil {
// 		fmt.Println(e.Error())
// 	}
// }

// func RSA_OAEP_Decrypt(cipherText string, privKey rsa.PrivateKey) string {
// 	ct, _ := base64.StdEncoding.DecodeString(cipherText)
// 	label := []byte("OAEP Encrypted")
// 	rng := rand.Reader
// 	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &privKey, ct, label)
// 	CheckError(err)
// 	fmt.Println("Plaintext:", string(plaintext))
// 	return string(plaintext)
// }

// privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
// if err != nil {
// 	logrus.Println(err)
// }

// publicKey := privateKey.PublicKey
// fmt.Println("Public Key  ", publicKey)
// fmt.Println("Private Key  ", privateKey)
// secretMessage := "This is super secret message!"

// // encryptedMessage := authentication.RSA_OAEP_Encrypt(secretMessage, publicKey)

// signature := authentication.CreateSignature(secretMessage, *privateKey)
// sEnc := base64.StdEncoding.EncodeToString(signature)
// fmt.Println("base64 signature ", sEnc)
// // fmt.Println(secretMessage, encryptedMessage)
// authentication.VerifySignature(secretMessage, signature, publicKey)
// // authentication.RSA_OAEP_Decrypt(encryptedMessage, *privateKey)
