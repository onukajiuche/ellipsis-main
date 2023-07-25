// Uncomment this package and import libraries if using mongo
package mongo

// import (
// 	"context"
// 	"fmt"
// 	"log"

// 	"unifylogic_backend/internal/config"
// 	"unifylogic_backend/utility"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// 	"go.mongodb.org/mongo-driver/mongo/readpref"
// )

// var (
// 	ctx         context.Context
// 	mongoclient *mongo.Client
// )

// func Connection() (db *mongo.Client) {
// 	return mongoclient
// }

// func ConnectToDB() *mongo.Client {
// 	logger := utility.NewLogger()
// 	uri := config.GetConfig().Mongodb.Url
// 	mongo_connection := options.Client().ApplyURI(uri)
// 	mongoclient, err := mongo.Connect(ctx, mongo_connection)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// PINGING THE CONNECTION
// 	err = mongoclient.Ping(ctx, readpref.Primary())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// IF EVERYTHING IS OKAY, THEN CONNECT
// 	fmt.Println("MONGO CONNECTION ESTABLISHED")
// 	logger.Info("MONGO CONNECTION ESTABLISHED")

// 	return mongoclient
// }
