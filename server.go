package main

import (
	"fmt"
	"log"
	"os"
    
    "github.com/tarosaiba/kafka-train-producer/handler"

	"github.com/Shopify/sarama"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	kafkaConn = "kafka:9092"
	//topic     = "test_topic"
)

func initProducer() (sarama.SyncProducer, error) {
	// setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// producer config
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	// async producer
	//prd, err := sarama.NewAsyncProducer([]string{kafkaConn}, config)

	// sync producer
	prd, err := sarama.NewSyncProducer([]string{kafkaConn}, config)

	return prd, err
}

///////////////////////////////////////////////////////////

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// kafka connection
	producer, err := initProducer()
	if err != nil {
		fmt.Println("Error producer: ", err.Error())
		os.Exit(1)
	}

	h := handler.NewHandler(producer)

	// Routes
	e.POST("/kafka", h.SendMessage)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
