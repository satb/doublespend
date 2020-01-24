package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/satb/doublespend/erc20"
	"log"
	"math/big"
	"testing"
	"time"
	//"time"
)

func deployContract(client *ethclient.Client, account Account) string {
	privateKey, err := crypto.HexToECDSA(account.PrivateKey)
	if err != nil {
		fmt.Println("Failed converting to ecdsca signature", err, account.PrivateKey)
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
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(8000000) // in units
	auth.GasPrice = gasPrice

	address, tx, instance, err := erc20.DeployErc20(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
	fmt.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0
	_ = instance
	return address.Hex()
}

//func TestDeployContract(t *testing.T) {
//	client2Url := "http://127.0.0.1:8101"
//	client2, err1 := ethclient.Dial(client2Url)
//	if err1 != nil {
//		t.Log(err1)
//		t.FailNow()
//	}
//	malice := CreateAccount()
//	fmt.Println("Created addresses for malice", malice.Address)
//	SetEtherBase(client2Url, malice.Address)
//	fmt.Println("Starting mining")
//	StartMining(client2Url, 1)
//	//sleep 20 seconds to accumulate eth and stop mining
//	fmt.Println("Waiting for 20 seconds")
//	time.Sleep(20 * time.Second)
//	fmt.Println("After 20 seconds...")
//	maliceBalance := GetBalance(client2, malice.Address)
//	if maliceBalance == 0 {
//		t.Log("Malice has no balance to do anything")
//		t.FailNow()
//	}
//	//eth.StopMining(client2Url)
//	//fmt.Println("Stopped mining on node1")
//
//	contractAddress := deployContract(client2, malice)
//
//	fmt.Println("Contract created, waiting for 20 seconds - ", contractAddress)
//	time.Sleep(20 * time.Second)
//
//}

//func TestGetTokenBalance(t *testing.T) {
//	client1Url := "http://127.0.0.1:8101"
//	client1, err1 := ethclient.Dial(client1Url)
//	if err1 != nil {
//		t.Log(err1)
//		t.FailNow()
//	}
//	bal, err := erc20.TokenBalance(client1, "0x7c43f397564d6d0183a163a97c9f68a9da22d68e", "0x71191A62829c797D88F918A4170624B137a404D8")
//	if err != nil {
//		t.Log("Cannot fetch balance on account")
//		t.FailNow()
//	}
//	fmt.Println("Balance is ", bal.String())
//}

func TestTokenTransfer(t *testing.T) {
	client1Url := "http://127.0.0.1:8101"
	client1, err1 := ethclient.Dial(client1Url)
	if err1 != nil {
		t.Log(err1)
		t.FailNow()
	}
	malice := CreateAccount()
	fmt.Println("Created addresses for malice address, privKey, pubKey", malice.Address, malice.PrivateKey, malice.PublicKey)
	SetEtherBase(client1Url, malice.Address)
	fmt.Println("Starting mining")
	StartMining(client1Url, 1)

	defer StopMining(client1Url)

	//sleep 20 seconds to accumulate eth and stop mining
	fmt.Println("Waiting for 20 seconds")
	time.Sleep(20 * time.Second)
	fmt.Println("After 20 seconds...")
	maliceBalance := GetBalance(client1, malice.Address)
	if maliceBalance == 0 {
		t.Log("Malice has no balance to do anything")
		t.FailNow()
	}
	//eth.StopMining(client1Url)
	//fmt.Println("Stopped mining on node1")

	contractAddress := deployContract(client1, malice)

	fmt.Println("Contract created, waiting for 20 seconds - ", contractAddress)
	time.Sleep(20 * time.Second)

	//Now create another address to transfer the tokens to
	john := CreateAccount()
	fmt.Println("Created addresses for John", john.Address)

	transferTxnId := erc20.Transfer(client1, contractAddress, malice.PrivateKey, john.Address, big.NewInt(1000))

	fmt.Println("Transferred from malice to john txnId=", transferTxnId)

	fmt.Println("Waiting 20 seconds for transaction ", transferTxnId, " to be mined")
	time.Sleep(20 * time.Second)

	newMaliceTokenBalance, _ := erc20.TokenBalance(client1, contractAddress, malice.Address)
	fmt.Println("Now malice token balance=", newMaliceTokenBalance.String())

	johnTokenBalance, _ := erc20.TokenBalance(client1, contractAddress, john.Address)

	if johnTokenBalance.Cmp(big.NewInt(0)) == 0 {
		t.Log("John token balance is 0. Transfer failed")
		t.FailNow()
	}

	fmt.Println("Now John token balance=", johnTokenBalance.String())
}
