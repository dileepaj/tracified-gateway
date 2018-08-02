package stellarExecuter

import (
	"log"

	"github.com/stellar/go/keypair"
)

func CreateAccount() bool {
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(pair.Seed())
	// SAV76USXIJOBMEQXPANUOQM6F5LIOTLPDIDVRJBFFE2MDJXG24TAPUU7
	log.Println(pair.Address())
	// GCFXHS4GXL6BVUCXBWXGTITROWLVYXQKQLF4YH5O5JT3YZXCYPAFBJZB
	return true
}
