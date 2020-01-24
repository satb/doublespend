package eth

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/satb/doublespend/cache"
	"github.com/satb/doublespend/erc20"
	"log"
	"strings"
)

var (
	logTransferSig     = []byte("Transfer(address,address,uint256)")
	logTransferSigHash = crypto.Keccak256Hash(logTransferSig)
	contractAbi, _     = abi.JSON(strings.NewReader(erc20.Erc20ABI))
)

func cacheNewTxn(tx TxnInfo) {
	log.Println("Caching new txn with id ", tx.txn.Hash().Hex())
	ethCache.AddItem(tx.from, cache.Item{
		Id:              tx.txn.Hash().Hex(),
		From:            tx.from,
		To:              tx.to,
		Amount:          *tx.amount,
		BlockNum:        tx.blockNum,
		Time:            tx.time,
		DoubleSpend:     false,
		Confirmations:   tx.numConfirmations,
		ContractAddress: tx.contractAddress,
	})
}

func purgeCache(addresses []string) {
	for _, address := range addresses {
		ethCache.Purge(address)
	}
}
