/**
This File Include the Encrypt and Decrypt method
**/
package accounts

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

/*Encrypt
@desc - Encrypt method is to encrypt or hide any classified text (secret key)
@params - Encrypted SK, decrypt password
*/
func Encrypt(text, mySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(mySecret))
	if err != nil {
		panic(err)
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

/*Decrypt
@desc - Decrypt method is to extract back the encrypted text
@params - Encrypted SK, decrypt password
*/
func Decrypt(text, mySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(mySecret))
	if err != nil {
		panic(err)
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
