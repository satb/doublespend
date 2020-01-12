package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TxnInfo struct {
	Txn      *types.Transaction
	BlockNum int64
}

func Scan(client *ethclient.Client, addresses []string) map[string][]*TxnInfo {
	topHeader := topHeader(client)
	headerNumber := topHeader.Number.Int64()
	return ScanBlocks(client, addresses, 0, headerNumber)
}

func ScanBlocks(client *ethclient.Client, addresses []string, fromBlock int64, toBlock int64) map[string][]*TxnInfo {
	txnMap := make(map[string][]*TxnInfo)
	var i int64
	for i = fromBlock; i < toBlock; i++ {
		blockNumber := big.NewInt(i)
		block, err := client.BlockByNumber(context.Background(), blockNumber)
		if err != nil {
			log.Fatal(err)
			continue
		}
		fmt.Println("Block - ", i, " has ", block.Transactions().Len(), " transactions")
		for _, tx := range block.Transactions() {
			chainID, err := client.NetworkID(context.Background())
			if err != nil {
				log.Fatal(err)
			}
			if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID)); err == nil {
				from := strings.ToLower(msg.From().Hex())
				var txns []*TxnInfo
				_, found := find(addresses, from)
				if found {
					fmt.Println("Found txn ", tx.Hash().Hex(), " during block scan for ", from)
					txnList := txnMap[msg.From().Hex()]
					if txnList == nil {
						txns = make([]*TxnInfo, 0)
					}
					txns = append(txns, &TxnInfo{Txn: tx, BlockNum: block.Number().Int64()})
					txnMap[from] = txns
				}
			}
		}
	}
	return txnMap
}

func GetBlock(client *ethclient.Client, blockNum int64) (block *types.Block, err error) {
	blockNumber := big.NewInt(blockNum)
	return client.BlockByNumber(context.Background(), blockNumber)
}

func Subscribe(client *ethclient.Client, channel chan<- *types.Block) {
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}
			go func() { channel <- block }()
		}
	}
}
