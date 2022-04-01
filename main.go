package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Predictions struct {
	Date       time.Time `json:"Date" bson:"Date"`
	Prediction float64   `json:"Prediction" bson:"Prediction"`
}

type Prediction struct {
	mgm.DefaultModel `bson:",inline"`
	ItemId           int64         `json:"Item_id" bson:"Item_id" xml:"Item_id" form:"Item_id"`
	Predictions      []Predictions `json:"predictions" bson:"predictions"`
}

type TempResponse struct {
	Success     bool         `json:"success" bson:"success" xml:"success" form:"success"`
	Collections []Prediction `json:"collections" bson:"collections"`
}

func CreateOrder(Item_id int64, Predictions []Predictions) *Prediction {
	return &Prediction{
		ItemId:      Item_id,
		Predictions: Predictions,
	}
}

func main() {

	app := fiber.New()
	app.Get("/getItem", Itemlist())
	app.Listen(":2303")
}

func Itemlist() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Append("Access-Control-Allow-Origin", "*")
		// var result Prediction
		var response = TempResponse{}
		var item_collection []Prediction

		clientOptions := options.Client().ApplyURI("mongodb+srv://appistock:oF07QuOZc57MXJQ9@cluster0.eqh4c.mongodb.net/admin?authSource=admin&replicaSet=atlas-sffvud-shard-0&readPreference=primary&appname=MongoDB%20Compass&ssl=true")
		// fmt.Println("clientOptions type:", reflect.TypeOf(clientOptions), "\n")
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			fmt.Println("mongo.Connect() ERROR:", err)
			os.Exit(1)
		}
		// Declare Context type object for managing multiple API requests
		c, _ := context.WithTimeout(context.Background(), 15*time.Second)
		// fmt.Println("this is c", c)

		// Access a MongoDB collection through a database
		col := client.Database("predictions").Collection("prediction")
		// fmt.Println("Collection type:", reflect.TypeOf(col), "\n")

		// var result Prediction

		cursor, err := col.Find(context.TODO(), bson.D{})

		// fmt.Println("this is cursor:  ", cursor)
		if err1 := cursor.All(c, &item_collection); err1 != nil {
			log.Fatal(err)

		}
		defer cursor.Close(c)
		// fmt.Println(item_collection)

		response.Success = true
		response.Collections = item_collection
		// fmt.Println("this is test  :", item_collection)
		// if result.ItemId != 0 {
		// 	response.Success = true
		// 	response.Collections = result
		// 	fmt.Println("resopnse", response)
		// }

		return ctx.Status(200).JSON(response)
	}

}
