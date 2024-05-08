package main

import (
	"context"
	"log"
	"log-service/data"
	"os"
	"time"
)

type RPCServcer struct {
}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServcer) LoginInfo(payload RPCPayload, resp *string) error {
	collection := mongoClient.Database(os.Getenv("MONGO_INITDB_DATABASE")).Collection("logs")
	var entry data.LogEntry
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()
	_, err := collection.InsertOne(context.TODO(), entry)
	if err != nil {
		log.Println("Error inserting into logs: ", err.Error())
		return err
	}

	return nil
}
