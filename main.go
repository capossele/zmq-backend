package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	
	"github.com/capossele/zmq-backend/handlers"
	"github.com/capossele/zmq-backend/scheme"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBNAME Database name
const DBNAME = "devnet"

// COLLECTION Collection name
const COLLECTION = "txs"

// CONNECTIONSTRING DB connection string
const CONNECTIONSTRING = "mongodb://localhost:27017"

func init() {
	// Populates database with dummy data

	var txs []models.Tx

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

	//fmt.Println("Connected to MongoDB!")

	// Collection types can be used to access the database
	db := client.Database(DBNAME)

	// Load values from JSON file to model
	byteValues, err := ioutil.ReadFile("tx_data.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(byteValues, &txs)

	// Insert txs into DB
	var transactions []interface{}
	for _, t := range txs {
		transactions = append(transactions, t)
	}
	_, err = db.Collection(COLLECTION).InsertMany(context.Background(), transactions)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("demo tx loaded correctly")
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/txs", handlers.GetAllTxsEndpoint).Methods("GET")
	router.HandleFunc("/txs/{hash}", handlers.GetTxEndpoint).Methods("GET")
	router.HandleFunc("/txs", handlers.CreateTxEndpoint).Methods("POST")
	router.HandleFunc("/txs", handlers.DeleteTxEndpoint).Methods("DELETE")
	fmt.Println("Starting server on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
