package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func transfer(client *ethclient.Client, from Account, to string, ethAmount int64) (txn *types.Transaction, err error) {
	privateKey, err := crypto.HexToECDSA(from.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	wei := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	value := new(big.Int).Mul(big.NewInt(ethAmount), wei)
	gasLimit := uint64(21000) // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress(to)
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Txn sent - txId=", signedTx.Hash().Hex(), " from=", strings.ToLower(fromAddress.Hex()))
	return signedTx, err
}

func getTxnReceipt(client *ethclient.Client, hash common.Hash) (receipt *types.Receipt, err error) {
	return client.TransactionReceipt(context.Background(), hash)
}

func TxnFrom(client *ethclient.Client, tx *types.Transaction) (from string, err error) {

	chainID, err := client.NetworkID(context.Background())

	if err != nil {
		return "", err
	}

	if msg, err := tx.AsMessage(types.NewEIP155Signer(chainID)); err == nil {
		return strings.ToLower(msg.From().Hex()), nil
	}

	return "", err
}
