package erc20

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

type ContractInfo struct {
	Name     string
	Symbol   string
	Decimals uint8
	Address  string
}

func GetContractInfo(client *ethclient.Client, contractAddress string) (contractInfo ContractInfo, err error) {
	tokenAddress := common.HexToAddress(contractAddress)
	instance, err := NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err, "Cannot create instance of erc20 instance for contract", contractAddress)
	}
	name, err := instance.Name(&bind.CallOpts{})
	symbol, err := instance.Symbol(&bind.CallOpts{})
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		return ContractInfo{}, err
	}
	return ContractInfo{
		Name:     name,
		Symbol:   symbol,
		Decimals: decimals,
	}, nil

}
