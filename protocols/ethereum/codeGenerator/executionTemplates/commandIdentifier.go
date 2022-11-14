package executionTemplates

import (
	"github.com/dileepaj/tracified-gateway/model"
)

// returns -> start variable declarations, setter list, built equation, error
func CommandBuilder(command model.Command) (string, string, error) {

	var commandStringStart string = ""
	var commandStringEnd string = ", "

	// check the command type and get the operator as a string
	if command.Ul_CommandType == 2100 {
		commandStringStart = "Add("
	} else if command.Ul_CommandType == 2101 {
		commandStringStart = "Subtract("
	} else if command.Ul_CommandType == 2102 {
		commandStringStart = "Multiply("
	} else if command.Ul_CommandType == 2103 {
		commandStringStart = "Divide("
	} else if command.Ul_CommandType == 10000 {
		commandStringStart = "Multiply("
	}

	// check the whether the command has argument or not and call relevant function
	if command.P_Arg.S_StartVarName != "" {
		if command.P_Arg.Lst_Commands != nil {
			strTemplate, _ := Template1Builder(command.P_Arg)
			commandStringEnd = commandStringEnd + strTemplate
		} else {
			strTemplate, _ := Template2Builder(command.P_Arg)
			commandStringEnd = commandStringEnd + strTemplate
		}
		commandStringEnd = commandStringEnd + ")"
	}

	return commandStringStart, commandStringEnd, nil
}