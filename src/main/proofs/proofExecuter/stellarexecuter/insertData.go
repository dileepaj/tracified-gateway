package stellarexecuter

import (
	"fmt"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
)

func InsertDataHash(hash string, secret string, public string) string {
	source := "SDL26B3CQN4AQHPV3MDRMUB5BXNMCQLHY3HVAD7ZOP4QACX2OL7V2IOW"
	destination := "GDLTXGZ2XHVDTNGKEXSRSEVNJEJJG6ENOBTE2CILFVWSSFV3ECYZDCTI"

	// Make sure destination account exists
	if _, err := horizon.DefaultTestNetClient.LoadAccount(destination); err != nil {
		panic(err)
	}

	//passphrase := network.TestNetworkPassphrase

	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{source},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.Payment(
			build.Destination{destination},
			build.NativeAmount{"10"},
		),
	)

	if err != nil {
		panic(err)
	}

	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(source)
	if err != nil {
		panic(err)
	}

	txeB64, err := txe.Base64()
	if err != nil {
		panic(err)
	}

	// And finally, send it off to Stellar!
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successful Transaction:")
	fmt.Println("Ledger:", resp.Ledger)
	fmt.Println("Hash:", resp.Hash)

	return "hellow"

}

// func AppendToTree(current string, pre string) string {

// }
