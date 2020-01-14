package eth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/google/uuid"
)

type EthError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ethResponse struct {
	ID      string          `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *EthError       `json:"error"`
}

func GetNodeInfo(clientUrl string) p2p.NodeInfo {
	id := uuid.New()
	message := map[string]interface{}{
		"id":      id.String(),
		"jsonrpc": "2.0",
		"method":  "admin_nodeInfo",
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
	var eresp ethResponse
	err = json.Unmarshal(body, &eresp)
	if err != nil {
		log.Println("Error unmarshaling", err)
	}
	var nodeInfo p2p.NodeInfo
	err = json.Unmarshal(eresp.Result, &nodeInfo)
	if err != nil {
		log.Fatalln("error:", err)
	}
	log.Println("Fetched node info with enode", nodeInfo.Enode)
	return nodeInfo
}

func SetEtherBase(clientUrl string, address string) {
	id := uuid.New()
	params := [...]string{address}
	message := map[string]interface{}{
		"id":      id.String(),
		"jsonrpc": "2.0",
		"method":  "miner_setEtherbase",
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

func AddPeer(clientUrl string, peerEnode string) {
	id := uuid.New()
	params := [...]string{peerEnode}
	message := map[string]interface{}{
		"id":      id.String(),
		"jsonrpc": "2.0",
		"method":  "admin_addPeer",
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

func RemovePeer(clientUrl string, peerEnode string) {
	id := uuid.New()
	params := [...]string{peerEnode}
	message := map[string]interface{}{
		"id":      id.String(),
		"jsonrpc": "2.0",
		"method":  "admin_removePeer",
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

func StartMining(clientUrl string, numThreads int) {
	id := uuid.New()
	params := [...]int{numThreads}
	message := map[string]interface{}{
		"id":      id.String(),
		"jsonrpc": "2.0",
		"method":  "miner_start",
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

func StopMining(clientUrl string) {
	id := uuid.New()
	params := [...]int{}
	message := map[string]interface{}{
		"id":      id.String(),
		"jsonrpc": "2.0",
		"method":  "miner_stop",
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
