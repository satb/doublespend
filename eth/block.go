package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strings"
)

func getBlock(client *ethclient.Client, blockNum int64) (block *types.Block, err error) {
	blockNumber := big.NewInt(blockNum)
	return client.BlockByNumber(context.Background(), blockNumber)
}

func scan(client *ethclient.Client, addresses []string) {
	log.Println("Started scan")
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	var i int64
	for i = 0; i < header.Number.Int64(); i++ {
		block, err := getBlock(client, i)
		if err != nil {
			log.Println(err)
			continue
		}
		for _, tx := range block.Transactions() {
			chainID, err := client.NetworkID(context.Background())
			if err != nil {
				log.Println(err)
				continue
			}
			if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID)); err == nil {
				from := strings.ToLower(msg.From().Hex())
				_, found := find(addresses, from)
				if found {
					cacheNewTxn(from, tx, block)
				}
			}
		}
	}
	log.Println("Completed scan")
}

func subscribe(client *ethclient.Client, channel chan<- *types.Block) {
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Println(err)
		case header := <-headers:
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Println(err)
			}
			go func() { channel <- block }()
		}
	}
}
