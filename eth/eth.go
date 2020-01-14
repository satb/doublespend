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
			processBlock(httpClient, addresses, block, ch)
		}
	}()
}

/**
On the arrival of the new block, this function is responsible for
1. Looping through the previously cached transactions
2. If the transaction receipt for any of them is NOT found
	a. if the transaction was not already marked as double spent
    b. update the cache marking the transaction as double spent
    c. send a message on the channel that the transaction was double spent
3. If the transaction receipt for any of them IS found and the transaction was marked double spent
	a. mark the transaction as NOT double spent
    b. update the cache marking the transaction as NOT double spent
    c. send a message on the channel that the transaction was not double spent
4. Add the new transactions, if any, that originated from the set of addresses of interest
*/
func processBlock(client *ethclient.Client, addresses []string, block *types.Block, ch chan<- cache.Item) {
	for _, address := range addresses {
		addr := strings.ToLower(address)
		if ethCache.Get(addr) != nil {
			for _, tx := range ethCache.Get(addr) {
				_, err := getTxnReceipt(client, common.HexToHash(tx.Id))
				if err != nil && err == ethereum.NotFound {
					if tx.DoubleSpend == false {
						tx.DoubleSpend = true
						log.Println("DOUBLE SPEND DETECTED - address=", addr, "txnId=", tx.Id, "block=", tx.BlockNum)
						ethCache.UpdateItem(addr, tx)
						ch <- tx
					}
				} else if err != nil {
					if tx.DoubleSpend {
						//Revert the transaction which was marked as double spend
						tx.DoubleSpend = false
						ethCache.UpdateItem(addr, tx)
						ch <- tx
					}
				} else {
					log.Println(err)
				}
			}
		}
	}
	purgeCache(addresses)
	cacheNewTxs(client, addresses, block)
}

func cacheNewTxs(client *ethclient.Client, addresses []string, block *types.Block) {
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
