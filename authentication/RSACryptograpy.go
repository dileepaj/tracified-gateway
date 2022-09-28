package authentication

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"

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