// SPDX-License-Identifier: MIT

pragma solidity >=0.7.0 <0.9.0;

contract UrbanWaterUsage {
	uint public UrbanWaterUsage;
	uint public WATER;
	uint public WATER_TO_ELECTRICITY_UNIT;
	uint public ELECTRICITY_UNIT_TO_CARBON_EMISSION;

	function Executor() public {
		UrbanWaterUsage = WATER * WATER_TO_ELECTRICITY_UNIT * ELECTRICITY_UNIT_TO_CARBON_EMISSION;
	}
}