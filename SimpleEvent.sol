// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleEvent {
    uint256 public number;
    event NumberSet(uint256 newNumber);

    function setNumber(uint256 _number) public {
        number = _number;
        emit NumberSet(_number);
    }
}