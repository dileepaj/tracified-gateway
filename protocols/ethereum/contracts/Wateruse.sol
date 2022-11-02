// SPDX-License-Identifier: MIT

pragma solidity ^0.8.7;

contract Wateruse {
	uint public Wateruse;
	uint public UrbanWaterUsage;
	uint public WATER;
	uint public WATER_TO_ELECTRICITY_UNIT;
	uint public ELECTRICITY_UNIT_TO_CARBON_EMISSION;

	function Executor() public {
		Wateruse = WATER * WATER_TO_ELECTRICITY_UNIT * ELECTRICITY_UNIT_TO_CARBON_EMISSION;
	}
}