package eth

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/satb/doublespend/cache"
	"log"
	"strings"
)

var (
	ethCache = cache.New()
)

func Monitor(wssClient *ethclient.Client, httpClient *ethclient.Client, addresses []string, ch chan<- cache.Item) {
	//1. full scan for addresses of interest
	scan(httpClient, addresses)
	//2. subscribe to new blocks
	messages := make(chan *types.Block, 20)
	go func() {
		subscribe(wssClient, messages)
	}()
	//3. on receiving new block, process
	go func() {
		for {
			block := <-messages
			go processBlock(httpClient, addresses, block, ch)
		}
	}()
}

func processBlock(client *ethclient.Client, addresses []string, block *types.Block, ch chan<- cache.Item) {
	for _, address := range addresses {
		addr := strings.ToLower(address)
		if ethCache.Get(addr) != nil {
			for _, tx := range ethCache.Get(addr) {
				_, err := GetTxnReceipt(client, common.HexToHash(tx.Id))
				if err != nil && err == ethereum.NotFound {
					tx.DoubleSpend = true
					log.Println("DOUBLE SPEND DETECTED - address=", addr)
					ethCache.UpdateItem(addr, tx.Id, tx)
					ch <- tx
				} else {
					if tx.DoubleSpend {
						//Revert the transaction which was marked as double spend
						tx.DoubleSpend = false
						ethCache.UpdateItem(addr, tx.Id, tx)
						ch <- tx
					}
				}
			}
		}
	}
	updateCache(client, addresses, block)
}

func updateCache(client *ethclient.Client, addresses []string, block *types.Block) {
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
				cacheTxn(from, tx, block)
			}
		}
	}
}

func cacheTxn(from string, tx *types.Transaction, block *types.Block) {
	log.Println("Caching new txn with id ", tx.Hash().Hex())
	blockNumber := block.Number()
	ethCache.AddItem(from, cache.Item{
		Id:          tx.Hash().Hex(),
		From:        from,
		To:          tx.To().Hex(),
		Amount:      *tx.Value(),
		BlockNum:    blockNumber.Int64(),
		Time:        block.Header().Time,
		DoubleSpend: false,
	})
}
