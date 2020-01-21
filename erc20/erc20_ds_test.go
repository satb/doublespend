package erc20

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/satb/doublespend/eth"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"testing"
	//"time"
)

const helloTokenByteCode = "0x60c0604052600560808190526448454c4c4f60d81b60a090815261002691600191906100cd565b5060408051808201909152600380825262484c4f60e81b6020909201918252610051916002916100cd565b5060038054601260ff19909116179081905560ff16600a0a620f42400260045534801561007d57600080fd5b5060045433600081815260208181526040808320859055805194855251929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a3610168565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061010e57805160ff191683800117855561013b565b8280016001018555821561013b579182015b8281111561013b578251825591602001919060010190610120565b5061014792915061014b565b5090565b61016591905b808211156101475760008155600101610151565b90565b61056f806101776000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063313ce56711610066578063313ce567146101a557806370a08231146101c357806395d89b41146101e9578063a9059cbb146101f1578063dd62ed3e1461021d57610093565b806306fdde0314610098578063095ea7b31461011557806318160ddd1461015557806323b872dd1461016f575b600080fd5b6100a061024b565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100da5781810151838201526020016100c2565b50505050905090810190601f1680156101075780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101416004803603604081101561012b57600080fd5b506001600160a01b0381351690602001356102d8565b604080519115158252519081900360200190f35b61015d61033e565b60408051918252519081900360200190f35b6101416004803603606081101561018557600080fd5b506001600160a01b03813581169160208101359091169060400135610344565b6101ad610421565b6040805160ff9092168252519081900360200190f35b61015d600480360360208110156101d957600080fd5b50356001600160a01b031661042a565b6100a061043c565b6101416004803603604081101561020757600080fd5b506001600160a01b038135169060200135610494565b61015d6004803603604081101561023357600080fd5b506001600160a01b038135811691602001351661051c565b60018054604080516020600284861615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156102d05780601f106102a5576101008083540402835291602001916102d0565b820191906000526020600020905b8154815290600101906020018083116102b357829003601f168201915b505050505081565b3360008181526005602090815260408083206001600160a01b038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a350600192915050565b60045481565b6001600160a01b03831660009081526020819052604081205482111561036957600080fd5b6001600160a01b038416600090815260056020908152604080832033845290915290205482111561039957600080fd5b6001600160a01b038085166000818152602081815260408083208054889003905593871680835284832080548801905583835260058252848320338452825291849020805487900390558351868152935191937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929081900390910190a35060019392505050565b60035460ff1681565b60006020819052908152604090205481565b6002805460408051602060018416156101000260001901909316849004601f810184900484028201840190925281815292918301828280156102d05780601f106102a5576101008083540402835291602001916102d0565b336000908152602081905260408120548211156104b057600080fd5b33600081815260208181526040808320805487900390556001600160a01b03871680845292819020805487019055805186815290519293927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929181900390910190a350600192915050565b60056020908152600092835260408084209091529082529020548156fea2646970667358221220cca447df38a8f7f87fb1022f4cb1d99efe983653c8c21881bb77124c98c569e864736f6c63430006010033"

func deployContract(client *ethclient.Client, account eth.Account) string {
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
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(600152) // in units
	auth.GasPrice = gasPrice

	address, tx, instance, err := DeployErc20(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
	fmt.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0
	_ = instance
	return address.Hex()
}

func getBalance(client *ethclient.Client, contractAddress string, account eth.Account, addr string) {
	tokenAddress := common.HexToAddress(contractAddress)
	instance, err := NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err, "Cannot create instance of erc20 instance for contract", contractAddress)
	}
	//privKey, err := crypto.HexToECDSA(account.PrivateKey)
	//chainID, err := client.NetworkID(context.Background())
	//nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(account.Address))
	//if err != nil {
	//	log.Fatal(err)
	//}

	// signerFn := func(signer types.Signer, addresses common.Address, tx *types.Transaction) (transaction *types.Transaction, e error) {
	//	return types.SignTx(tx, types.NewEIP155Signer(chainID), privKey)
	//}

	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err, " .Can't get name ", contractAddress)
	}

	//contractAbi, err := abi.JSON(strings.NewReader(Erc20ABI))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//var nameStr string
	//err = contractAbi.Unpack(&nameStr, "name", name.Data())
	//if err != nil {
	//	log.Fatal(err)
	//}

	fmt.Println("************ name ***************", name)

	//nonce, _ = client.PendingNonceAt(context.Background(), common.HexToAddress(account.Address))
	address := common.HexToAddress(addr)
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err, ". Cannot fetch balance from erc20 instance for contractAddress ", contractAddress, "at address ", addr)
	}

	fmt.Printf("name: %s\n", name) // "name: Golem Network"
	//fmt.Printf("symbol: %s\n", symbol)     // "symbol: GNT"
	fmt.Printf(" balance: %v\n", bal)

}

//func createToken(client *ethclient.Client, account eth.Account, data string) {
//	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(account.Address))
//	if err != nil {
//		log.Fatal(err)
//	}
//	gasLimit := uint64(489152)
//	dataBytes := []byte(data)
//	tx := types.NewContractCreation(nonce, nil,  gasLimit, nil, dataBytes)
//	chainID, err := client.NetworkID(context.Background())
//	if err != nil {
//		fmt.Println("Failed creating contract", err)
//		log.Fatal(err)
//	}
//	privKey, err := crypto.HexToECDSA(account.PrivateKey)
//	if err != nil {
//		fmt.Println("Failed converting to ecdsca signature", err, account.PrivateKey)
//		log.Fatal(err)
//	}
//
//	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privKey)
//
//	if err != nil {
//		fmt.Println("failed to sign transaction", err)
//		log.Fatal(err)
//	}
//
//	ts := types.Transactions{signedTx}
//	rawTxBytes := ts.GetRlp(0)
//	rawTxHex := hex.EncodeToString(rawTxBytes)
//
//	fmt.Printf(rawTxHex) // f86...772
//
//	err = client.SendTransaction(context.Background(), tx)
//	if err != nil {
//		fmt.Println("Failed to send transaction", err)
//		log.Fatal(err)
//	}
//
//	fmt.Printf("tx sent: %s", tx.Hash().Hex())
//}

/////
func newAccount(clientUrl string, passphrase string) string {
	id := uuid.New()
	params := [...]interface{}{passphrase}
	message := map[string]interface{}{
		"id":      id.String(),
		"jsonrpc": "2.0",
		"method":  "personal_newAccount",
		"params":  params,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(clientUrl, "application/json", bytes.NewBuffer(bytesRepresentation))

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))

	var eresp eth.EthResponse
	err = json.Unmarshal(body, &eresp)
	if err != nil {
		log.Println("Error unmarshaling", err)
	}
	var address string
	err = json.Unmarshal(eresp.Result, &address)
	if err != nil {
		log.Fatal("Error extrating address")
	}
	return address
}

func unlockAccount(clientUrl string, from string, passphrase string) {
	id := uuid.New()
	params := [...]interface{}{from, passphrase, 0}
	message := map[string]interface{}{
		"id":      id.String(),
		"jsonrpc": "2.0",
		"method":  "personal_unlockAccount",
		"params":  params,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := http.Post(clientUrl, "application/json", bytes.NewBuffer(bytesRepresentation))

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))
}

//func getTokenBalance(client *ethclient.Client) {
//	tokenAddress := common.HexToAddress("0xFdC4390f63cC385Ef34B52d5f780dd671B1d067B")
//	instance, err := NewToken(tokenAddress, client)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	address := common.HexToAddress("0x0536806df512d6cdde913cf95c9886f65b1d3462")
//	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	name, err := instance.Name(&bind.CallOpts{})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	symbol, err := instance.Symbol(&bind.CallOpts{})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	decimals, err := instance.Decimals(&bind.CallOpts{})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("name: %s\n", name)         // "name: Golem Network"
//	fmt.Printf("symbol: %s\n", symbol)     // "symbol: GNT"
//	fmt.Printf("decimals: %v\n", decimals) // "decimals: 18"
//
//	fmt.Printf("wei: %s\n", bal) // "wei: 74605500647408739782407023"
//
//	fbal := new(big.Float)
//	fbal.SetString(bal.String())
//	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
//
//	fmt.Printf("balance: %f", value) // "balance: 74605500.647409"
//}

////

func createToken1(clientUrl string, from string, data string) {

	id := uuid.New()
	params := []map[string]string{{
		"from": from,
		"gas":  "0x776c0",
		"data": data,
	}}

	message := map[string]interface{}{
		"id":      id.String(),
		"jsonrpc": "2.0",
		"method":  "eth_sendTransaction",
		"params":  params,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := http.Post(clientUrl, "application/json", bytes.NewBuffer(bytesRepresentation))

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))
}

//func TestAddItems(t *testing.T) {
//	client1Url := "http://127.0.0.1:8101"
//	client1, err1 := ethclient.Dial(client1Url)
//	if err1 != nil {
//		t.Log(err1)
//		t.FailNow()
//	}
//	//malice := eth.CreateAccount()
//	//fmt.Println("Created addresses for malice", malice.Address)
//	//eth.SetEtherBase(client1Url, malice.Address)
//	//fmt.Println("Starting mining")
//	//eth.StartMining(client1Url, 1)
//	////sleep 20 seconds to accumulate eth and stop mining
//	//fmt.Println("Waiting for 20 seconds")
//	//time.Sleep(20 * time.Second)
//	//fmt.Println("After 20 seconds...")
//	//maliceBalance := eth.GetBalance(client1, malice.Address)
//	//if maliceBalance == 0 {
//	//	t.Log("Malice has no balance to do anything")
//	//	t.FailNow()
//	//}
//	////eth.StopMining(client1Url)
//	////fmt.Println("Stopped mining on node1")
//	//
//	//createToken(client1, malice, helloTokenByteCode)
//
//	//*********** WORKING BELOW
//
//	passphrase := "john3"
//	address := newAccount(client1Url, passphrase)
//	unlockAccount(client1Url, address, passphrase)
//
//	//Set etherbase and start mining
//	eth.SetEtherBase(client1Url, address)
//	fmt.Println("Starting mining")
//	eth.StartMining(client1Url, 1)
//	//sleep 20 seconds to accumulate eth and stop mining
//	fmt.Println("Waiting for 20 seconds")
//	time.Sleep(20 * time.Second)
//	fmt.Println("After 20 seconds...")
//	balance := eth.GetBalance(client1, address)
//	if balance == 0 {
//		t.Log("address has no balance")
//		t.FailNow()
//	}
//
//	//stop mining
//	//eth.StopMining(client1Url)
//
//	createToken1(client1Url, address, helloTokenByteCode)
//
//	//Create malice's account and transfer some of the tokens
//	malice := eth.CreateAccount()
//	fmt.Println("Created addresses for malice", malice.Address)
//
//}

//func TestDeployContract(t *testing.T) {
//	client2Url := "http://127.0.0.1:8102"
//	client2, err1 := ethclient.Dial(client2Url)
//	if err1 != nil {
//		t.Log(err1)
//		t.FailNow()
//	}
//	malice := eth.CreateAccount()
//	fmt.Println("Created addresses for malice", malice.Address)
//	eth.SetEtherBase(client2Url, malice.Address)
//	fmt.Println("Starting mining")
//	eth.StartMining(client2Url, 1)
//	//sleep 20 seconds to accumulate eth and stop mining
//	fmt.Println("Waiting for 20 seconds")
//	time.Sleep(20 * time.Second)
//	fmt.Println("After 20 seconds...")
//	maliceBalance := eth.GetBalance(client2, malice.Address)
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
//	getBalance(client2, contractAddress, malice, malice.Address)
//}

func TestGetBalance(t *testing.T) {
	client2Url := "http://127.0.0.1:8102"
	client2, err1 := ethclient.Dial(client2Url)
	if err1 != nil {
		t.Log(err1)
		t.FailNow()
	}
	account := eth.Account{
		PrivateKey: "e87f0151e23b607fdce450aacee56f1544287020376ba43dcc7714b5ccd43d7b",
		PublicKey:  "c4f06f79b2b409761a925484a968ff136876eeb406373410c10a2d08e6b3d6efc19a2a402cf80ef297a1bcc3d9f24af590c05a9bc61e2661a82c1c8893878382",
		Address:    "0xD5fB2a241782b7016BB2792ffbacd9dd691fb03a",
	}
	getBalance(client2, "0x1824C0d688b51e68199EFe2330Ff8482FA6f1A7b", account, "0x63ac173503f0565065f15b62ccc835f21ae23ebd")
	//getBalance(client1, "0x007f17daC8e1e81b1ccB66E2692D7116b5Efc317", "0x5cb2DAd087D37EAA788F4109261570b03f7EF003")
	//getBalance(client1, "0x12CF63a69B337aaD94beD2339fd5D8B29aE5Ddc0", "0xb16910653dDAc532d4D67729377B9E895C7f361b")
}

//func TestStoreDeploy(t *testing.T) {
//		client1Url := "http://127.0.0.1:8101"
//		client, err1 := ethclient.Dial(client1Url)
//		if err1 != nil {
//			t.Log(err1)
//			t.FailNow()
//		}
//
//		malice := eth.CreateAccount()
//		fmt.Println("Created addresses for malice", malice.Address)
//		eth.SetEtherBase(client1Url, malice.Address)
//		fmt.Println("Starting mining")
//		eth.StartMining(client1Url, 1)
//		//sleep 20 seconds to accumulate eth and stop mining
//		fmt.Println("Waiting for 20 seconds")
//		time.Sleep(20 * time.Second)
//		fmt.Println("After 20 seconds...")
//		maliceBalance := eth.GetBalance(client, malice.Address)
//		if maliceBalance == 0 {
//			t.Log("Malice has no balance to do anything")
//			t.FailNow()
//		}
//	privateKey, err := crypto.HexToECDSA(malice.PrivateKey)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	publicKey := privateKey.Public()
//	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
//	if !ok {
//		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
//	}
//
//	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
//	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	gasPrice, err := client.SuggestGasPrice(context.Background())
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	auth := bind.NewKeyedTransactor(privateKey)
//	auth.Nonce = big.NewInt(int64(nonce))
//	auth.Value = big.NewInt(0)     // in wei
//	auth.GasLimit = uint64(300000) // in units
//	auth.GasPrice = gasPrice
//
//	input := "1.0"
//	address, tx, instance, err := DeployStore(auth, client, input)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
//	fmt.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0
//
//	_ = instance
//}

//func TestTokenDeploy(t *testing.T) {
//	client1Url := "http://127.0.0.1:8101"
//	client, err1 := ethclient.Dial(client1Url)
//	if err1 != nil {
//		t.Log(err1)
//		t.FailNow()
//	}
//
//	malice := eth.CreateAccount()
//	fmt.Println("Created addresses for malice", malice.Address)
//	eth.SetEtherBase(client1Url, malice.Address)
//	fmt.Println("Starting mining")
//	eth.StartMining(client1Url, 1)
//	//sleep 20 seconds to accumulate eth and stop mining
//	fmt.Println("Waiting for 20 seconds")
//	time.Sleep(20 * time.Second)
//	fmt.Println("After 20 seconds...")
//	maliceBalance := eth.GetBalance(client, malice.Address)
//	if maliceBalance == 0 {
//		t.Log("Malice has no balance to do anything")
//		t.FailNow()
//	}
//	privateKey, err := crypto.HexToECDSA(malice.PrivateKey)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	publicKey := privateKey.Public()
//	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
//	if !ok {
//		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
//	}
//
//	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
//	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	gasPrice, err := client.SuggestGasPrice(context.Background())
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	auth := bind.NewKeyedTransactor(privateKey)
//	auth.Nonce = big.NewInt(int64(nonce))
//	auth.Value = big.NewInt(0)     // in wei
//	auth.GasLimit = uint64(3000000) // in units
//	auth.GasPrice = gasPrice
//
//	address, tx, instance, err := DeployErc20(auth, client)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
//	fmt.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0
//
//	_ = instance
//}
