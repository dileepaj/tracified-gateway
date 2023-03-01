package services

func CheckContractStatus() {
	//TODO call the DB and filter out the transaction with pending status
	//TODO loop through the transactions and call the ethereum endpoint to check the transaction status
	//TODO check the status
	/*
		TODO
			if pending - Update the index
			if success - update in the DB collection as completed
			if failed - log the error and mark the status as failed
	*/
}
