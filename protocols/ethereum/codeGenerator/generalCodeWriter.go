package codeGenerator

import "github.com/dileepaj/tracified-gateway/model"

func GeneralCodeWriter(contractName string) model.ContractGeneral {
	generalBuilder := model.ContractGeneral{
		License:          `// SPDX-License-Identifier: MIT`,
		StartingCodeLine: `pragma solidity ^0.8.7;`,
		ContractStart:    `contract ` + contractName + ` {`,
		ContractEnd:      `}`,
	}

	return generalBuilder
}
