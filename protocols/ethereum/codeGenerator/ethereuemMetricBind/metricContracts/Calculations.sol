//SPDX-License-Identifier: Unlicensed

pragma solidity ^0.8.7;

contract Calculations {
    int256 exponent = 0;
    //Multiplication
    function Multiply(
        int256 _opOneValue,
        int256 _opOneExponent,
        int256 _opTwoValue,
        int256 _opTwoExponent
    ) public returns (int256) {
        //multiply values
        int256 valueRes = _opOneValue * _opTwoValue;

        //add exponents
        int256 exponentRes = _opOneExponent + _opTwoExponent;

        exponent = exponentRes;

        return valueRes;
    }

    //Division
    //opOneValue - numerator value
    //opOneExponent - numerator exponent
    //opTwoValue - denominator value
    //opTwoExponent - denominator exponent
    function Divide(
        int256 _opOneValue,
        int256 _opOneExponent,
        int256 _opTwoValue,
        int256 _opTwoExponent
    ) public returns (int256) {
        //multiply numerator value from 10^6
        int256 opOneValuePowValue = _opOneValue * 1000000;

        //divide the values
        int256 valueRes = opOneValuePowValue / _opTwoValue;

        //calculate the exponent result
        //6 represents the value value of 1000000
        //!Note - this will give the value outcome with the 10^6 multiplication involved, should be handled in the calling contract
        int256 exponentRes = _opOneExponent - _opTwoExponent - 6;

        exponent = exponentRes;

        return valueRes;
    }

    //Subtract
    function Subtract(
        int256 _opOneValue,
        int256 _opOneExponent,
        int256 _opTwoValue,
        int256 _opTwoExponent
    ) public returns (int256) {
        int256 valueRes;
        int256 valueConvertedInBetween;
        int256 exponentRes;
        int256 exponentDif;
        //check the maximum exponent value
        if (_opOneExponent > _opTwoExponent) {
            exponentDif = _opOneExponent - _opTwoExponent;
            valueConvertedInBetween =
                _opOneValue *
                int256(10**uint256(exponentDif));
            valueRes = valueConvertedInBetween - _opTwoValue;
            exponentRes = _opTwoExponent;
        } else if (_opOneExponent < _opTwoExponent) {
            exponentDif = _opTwoExponent - _opOneExponent;
            valueConvertedInBetween =
                _opTwoValue *
                int256(10**uint256(exponentDif));
            valueRes = _opOneValue - valueConvertedInBetween;
            exponentRes = _opOneExponent;
        } else {
            valueRes = _opOneValue - _opTwoValue;
            exponentRes = _opOneExponent;
        }

        exponent = exponentRes;

        return valueRes;
    }

    //Addition
    function Add(
        int256 _opOneValue,
        int256 _opOneExponent,
        int256 _opTwoValue,
        int256 _opTwoExponent
    ) public returns (int256) {
        int256 valueRes;
        int256 valueConvertedInBetween;
        int256 exponentRes;
        int256 exponentDif;
        //check the maximum exponent value
        if (_opOneExponent > _opTwoExponent) {
            exponentDif = _opOneExponent - _opTwoExponent;
            valueConvertedInBetween =
                _opOneValue *
                int256(10**uint256(exponentDif));
            valueRes = valueConvertedInBetween + _opTwoValue;
            exponentRes = _opTwoExponent;
        } else if (_opOneExponent < _opTwoExponent) {
            exponentDif = _opTwoExponent - _opOneExponent;
            valueConvertedInBetween =
                _opTwoValue *
                int256(10**uint256(exponentDif));
            valueRes = _opOneValue + valueConvertedInBetween;
            exponentRes = _opOneExponent;
        } else {
            valueRes = _opOneValue + _opTwoValue;
            exponentRes = _opOneExponent;
        }

        exponent = exponentRes;

        return valueRes;
    }

    function GetExponent() public view returns (int256) {
        return exponent;
    }
}