package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx context.Context
var err error
var client *mongo.Client

func init() {

	ctx = context.Background()
	client, err = mongo.Connect(ctx,
		options.Client().ApplyURI("mongodb+srv://kirankalidass:kalidasskiran@cluster0.xe3j2.mongodb.net/test"))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

}

type Table struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	TableNumber string             `json:"tnumber" bson:"tnumber"`
	Color       string             `json:"color" bson:"color"`
	Location    string             `json:"location" bson:"location"`
	Waiter      string             `json:"waiter" bson:"waiter"`
}

func NewTableHandler(c *gin.Context) {
	var table Table
	collection := client.Database("restaurant").Collection("Tables")
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	table.ID = primitive.NewObjectID()
	_, err = collection.InsertOne(ctx, table)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Error while inserting a new Table"})
		return
	}
	c.JSON(http.StatusOK, table)
}

func ListTableHandler(c *gin.Context) {
	collection := client.Database("restaurant").Collection("Tables")
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)
	tables := make([]Table, 0)
	for cur.Next(ctx) {
		var table Table
		cur.Decode(&table)
		tables = append(tables, table)
	}
	c.JSON(http.StatusOK, tables)
}

func UpdateTableHandler(c *gin.Context) {
	id := c.Param("id")
	var table Table
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	objectId, _ := primitive.ObjectIDFromHex(id)
	collection := client.Database("restaurant").Collection("Tables")
	_, err = collection.UpdateOne(ctx, bson.M{
		"_id": objectId,
	}, bson.D{{"$set", bson.D{
		{"tnumber", table.TableNumber},
		{"color", table.Color},
		{"location", table.Location},
		{"waiter", table.Waiter},
	}}})
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "table has been updated"})
}
func DeleteTableHandler(c *gin.Context) {

	collection := client.Database("restaurant").Collection("Tables")
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := collection.DeleteOne(ctx,
		bson.M{
			"_id": objectId,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Table has been deleted"})
}

func GetTableHandler(c *gin.Context) {

	collection := client.Database("restaurant").Collection("Tables")
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)
	cur := collection.FindOne(ctx, bson.M{
		"_id": objectId,
	})
	var table Table
	err := cur.Decode(&table)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, table)
}

//tables end

func main() {
	router := gin.Default()
	//tables
	router.GET("/table", ListTableHandler)
	router.POST("/table", NewTableHandler)
	router.PUT("/table/:id", UpdateTableHandler)
	router.DELETE("/table/:id", DeleteTableHandler)
	router.GET("/table/:id", GetTableHandler)
	//tables end
	router.Run()
}
