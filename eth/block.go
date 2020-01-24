package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func getBlock(client *ethclient.Client, blockNum int64) (block *types.Block, err error) {
	blockNumber := big.NewInt(blockNum)
	return client.BlockByNumber(context.Background(), blockNumber)
}

func startFullScan(client *ethclient.Client, addresses []string) {
	log.Println("Started startFullScan")
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
		cacheNewTxs(client, addresses, block)
	}
	log.Println("Completed startFullScan")
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
