// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package build

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// BuildMetaData contains all meta data concerning the Build contract.
var BuildMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_opOneValue\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"_opOneExponent\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"_opTwoValue\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"_opTwoExponent\",\"type\":\"int256\"}],\"name\":\"Add\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_opOneValue\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"_opOneExponent\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"_opTwoValue\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"_opTwoExponent\",\"type\":\"int256\"}],\"name\":\"Divide\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetExponent\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_opOneValue\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"_opOneExponent\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"_opTwoValue\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"_opTwoExponent\",\"type\":\"int256\"}],\"name\":\"Multiply\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"_opOneValue\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"_opOneExponent\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"_opTwoValue\",\"type\":\"int256\"},{\"internalType\":\"int256\",\"name\":\"_opTwoExponent\",\"type\":\"int256\"}],\"name\":\"Subtract\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526000805534801561001457600080fd5b50610798806100246000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c8063107dd6271461005c5780631bcbc5971461008c578063355e7a8e146100bc5780637815e9e0146100ec578063ad170b501461011c575b600080fd5b61007660048036038101906100719190610382565b61013a565b60405161008391906103f8565b60405180910390f35b6100a660048036038101906100a19190610382565b6101f7565b6040516100b391906103f8565b60405180910390f35b6100d660048036038101906100d19190610382565b61022c565b6040516100e391906103f8565b60405180910390f35b61010660048036038101906101019190610382565b610281565b60405161011391906103f8565b60405180910390f35b61012461033e565b60405161013191906103f8565b60405180910390f35b6000806000806000858813156101885785886101569190610442565b905080600a61016591906105c2565b89610170919061060d565b9250868361017e9190610442565b93508591506101e1565b858812156101ce57878661019c9190610442565b905080600a6101ab91906105c2565b876101b6919061060d565b925082896101c49190610442565b93508791506101e0565b86896101da9190610442565b93508791505b5b8160008190555083945050505050949350505050565b6000808386610206919061060d565b9050600083866102169190610685565b9050806000819055508192505050949350505050565b600080620f42408661023e919061060d565b90506000848261024e91906106f8565b90506000600685886102609190610442565b61026a9190610442565b905080600081905550819350505050949350505050565b6000806000806000858813156102cf57858861029d9190610442565b905080600a6102ac91906105c2565b896102b7919061060d565b925086836102c59190610685565b9350859150610328565b858812156103155787866102e39190610442565b905080600a6102f291906105c2565b876102fd919061060d565b9250828961030b9190610685565b9350879150610327565b86896103219190610685565b93508791505b5b8160008190555083945050505050949350505050565b60008054905090565b600080fd5b6000819050919050565b61035f8161034c565b811461036a57600080fd5b50565b60008135905061037c81610356565b92915050565b6000806000806080858703121561039c5761039b610347565b5b60006103aa8782880161036d565b94505060206103bb8782880161036d565b93505060406103cc8782880161036d565b92505060606103dd8782880161036d565b91505092959194509250565b6103f28161034c565b82525050565b600060208201905061040d60008301846103e9565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061044d8261034c565b91506104588361034c565b925082820390508181126000841216828213600085121516171561047f5761047e610413565b5b92915050565b60008160011c9050919050565b6000808291508390505b60018511156104dc578086048111156104b8576104b7610413565b5b60018516156104c75780820291505b80810290506104d585610485565b945061049c565b94509492505050565b6000826104f557600190506105b1565b8161050357600090506105b1565b8160018114610519576002811461052357610552565b60019150506105b1565b60ff84111561053557610534610413565b5b8360020a91508482111561054c5761054b610413565b5b506105b1565b5060208310610133831016604e8410600b84101617156105875782820a90508381111561058257610581610413565b5b6105b1565b6105948484846001610492565b925090508184048111156105ab576105aa610413565b5b81810290505b9392505050565b6000819050919050565b60006105cd826105b8565b91506105d8836105b8565b92506106057fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff84846104e5565b905092915050565b60006106188261034c565b91506106238361034c565b92508282026106318161034c565b91507f8000000000000000000000000000000000000000000000000000000000000000841460008412161561066957610668610413565b5b828205841483151761067e5761067d610413565b5b5092915050565b60006106908261034c565b915061069b8361034c565b9250828201905082811215600083121683821260008412151617156106c3576106c2610413565b5b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60006107038261034c565b915061070e8361034c565b92508261071e5761071d6106c9565b5b600160000383147f80000000000000000000000000000000000000000000000000000000000000008314161561075757610756610413565b5b82820590509291505056fea264697066735822122074a84dc35e92a771ee2ff003f000e82e551ea538ee87db96ec93767fba48630f64736f6c63430008110033",
}

// BuildABI is the input ABI used to generate the binding from.
// Deprecated: Use BuildMetaData.ABI instead.
var BuildABI = BuildMetaData.ABI

// BuildBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BuildMetaData.Bin instead.
var BuildBin = BuildMetaData.Bin

// DeployBuild deploys a new Ethereum contract, binding an instance of Build to it.
func DeployBuild(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Build, error) {
	parsed, err := BuildMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BuildBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Build{BuildCaller: BuildCaller{contract: contract}, BuildTransactor: BuildTransactor{contract: contract}, BuildFilterer: BuildFilterer{contract: contract}}, nil
}

// Build is an auto generated Go binding around an Ethereum contract.
type Build struct {
	BuildCaller     // Read-only binding to the contract
	BuildTransactor // Write-only binding to the contract
	BuildFilterer   // Log filterer for contract events
}

// BuildCaller is an auto generated read-only Go binding around an Ethereum contract.
type BuildCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BuildTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BuildTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BuildFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BuildFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BuildSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BuildSession struct {
	Contract     *Build            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BuildCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BuildCallerSession struct {
	Contract *BuildCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// BuildTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BuildTransactorSession struct {
	Contract     *BuildTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BuildRaw is an auto generated low-level Go binding around an Ethereum contract.
type BuildRaw struct {
	Contract *Build // Generic contract binding to access the raw methods on
}

// BuildCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BuildCallerRaw struct {
	Contract *BuildCaller // Generic read-only contract binding to access the raw methods on
}

// BuildTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BuildTransactorRaw struct {
	Contract *BuildTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBuild creates a new instance of Build, bound to a specific deployed contract.
func NewBuild(address common.Address, backend bind.ContractBackend) (*Build, error) {
	contract, err := bindBuild(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Build{BuildCaller: BuildCaller{contract: contract}, BuildTransactor: BuildTransactor{contract: contract}, BuildFilterer: BuildFilterer{contract: contract}}, nil
}

// NewBuildCaller creates a new read-only instance of Build, bound to a specific deployed contract.
func NewBuildCaller(address common.Address, caller bind.ContractCaller) (*BuildCaller, error) {
	contract, err := bindBuild(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BuildCaller{contract: contract}, nil
}

// NewBuildTransactor creates a new write-only instance of Build, bound to a specific deployed contract.
func NewBuildTransactor(address common.Address, transactor bind.ContractTransactor) (*BuildTransactor, error) {
	contract, err := bindBuild(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BuildTransactor{contract: contract}, nil
}

// NewBuildFilterer creates a new log filterer instance of Build, bound to a specific deployed contract.
func NewBuildFilterer(address common.Address, filterer bind.ContractFilterer) (*BuildFilterer, error) {
	contract, err := bindBuild(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BuildFilterer{contract: contract}, nil
}

// bindBuild binds a generic wrapper to an already deployed contract.
func bindBuild(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BuildABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Build *BuildRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Build.Contract.BuildCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Build *BuildRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Build.Contract.BuildTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Build *BuildRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Build.Contract.BuildTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Build *BuildCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Build.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Build *BuildTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Build.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Build *BuildTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Build.Contract.contract.Transact(opts, method, params...)
}

// GetExponent is a free data retrieval call binding the contract method 0xad170b50.
//
// Solidity: function GetExponent() view returns(int256)
func (_Build *BuildCaller) GetExponent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Build.contract.Call(opts, &out, "GetExponent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetExponent is a free data retrieval call binding the contract method 0xad170b50.
//
// Solidity: function GetExponent() view returns(int256)
func (_Build *BuildSession) GetExponent() (*big.Int, error) {
	return _Build.Contract.GetExponent(&_Build.CallOpts)
}

// GetExponent is a free data retrieval call binding the contract method 0xad170b50.
//
// Solidity: function GetExponent() view returns(int256)
func (_Build *BuildCallerSession) GetExponent() (*big.Int, error) {
	return _Build.Contract.GetExponent(&_Build.CallOpts)
}

// Add is a paid mutator transaction binding the contract method 0x7815e9e0.
//
// Solidity: function Add(int256 _opOneValue, int256 _opOneExponent, int256 _opTwoValue, int256 _opTwoExponent) returns(int256)
func (_Build *BuildTransactor) Add(opts *bind.TransactOpts, _opOneValue *big.Int, _opOneExponent *big.Int, _opTwoValue *big.Int, _opTwoExponent *big.Int) (*types.Transaction, error) {
	return _Build.contract.Transact(opts, "Add", _opOneValue, _opOneExponent, _opTwoValue, _opTwoExponent)
}

// Add is a paid mutator transaction binding the contract method 0x7815e9e0.
//
// Solidity: function Add(int256 _opOneValue, int256 _opOneExponent, int256 _opTwoValue, int256 _opTwoExponent) returns(int256)
func (_Build *BuildSession) Add(_opOneValue *big.Int, _opOneExponent *big.Int, _opTwoValue *big.Int, _opTwoExponent *big.Int) (*types.Transaction, error) {
	return _Build.Contract.Add(&_Build.TransactOpts, _opOneValue, _opOneExponent, _opTwoValue, _opTwoExponent)
}

// Add is a paid mutator transaction binding the contract method 0x7815e9e0.
//
// Solidity: function Add(int256 _opOneValue, int256 _opOneExponent, int256 _opTwoValue, int256 _opTwoExponent) returns(int256)
func (_Build *BuildTransactorSession) Add(_opOneValue *big.Int, _opOneExponent *big.Int, _opTwoValue *big.Int, _opTwoExponent *big.Int) (*types.Transaction, error) {
	return _Build.Contract.Add(&_Build.TransactOpts, _opOneValue, _opOneExponent, _opTwoValue, _opTwoExponent)
}

// Divide is a paid mutator transaction binding the contract method 0x355e7a8e.
//
// Solidity: function Divide(int256 _opOneValue, int256 _opOneExponent, int256 _opTwoValue, int256 _opTwoExponent) returns(int256)
func (_Build *BuildTransactor) Divide(opts *bind.TransactOpts, _opOneValue *big.Int, _opOneExponent *big.Int, _opTwoValue *big.Int, _opTwoExponent *big.Int) (*types.Transaction, error) {
	return _Build.contract.Transact(opts, "Divide", _opOneValue, _opOneExponent, _opTwoValue, _opTwoExponent)
}

// Divide is a paid mutator transaction binding the contract method 0x355e7a8e.
//
// Solidity: function Divide(int256 _opOneValue, int256 _opOneExponent, int256 _opTwoValue, int256 _opTwoExponent) returns(int256)
func (_Build *BuildSession) Divide(_opOneValue *big.Int, _opOneExponent *big.Int, _opTwoValue *big.Int, _opTwoExponent *big.Int) (*types.Transaction, error) {
	return _Build.Contract.Divide(&_Build.TransactOpts, _opOneValue, _opOneExponent, _opTwoValue, _opTwoExponent)
}

// Divide is a paid mutator transaction binding the contract method 0x355e7a8e.
//
// Solidity: function Divide(int256 _opOneValue, int256 _opOneExponent, int256 _opTwoValue, int256 _opTwoExponent) returns(int256)
func (_Build *BuildTransactorSession) Divide(_opOneValue *big.Int, _opOneExponent *big.Int, _opTwoValue *big.Int, _opTwoExponent *big.Int) (*types.Transaction, error) {
	return _Build.Contract.Divide(&_Build.TransactOpts, _opOneValue, _opOneExponent, _opTwoValue, _opTwoExponent)
}

// Multiply is a paid mutator transaction binding the contract method 0x1bcbc597.
//
// Solidity: function Multiply(int256 _opOneValue, int256 _opOneExponent, int256 _opTwoValue, int256 _opTwoExponent) returns(int256)
func (_Build *BuildTransactor) Multiply(opts *bind.TransactOpts, _opOneValue *big.Int, _opOneExponent *big.Int, _opTwoValue *big.Int, _opTwoExponent *big.Int) (*types.Transaction, error) {
	return _Build.contract.Transact(opts, "Multiply", _opOneValue, _opOneExponent, _opTwoValue, _opTwoExponent)
}

// Multiply is a paid mutator transaction binding the contract method 0x1bcbc597.
//
// Solidity: function Multiply(int256 _opOneValue, int256 _opOneExponent, int256 _opTwoValue, int256 _opTwoExponent) returns(int256)
func (_Build *BuildSession) Multiply(_opOneValue *big.Int, _opOneExponent *big.Int, _opTwoValue *big.Int, _opTwoExponent *big.Int) (*types.Transaction, error) {
	return _Build.Contract.Multiply(&_Build.TransactOpts, _opOneValue, _opOneExponent, _opTwoValue, _opTwoExponent)
}

// Multiply is a paid mutator transaction binding the contract method 0x1bcbc597.
//
// Solidity: function Multiply(int256 _opOneValue, int256 _opOneExponent, int256 _opTwoValue, int256 _opTwoExponent) returns(int256)
func (_Build *BuildTransactorSession) Multiply(_opOneValue *big.Int, _opOneExponent *big.Int, _opTwoValue *big.Int, _opTwoExponent *big.Int) (*types.Transaction, error) {
	return _Build.Contract.Multiply(&_Build.TransactOpts, _opOneValue, _opOneExponent, _opTwoValue, _opTwoExponent)
}

// Subtract is a paid mutator transaction binding the contract method 0x107dd627.
//
// Solidity: function Subtract(int256 _opOneValue, int256 _opOneExponent, int256 _opTwoValue, int256 _opTwoExponent) returns(int256)
func (_Build *BuildTransactor) Subtract(opts *bind.TransactOpts, _opOneValue *big.Int, _opOneExponent *big.Int, _opTwoValue *big.Int, _opTwoExponent *big.Int) (*types.Transaction, error) {
	return _Build.contract.Transact(opts, "Subtract", _opOneValue, _opOneExponent, _opTwoValue, _opTwoExponent)
}

// Subtract is a paid mutator transaction binding the contract method 0x107dd627.
//
// Solidity: function Subtract(int256 _opOneValue, int256 _opOneExponent, int256 _opTwoValue, int256 _opTwoExponent) returns(int256)
func (_Build *BuildSession) Subtract(_opOneValue *big.Int, _opOneExponent *big.Int, _opTwoValue *big.Int, _opTwoExponent *big.Int) (*types.Transaction, error) {
	return _Build.Contract.Subtract(&_Build.TransactOpts, _opOneValue, _opOneExponent, _opTwoValue, _opTwoExponent)
}

// Subtract is a paid mutator transaction binding the contract method 0x107dd627.
//
// Solidity: function Subtract(int256 _opOneValue, int256 _opOneExponent, int256 _opTwoValue, int256 _opTwoExponent) returns(int256)
func (_Build *BuildTransactorSession) Subtract(_opOneValue *big.Int, _opOneExponent *big.Int, _opTwoValue *big.Int, _opTwoExponent *big.Int) (*types.Transaction, error) {
	return _Build.Contract.Subtract(&_Build.TransactOpts, _opOneValue, _opOneExponent, _opTwoValue, _opTwoExponent)
}
