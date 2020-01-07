package eth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
)

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
