package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/capossele/zmq-backend/handlers"
	"github.com/gorilla/mux"
	"github.com/iotaledger/iota.go/transaction"
	"github.com/iotaledger/iota.go/trinary"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	czmq "github.com/zeromq/goczmq"
)

// DBNAME Database name
const DBNAME = "devnet"

// COLLECTION Collection name
const COLLECTION = "txs"

// CONNECTIONSTRING DB connection string
const CONNECTIONSTRING = "mongodb://localhost:27017"

func init() {
	// Populates database with dummy data

	//var txs []models.Tx

	clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Collection types can be used to access the database
	//db := client.Database(DBNAME)

	// Load values from JSON file to model
	// byteValues, err := ioutil.ReadFile("tx_data.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// json.Unmarshal(byteValues, &txs)

	// // Insert txs into DB
	// var transactions []interface{}
	// for _, t := range txs {
	// 	transactions = append(transactions, t)
	// }
	// _, err = db.Collection(COLLECTION).InsertMany(context.Background(), transactions)
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println("demo tx loaded correctly")
	// }
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/txs", handlers.GetAllTxsEndpoint).Methods("GET")
	router.HandleFunc("/txs/{hash}", handlers.GetTxEndpoint).Methods("GET")
	router.HandleFunc("/txs", handlers.CreateTxEndpoint).Methods("POST")
	router.HandleFunc("/txs", handlers.DeleteAllTxsEndpoints).Methods("DELETE")
	fmt.Println("Starting server on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))

}

type Tx struct {
	Hash               trinary.Hash
	id                 int
	time               int64
	cw                 int
	cw2                int // TODO: to remove, used only to compare different CW update mechanisms
	ref                []int
	refHash            []trinary.Hash
	app                []int
	appHash            []trinary.Hash
	bundle             trinary.Hash
	firstApproval      float64
	bundleCurrentIndex uint64
	bundleLastIndex    uint64
	isTip              bool
	//trunkHash          trinary.Hash
	//branchHash         trinary.Hash
}

func zmqService() {
	pubEndpoint := "tcp://35.246.106.61:5556"

	topics := "tx"

	subSock, err := czmq.NewSub(pubEndpoint, topics)
	if err != nil {
		panic(err)
	}

	defer subSock.Destroy()

	fmt.Printf("Collecting updates from IRI Node for %s…\n", topics)
	subSock.Connect(pubEndpoint)

	timestampMap := make(map[trinary.Hash]int64)
	trytesMap := make(map[trinary.Hash]trinary.Trytes)
	var timeOffset int64
	//for i := 0; i < 1000; {
	for {
		msg, _, err := subSock.RecvFrame()
		if err != nil {
			panic(err)
		}

		data := strings.Split(string(msg), " ")
		now := time.Now()
		nowMillis := now.UnixNano() /// 1000000
		if timeOffset == 0 {
			timeOffset = nowMillis
		}

		if data[0] == "tx" {
			hash := data[1]
			timestampMap[hash] = nowMillis - timeOffset
			fmt.Println(data[0], hash, timestampMap[hash])
			//i++
		} else if data[0] == "tx_trytes" {
			txObject, err := transaction.AsTransactionObject(data[1])
			if err != nil {
				//log.Fatal(err)
				fmt.Println("ERROR on AsTransactionObject", err)
			}
			txObject.AttachmentTimestamp = timestampMap[txObject.Hash]
			trytes, err := transaction.TransactionToTrytes(txObject)
			if err != nil {
				//log.Fatal(err)
				fmt.Println("ERROR on TransactionToTrytes", err)
			}
			trytesMap[txObject.Hash] = trytes
			//i++
		}
	}
	for k, v := range timestampMap {
		fmt.Println(k, v)
	}
	fmt.Println(trytesMap)
}
