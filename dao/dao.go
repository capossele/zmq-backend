package dao

import (
	"context"
	"fmt"
	"log"
	"zmq-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CONNECTIONSTRING DB connection string
const CONNECTIONSTRING = "mongodb://localhost:27017"

// DBNAME Database name
const DBNAME = "devnet"

// COLLNAME Collection name
const COLLNAME = "txs"

var db *mongo.Database

// Connect establish a connection to database
func init() {
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
	db = client.Database(DBNAME)
}

// InsertManyValues inserts many items from byte slice
func InsertManyValues(txs []models.Tx) {
	var ppl []interface{}
	for _, p := range txs {
		ppl = append(ppl, p)
	}
	_, err := db.Collection(COLLNAME).InsertMany(context.Background(), ppl)
	if err != nil {
		log.Fatal(err)
	}
}

// InsertOneValue inserts one item from Tx model
func InsertOneValue(tx models.Tx) {
	fmt.Println(tx)
	_, err := db.Collection(COLLNAME).InsertOne(context.Background(), tx)
	if err != nil {
		log.Fatal(err)
	}
}

// GetAllTxs returns all txs from DB
func GetAllTxs() []models.Tx {
	cur, err := db.Collection(COLLNAME).Find(context.Background(), nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	var elements []models.Tx
	var elem models.Tx
	// Get the next result from the cursor
	for cur.Next(context.Background()) {
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		elements = append(elements, elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.Background())
	return elements
}

// GetAllTxs returns all txs from DB
func GetTx(hash string) models.Tx {
	var result models.Tx
	filter := bson.D{{"hash", hash}}
	fmt.Println("Looking for", hash, filter)
	err := db.Collection(COLLNAME).FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result

}

// DeleteTx deletes an existing tx
func DeleteTx(tx models.Tx) {
	_, err := db.Collection(COLLNAME).DeleteOne(context.Background(), tx, nil)
	if err != nil {
		log.Fatal(err)
	}
}