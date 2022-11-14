package stellar

import (
	log "github.com/sirupsen/logrus"

	"github.com/dileepaj/tracified-gateway/commons"
	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/dao"
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/sirupsen/logrus"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

/*IssueNft
@desc - Issue NFT with the newly generated issuer credentials
@params - Current Issuer public key, Distributor public key, asset code, TDP hash
*/
func IssueNft(CurrentIssuerPK string, distributerPK string, assetcode string, svg string) (string, string, error) {
	object := dao.Connection{}
	//getting the issuer secret key
	data, err := object.GetNFTIssuerSK(CurrentIssuerPK).Then(func(data interface{}) interface{} {
		return data
	}).Await()
	if err != nil {
		log.Fatal(err)
		return "", "", err
	} else {
		nftKeys := data.([]model.NFTKeys)
		//decrypt the secret key
		decrpytNftissuerSecretKey := commons.Decrypt(nftKeys[0].SecretKey)
		if data == nil {
			logrus.Error("PublicKey is not found in gateway datastore")
			panic(data)
		}
		var nftConten = []byte(svg)

		request := horizonclient.AccountRequest{AccountID: CurrentIssuerPK}
		issuerAccount, err := commons.GetHorizonNetwork().AccountDetail(request)
		if err != nil {
			return "", "", err
		}
		issuerSign, err := keypair.ParseFull(decrpytNftissuerSecretKey)
		if err != nil {
			return "", "", err
		}

		var homeDomain string = commons.GoDotEnvVariable("HOMEDOMAIN")

		payments := []txnbuild.Operation{
			&txnbuild.Payment{Destination: distributerPK, Asset: txnbuild.CreditAsset{Code: assetcode,
				Issuer: CurrentIssuerPK}, Amount: "1"},
			&txnbuild.ManageData{
				Name:  assetcode,
				Value: nftConten,
			},
			&txnbuild.SetOptions{
				MasterWeight: txnbuild.NewThreshold(0),
				HomeDomain:   &homeDomain,
			},
		}
		tx, err := txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        &issuerAccount,
				IncrementSequenceNum: true,
				Operations:           payments,
				BaseFee:              constants.MinBaseFee,
				Memo:                 nil,
				Preconditions:        txnbuild.Preconditions{TimeBounds: constants.TransactionTimeOut},
			},
		)
		if err != nil {
			log.Fatal(err)
			return "", "", err
		}

		signedTx, err := tx.Sign(network.TestNetworkPassphrase, issuerSign)
		if err != nil {
			return "", "", err
		}
		// submit transaction
		respn, err := commons.GetHorizonClient().SubmitTransaction(signedTx)
		if err != nil {
			log.Fatal("Error submitting transaction:", err)
			panic(err)
		}
		return respn.Hash, string(nftConten), nil
	}

}
