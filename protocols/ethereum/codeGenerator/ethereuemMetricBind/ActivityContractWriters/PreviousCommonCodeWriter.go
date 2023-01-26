package activitywriters

import "github.com/dileepaj/tracified-gateway/model"

// PreviousCommonCodeWriter writes solidity code for variable declaration and getter for the previous contract address
// The previous contract address is the address of the contract that was deployed before the current contract with the same metric name (it will be retrieved from the DB collection EthMetricLatest)

func WritePreviousCommonCode(metricID string) (model.PreviousCode, error) {

	// TODO: get the previous contract address from the DB
	previousContractAddress, err := getPreviousContractAddress(metricID)
	if err != nil {
		return model.PreviousCode{}, err
	}

	// variable declaration and initialization
	deceleration := "\t" + `string private previousContractAddress = "` + previousContractAddress + `";	//previous contract address` + "\n"

	// getter
	getter := "\t" + `function getPreviousContractAddress() public view returns (string memory) {` + "\n"
	getter += "\t\t" + `return previousContractAddress;` + "\n"
	getter += "\t" + `}` + "\n"

	codesForPreviousContract := model.PreviousCode{
		Setter: deceleration,
		Getter: getter,
	}

	return codesForPreviousContract, nil
}
