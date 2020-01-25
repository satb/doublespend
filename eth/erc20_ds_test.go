package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/satb/doublespend/cache"
	"github.com/satb/doublespend/erc20"
	"log"
	"math/big"
	"testing"
	"time"
	//"time"
)

func deployContract(client *ethclient.Client, account Account) (txnId string, contractAddress string) {
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
	return tx.Hash().Hex(), address.Hex()
}

func TestErc20DoubleSpend(t *testing.T) {
	client1Url := "http://127.0.0.1:8101"
	client1, err1 := ethclient.Dial(client1Url)
	if err1 != nil {
		t.Log(err1)
		t.FailNow()
	}

	client2Url := "http://127.0.0.1:8102"
	client2, err2 := ethclient.Dial(client2Url)
	if err2 != nil {
		t.Log(err2)
		t.FailNow()
	}

	wssClient2, wssErr2 := ethclient.Dial("ws://127.0.0.1:4102")
	if wssErr2 != nil {
		t.Log(wssErr2)
		t.FailNow()
	}

	malice := CreateAccount()
	bob := CreateAccount()
	jane := CreateAccount()
	eb1 := CreateAccount()
	eb2 := CreateAccount()

	defer StopMining(client1Url)
	defer StopMining(client2Url)

	fmt.Println("Created addresses for malice, bob, jane, eb1, eb2", malice.Address, bob.Address, jane.Address, eb1.Address, eb2.Address)

	SetEtherBase(client1Url, malice.Address)
	SetEtherBase(client2Url, bob.Address)

	fmt.Println("Fetching node information for client2")
	nodeInfo2 := getNodeInfo(client2Url)
	fmt.Println("Adding client2 as peer of client1")
	addPeer(client1Url, nodeInfo2.Enode)

	fmt.Println("Starting mining")
	StartMining(client1Url, 1)
	StartMining(client2Url, 1)
	fmt.Println("Started mining on both clients")

	//sleep 20 seconds to accumulate eth and stop mining
	fmt.Println("Waiting for 30 seconds")
	time.Sleep(30 * time.Second)

	fmt.Println("After 20 seconds...")
	maliceNode2EthBalance := GetBalance(client2, malice.Address)
	fmt.Println("On node2 malice balance=", maliceNode2EthBalance)

	maliceNode1EthBalance := GetBalance(client1, malice.Address)
	fmt.Println("On node1 malice balance=", maliceNode1EthBalance)

	//Deploy ERC20 token to node1
	contractTxnId, contractAddress := deployContract(client1, malice)
	//sleep 20 seconds to accumulate eth and stop mining
	fmt.Println("Deployed contract. Waiting for 20 seconds")
	time.Sleep(20 * time.Second)

	StopMining(client1Url)
	StopMining(client2Url)

	fmt.Println("Stopped mining on both nodes")

	//set etherbase to eb1 and eb2 before mining again
	SetEtherBase(client1Url, eb1.Address)
	SetEtherBase(client2Url, eb2.Address)

	fmt.Println("Changed etherbase on both nodes so malice and bob balances remain constant for the rest of the test")

	maliceBalance := GetBalance(client1, malice.Address)
	if maliceBalance <= 0 {
		t.Log("Malice account balance is 0")
		t.FailNow()
	}

	fmt.Println("Malice has an eth balance of ", maliceBalance)

	//Ensure contract has been synced to the other node
	_, err := getTxnReceipt(client2, common.HexToHash(contractTxnId))
	if err != nil {
		t.Log("Failed to retrieve contract transaction from the other node. Failing")
		t.FailNow()
	}

	//Get the token balance for malice now on node1 - it should show the full balance of the newly deployed contract
	maliceTokenBalance, err := erc20.TokenBalance(client1, contractAddress, malice.Address)
	if err != nil {
		t.Log("Could not fetch token balance for malice on node1. Failing")
		t.FailNow()
	}
	if maliceTokenBalance.Cmp(big.NewInt(0)) == 0 {
		t.Log("Malice has no token balance on node1. Failing")
		t.FailNow()
	}
	fmt.Println("On node 1, malice has a token balance of ", maliceTokenBalance.String())

	//Get the token balance for malice on node2 - it should show the full balance of the newly deployed contract
	//Get the token balance for malice now on node1 - it should show the full balance of the newly deployed contract
	maliceTokenBalance, err = erc20.TokenBalance(client2, contractAddress, malice.Address)
	if err != nil {
		t.Log("Could not fetch token balance for malice on node2. Failing")
		t.FailNow()
	}
	if maliceTokenBalance.Cmp(big.NewInt(0)) == 0 {
		t.Log("Malice has no token balance on node2. Failing")
		t.FailNow()
	}
	fmt.Println("On node2, malice has a token balance of ", maliceTokenBalance.String())

	fmt.Println("Removing client2 as the peer of client1")
	//Remove peers now
	removePeer(client1Url, nodeInfo2.Enode)
	fmt.Println("Removed client2 as the peer of client1")

	//wait 2 second
	time.Sleep(2 * time.Second)

	//transfer 70% of the tokens
	maliceTokenTxAmount, _ := new(big.Int).SetString("400000000000000000000000", 10)

	fmt.Println("Transferring from malice to bob on node2 amount of ", maliceTokenTxAmount.String())
	txnId := erc20.Transfer(client2, contractAddress, malice.PrivateKey, bob.Address, maliceTokenTxAmount)
	fmt.Println("Transferred from malice to bob amount of ", maliceTokenTxAmount.String(), " txnId=", txnId)

	//Now start mining on node2 only so the transaction is picked up
	fmt.Println("Started mining on node2. Waiting for 10 seconds for more blocks to be mined")
	StartMining(client2Url, 2)
	time.Sleep(20 * time.Second)
	StopMining(client2Url)
	fmt.Println("Stopped mining after 20 seconds on node2...")

	maliceNode2TokenBalance, err := erc20.TokenBalance(client2, contractAddress, malice.Address)
	if err != nil {
		t.Log("Could not fetch token balance for malice. Failing")
		t.FailNow()
	}
	fmt.Println("On node2 malice balance=", maliceNode2TokenBalance.String())

	if maliceTokenBalance.Cmp(maliceNode2TokenBalance) == 0 {
		t.Log("Balance has not reduced after transfer for malice. Failing test")
		t.FailNow()
	}

	maliceTokenTxnReceipt, err := getTxnReceipt(client2, common.HexToHash(txnId))
	if err != nil || maliceTokenTxnReceipt == nil {
		t.Log("Cannot find transaction receipt for ", txnId)
		t.FailNow()
	}
	fmt.Println("Recorded txn ", txnId, " in blockNumber=", maliceTokenTxnReceipt.BlockNumber.Int64())

	maliceNode1TokenBalance, err := erc20.TokenBalance(client1, contractAddress, malice.Address)
	if err != nil {
		t.Log("Cannot fetch token balance for malice on node1", err)
	}
	fmt.Println("On node1 malice token balance=", maliceNode1TokenBalance.String())

	fmt.Println("Starting monitoring on node2 for malice address ", malice.Address)

	//start monitorning on node 2 malice's address
	ch := make(chan<- cache.Item)
	Monitor(wssClient2, client2, []string{malice.Address}, ch)

	fmt.Println("Now monitoring ", malice.Address)

	//transfer all but 1000 tokens as the double spend
	var doubleSpendMaliceTxAmount = maliceNode1TokenBalance.Sub(maliceNode1TokenBalance, big.NewInt(1000))
	fmt.Println("Now sending the double spend to node1 txn of ", doubleSpendMaliceTxAmount)
	//*******Double spent with this transfer now*******
	txnId = erc20.Transfer(client1, contractAddress, malice.PrivateKey, jane.Address, doubleSpendMaliceTxAmount)

	//See the balances on both nodes one last time before starting to mine
	maliceNode2TokenBalance, err = erc20.TokenBalance(client2, contractAddress, malice.Address)
	if err != nil {
		t.Log("Cannot get malice token balance on node2")
		t.FailNow()
	}
	fmt.Println("On node2 malice balance=", maliceNode2TokenBalance.String())

	maliceNode1TokenBalance, err = erc20.TokenBalance(client1, contractAddress, malice.Address)
	if err != nil {
		t.Log("Cannot get malice token balance on node2")
		t.FailNow()
	}
	fmt.Println("On node1 malice token balance=", maliceNode1TokenBalance.String())

	//start mining again - very fast on node 1 so it can add more blocks and overpower node 2
	StartMining(client1Url, 10)
	StartMining(client2Url, 1)

	fmt.Println("Started mining on both nodes....very fast on node1 and someone slowly on node2")

	fmt.Println("Waiting for 30 seconds for more blocks to be mined")
	//Wait for 30 seconds so node 1 adds more blocks
	time.Sleep(30 * time.Second)

	maliceTokenTxnReceipt, err = getTxnReceipt(client1, common.HexToHash(txnId))
	if err != nil || maliceTokenTxnReceipt == nil {
		t.Log("Could not fetch txn receipt for txn=", txnId)
		t.FailNow()
	}

	fmt.Println("Recorded double spend txn ", txnId, " in blockNumber=", maliceTokenTxnReceipt.BlockNumber.Int64())

	maliceNode2TokenBalance, err = erc20.TokenBalance(client2, contractAddress, malice.Address)
	if err != nil {
		t.Log("Cannot get malice token balance on node2")
		t.FailNow()
	}
	fmt.Println("Before adding peer of node2 on node1, on node2 malice token balance=", maliceNode2TokenBalance.String())

	maliceNode1TokenBalance, err = erc20.TokenBalance(client1, contractAddress, malice.Address)
	if err != nil {
		t.Log("Cannot get malice token balance on node1")
		t.FailNow()
	}
	fmt.Println("Before adding peer of node1 on node2, on node1 malice token balance=", maliceNode1TokenBalance.String())

	fmt.Println("Adding node2 as peer of node1 again")
	//now add the node 2 as peer to node 1 so the network syncs again
	addPeer(client1Url, nodeInfo2.Enode)

	//we should see the double spend being logged
	var b1 *types.Block
	var b2 *types.Block
	for i := 0; i < 5; i++ {
		b1, _ = client1.BlockByNumber(context.Background(), nil)
		fmt.Println("Blocks mined on node1 so far - ", b1.Number().String())
		b2, _ = client2.BlockByNumber(context.Background(), nil)
		fmt.Println("Blocks mined on node2 so far - ", b2.Number().String())
		time.Sleep(6 * time.Second)
	}

	fmt.Println("After 30 sec...")

	maliceNode2TokenBalance, err = erc20.TokenBalance(client2, contractAddress, malice.Address)
	if err != nil {
		t.Log("Cannot get malice token balance on node2")
		t.FailNow()
	}
	fmt.Println("On node2 malice balance=", maliceNode2TokenBalance.String())

	maliceNode1TokenBalance, err = erc20.TokenBalance(client1, contractAddress, malice.Address)
	if err != nil {
		t.Log("Cannot get malice token balance on node1")
		t.FailNow()
	}
	fmt.Println("On node1 malice balance=", maliceNode1TokenBalance.String())

	fmt.Println("Removing node2 as peer of node1 to reset to original state")
	removePeer(client1Url, nodeInfo2.Enode)

	fmt.Println("Stopping mining on all nodes")
}

func TestLogFetch(t *testing.T) {
	client1Url := "http://127.0.0.1:8101"
	client1, err1 := ethclient.Dial(client1Url)
	if err1 != nil {
		t.Log(err1)
		t.FailNow()
	}
	receipt, _ := getTxnReceipt(client1, common.HexToHash("0x06a77154b8f30a4924b9123e194ff1e0d017d6f54e642bbe8b3ad0ecd8acfaf5"))
	txLogs := ExtractTransferLog(receipt)
	fmt.Println("size - ", len(txLogs))
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

//func TestTokenTransfer(t *testing.T) {
//	client1Url := "http://127.0.0.1:8101"
//	client1, err1 := ethclient.Dial(client1Url)
//	if err1 != nil {
//		t.Log(err1)
//		t.FailNow()
//	}
//	malice := CreateAccount()
//	fmt.Println("Created addresses for malice address, privKey, pubKey", malice.Address, malice.PrivateKey, malice.PublicKey)
//	SetEtherBase(client1Url, malice.Address)
//	fmt.Println("Starting mining")
//	StartMining(client1Url, 1)
//
//	defer StopMining(client1Url)
//
//	//sleep 20 seconds to accumulate eth and stop mining
//	fmt.Println("Waiting for 20 seconds")
//	time.Sleep(20 * time.Second)
//	fmt.Println("After 20 seconds...")
//	maliceBalance := GetBalance(client1, malice.Address)
//	if maliceBalance == 0 {
//		t.Log("Malice has no balance to do anything")
//		t.FailNow()
//	}
//	//eth.StopMining(client1Url)
//	//fmt.Println("Stopped mining on node1")
//
//	contractAddress := deployContract(client1, malice)
//
//	fmt.Println("Contract created, waiting for 20 seconds - ", contractAddress)
//	time.Sleep(20 * time.Second)
//
//	//Now create another address to transfer the tokens to
//	john := CreateAccount()
//	fmt.Println("Created addresses for John", john.Address)
//
//	transferTxnId := erc20.Transfer(client1, contractAddress, malice.PrivateKey, john.Address, big.NewInt(1000))
//
//	fmt.Println("Transferred from malice to john txnId=", transferTxnId)
//
//	fmt.Println("Waiting 20 seconds for transaction ", transferTxnId, " to be mined")
//	time.Sleep(20 * time.Second)
//
//	newMaliceTokenBalance, _ := erc20.TokenBalance(client1, contractAddress, malice.Address)
//	fmt.Println("Now malice token balance=", newMaliceTokenBalance.String())
//
//	johnTokenBalance, _ := erc20.TokenBalance(client1, contractAddress, john.Address)
//
//	if johnTokenBalance.Cmp(big.NewInt(0)) == 0 {
//		t.Log("John token balance is 0. Transfer failed")
//		t.FailNow()
//	}
//
//	fmt.Println("Now John token balance=", johnTokenBalance.String())
//}
