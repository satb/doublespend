package erc20

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func TokenBalance(client *ethclient.Client, contractAddr string, addr string) (*big.Int, error) {
	tokenAddress := common.HexToAddress(contractAddr)
	instance, err := NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	address := common.HexToAddress(addr)
	return instance.BalanceOf(&bind.CallOpts{}, address)
}
