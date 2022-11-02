// SPDX-License-Identifier: MIT

pragma solidity ^0.8.7;

contract UrbanWaterUsage {
	uint public Result;
	uint public WATER;
	uint public WATER_TO_ELECTRICITY_UNIT;
	uint public ELECTRICITY_UNIT_TO_CARBON_EMISSION;

	function setWATER(uint _WATER) public {
		WATER = _WATER;
	}

	function setWATER_TO_ELECTRICITY_UNIT(uint _WATER_TO_ELECTRICITY_UNIT) public {
		WATER_TO_ELECTRICITY_UNIT = _WATER_TO_ELECTRICITY_UNIT;
	}

	function setELECTRICITY_UNIT_TO_CARBON_EMISSION(uint _ELECTRICITY_UNIT_TO_CARBON_EMISSION) public {
		ELECTRICITY_UNIT_TO_CARBON_EMISSION = _ELECTRICITY_UNIT_TO_CARBON_EMISSION;
	}


	function Executor() public {
		Result = WATER * WATER_TO_ELECTRICITY_UNIT * ELECTRICITY_UNIT_TO_CARBON_EMISSION;
	}

	function getResult() public view returns (uint) {
		return Result;
	}
}