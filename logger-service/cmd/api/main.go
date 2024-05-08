package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"log-service/data"
	"log-service/logs/logs"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	PORT = "80"
)

var mongoClient *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	var err error

	//os.Setenv("MONGO_URL", "mongodb://mongo:5003")
	//os.Setenv("MONGO_INITDB_DATABASE", "logs")
	//os.Setenv("MONGO_INITDB_ROOT_USERNAME", "admin")
	//os.Setenv("MONGO_INITDB_ROOT_PASSWORD", "password")

	mongoClient, err = connectToMongo()
	if err != nil {
		log.Panic(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err.Error())
		}
	}()

	app := Config{Models: data.New(mongoClient)}
	go app.gRPCListen()

	app.serve()

}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	log.Println("Server running at port", PORT)
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URL"))
	clientOptions.SetAuth(options.Credential{
		Username: os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
		Password: os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting ", err.Error())
	}
	log.Println("Connected to mongodb")
	return c, nil
}

func (app Config) gRPCListen() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("RPC_PORT")))
	if err != nil {
		log.Fatalf("Failed to listne for gRPC: %v", err.Error())
	}

	s := grpc.NewServer()
	logs.RegisterLogServiceServer(s, LogServer{Model: app.Models})

	log.Printf("gRPC serverted on port %s", os.Getenv("RPC_PORT"))

	if err = s.Serve(listen); err != nil {
		log.Fatalf("Failed to listne for gRPC: %v", err.Error())
	}
}
