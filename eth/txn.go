package eth

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// var n = eth.blocknumber;

// var txs = [];
// for(var i = 0; i < n; i++) {
//     var block = eth.getBlock(i, true);
//     for(var j = 0; j < block.transactions; j++) {
//         if( block.transactions[j].to == the_address )
//             txs.push(block.transactions[j]);
//     }
// }

func Scan(client *ethclient.Client, address string) []*types.Transaction {
	topHeader := topHeader(client)
	headerNumber := topHeader.Number.Int64()
	return ScanBlocks(client, address, 0, headerNumber)
}

func ScanBlocks(client *ethclient.Client, address string, fromBlock int64, toBlock int64) []*types.Transaction {
	txns := make([]*types.Transaction, 5)
	var i int64
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for i = fromBlock; i < toBlock; i++ {
		blockNumber := big.NewInt(i)
		block, err := client.BlockByNumber(context.Background(), blockNumber)
		if err != nil {
			log.Fatal(err)
		}
		var j int
		for j = 0; j < block.Transactions().Len(); j++ {
			tx := block.Transactions()[j]
			if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID)); err != nil {
				if msg.From().Hex() == address {
					txns = append(txns, tx)
				}
			}
		}
	}
	return txns
}

func topHeader(client *ethclient.Client) *types.Header {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	return header
}
