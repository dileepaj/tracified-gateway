package authentication

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"github.com/sirupsen/logrus"
)

func PGPValidator(sha256hash string, signature []byte, originalMsg string) (error, bool) {
	//object := dao.Connection{}

	//TODO: get the public key from the DB
	// publicKeyDet, errWhenGettingPublicKey := object.GetRSAPublicKeyBySHA256PK(sha256hash).Then(func(data interface{}) interface{} {
	// 	return data
	// }).Await()
	// if errWhenGettingPublicKey != nil {
	// 	logrus.Error("Error when getting the public key from gateway datastore " + errWhenGettingPublicKey.Error())
	// 	return errors.New("Error when getting the public key from gateway datastore " + errWhenGettingPublicKey.Error()), false
	// }
	// if publicKeyDet == nil {
	// 	logrus.Error("Public key does not exist in the gateway datastore for the hash ", sha256hash)
	// 	return errors.New("Public key does not exist in the gateway datastore for the hash " + sha256hash), false
	// }
	// publicKeyData := publicKeyDet.(model.RSAPublickey)
	// publicKey = publicKeyData.Publickey
	// logrus.Info("Public key : ", publicKey)
	// pgpPubKey := "LS0tLS1CRUdJTiBQR1AgUFVCTElDIEtFWSBCTE9DSy0tLS0tClZlcnNpb246IEtleWJhc2UgT3BlblBHUCB2Mi4wLjc2CkNvbW1lbnQ6IGh0dHBzOi8va2V5YmFzZS5pby9jcnlwdG8KCnhzQk5CR05XZGNRQkNBRGR4ZzZRNFA1QzhyeUxhc1Z0UE5IcUFCcUF5YzZkNkxjcUgzNGplWFFWRlNKM0JsSjYKdkE4SEVNRzRNdkhKZnRoM2dETzJ1c1hnRUFEMExkVG0vMTl3bFVWK3Q0QWNFd0R1UkVzY2ZTUmRjUHZDUnpvZwp1Y3VtcTI1Y3l1ZXRBK05MYzlnU2pkMHkvRjIyZ2hSVlQ2MGVXQnRVQWRGZkhQdHZ1RTYvK1pUeDlBeUdhUTVBCm9kQVh1ZkRBVjlmdWFKaFlQQU9zMVMvckwySUpMTHRQZW10L1NuZU42VElhbW5lSFl4d29TS1pZZDNlVTd1QTEKZ2J4NW9YUG12LzBwbWJRYU9LWXByMFcrMjlRTTVjWllGV0w1ck5FeHVKRzJPdnBRN09xcUNGdDlBSG1UY0VENgpWeTdUb1g5aEdYd00yK0V6N0N4elBZKzgxZVl4VmhmM09FSVBBQkVCQUFITkowSm9ZWGRoYm5Sb1lTQndZWE5oCmJpQThjMmhoYm5WcllTNWljSE5BWjIxaGFXd3VZMjl0UHNMQWVnUVRBUW9BSkFVQ1kxWjF4QUliTHdNTENRY0QKRlFvSUFoNEJBaGVBQXhZQ0FRSVpBUVVKQUFBQUFBQUtDUkJjdXQzYVNHR3RWOWw0Qi8wUkh4bWFuT1g1Tk9GZgp4Um1WNFYzQmxuOW1qbjBrbDJQaTFrK29nSU5KOVU4bDVnS3hHclJFN2lMbElUR1BiYXZCYUVFMTFwQ0prOUFICmFkZHlacXZId1ltb3p6bmtEeUNwYW1CdjFlc21ObTNkd01iNkVMVE1USFFmcmFjZ1ZQNHJMeE5ybFQ2R2RwMkEKdnNhd3ZIQXYvSlpvUDhQMEUrQ2FQM081VCtXVTcyWE1vVzNkU3A2K3l5ZHlMNUpyZUo3UnhlUGhNL3Y3SERMRQpGZDlGS3JKR01DNVd1RWg0TWxVUjRMd0xjWjU3QWwzMWljR1ZLZnhDZnFLSjlCMzNqZWN2SWY3ZGxrcjl1TWRkCmU2LzYxNlJDa3dSODhvb2VqL0VETFlpcm9DcTNaUTAramJjSFVRWmVnR0hRaTZlbEF0UE1TMzY3dkdaRkdTeTEKbDhaSlVlOVJ6c0JOQkdOV2RjUUJDQURYdlpTdlgxakRHMlVaN0IydVl1NmVDcjNBcXlwd3REbVVYbjBwQkkyUQpyTTV5L1VnRzlaVHVQWnpIUGVBTEdnUWl5WXlsMEZncFdpQXJaRGhPVWFTeFVGdWtqbDRvTTNKMURIMlI0cjMzCk1QS1ZxMnF1dWlRWVZDS2djcHErZDRlYllldzViSUlSS2VtVFJlYTVXZDIrU1VSWTBFTG1ZNlNWa0gvRFJ3VkMKQmVCMGtRZGp0TVhSRi9aalMwbG1HeWxTa2Y4a1BJODdIajk4d1ZFKzl4WGxJZWlGRDhzV1FpanJ2QWFkUHNFNQpWelFMTDBvb3F2Sys4U2JIR2NkZlRHVG9IZ3VMb08xMXFJSk1nanl6S2s3b0dma1NEZU8wdnpGMzdZaVduNXVmCjV0Tk9yUTRmZVJQQ1VBTTA0T253dkZRbDZMZldZeGg2UkR1eXNiVnRhYm9OQUJFQkFBSEN3WVFFR0FFS0FBOEYKQW1OV2RjUUZDUUFBQUFBQ0d5NEJLUWtRWExyZDJraGhyVmZBWFNBRUdRRUtBQVlGQW1OV2RjUUFDZ2tRTkR3ZQpwWFd2U0tteXNBZitPSzIrbi90TURBcnNTNEIyU0ZQZ1YrbmxpRVhOYU9HL2tTTW1iUjBvbW9VTG50K0tFRnRECldodnNKR21obXdjeFJJR3YwNC9aUmoxVk1sZlNkWEtCOGlyYWhjQ0tXTmQ2V29LcXJWK0RsbmJOZGtZQkovM3UKdVdHNzRpazAyOWRrdjFnYUtGTnlyRHNSOTdyLzd5bzFRTDJ1TEYvaEdVY0J3ZVNnWUMyN0pFMXpWTjQ3TTF1RApZZEFCcmw4QkQxeWZFVlAvS01na1V1bDFmait5blU1NkpvY2kvUkx1K3Z5djFQemtFUUhjTkFtSmRaSENjbWZKCm9XM3F2am96ZkQ3NGQ4WXRNWmlPaXowSFQ3TStjdHBJcDUxR3YzY1dZVnVoK1NNOCtNL2dScnpMaTFCZTNsYXIKNExwNVZ1RENkM2I1dzZlaGhGdFdIenJpWGZIRjUrWE5ITEtIQ0FDRFdxbFFHeGlHWGZlamhObDN5UHVaRHlsagpTUFVVcHlaZ2k4TDMzS1FzNXo4Y1NWU2szWjN2MXhHOUFtcHNpdXR6TE1RYkdFeE01clFOUGQ3cHo1UnJaV1lBCnE4Tm9ueStBa20vQzFLOUZlVkNiYUF3cE15RWtjbkxKeCtETFRGQ0FVRnFDQ2w3b3BiVjI0Y2JiSlRlcUEvMkoKNzZDck9JWmY1bitFaC83WXJ0RFZrY3V0Z3pGbG8vSmo1dno5dEVlQnB2ajNNVUNtRWJzb3JZb3FUbHpRQjVpeAp0MW4yTkY1RXJib0t6SHN0cEhsaWxoTlN1QUs1SXkxaWJjWWMvdlBpU2lQWVlhNFhHYmFQWnJoMWNtZWlKUSs1Clk2dWdBT2FjS0JHTUlPUDJCeXJHUFgzditKU1VUdEVBU2Zrd1dLRnJmQWoyeVRPQWlYcGt5Qm13VGhFYXpzQk4KQkdOV2RjUUJDQUMwZUxXTi9wOE9JbkwzQXErckRicWZkSWJkdXRmUDdZNUY3TWhKYnVrTEpLVk00ZW1NMTdKVAoyMDNSUmxDalJxY0VUeklaRHE0N0hLdGNSMVJPYmVxYTJXY0poMURoY3JqSms5SlZSNzNsVzdDWno2Z0dmZllsCjZXcXpnZU1Ra3o0bVJEOUR0L3ZVeHV4MTZFb1g1TXl1UittaGZCdlVPYnJkczhQQWZtS0wrWHBRTlBBcXVIaTYKVGU5N0VvMW9teWwzcTlweiswcHljaGFzN1REMThrUjVVc24zRldkTTE3cmRDaWZRZTBXS1Yyd3p2RmxtcEVPTApvdG9aZFZMY01MVE96Z0dISi9EUXpsOHk2dkQ5MlRMMy9FbHhIM3RvU3pCZHA1K3Vwa0o4aFBzVENMWmVoWUY0CmFNeVk2ekNVZFdPTEdEdlRzYkFjMnpaSWFxMjRkNXZqQUJFQkFBSEN3WVFFR0FFS0FBOEZBbU5XZGNRRkNRQUEKQUFBQ0d5NEJLUWtRWExyZDJraGhyVmZBWFNBRUdRRUtBQVlGQW1OV2RjUUFDZ2tRb1NLa3hYUDRHb1lDM1FnQQprVnZ6K2g1Uy9yOUpSZUdFTk9FU05WdXl1UmMwWkNWOWtvR205NFExV3JIS1ZTbWp1TEE0SG4rQ3d5TDYraFNwCnVUVjV3M0FTK2FBdXhCbzZHYVExZFJ2c1RSc0x4TmhIQ3RQQTJ3bG9LTjJHTVB0SmEzN251cmt6VEdIWHBuTDIKcmYvRy8yZC9qTnB0WEIydHNzRWQyb3VYMEJ6ZVdDRlA5NksrTStSQzhEZ3hUZk1mUnlHblJ5TUFmdy9oclJicwpiL055dGIvV2l1VC9DMEJxM1laRzFibTF2SFN1SGI4WkxrRTFMY1VJREtVaUx1bkx4YUFRWGFCZWl2YVNlNlhvCk85anZWTGl0MWdrMlVNS1RJMVpOZmt2dVZzZkx1SE1FbHczUklqOEY0djZUQlptd1ZybkJpQzBpRWpkUHE4SlkKK3lRUHovSEJZMHJKdDEra2tQT2NVZXdGQ0FDV3dRaHd6RXZHaEtNOW1PTWgzbEl1RURqamd1bGpZVjNKYzlDOApCbitZVGpPU01ndkEycEY3UjIrVzBjSFhnN25wRHNyMjR0M0RYeVBXdFhyQWVJYVJ0ZWpKdGhnT2pSUlVTR0o4CnBFZDI3enhjLzIwQUN0Mm1RL282K1g1QzlwWmR5L3dyZGtkY1lURm4xVXREVkhucEdXS0Vkd1dRN0FNTEU2Z3kKWmxtQXhTRUprdmovT2x2Zy9uYk9XVGYzL012TTJiODNQcy9NVWpEWWEyRzU2dEt6K2lsckYxbTFsK0RsN0lJegpMdFRybndYenhzOEpTNlZ0OFlzcmI1eXM5aWQzV29XY0FyN24xSGpRZGUzRkg1R0pTMGJmRDJpVjU2aHBmOHBNCmJHekdHTURoa1Y1MjRXQ0YraWVQY3psM1Q3Q3lzaVJzVnBWUzQ0N1JDdjM1cTV0VAo9VEY3OQotLS0tLUVORCBQR1AgUFVCTElDIEtFWSBCTE9DSy0tLS0t"
	// b64, _ := base64.StdEncoding.DecodeString(pgpPubKey)

	//Generate key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		logrus.Error(err)
	}
	secretMessage := "Superrrrrrrrrr secret"
	signatureTxt := CreateSignature(secretMessage, *privateKey)

	//convert to pem string
	publicKeyPemString, _ := ExportRsaPublicKeyAsPemStr(&privateKey.PublicKey)

	//base64 convert the public key
	base64PublicKey := base64.StdEncoding.EncodeToString([]byte(publicKeyPemString))

	//-------------------------------------------------------------------------------
	decodedPublicKey, errWhenDecodongBase64 := base64.StdEncoding.DecodeString(base64PublicKey)

	if errWhenDecodongBase64 != nil {
		logrus.Error()
	}

	parsingPublicKey := string(decodedPublicKey[:])

	pub_parsed, _ := ParseRsaPublicKeyFromPemStr(parsingPublicKey) //pass the decoded base 64 value to this method

	hashed := sha256.Sum256([]byte(secretMessage))

	errInVerifier := rsa.VerifyPKCS1v15(pub_parsed, crypto.SHA256, hashed[:], signatureTxt)
	if errInVerifier != nil {
		logrus.Error("Verification failed " + errInVerifier.Error())
		return errors.New("Verification failed :" + errInVerifier.Error()), false
	}

	return nil, true

}