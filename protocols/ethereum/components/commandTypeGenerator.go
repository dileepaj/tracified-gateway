package components

func GenerateCommandType(commandType uint32) (string, error) {

	commandTypeString := ""

	// 2100-add, 2101-subtract, 2102-multiply, 2103-divide
	if commandType == 2100 {
		commandTypeString = " + "
	} else if commandType == 2101 {
		commandTypeString = " - "
	} else if commandType == 2102 {
		commandTypeString = " * "
	} else if commandType == 2103 {
		commandTypeString = " / "
	} else if commandType == 10000 {
		commandTypeString = " * "
	}

	return commandTypeString, nil
}