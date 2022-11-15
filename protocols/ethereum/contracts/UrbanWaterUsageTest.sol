// SPDX-License-Identifier: MIT

pragma solidity ^0.8.7;

import './Calculations.sol';

contract UrbanWaterUsageTest {
	
	// Result structure
	struct Result {
		int256 value;
		int256 exponent;
	}
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
		int256 value; //initialize at start, added using setter
		int256 exponent; //initialize at start, added using setter
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
	Metadata metadata = Metadata("631a0b4ad9241a9374f9f035","Urban water usage test","GBOK5LA4FABVOK76XOZGYHXDF5XYMS5FQJ7XZNP6TTAG6NYS766TO7EF");

	// Result initialization
	Result result = Result(0, 0);

	// Calculation object creation
	Calculations calculations = new Calculations();

	// Value initializations
	// value initialization for VARIABLE -> WATER
	Variable WATER = Variable(Value("VARIABLE", "e1223109481cf739d19e6735d0236577", "water", 0, 0, ""), "l", 0);
	// value initialization for REFERREDCONSTANT -> WATER_TO_ELECTRICITY_UNIT
	ReferredConstant WATER_TO_ELECTRICITY_UNIT = ReferredConstant(Value("REFERREDCONSTANT", "faa0754b7eddd458180b1d36ff9c0494", "water to e. unit", 20052, -2, ""), "kWh", "Reference URL");
	// value initialization for SEMANTICCONSTANT -> ELECTRICITY_UNIT_TO_CARBON_EMISSION
	SemanticConstant ELECTRICITY_UNIT_TO_CARBON_EMISSION = SemanticConstant(Value("SEMANTICCONSTANT", "faa0754b7eddd458180b1d36ff4c0492", "e. unit to carb. em.", 15, 0, ""));

	// value setter for VARIABLE WATER
	function setWATER(int256 _WATER, int256 _EXPONENT) public {
		WATER.value.value = _WATER;
		WATER.value.exponent = _EXPONENT;
	}

	// method to get the result of the calculation
	function executeCalculation() public returns (int256, int256) {
		result.value = calculations.Multiply(calculations.Multiply(WATER.value.value, WATER.value.exponent, WATER_TO_ELECTRICITY_UNIT.value.value, WATER_TO_ELECTRICITY_UNIT.value.exponent), calculations.GetExponent(), calculations.Multiply(ELECTRICITY_UNIT_TO_CARBON_EMISSION.value.value, ELECTRICITY_UNIT_TO_CARBON_EMISSION.value.exponent, calculations.Add(WATER_TO_ELECTRICITY_UNIT.value.value, WATER_TO_ELECTRICITY_UNIT.value.exponent, ELECTRICITY_UNIT_TO_CARBON_EMISSION.value.value, ELECTRICITY_UNIT_TO_CARBON_EMISSION.value.exponent), calculations.GetExponent()), calculations.GetExponent());
		result.exponent = calculations.GetExponent();
		
		return (result.value, result.exponent);
	}
}