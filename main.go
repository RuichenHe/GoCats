package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg mongoInstance

const dbName = "catsManagement"
const mongoURI = "mongodb://localhost:27017/" + dbName

func Connect() error {
	//Create a new client to connect to a deployment specified by the uri
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}
	//Use context.WithTimeout to handle blocked operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	//Access the database, if not exist, create one
	db := client.Database(dbName)
	//Update the mg with the client and db
	mg = mongoInstance{
		Client: client,
		Db:     db,
	}
	return nil
}
