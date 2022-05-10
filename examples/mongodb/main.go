package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cookie-for-pres/chc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DB *mongo.Client
)

func Controller(req *chc.Request) *chc.Response {
	databases, err := DB.ListDatabaseNames(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}

	jsonData, err := json.Marshal(databases)
	if err != nil {
		panic(err)
	}

	res := req.NewResponse()
	res.SetStatusCode(200)
	res.SetHeader("Content-Type", "application/json")
	res.SetStringBody(`{"databases":` + string(jsonData) + `}`)

	return res
}

func main() {
	CHC := chc.NewCHC()
	CHC.LoadEnv(".env")

	mongoUri := CHC.GetEnv("MONGO_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	DB = client

	CHC.AddRoute(&chc.Route{
		Path:       "/",
		Methods:    []string{"GET"},
		Controller: Controller,
	})

	CHC.Listen("localhost", 8080)
}
