package main

import (
	"broker/event"
	"broker/logs/logs"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

const PORT = 8081

type Config struct {
	Emitter    *event.Emitter
	GRPCCLient GRPCCLient
}

type GRPCCLient struct {
	LogClient logs.LogServiceClient
}

func main() {
	var err error
	var app = Config{}
	//var grpcConn *grpc.ClientConn

	// rabbitmq client
	//amqpConnection, err := connect()
	//if err != nil {
	//	panic(err.Error())
	//}
	//emitter, err := event.NewEventEmitter(amqpConnection)
	//if err != nil {
	//	panic(err.Error())
	//}
	//app.Emitter = &emitter

	// setup grpc
	//grpcConn = app.ConnectGRPC()
	//defer grpcConn.Close()
	//app.GRPCCLient.LogClient = logs.NewLogServiceClient(grpcConn)

	var srv = &http.Server{
		Addr:    fmt.Sprintf(":%v", PORT),
		Handler: app.routes(),
	}

	log.Println("starting broker service on port ", PORT)

	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		c, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s", os.Getenv("BASE_URL")))
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}

func (app Config) ConnectGRPC() *grpc.ClientConn {
	var err error
	var connection *grpc.ClientConn

	path := fmt.Sprintf("%s:%s", os.Getenv("BASE_URL"), os.Getenv("RPC_PORT"))
	connection, err = grpc.Dial(path, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("gprc client connected")

	return connection

}
