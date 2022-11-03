package executionTemplates

import (
	"github.com/dileepaj/tracified-gateway/model"
)

// returns -> start variable declarations, setter list, built equation, error
func CommandBuilder(command model.Command) (string, error) {

	var commandString string = ""

	// check the command type and get the operator as a string
	if command.Ul_CommandType == 2100 {
		commandString = " + "
	} else if command.Ul_CommandType == 2101 {
		commandString = " - "
	} else if command.Ul_CommandType == 2102 {
		commandString = " * "
	} else if command.Ul_CommandType == 2103 {
		commandString = " / "
	} else if command.Ul_CommandType == 10000 {
		commandString = " * "
	}

	// check the whether the command has argument or not and call relevant function
	if command.P_Arg.S_StartVarName != "" {
		if command.P_Arg.Lst_Commands != nil {
			strTemplate, _ := Template1Builder(command.P_Arg)
			commandString += strTemplate
		} else {
			strTemplate, _ := Template2Builder(command.P_Arg)
			commandString += strTemplate
		}
	}

	return commandString, nil
}