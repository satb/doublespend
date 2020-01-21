// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package erc20

//import (
//	"math/big"
//	"strings"
//
//	ethereum "github.com/ethereum/go-ethereum"
//	"github.com/ethereum/go-ethereum/accounts/abi"
//	"github.com/ethereum/go-ethereum/accounts/abi/bind"
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/core/types"
//	"github.com/ethereum/go-ethereum/event"
//)
//
//// Reference imports to suppress errors if they are not otherwise used.
//var (
//	_ = big.NewInt
//	_ = strings.NewReader
//	_ = ethereum.NotFound
//	_ = abi.U256
//	_ = bind.Bind
//	_ = common.Big1
//	_ = types.BloomLookup
//	_ = event.NewSubscription
//)
//
//// Erc20ABI is the input ABI used to generate the binding from.
//const Erc20ABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
//
//// Erc20Bin is the compiled bytecode used for deploying new contracts.
//var Erc20Bin = "0x60806040526040518060400160405280600581526020017f48454c4c4f0000000000000000000000000000000000000000000000000000008152506000908051906020019061004f929190610196565b506040518060400160405280600381526020017f484c4f00000000000000000000000000000000000000000000000000000000008152506001908051906020019061009b929190610196565b506012600260006101000a81548160ff021916908360ff160217905550600260009054906101000a900460ff1660ff16600a0a620f4240026003553480156100e257600080fd5b50600354600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055503373ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef6003546040518082815260200191505060405180910390a361023b565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106101d757805160ff1916838001178555610205565b82800160010185558215610205579182015b828111156102045782518255916020019190600101906101e9565b5b5090506102129190610216565b5090565b61023891905b8082111561023457600081600090555060010161021c565b5090565b90565b610acc8061024a6000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c806370a082311161005b57806370a08231146101d857806395d89b4114610230578063dd62ed3e146102b3578063e1f21c671461032b57610088565b806306fdde031461008d57806318160ddd1461011057806323b872dd1461012e578063313ce567146101b4575b600080fd5b6100956103b1565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100d55780820151818401526020810190506100ba565b50505050905090810190601f1680156101025780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61011861044f565b6040518082815260200191505060405180910390f35b61019a6004803603606081101561014457600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610455565b604051808215151515815260200191505060405180910390f35b6101bc610479565b604051808260ff1660ff16815260200191505060405180910390f35b61021a600480360360208110156101ee57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061048c565b6040518082815260200191505060405180910390f35b6102386104d5565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561027857808201518184015260208101905061025d565b50505050905090810190601f1680156102a55780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b610315600480360360408110156102c957600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610573565b6040518082815260200191505060405180910390f35b6103976004803603606081101561034157600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506105fa565b604051808215151515815260200191505060405180910390f35b60008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156104475780601f1061041c57610100808354040283529160200191610447565b820191906000526020600020905b81548152906001019060200180831161042a57829003601f168201915b505050505081565b60035481565b60006104628484846107f8565b61046d8433846105fa565b50600190509392505050565b600260009054906101000a900460ff1681565b6000600460008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b60018054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561056b5780601f106105405761010080835404028352916020019161056b565b820191906000526020600020905b81548152906001019060200180831161054e57829003601f168201915b505050505081565b6000600560008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b60008073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161415610681576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526024815260200180610a736024913960400191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610707576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526022815260200180610a2c6022913960400191505060405180910390fd5b81600560008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040518082815260200191505060405180910390a3600190509392505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16141561087e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180610a4e6025913960400191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610904576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526023815260200180610a096023913960400191505060405180910390fd5b80600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254039250508190555080600460008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040518082815260200191505060405180910390a350505056fe45524332303a207472616e7366657220746f20746865207a65726f206164647265737345524332303a20617070726f766520746f20746865207a65726f206164647265737345524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f2061646472657373a2646970667358221220164fbd7a8d67fce2b5ead914c043a7cbe210fc6a01496f8bf83a987184f21b5a64736f6c63430006010033"
//
//// DeployErc20 deploys a new Ethereum contract, binding an instance of Erc20 to it.
//func DeployErc20(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Erc20, error) {
//	parsed, err := abi.JSON(strings.NewReader(Erc20ABI))
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//
//	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(Erc20Bin), backend)
//	if err != nil {
//		return common.Address{}, nil, nil, err
//	}
//	return address, tx, &Erc20{Erc20Caller: Erc20Caller{contract: contract}, Erc20Transactor: Erc20Transactor{contract: contract}, Erc20Filterer: Erc20Filterer{contract: contract}}, nil
//}
//
//// Erc20 is an auto generated Go binding around an Ethereum contract.
//type Erc20 struct {
//	Erc20Caller     // Read-only binding to the contract
//	Erc20Transactor // Write-only binding to the contract
//	Erc20Filterer   // Log filterer for contract events
//}
//
//// Erc20Caller is an auto generated read-only Go binding around an Ethereum contract.
//type Erc20Caller struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// Erc20Transactor is an auto generated write-only Go binding around an Ethereum contract.
//type Erc20Transactor struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// Erc20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
//type Erc20Filterer struct {
//	contract *bind.BoundContract // Generic contract wrapper for the low level calls
//}
//
//// Erc20Session is an auto generated Go binding around an Ethereum contract,
//// with pre-set call and transact options.
//type Erc20Session struct {
//	Contract     *Erc20            // Generic contract binding to set the session for
//	CallOpts     bind.CallOpts     // Call options to use throughout this session
//	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
//}
//
//// Erc20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
//// with pre-set call options.
//type Erc20CallerSession struct {
//	Contract *Erc20Caller  // Generic contract caller binding to set the session for
//	CallOpts bind.CallOpts // Call options to use throughout this session
//}
//
//// Erc20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
//// with pre-set transact options.
//type Erc20TransactorSession struct {
//	Contract     *Erc20Transactor  // Generic contract transactor binding to set the session for
//	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
//}
//
//// Erc20Raw is an auto generated low-level Go binding around an Ethereum contract.
//type Erc20Raw struct {
//	Contract *Erc20 // Generic contract binding to access the raw methods on
//}
//
//// Erc20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
//type Erc20CallerRaw struct {
//	Contract *Erc20Caller // Generic read-only contract binding to access the raw methods on
//}
//
//// Erc20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
//type Erc20TransactorRaw struct {
//	Contract *Erc20Transactor // Generic write-only contract binding to access the raw methods on
//}
//
//// NewErc20 creates a new instance of Erc20, bound to a specific deployed contract.
//func NewErc20(address common.Address, backend bind.ContractBackend) (*Erc20, error) {
//	contract, err := bindErc20(address, backend, backend, backend)
//	if err != nil {
//		return nil, err
//	}
//	return &Erc20{Erc20Caller: Erc20Caller{contract: contract}, Erc20Transactor: Erc20Transactor{contract: contract}, Erc20Filterer: Erc20Filterer{contract: contract}}, nil
//}
//
//// NewErc20Caller creates a new read-only instance of Erc20, bound to a specific deployed contract.
//func NewErc20Caller(address common.Address, caller bind.ContractCaller) (*Erc20Caller, error) {
//	contract, err := bindErc20(address, caller, nil, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &Erc20Caller{contract: contract}, nil
//}
//
//// NewErc20Transactor creates a new write-only instance of Erc20, bound to a specific deployed contract.
//func NewErc20Transactor(address common.Address, transactor bind.ContractTransactor) (*Erc20Transactor, error) {
//	contract, err := bindErc20(address, nil, transactor, nil)
//	if err != nil {
//		return nil, err
//	}
//	return &Erc20Transactor{contract: contract}, nil
//}
//
//// NewErc20Filterer creates a new log filterer instance of Erc20, bound to a specific deployed contract.
//func NewErc20Filterer(address common.Address, filterer bind.ContractFilterer) (*Erc20Filterer, error) {
//	contract, err := bindErc20(address, nil, nil, filterer)
//	if err != nil {
//		return nil, err
//	}
//	return &Erc20Filterer{contract: contract}, nil
//}
//
//// bindErc20 binds a generic wrapper to an already deployed contract.
//func bindErc20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
//	parsed, err := abi.JSON(strings.NewReader(Erc20ABI))
//	if err != nil {
//		return nil, err
//	}
//	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_Erc20 *Erc20Raw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _Erc20.Contract.Erc20Caller.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_Erc20 *Erc20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _Erc20.Contract.Erc20Transactor.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_Erc20 *Erc20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _Erc20.Contract.Erc20Transactor.contract.Transact(opts, method, params...)
//}
//
//// Call invokes the (constant) contract method with params as input values and
//// sets the output to result. The result type might be a single field for simple
//// returns, a slice of interfaces for anonymous returns and a struct for named
//// returns.
//func (_Erc20 *Erc20CallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
//	return _Erc20.Contract.contract.Call(opts, result, method, params...)
//}
//
//// Transfer initiates a plain transaction to move funds to the contract, calling
//// its default method if one is available.
//func (_Erc20 *Erc20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _Erc20.Contract.contract.Transfer(opts)
//}
//
//// Transact invokes the (paid) contract method with params as input values.
//func (_Erc20 *Erc20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
//	return _Erc20.Contract.contract.Transact(opts, method, params...)
//}
//
//// Allowance is a paid mutator transaction binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(address owner, address spender) returns(uint256)
//func (_Erc20 *Erc20Transactor) Allowance(opts *bind.TransactOpts, owner common.Address, spender common.Address) (*types.Transaction, error) {
//	return _Erc20.contract.Transact(opts, "allowance", owner, spender)
//}
//
//// Allowance is a paid mutator transaction binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(address owner, address spender) returns(uint256)
//func (_Erc20 *Erc20Session) Allowance(owner common.Address, spender common.Address) (*types.Transaction, error) {
//	return _Erc20.Contract.Allowance(&_Erc20.TransactOpts, owner, spender)
//}
//
//// Allowance is a paid mutator transaction binding the contract method 0xdd62ed3e.
////
//// Solidity: function allowance(address owner, address spender) returns(uint256)
//func (_Erc20 *Erc20TransactorSession) Allowance(owner common.Address, spender common.Address) (*types.Transaction, error) {
//	return _Erc20.Contract.Allowance(&_Erc20.TransactOpts, owner, spender)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0xe1f21c67.
////
//// Solidity: function approve(address owner, address spender, uint256 amount) returns(bool)
//func (_Erc20 *Erc20Transactor) Approve(opts *bind.TransactOpts, owner common.Address, spender common.Address, amount *big.Int) (*types.Transaction, error) {
//	return _Erc20.contract.Transact(opts, "approve", owner, spender, amount)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0xe1f21c67.
////
//// Solidity: function approve(address owner, address spender, uint256 amount) returns(bool)
//func (_Erc20 *Erc20Session) Approve(owner common.Address, spender common.Address, amount *big.Int) (*types.Transaction, error) {
//	return _Erc20.Contract.Approve(&_Erc20.TransactOpts, owner, spender, amount)
//}
//
//// Approve is a paid mutator transaction binding the contract method 0xe1f21c67.
////
//// Solidity: function approve(address owner, address spender, uint256 amount) returns(bool)
//func (_Erc20 *Erc20TransactorSession) Approve(owner common.Address, spender common.Address, amount *big.Int) (*types.Transaction, error) {
//	return _Erc20.Contract.Approve(&_Erc20.TransactOpts, owner, spender, amount)
//}
//
//// BalanceOf is a paid mutator transaction binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(address _owner) returns(uint256 balance)
//func (_Erc20 *Erc20Transactor) BalanceOf(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
//	return _Erc20.contract.Transact(opts, "balanceOf", _owner)
//}
//
//// BalanceOf is a paid mutator transaction binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(address _owner) returns(uint256 balance)
//func (_Erc20 *Erc20Session) BalanceOf(_owner common.Address) (*types.Transaction, error) {
//	return _Erc20.Contract.BalanceOf(&_Erc20.TransactOpts, _owner)
//}
//
//// BalanceOf is a paid mutator transaction binding the contract method 0x70a08231.
////
//// Solidity: function balanceOf(address _owner) returns(uint256 balance)
//func (_Erc20 *Erc20TransactorSession) BalanceOf(_owner common.Address) (*types.Transaction, error) {
//	return _Erc20.Contract.BalanceOf(&_Erc20.TransactOpts, _owner)
//}
//
//// Decimals is a paid mutator transaction binding the contract method 0x313ce567.
////
//// Solidity: function decimals() returns(uint8)
//func (_Erc20 *Erc20Transactor) Decimals(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _Erc20.contract.Transact(opts, "decimals")
//}
//
//// Decimals is a paid mutator transaction binding the contract method 0x313ce567.
////
//// Solidity: function decimals() returns(uint8)
//func (_Erc20 *Erc20Session) Decimals() (*types.Transaction, error) {
//	return _Erc20.Contract.Decimals(&_Erc20.TransactOpts)
//}
//
//// Decimals is a paid mutator transaction binding the contract method 0x313ce567.
////
//// Solidity: function decimals() returns(uint8)
//func (_Erc20 *Erc20TransactorSession) Decimals() (*types.Transaction, error) {
//	return _Erc20.Contract.Decimals(&_Erc20.TransactOpts)
//}
//
//// Name is a paid mutator transaction binding the contract method 0x06fdde03.
////
//// Solidity: function name() returns(string)
//func (_Erc20 *Erc20Transactor) Name(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _Erc20.contract.Transact(opts, "name")
//}
//
//// Name is a paid mutator transaction binding the contract method 0x06fdde03.
////
//// Solidity: function name() returns(string)
//func (_Erc20 *Erc20Session) Name() (*types.Transaction, error) {
//	return _Erc20.Contract.Name(&_Erc20.TransactOpts)
//}
//
//// Name is a paid mutator transaction binding the contract method 0x06fdde03.
////
//// Solidity: function name() returns(string)
//func (_Erc20 *Erc20TransactorSession) Name() (*types.Transaction, error) {
//	return _Erc20.Contract.Name(&_Erc20.TransactOpts)
//}
//
//// Symbol is a paid mutator transaction binding the contract method 0x95d89b41.
////
//// Solidity: function symbol() returns(string)
//func (_Erc20 *Erc20Transactor) Symbol(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _Erc20.contract.Transact(opts, "symbol")
//}
//
//// Symbol is a paid mutator transaction binding the contract method 0x95d89b41.
////
//// Solidity: function symbol() returns(string)
//func (_Erc20 *Erc20Session) Symbol() (*types.Transaction, error) {
//	return _Erc20.Contract.Symbol(&_Erc20.TransactOpts)
//}
//
//// Symbol is a paid mutator transaction binding the contract method 0x95d89b41.
////
//// Solidity: function symbol() returns(string)
//func (_Erc20 *Erc20TransactorSession) Symbol() (*types.Transaction, error) {
//	return _Erc20.Contract.Symbol(&_Erc20.TransactOpts)
//}
//
//// TotalSupply is a paid mutator transaction binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() returns(uint256)
//func (_Erc20 *Erc20Transactor) TotalSupply(opts *bind.TransactOpts) (*types.Transaction, error) {
//	return _Erc20.contract.Transact(opts, "totalSupply")
//}
//
//// TotalSupply is a paid mutator transaction binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() returns(uint256)
//func (_Erc20 *Erc20Session) TotalSupply() (*types.Transaction, error) {
//	return _Erc20.Contract.TotalSupply(&_Erc20.TransactOpts)
//}
//
//// TotalSupply is a paid mutator transaction binding the contract method 0x18160ddd.
////
//// Solidity: function totalSupply() returns(uint256)
//func (_Erc20 *Erc20TransactorSession) TotalSupply() (*types.Transaction, error) {
//	return _Erc20.Contract.TotalSupply(&_Erc20.TransactOpts)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
//func (_Erc20 *Erc20Transactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
//	return _Erc20.contract.Transact(opts, "transferFrom", sender, recipient, amount)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
//func (_Erc20 *Erc20Session) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
//	return _Erc20.Contract.TransferFrom(&_Erc20.TransactOpts, sender, recipient, amount)
//}
//
//// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
////
//// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
//func (_Erc20 *Erc20TransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
//	return _Erc20.Contract.TransferFrom(&_Erc20.TransactOpts, sender, recipient, amount)
//}
//
//// Erc20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Erc20 contract.
//type Erc20ApprovalIterator struct {
//	Event *Erc20Approval // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *Erc20ApprovalIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(Erc20Approval)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(Erc20Approval)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *Erc20ApprovalIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *Erc20ApprovalIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// Erc20Approval represents a Approval event raised by the Erc20 contract.
//type Erc20Approval struct {
//	Owner   common.Address
//	Spender common.Address
//	Value   *big.Int
//	Raw     types.Log // Blockchain specific contextual infos
//}
//
//// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
////
//// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
//func (_Erc20 *Erc20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*Erc20ApprovalIterator, error) {
//
//	var ownerRule []interface{}
//	for _, ownerItem := range owner {
//		ownerRule = append(ownerRule, ownerItem)
//	}
//	var spenderRule []interface{}
//	for _, spenderItem := range spender {
//		spenderRule = append(spenderRule, spenderItem)
//	}
//
//	logs, sub, err := _Erc20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
//	if err != nil {
//		return nil, err
//	}
//	return &Erc20ApprovalIterator{contract: _Erc20.contract, event: "Approval", logs: logs, sub: sub}, nil
//}
//
//// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
////
//// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
//func (_Erc20 *Erc20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *Erc20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {
//
//	var ownerRule []interface{}
//	for _, ownerItem := range owner {
//		ownerRule = append(ownerRule, ownerItem)
//	}
//	var spenderRule []interface{}
//	for _, spenderItem := range spender {
//		spenderRule = append(spenderRule, spenderItem)
//	}
//
//	logs, sub, err := _Erc20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(Erc20Approval)
//				if err := _Erc20.contract.UnpackLog(event, "Approval", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
////
//// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
//func (_Erc20 *Erc20Filterer) ParseApproval(log types.Log) (*Erc20Approval, error) {
//	event := new(Erc20Approval)
//	if err := _Erc20.contract.UnpackLog(event, "Approval", log); err != nil {
//		return nil, err
//	}
//	return event, nil
//}
//
//// Erc20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Erc20 contract.
//type Erc20TransferIterator struct {
//	Event *Erc20Transfer // Event containing the contract specifics and raw log
//
//	contract *bind.BoundContract // Generic contract to use for unpacking event data
//	event    string              // Event name to use for unpacking event data
//
//	logs chan types.Log        // Log channel receiving the found contract events
//	sub  ethereum.Subscription // Subscription for errors, completion and termination
//	done bool                  // Whether the subscription completed delivering logs
//	fail error                 // Occurred error to stop iteration
//}
//
//// Next advances the iterator to the subsequent event, returning whether there
//// are any more events found. In case of a retrieval or parsing error, false is
//// returned and Error() can be queried for the exact failure.
//func (it *Erc20TransferIterator) Next() bool {
//	// If the iterator failed, stop iterating
//	if it.fail != nil {
//		return false
//	}
//	// If the iterator completed, deliver directly whatever's available
//	if it.done {
//		select {
//		case log := <-it.logs:
//			it.Event = new(Erc20Transfer)
//			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//				it.fail = err
//				return false
//			}
//			it.Event.Raw = log
//			return true
//
//		default:
//			return false
//		}
//	}
//	// Iterator still in progress, wait for either a data or an error event
//	select {
//	case log := <-it.logs:
//		it.Event = new(Erc20Transfer)
//		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
//			it.fail = err
//			return false
//		}
//		it.Event.Raw = log
//		return true
//
//	case err := <-it.sub.Err():
//		it.done = true
//		it.fail = err
//		return it.Next()
//	}
//}
//
//// Error returns any retrieval or parsing error occurred during filtering.
//func (it *Erc20TransferIterator) Error() error {
//	return it.fail
//}
//
//// Close terminates the iteration process, releasing any pending underlying
//// resources.
//func (it *Erc20TransferIterator) Close() error {
//	it.sub.Unsubscribe()
//	return nil
//}
//
//// Erc20Transfer represents a Transfer event raised by the Erc20 contract.
//type Erc20Transfer struct {
//	From  common.Address
//	To    common.Address
//	Value *big.Int
//	Raw   types.Log // Blockchain specific contextual infos
//}
//
//// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
//func (_Erc20 *Erc20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*Erc20TransferIterator, error) {
//
//	var fromRule []interface{}
//	for _, fromItem := range from {
//		fromRule = append(fromRule, fromItem)
//	}
//	var toRule []interface{}
//	for _, toItem := range to {
//		toRule = append(toRule, toItem)
//	}
//
//	logs, sub, err := _Erc20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
//	if err != nil {
//		return nil, err
//	}
//	return &Erc20TransferIterator{contract: _Erc20.contract, event: "Transfer", logs: logs, sub: sub}, nil
//}
//
//// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
//func (_Erc20 *Erc20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *Erc20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {
//
//	var fromRule []interface{}
//	for _, fromItem := range from {
//		fromRule = append(fromRule, fromItem)
//	}
//	var toRule []interface{}
//	for _, toItem := range to {
//		toRule = append(toRule, toItem)
//	}
//
//	logs, sub, err := _Erc20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
//	if err != nil {
//		return nil, err
//	}
//	return event.NewSubscription(func(quit <-chan struct{}) error {
//		defer sub.Unsubscribe()
//		for {
//			select {
//			case log := <-logs:
//				// New log arrived, parse the event and forward to the user
//				event := new(Erc20Transfer)
//				if err := _Erc20.contract.UnpackLog(event, "Transfer", log); err != nil {
//					return err
//				}
//				event.Raw = log
//
//				select {
//				case sink <- event:
//				case err := <-sub.Err():
//					return err
//				case <-quit:
//					return nil
//				}
//			case err := <-sub.Err():
//				return err
//			case <-quit:
//				return nil
//			}
//		}
//	}), nil
//}
//
//// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
////
//// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
//func (_Erc20 *Erc20Filterer) ParseTransfer(log types.Log) (*Erc20Transfer, error) {
//	event := new(Erc20Transfer)
//	if err := _Erc20.contract.UnpackLog(event, "Transfer", log); err != nil {
//		return nil, err
//	}
//	return event, nil
//}
