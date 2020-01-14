package eth

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/satb/doublespend/cache"
	"log"
)

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

func purgeCache(addresses []string) {
	for _, address := range addresses {
		ethCache.Purge(address)
	}
}
