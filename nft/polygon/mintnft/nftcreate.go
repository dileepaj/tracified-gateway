package mintnft

// import (
// 	"context"
// 	"crypto/ecdsa"
// 	"fmt"

// 	"log"
// 	"math/big"

// 	"github.com/dileepaj/tracified-gateway/nft/polygon/TracifiedNFT"
// 	//"github.com/dileepaj/tracified-gateway/nft/polygon/polygon"

// 	"github.com/dileepaj/tracified-gateway/nft/polygon/polygonconnection"
// 	"github.com/ethereum/go-ethereum/accounts/abi/bind"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/crypto"
// 	//"github.com/jaynagpaul/go-web3"
// 	//"github.com/hrharder/go-gas"
// )

// func ecdsaPrivateKeyFromHex(privateKeyHex string) *ecdsa.PrivateKey {
// 	ecdsaKey, err := crypto.HexToECDSA(privateKeyHex)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return ecdsaKey
// }

// func MintNFT() {
// 	// println("lesssssss ggooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooo")
// 	// //w3 := web3.NewWeb3()
// 	// account1 := common.HexToAddress("0x7d612CA535521777895eA3a5F73F38e06d5E7549")
// 	// account2 := common.HexToAddress("0xb0f63fbB5c24dbAe40bc5cDF8E823226aC66E717")

// 	// /*b1, err := polygonconnection.CallPolygonClient().BalanceAt(context.TODO(), account1, nil)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// b2, err := polygonconnection.CallPolygonClient().BalanceAt(context.Background(), account2, nil)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// fmt.Println("Balance 1:", b1)
// 	// fmt.Println("Balance 2:", b2)
// 	// */
// 	// nonce, err := polygonconnection.CallPolygonClient().PendingNonceAt(context.Background(), account1)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// amount := big.NewInt(1000000000000000000) // wei
// 	// //gasprice := big.NewInt(20000000000)

// 	// gasPrice, err := polygonconnection.CallPolygonClient().SuggestGasPrice(context.Background())
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// /*gasPrice, err := gas.SuggestGasPrice(gas.GasPriorityFastest)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }*/

// 	// contractAddress := common.HexToAddress("0x1b43258Ceb9e9cB40aa8c9cbc158205C22a6F45a")

// 	// t, err := TracifiedNFT.NewTracifiedNFT(contractAddress, polygonconnection.CallPolygonClient())

// 	// instance, err := TracifiedNFT.NewTracifiedNFT(contractAddress, polygonconnection.CallPolygonClient())
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// tokenURI := "Test"

// 	// tx := types.NewTransaction(nonce, account2, amount, 21000, gasPrice, nil)

// 	// chainID, err := polygonconnection.CallPolygonClient().NetworkID(context.Background())
// 	// if err != nil {
// 	// 	log.Fatal("Error : ", err)
// 	// 	log.Fatal(t, tx)
// 	// }

// 	// privateKey1, err := crypto.HexToECDSA("18ad5d8357e418199fcd9dbdfa810380441c6dff269a7a516f8b515e3d0c2aa2")
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// privateKey := "18ad5d8357e418199fcd9dbdfa810380441c6dff269a7a516f8b515e3d0c2aa2"
// 	// x, err := bind.NewKeyedTransactorWithChainID(ecdsaPrivateKeyFromHex(privateKey), chainID)
// 	// if err != nil {
// 	// 	log.Fatal("Error : ", err)
// 	// 	log.Fatal(t, tx)
// 	// }

// 	// tx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey1)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// err = polygonconnection.CallPolygonClient().SendTransaction(context.Background(), tx)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// fmt.Printf("tx sent: %s\n", tx.Hash().Hex())

// 	// x.GasLimit = 500000
// 	// x.GasPrice = gasPrice

// 	// version, err := instance.MintNFT(&bind.TransactOpts{}, account1, tokenURI)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// fmt.Println("Version", version)
// }
