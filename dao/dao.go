package dao

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/capossele/zmq-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
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

	//populate index
	PopulateIndex(DBNAME, COLLNAME, client)

}

// InsertManyValues inserts many items from byte slice
// func InsertManyValues(txs []models.Tx) {
// 	var ppl []interface{}
// 	for _, p := range txs {
// 		ppl = append(ppl, p)
// 	}
// 	_, err := db.Collection(COLLNAME).InsertMany(context.Background(), ppl)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// InsertOneValue inserts one item from Tx model
func InsertOneValue(tx models.Tx) {
	fmt.Println(tx)
	filter := bson.D{
		{"hash", tx.Hash},
	}
	update := bson.D{
		{"$set", bson.D{
			{"hash", tx.Hash},
			{"timestamp", tx.Timestamp},
		}},
	}
	updateOptions := options.Update()
	updateOptions.SetUpsert(true)
	//findOptions.SetLimit(2)
	//_, err := db.Collection(COLLNAME).InsertOne(context.Background(), tx)
	_, err := db.Collection(COLLNAME).UpdateOne(context.Background(), filter, update, updateOptions)
	if err != nil {
		log.Fatal(err)
	}
}

//GetAllTxs returns all txs from DB
func GetAllTxs() []models.Tx {

	// Pass these options to the Find method
	//findOptions := options.Find()
	//findOptions.SetLimit(2)
	cur, err := db.Collection(COLLNAME).Find(context.TODO(), bson.D{})
	//fmt.Println(cur)
	if err != nil {
		fmt.Println("first line")
		log.Fatal(err)
	}
	var elements []models.Tx
	var elem models.Tx
	// Get the next result from the cursor
	for cur.Next(context.TODO()) {
		err := cur.Decode(&elem)
		if err != nil {
			fmt.Println("second line")
			log.Fatal(err)
		}
		elements = append(elements, elem)
	}
	if err := cur.Err(); err != nil {
		fmt.Println("third line")
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	return elements
}

// GetTx returns a given tx from DB
func GetTx(hash string) models.Tx {
	var result models.Tx
	filter := bson.D{{"hash", hash}}
	//fmt.Println("Looking for", hash, filter)
	err := db.Collection(COLLNAME).FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		//log.Fatal(err)
	}

	return result

}

// DeleteAllTxs deletes all existing txs
func DeleteAllTxs() {
	_, err := db.Collection(COLLNAME).DeleteMany(context.Background(), bson.D{}, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// DeleteTx deletes an existing tx
// func DeleteTx(tx models.Tx) {
// 	_, err := db.Collection(COLLNAME).DeleteOne(context.Background(), tx, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func PopulateIndex(database, collection string, client *mongo.Client) {
	c := client.Database(database).Collection(collection)
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	index := yieldIndexModel()
	c.Indexes().CreateOne(context.Background(), index, opts)
	log.Println("Successfully created hash index")
}

func yieldIndexModel() mongo.IndexModel {
	keys := bsonx.Doc{{Key: "hash", Value: bsonx.Int32(1)}}
	index := mongo.IndexModel{}
	index.Keys = keys
	//index.Options = bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}}
	return index
}
