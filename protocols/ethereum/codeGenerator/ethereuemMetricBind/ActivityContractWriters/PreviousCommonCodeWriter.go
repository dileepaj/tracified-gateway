package activitywriters

// PreviousCommonCodeWriter writes solidity code for variable declaration and getter for the previous contract address
// The previous contract address is the address of the contract that was deployed before the current contract with the same metric name (it will be retrieved from the DB collection EthMetricLatest)

func WritePreviousCommonCode(metricID string) (string, error) {
	codesForPreviousContract := ""

	// TODO: get the previous contract address from the DB
	previousContractAddress, err := getPreviousContractAddress(metricID)
	if err != nil {
		return "", err
	}

	// variable declaration and initialization
	codesForPreviousContract += "\t" + `string private previousContractAddress = "`+ previousContractAddress + `";	//previous contract address` + "\n"

	// getter
	codesForPreviousContract += "\t" + `function getPreviousContractAddress() public view returns (string memory) {` + "\n"
	codesForPreviousContract += "\t\t" + `return previousContractAddress;` + "\n"
	codesForPreviousContract += "\t" + `}` + "\n"

	return codesForPreviousContract, nil
}