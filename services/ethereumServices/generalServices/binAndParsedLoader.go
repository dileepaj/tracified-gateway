package generalservices

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/sirupsen/logrus"
)

//Load the contract BIN and ABI as an interface
func LoadContractBinAndParsedAbi(bin string, abi string) (string, *abi.ABI, error) {

	var ContractBIN string

	//assign metadata for the contract
	var BuildData = &bind.MetaData{
		ABI: abi,
		Bin: bin,
	}
	ContractBIN = BuildData.Bin

	parsed, errWhenGettingABI := BuildData.GetAbi()
	if errWhenGettingABI != nil {
		logrus.Error("Error when getting abi from passed ABI string " + errWhenGettingABI.Error())
		return ContractBIN, parsed, errWhenGettingABI
	}

	return ContractBIN, parsed, nil

}
