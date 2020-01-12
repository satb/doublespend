package main

func main() {
	// client, err := ethclient.Dial("http://127.0.0.1:8101")
	// client, err := ethclient.Dial("/tmp/eth/60/01/geth.ipc")
	// account := common.HexToAddress("0x71c7656ec7ab88b098defb751b7401b5f6d8976f")

	// client, err := ethclient.Dial("https://mainnet.infura.io")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// wssClient, err := ethclient.Dial("ws://127.0.0.1:4101")
	// client, err := ethclient.Dial("/tmp/eth/60/01/geth.ipc")
	// account := common.HexToAddress("0x71c7656ec7ab88b098defb751b7401b5f6d8976f")

	// client, err := ethclient.Dial("https://mainnet.infura.io")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("we have a connection")

	// balance, err := client.BalanceAt(context.Background(), account, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(balance)

	// id := uuid.New()
	// message := map[string]interface{}{
	// 	"id":      id.String(),
	// 	"jsonrpc": "2.0",
	// 	"method":  "admin_datadir",
	// }

	// bytesRepresentation, err := json.Marshal(message)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// resp, err := http.Post("http://127.0.0.1:8101", "application/json", bytes.NewBuffer(bytesRepresentation))

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// log.Println(string(body))

	// txns := Scan(client, "0x69259375afd0944fa345f9430b3cde078ab0eb76")
	// log.Fatal(txns)

	// fmt.Println("Subscribing to blocks")
	// Subscribe(wssClient)
}
