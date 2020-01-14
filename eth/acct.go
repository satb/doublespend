package eth

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"log"
	"math"
	"math/big"
)

type Account struct {
	PrivateKey string
	PublicKey  string
	Address    string
}

func CreateAccount() Account {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	return Account{
		hexutil.Encode(privateKeyBytes)[2:],
		hexutil.Encode(publicKeyBytes)[4:],
		address,
	}
}

/*
Returns ETH balance
*/
func GetBalance(client *ethclient.Client, address string) int64 {
	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	ethBalance := toEth(balance)
	result, _ := ethBalance.Int(new(big.Int))
	return result.Int64()
}

func toEth(value *big.Int) *big.Float {
	fbalance := new(big.Float)
	fbalance.SetString(value.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return ethValue
}
