package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/capossele/zmq-backend/models"
)

// GetTimeOfArrival returns the time of arrival (int64) of a given tx
// returns 0 if the tx was not recored
// optionally, you can give a different API uri as the second parameter
func GetTimeOfArrival(txHash string, uri ...string) (int64, error) {
	APIurl := "http://node1.iota.capossele.org:8000/txs/"
	if uri != nil {
		APIurl = uri[0] + "/txs/"
	}
	resp, err := http.Get(APIurl + txHash)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}

	var tx models.Tx
	err = json.Unmarshal(body, &tx)
	return tx.Timestamp, err
}
