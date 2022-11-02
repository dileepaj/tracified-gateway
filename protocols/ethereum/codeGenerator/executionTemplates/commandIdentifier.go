package executionTemplates

import (
	"github.com/dileepaj/tracified-gateway/model"
	"github.com/dileepaj/tracified-gateway/protocols/ethereum/components"
)

func CommandBuilder(command model.Command) (string, error) {
	var commandString string = ""
	
	str, errInCommandString := components.GenerateCommandType(command.Ul_CommandType) 
	if errInCommandString != nil {
		return "", errInCommandString
	}
	commandString = commandString + str

	// check the command type and generate the command type

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