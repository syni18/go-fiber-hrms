package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client   *mongo.Client
	Db       *mongo.Database
}

var mg MongoInstance

const dbName = "fiber-hrms"
const mongoURI = "mongodb://localhost:27017" + dbName

type Employee struct {
	ID     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	Salary float64 `json:"salary"`
	Age    int     `json:"age"`
}

func Connect() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		return err
	}

	mg = MongoInstance{
		Client: client,
		Db: db,
	}
	return nil
}

func main() {

	if err := Connect(); err != nil {
		log.Fatal(err)
	}
	app := fiber.New()

	app.Get("/employee", func(c *fiber.Ctx) error {

		query := bson.D{{}}
		cursor, err := mg.Db.Collection("employees").Find(c.Context(),query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		var employees []Employee = make([]Employee, 0)
		if err := cursor.All(c.Context(), &employees); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(employees)
	})
	
	app.Post("/employee")
	app.Put("/employee/:id")
	app.Delete("/employee/:id")
}
