package main

import (
	"context"
	"log"
	"time"

	"github.com/RuichenHe/GoCats/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	//Second op: get all cats
	app.Get("/cats", func(c *fiber.Ctx) error {
		query := bson.D{{}}

		cursor, err := mg.Db.Collection("cats").Find(c.Context(), query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		//Create a cats slice, and use All func to populate cats with all of the query results
		var cats []models.Cat = make([]models.Cat, 0)
		if err := cursor.All(c.Context(), &cats); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(cats)
	})
	//Third op: get one cat based on cat name
	app.Get("/cat/:name", func(c *fiber.Ctx) error {
		catName := c.Params("name")
		query := bson.M{"name": catName}
		cat := new(models.Cat)
		err := mg.Db.Collection("cats").FindOne(c.Context(), query).Decode(&cat)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		return c.Status(200).JSON(cat)
	})
	//Fourth op: update a cat info based on id
	app.Put("/cat/:id", func(c *fiber.Ctx) error {
		inputId := c.Params("id")
		catId, err := primitive.ObjectIDFromHex(inputId)
		if err != nil {
			return c.SendStatus(400)
		}
		query := bson.M{"_id": catId}
		cat := new(models.Cat)
		if err := c.BodyParser(cat); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		originalCat := new(models.Cat)
		err = mg.Db.Collection("cats").FindOne(c.Context(), query).Decode(&originalCat)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
		cat.CheckField(originalCat)
		newInfo := bson.D{
			{
				Key: "$set",
				Value: bson.D{
					{Key: "name", Value: cat.Name},
					{Key: "brand", Value: cat.Brand},
					{Key: "age", Value: cat.Age},
					{Key: "gender", Value: cat.Gender},
					{Key: "color", Value: cat.Color},
					{Key: "weight", Value: cat.Weight},
				},
			},
		}
		err = mg.Db.Collection("cats").FindOneAndUpdate(c.Context(), query, newInfo).Err()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.SendStatus(400)
			}
			return c.SendStatus(500)
		}
		cat.ID = inputId
		return c.Status(200).JSON(cat)
	})
	//Listen at localhost:3000
	log.Fatal(app.Listen(":3000"))
}
