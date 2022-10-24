package fclqueryexecuter

import (
	"errors"

	fcl "github.com/shanukabps/FCL-GO_test"
)

func FCLQueryToExecutionTempalteJsonString(query string) (string, error) {
	if query == "" {
		return "", errors.New("query can not be empty")
	}
	 var executionTemplateString string= fcl.NewFCLWrapper().GetExecutionTemplateJSONString("protocols/stellarprotocols/FCLQueryExecuter/Defs.txt",query)
	return executionTemplateString, nil
}
