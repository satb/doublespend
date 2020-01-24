package eth

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
)

type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

func ExtractTransferLog(receipt *types.Receipt) []LogTransfer {
	txLogs := make([]LogTransfer, 0)
	logs := receipt.Logs
	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			fmt.Printf("Log Name: Transfer\n")

			var transferEvent LogTransfer

			err := contractAbi.Unpack(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			txLogs = append(txLogs, LogTransfer{
				From:   transferEvent.From,
				To:     transferEvent.To,
				Tokens: transferEvent.Tokens,
			})
			fmt.Printf("From: %s\n", transferEvent.From.Hex())
			fmt.Printf("To: %s\n", transferEvent.To.Hex())
			fmt.Printf("NumTokens: %s\n", transferEvent.Tokens.String())
		}
	}
	return txLogs
}
