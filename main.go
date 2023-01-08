package main

import (
	"context"
	"log"
	"time"

	"github.com/RuichenHe/GoCats/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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

func main() {
	//Connect to database
	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	//Init fiber
	app := fiber.New()
	//First op: add one cat
	app.Post("/cat", func(c *fiber.Ctx) error {
		//Obtain the cats collection from the database
		collection := mg.Db.Collection("cats")
		cat := new(models.Cat)
		if err := c.BodyParser(cat); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		cat.ID = ""
		//Insert the recieved json document to the database collection
		insertedId, err := collection.InsertOne(c.Context(), cat)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		//Use the insertedId to access the inserted cat info
		query := bson.D{{Key: "_id", Value: insertedId.InsertedID}}
		createdRecord := collection.FindOne(c.Context(), query)
		createdCat := &models.Cat{}
		createdRecord.Decode(createdCat)
		return c.Status(201).JSON(createdCat)
	})

	//Listen at localhost:3000
	log.Fatal(app.Listen(":3000"))
}
