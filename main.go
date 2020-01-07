package main

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/satb/doublespend/eth"
)

func main() {
	client, err := ethclient.Dial("http://127.0.0.1:8101")
	// client, err := ethclient.Dial("/tmp/eth/60/01/geth.ipc")
	// account := common.HexToAddress("0x71c7656ec7ab88b098defb751b7401b5f6d8976f")

	if err != nil {
		log.Fatal(err)
	}

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

	eth.ScanBlocks(client, "0xFB9931c331f8119E406E1307A896a09029151277", 9231655, 9231658)
}
