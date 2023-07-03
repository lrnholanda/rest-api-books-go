package main

import (
	"context"
	"fmt"
	"gin"

	"github.com/qiniu/qmgo"
)

var database *qmgo.Database
var collection *qmgo.Collection

func main() {

	// create new Client
	const databaseURI = "mongodb://localhost:27017"
	fmt.Println("Connecting to database", databaseURI)
	ctx := context.Background()
	connecting, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: databaseURI})

	database = connecting.Database("test")    // creating Database connection
	collection = database.Collection("books") // get the collection
	defer func() {
		if err = connecting.Close(ctx); err != nil {
			fmt.Println("Closing Connection to database", databaseURI)
			panic(err)
		}
	}()

	router := gin.Default() // create router using gin
}
