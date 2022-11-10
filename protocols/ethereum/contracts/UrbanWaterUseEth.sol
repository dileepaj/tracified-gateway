// SPDX-License-Identifier: MIT

pragma solidity ^0.8.7;

contract UrbanWaterUseEth {
	int public result = -9999;
	//Metadata structure
	struct Metadata {
		bytes32 formulaID; //initialize at start
		string formulaName; //initialize at start
		string expertPK; //initialize at start
	}
	//Parent value structure
	struct Value {
		string valueType; //initialize at start
		bytes32 valueID; //initialize at start
		string valueName; //initialize at start
		int value; //initialize at start, added using setter
		string description; //initialize at start
	}
	//Variable structure, child of Value
	struct Variable {
		Value value; //initialize at start
		bytes32 unit; //initialize at start
		bytes32 precision; //initialize at start
	}
	//Semantic constant structure, child of Value
	struct SemanticConstant {
		Value value; //initialize at start
	}
	//Referred constant structure, child of Value
	struct ReferredConstant {
		Value value; //initialize at start
		bytes32 unit; //initialize at start
		string refUrl; //initialize at start
	}
	//Metadata declaration
	Metadata metadata = Metadata("631a0b4ad9241a9374f9f001EthTest1","Urban water use ETH","GBOK5LA4FABVOK76XOZGYHXDF5XYMS5FQJ7XZNP6TTAG6NYS766TO7EF");

	// Value initializations
	// value initialization for VARIABLE -> WATER
	Variable WATER = Variable(Value("VARIABLE", "e1223109481cf739d19e6735d0236577", "water", 0, ""), "l", 0);
	// value initialization for REFERREDCONSTANT -> WATER_TO_ELECTRICITY_UNIT
	ReferredConstant WATER_TO_ELECTRICITY_UNIT = ReferredConstant(Value("REFERREDCONSTANT", "faa0754b7eddd458180b1d36ff9c0494", "water to e. unit", 200.000000, ""), "kWh", "Reference URL");
	// value initialization for SEMANTICCONSTANT -> ELECTRICITY_UNIT_TO_CARBON_EMISSION
	SemanticConstant ELECTRICITY_UNIT_TO_CARBON_EMISSION = SemanticConstant(Value("SEMANTICCONSTANT", "faa0754b7eddd458180b1d36ff4c0492", "e. unit to carb. em.", 15.000000, ""));

	// value setter for VARIABLE WATER
	function setWATER(int _WATER) public {
		WATER.value.value = _WATER;
	}

	// method to get the result of the calculation
	function executeCalculation() public returns (int) {	
		if (result == -9999) {
			result = (WATER.value.value * WATER_TO_ELECTRICITY_UNIT.value.value * (ELECTRICITY_UNIT_TO_CARBON_EMISSION.value.value * (WATER_TO_ELECTRICITY_UNIT.value.value + ELECTRICITY_UNIT_TO_CARBON_EMISSION.value.value)));
		}
		return result;
	}
}