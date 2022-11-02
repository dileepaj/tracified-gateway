// SPDX-License-Identifier: MIT

pragma solidity ^0.8.7;

contract UrbanWaterUsage {
	uint public Result;
	bytes32 formulaID;
	bytes32 expertPK;
	bytes32 tenetID;
	uint public WATER;
	uint public WATER_TO_ELECTRICITY_UNIT;
	uint public ELECTRICITY_UNIT_TO_CARBON_EMISSION;

	function setMetadata(bytes32 _formulaID, bytes32 _expertPK, bytes32 _tenetID) public {
		formulaID = _formulaID;
		expertPK = _expertPK;
		tenetID = _tenetID;
	}

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
		Result = (WATER * WATER_TO_ELECTRICITY_UNIT * (ELECTRICITY_UNIT_TO_CARBON_EMISSION * (WATER_TO_ELECTRICITY_UNIT + ELECTRICITY_UNIT_TO_CARBON_EMISSION)));
	}

	function getResult() public view returns (uint) {
		return Result;
	}
	function getFormulaID() public view returns (bytes32) {
		return formulaID;
	}

	function getExpertPK() public view returns (bytes32) {
		return expertPK;
	}

	function getTenetID() public view returns (bytes32) {
		return tenetID;
	}

}