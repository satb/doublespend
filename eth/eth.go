package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/satb/doublespend/cache"
	"github.com/savier89/ethunitconv"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"
)

var (
	ethCache = cache.New()
)

type TxnInfo struct {
	txn              *types.Transaction
	receipt          *types.Receipt
	from             string
	blockNum         int64
	numConfirmations int
	blockTxnValue    *big.Int
	contractAddress  string
	amount           *big.Int
	time             uint64
	to               string
}

func Monitor(wssClient *ethclient.Client, httpClient *ethclient.Client, addresses []string, ch chan<- cache.Item) {
	//1. full startFullScan for addresses of interest
	startFullScan(httpClient, addresses)
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
				receipt, err := getTxnReceipt(client, common.HexToHash(tx.Id))
				if err != nil && err == ethereum.NotFound {
					if tx.DoubleSpend == false {
						tx.DoubleSpend = true
						log.Println("DOUBLE SPEND DETECTED - address=", addr, "txnId=", tx.Id, "block=", tx.BlockNum)
						ethCache.UpdateItem(addr, tx)
						ch <- tx
					}
				} else if err == nil && tx.Amount.Cmp(big.NewInt(0)) == 0 && tx.ContractAddress != "" {
					//erc token - see txn receipt has the token transfer previously done
					transferLogs := ExtractTransferLog(receipt)
					erc20Found := false
					for _, txLog := range transferLogs {
						if strings.ToLower(txLog.From.Hex()) == tx.From && strings.ToLower(txLog.To.Hex()) == tx.To && strings.ToLower(receipt.ContractAddress.Hex()) == tx.ContractAddress {
							erc20Found = true
							break
						}
					}
					if !erc20Found {
						log.Println("ERC20 token not found for txnId=", tx.Id, " from=", tx.From, " to=", tx.To, " contractAddress=", tx.ContractAddress, " confirmations=", tx.Confirmations)
						if int(block.Number().Int64()-tx.BlockNum) > tx.Confirmations {
							if tx.DoubleSpend != true {
								log.Println("******ERC20 DOUBLE SPEND DETECTED*****")
								tx.DoubleSpend = true
								ethCache.UpdateItem(addr, tx)
								ch <- tx
							}
						}
					} else {
						//erc20 found...if double spent signaled earlier, signal not double spent
						if tx.DoubleSpend == true {
							log.Println("******ERC20 DOUBLE SPEND NEEDS TO BE REVERTED*****")
							tx.DoubleSpend = false
							ethCache.UpdateItem(addr, tx)
							ch <- tx
						}
					}
				} else if err == nil {
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

func processEth(client *ethclient.Client, addresses []string, block *types.Block, totalValue *big.Int) []TxnInfo {
	blockTxns := make([]TxnInfo, 0)
	for _, tx := range block.Transactions() {
		from, err := TxnFrom(client, tx)
		if err == nil {
			_, found := find(addresses, from)
			if found == true && tx.Value().Cmp(big.NewInt(0)) > 0 {
				blockTxns = append(blockTxns, TxnInfo{
					txn:              tx,
					blockNum:         block.Number().Int64(),
					from:             strings.ToLower(from),
					to:               strings.ToLower(tx.To().Hex()),
					blockTxnValue:    totalValue,
					numConfirmations: numConfirmations(totalValue),
					amount:           tx.Value(),
					time:             block.Time(),
				})
			}
		} else {
			log.Println("Cannot get transaction from for ", tx.Hash().Hex())
		}
	}
	return blockTxns
}

func processErc(client *ethclient.Client, addresses []string, block *types.Block, totalValue *big.Int) []TxnInfo {
	blockTxns := make([]TxnInfo, 0)
	for _, tx := range block.Transactions() {
		from, err := TxnFrom(client, tx)
		if err == nil {
			_, found := find(addresses, from)
			if found && tx.Value().Cmp(big.NewInt(0)) == 0 {
				receipt, err := getTxnReceipt(client, tx.Hash())
				if err != nil {
					log.Println("Could not fetch transaction receipt for ", tx.Hash().Hex())
					continue
				}
				transferLogs := ExtractTransferLog(receipt)
				fmt.Println("Found for txId ane extracted logs of log size", tx.Hash().Hex(), len(transferLogs))

				for _, txLog := range transferLogs {
					if strings.ToLower(txLog.From.Hex()) == strings.ToLower(from) {
						blockTxns = append(blockTxns, TxnInfo{
							txn:              tx,
							receipt:          receipt,
							from:             strings.ToLower(from),
							to:               strings.ToLower(txLog.To.Hex()),
							blockNum:         block.Number().Int64(),
							blockTxnValue:    totalValue,
							numConfirmations: numConfirmations(totalValue),
							contractAddress:  strings.ToLower(tx.To().Hex()),
							amount:           txLog.Tokens,
							time:             block.Time(),
						})
					}
				}
			}
		} else {
			log.Println("Cannot get transaction from for ", tx.Hash().Hex())
		}
	}
	return blockTxns
}

func cacheNewTxs(client *ethclient.Client, addresses []string, block *types.Block) {
	blockTxns := make([]TxnInfo, 0)
	totalValue := blockTxnValue(block.Transactions())
	blockTxns = append(blockTxns, processEth(client, addresses, block, totalValue)...)
	blockTxns = append(blockTxns, processErc(client, addresses, block, totalValue)...)
	for _, tx := range blockTxns {
		cacheNewTxn(tx)
	}
}

func blockTxnValue(txns types.Transactions) *big.Int {
	total := big.NewInt(0)
	for _, tx := range txns {
		v := tx.Value()
		total = total.Add(total, v)
	}
	return total
}

/**
Today, the total block txn value is how many confirmations are required. It could change later
*/
func numConfirmations(blockTxnValue *big.Int) int {
	ether := ethunitconv.FromWei(blockTxnValue.String(), "Ether")
	f, err := strconv.ParseFloat(ether, 64)
	if err != nil {
		log.Println("Could not convert wei to eth and parse into float for rounding")
	}
	return int(math.Ceil(f))
}
