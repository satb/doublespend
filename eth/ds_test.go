package eth

import (
	"fmt"
	"github.com/satb/doublespend/cache"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

/*
* Create a new address for Malice - the bad one
* Create a new address for Bob
* Create a new address for Jane
* Create new etherbase accounts eb1 and eb2 on node1 and node2 respectively
* Set etherbase on node 1 to Malice's address
* Set etherbase on node 2 to Bob's address
* Add peer of node 2 on node 1 so they sync
* Start mining for 10 seconds so Malice and Bob have some ETH
* Stop mining
* Set etherbase accounts for node 1 and node 2 to the newly created accounts
* Start full blockchain startFullScan for Malice's address
* subscribe to new blocks from node2
* Remove node 2 peer from node 1
* Transfer 10% of the ETH from Malice to Jane and send transaction to node 2 only (good one)
* After 1 seconds, get Malice's balance from node 2 and ensure it reads 10% less
* Send 100% of Malice's balance to Jane on node 1 only
* Start mining on node 1 and node 2
* Double spend should log becasue of the block subscription
 */

func TestDoubleSpend(t *testing.T) {
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
	fmt.Println("Waiting for 20 seconds")
	time.Sleep(20 * time.Second)

	fmt.Println("After 20 seconds...")
	maliceNode2Balance := GetBalance(client2, malice.Address)
	fmt.Println("On node2 malice balance=", maliceNode2Balance)

	maliceNode1Balance := GetBalance(client1, malice.Address)
	fmt.Println("On node1 malice balance=", maliceNode1Balance)

	StopMining(client1Url)
	StopMining(client2Url)

	fmt.Println("Stopped mining on both nodes")

	//set etherbase to eb1 and eb2 before mining again
	SetEtherBase(client1Url, eb1.Address)
	SetEtherBase(client2Url, eb2.Address)

	maliceBalance := GetBalance(client1, malice.Address)
	if maliceBalance <= 0 {
		t.Log("Malice account balance is 0")
		t.FailNow()
	}

	fmt.Println("Malice has a balance of ", maliceBalance)

	fmt.Println("Removing client2 as the peer of client1")
	//Remove peers now
	removePeer(client1Url, nodeInfo2.Enode)
	fmt.Println("Removed client2 as the peer of client1")

	//wait 2 second
	time.Sleep(2 * time.Second)

	//transfer 50% of the balance to the good client (client2)
	maliceTxAmount := maliceBalance / 2

	maliceNode2BalanceBeforeMiningStart := GetBalance(client2, malice.Address)
	fmt.Println("Before transferring maliceBalance from node2 ", maliceNode2BalanceBeforeMiningStart)
	fmt.Println("Transferring from malice to bob on node2 amount of ", maliceTxAmount)
	tx, err := transfer(client2, malice, bob.Address, maliceTxAmount)
	if err != nil {
		t.Log("Signing of txn failed")
		t.Fail()
	}
	fmt.Println("Transferred from malice to bob amount of ", maliceTxAmount)
	fmt.Println("Started mining on node2. Waiting for 10 seconds for more blocks to be mined")
	StartMining(client2Url, 2)
	time.Sleep(10 * time.Second)
	StopMining(client2Url)
	fmt.Println("Stopped mining after 10 seconds on node2...")

	maliceNode2Balance = GetBalance(client2, malice.Address)
	fmt.Println("On node2 malice balance=", maliceNode2Balance)

	if maliceNode2BalanceBeforeMiningStart == maliceNode2Balance {
		t.Log("Balance has not reduced after transfer for malice")
		t.FailNow()
	}
	maliceTxnReceipt, err := getTxnReceipt(client2, tx.Hash())
	if err != nil || maliceTxnReceipt == nil {
		t.Log("Cannot find transaction receipt for ", tx.Hash().Hex())
		t.FailNow()
	}
	fmt.Println("Recorded txn ", tx.Hash().Hex(), " in blockNumber=", maliceTxnReceipt.BlockNumber.Int64())

	b, err := getBlock(client2, maliceTxnReceipt.BlockNumber.Int64())
	if err != nil {
		t.Log("Cannot find transaction", tx.Hash().Hex(), "in block", b.Number().Int64())
		t.FailNow()
	} else {
		fmt.Println("Found transaction", tx.Hash().Hex(), "in block", b.Number().Int64())
	}
	maliceNode1Balance = GetBalance(client1, malice.Address)
	fmt.Println("On node1 malice balance=", maliceNode1Balance)

	fmt.Println("Starting monitoring on node2 for malice address ", malice.Address)

	//start monitorning on node 2 malice's address
	ch := make(chan<- cache.Item)
	Monitor(wssClient2, client2, []string{malice.Address}, ch)

	fmt.Println("Now monitoring ", malice.Address)

	//transfer 80% of malice's balance to Jane on node 1 (malicious node)
	var doubleSpendMaliceTxAmount = int64(float64(maliceBalance) * 0.8)
	fmt.Println("Now sending the double spend to node1 txn of ", doubleSpendMaliceTxAmount)
	tx, err = transfer(client1, malice, jane.Address, doubleSpendMaliceTxAmount)

	maliceNode2Balance = GetBalance(client2, malice.Address)
	fmt.Println("On node2 malice balance=", maliceNode2Balance)

	maliceNode1Balance = GetBalance(client1, malice.Address)
	fmt.Println("On node1 malice balance=", maliceNode1Balance)

	//start mining again - very fast on node 1 so it can add more blocks and overpower node 2
	StartMining(client1Url, 10)
	StartMining(client2Url, 1)

	fmt.Println("Started mining on both nodes....very fast on node1 and someone slowly on node2")

	fmt.Println("Waiting for 30 seconds for more blocks to be mined")
	//Wait for 30 seconds so node 1 adds more blocks because it has less difficulty
	time.Sleep(30 * time.Second)

	maliceTxnReceipt, err = getTxnReceipt(client1, tx.Hash())
	if err != nil || maliceTxnReceipt == nil {
		t.Log("Could not fetch txn receipt for txn=", tx.Hash().Hex())
		t.FailNow()
	}

	fmt.Println("Recorded double spend txn ", tx.Hash().Hex(), " in blockNumber=", maliceTxnReceipt.BlockNumber.Int64())

	maliceNode2Balance = GetBalance(client2, malice.Address)
	fmt.Println("Before adding peer of node2 on node1, on node2 malice balance=", maliceNode2Balance)

	maliceNode1Balance = GetBalance(client1, malice.Address)
	fmt.Println("Before adding peer of nod1 on node2, on node1 malice balance=", maliceNode1Balance)

	fmt.Println("Adding node2 as peer of node1 again")
	//now add the node 2 as peer to node 1 so the network syncs again
	addPeer(client1Url, nodeInfo2.Enode)

	//we should see the double spend being logged
	time.Sleep(20 * time.Second)

	fmt.Println("After 20 sec...")

	maliceNode2Balance = GetBalance(client2, malice.Address)
	fmt.Println("On node2 malice balance=", maliceNode2Balance)

	maliceNode1Balance = GetBalance(client1, malice.Address)
	fmt.Println("On node1 malice balance=", maliceNode1Balance)

	fmt.Println("Removing node2 as peer of node1 to reset to original state")
	removePeer(client1Url, nodeInfo2.Enode)

	StopMining(client1Url)
	StopMining(client2Url)
	fmt.Println("Stopped mining on all nodes")
}
