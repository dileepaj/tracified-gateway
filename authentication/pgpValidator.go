package authentication

import (
	"errors"

	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/jchavannes/go-pgp/pgp"
	"github.com/sirupsen/logrus"
)

func PGPValidator(sha256hash string, encryptedMsg string, originalMsg string) (error, bool) {
	object := dao.Connection{}
	var publicKey string

	//TODO: get the public key from the DB
	publicKeyDet, errWhenGettingPublicKey := object.GetRSAPublicKeyBySHA256PK(sha256hash).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if errWhenGettingPublicKey != nil {
		logrus.Error("Error when getting the public key from gateway datastore " + errWhenGettingPublicKey.Error())
		return errors.New("Error when getting the public key from gateway datastore " + errWhenGettingPublicKey.Error()), false
	}
	if publicKeyDet == nil {
		logrus.Error("Public key does not exist in the gateway datastore for the hash ", sha256hash)
		return errors.New("Public key does not exist in the gateway datastore for the hash " + sha256hash), false
	}
	publicKeyData := publicKeyDet.(model.RSAPublickey)
	publicKey = publicKeyData.Publickey
	logrus.Info("Public key : ", publicKey)

	pgpPubKey := `----BEGIN PGP PUBLIC KEY BLOCK----- Version: Keybase OpenPGP v2.0.76 Comment: https://keybase.io/crypto xo0EY1Y1ngEEAMsCN9dUlhpn4tdqxUWl5oeAgxZ5HuD+IiP3UQb1mRLU/0lq9wQrcCSAM5qN3eP6hhEbfUOe4qIJJuLI/orsaGziYCVsYU45OKPeimPP7gzur8wynFp6YebupyRkl04oRqwsV2T9mI18PYtzCpeuAvUo9sRYLAYvr2fRdS8O0e6tABEBAAHNI0lzaGluaSAgPGlzaGluaS5raXJpZGVuYUBnbWFpbC5jb20+wroEEwEKACQFAmNWNZ4CGy8DCwkHAxUKCAIeAQIXgAMWAgECGQEFCQHhM4AACgkQCkY1a2czFk1IUQP/ZL19OnAbTeS77wocLzoV3F6e35LaT+EUIHf7n70aS74oXp+MqHxHFrfx6PCiChBlSkoYnROAW4BzjgK3kvw3KBMBGCju4Th7vavwPBRad4LFDqLkxWF3BG77nMI4oDbGY/ehuxyXRSTjdGyXyO6WeHqUB+0tW921x6TjzS1FIvDOjQRjVjWeAQQAzBWTDhA6qZ9WCdTC37VgXZ/yNWiBFlSsb/cUDqjVkOY8mab78fmHh9sHwv8IxAti/yMBa93uxOOzadf0eXops3fGADi815LaNUROfqiqnKOkdd+jDv1p9iEcIXEpfy6plyzqgKgZdWFFDxCxu6yXXMLRSXDcjXiYjPZ6Qw6eMRUAEQEAAcLAgwQYAQoADwUCY1Y1ngUJAeEzgAIbLgCoCRAKRjVrZzMWTZ0gBBkBCgAGBQJjVjWeAAoJEFuAhL/MU7/i7TIEAL1G+vt8RJlPK0WIXva0rwrCILy8eSwcopCkB8xXpERzCbe2DTNQVINPRYrK3+IRuERAQ09j7ZVzwoYmgWMLefMNVoGgbSq5b51pPkdbjOdtr04HN+61jmpmeHuO2M05ru0fChRLG9fOfHwLyUJvS8eFRbkQ3YsqNpbmq6/y3bBHSpgEAJ3SIHTiza8mZjNy/DHU81HDg5R3hf5SeBZxRjMRaAr3xTncSHaBDJJZFMmHodqixZ/LAN+/enQy7HXUqrUb34zDcVi3Kk3mDUCtpGQMTtm/GLuJc4RjKulEm0OPz69EsbiVuibjOp54g1A5mFjlRbUizrvhrHi3M8mp3BQaqSWMzo0EY1Y1ngEEAJ0VMhLwNEYNYPxSMghBd7K1njOl+3N2NGtByrXDg/Nng70omwjBw1aLzTsOG9pV6AuTdSDlBGh1FNRpFKr9KAYmle8TpLYfK64ivgcZpr6Q0gkEA3jnW1+9PmI46/sbL1QSzZLWuoJQV3utKw3AnHJeZ1uP+j5cayIvFqs/Ph7JABEBAAHCwIMEGAEKAA8FAmNWNZ4FCQHhM4ACGy4AqAkQCkY1a2czFk2dIAQZAQoABgUCY1Y1ngAKCRA3ddcv/y/mSwDQA/9SjHleayhMc3SuQdB6oBHjJ6Z4mCzdjQeWWx8nQ6nHPniDXhhfVaOiB05DM3Iz8Nq7cGcWeXJ3o+cx+O0B1nSj6nD1U9IgfWozmGCZliKX8pn57nQHk/p65yBv62kV17fOZVVsERsxxpDT7TdSNMaMt16WhoOB1EDjZnkNGpnJ8tyKA/0QT5fm3IYq5F8nUeUiCHBZ+j9M4aMBLyIJee6dAAoAa9QM5o/F1+okjwzEVOU0nmf/YE6pdkMNSYs2aVbOn3mDWLIBRNI7kHvdEm3uKO2MjOweVJjbumRjJFtg4queOJpz/+fZPwmnLCDOQhYQ/1a1GA68WzpN0zPvu+09FCv3Yw===wmKm -----END PGP PUBLIC KEY BLOCK----`

	//export the keys to pem string
	// pub_parsed, errWhenGettingpubparsed := ParseRsaPublicKeyFromPemStr(pgpPubKey)
	// if errWhenGettingpubparsed != nil {
	// 	logrus.Error("Error whrn getting the public key parser " + errWhenGettingpubparsed.Error())
	// 	return errors.New("Error whrn getting the public key parser " + errWhenGettingpubparsed.Error()), false
	// }

	//create signature
	// signature := CreateSignature(encryptedMsg, pub_parsed)

	//-------------------------------------------------------------------------------------------------------------------

	//TODO: decrypt the string
	publicEntity, errWhenGettingEntity := pgp.GetEntity([]byte(pgpPubKey), []byte{})
	if errWhenGettingEntity != nil {
		logrus.Error("Error when getting the entitiy from public key " + errWhenGettingEntity.Error())
		return errors.New("Error when getting the entitiy from public key " + errWhenGettingEntity.Error()), false
	}

	decryptedMsg, errIndecrypt := pgp.Decrypt(publicEntity, []byte(encryptedMsg))
	if errIndecrypt != nil {
		logrus.Error("Error when decrytping the message " + errIndecrypt.Error())
		return errors.New("Error when decrytping the message " + errIndecrypt.Error()), false
	}

	decryptedMessageString := string(decryptedMsg)
	logrus.Info("Decrypted Message : ", decryptedMessageString)

	//TODO: verify the string
	if decryptedMessageString != originalMsg {
		logrus.Error("PGP validation failed, incorrect key pairs")
		return errors.New("PGP validation failed, incorrect key pairs"), false
	}

	return nil, true

}
