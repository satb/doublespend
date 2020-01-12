package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

const CONFIRMATIONS int = 6

var suspiciousTxnMap map[string][]*types.Transaction = make(map[string][]*types.Transaction)

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
	ethBalance := ToEth(balance)
	result, _ := ethBalance.Int(new(big.Int))
	return result.Int64()
}

func Monitor(wssClient *ethclient.Client, httpClient *ethclient.Client, addresses []string) {
	txnMap := Scan(httpClient, addresses)
	printTxnMap(txnMap)
	messages := make(chan *types.Block, 20)
	go func() {
		Subscribe(wssClient, messages)
	}()

	go func() {
		for {
			block := <-messages
			go processNewBlock(httpClient, addresses, txnMap, block)
		}
	}()
}

func ToEth(value *big.Int) *big.Float {
	fbalance := new(big.Float)
	fbalance.SetString(value.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return ethValue
}

func processNewBlock(client *ethclient.Client, addresses []string, txnMap map[string][]*TxnInfo, block *types.Block) {
	log.Println("Processing new block with ", block.Transactions().Len(), " transactions in blockNumber", block.Number().Int64())
	//Get the txn receipt for the most recent 2 tx in the map, get the block hash, check if it is still valid and if not, log double spend
	for _, address := range addresses {
		addr := strings.ToLower(address)
		if txnMap[addr] != nil {
			//grab the most recent txn
			tx := txnMap[addr][len(txnMap[addr])-1]
			_, err := GetTxnReceipt(client, tx.Txn.Hash())
			if err != nil && err == ethereum.NotFound {
				log.Println("**********DOUBLE SPEND DETECTED FOR ", addr, " ********")
			}
		}
	}
	appendTxn(client, addresses, txnMap, block)
}

func appendTxn(client *ethclient.Client, addresses []string, txnMap map[string][]*TxnInfo, block *types.Block) {
	for _, tx := range block.Transactions() {
		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
			return
		}
		if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID)); err == nil {
			from := strings.ToLower(msg.From().Hex())
			_, found := find(addresses, from)
			if found {
				txnInfo := TxnInfo{Txn: tx, BlockNum: block.Number().Int64()}
				if txnMap[from] != nil {
					txnMap[from] = append(txnMap[from], &txnInfo)
				} else {
					txnMap[from] = []*TxnInfo{&txnInfo}
				}
			}
		}
	}
}

func printTxnMap(txnMap map[string][]*TxnInfo) {
	for k, v := range txnMap {
		fmt.Println("txnMap has address ", k, " and transactions ", printTxns(v))
	}
}

func printTxns(txns []*TxnInfo) string {
	result := ""
	if txns != nil {
		for _, txn := range txns {
			result += txn.Txn.Hash().Hex()
			result += ","
		}
	}
	return result
}
