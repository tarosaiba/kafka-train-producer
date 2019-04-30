package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tarosaiba/kafka-train-producer/handler"

	"github.com/Shopify/sarama"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//const (
//	kafkaConn = "kafka:9092"
//)

var kafkaConn string

func initProducer() (sarama.SyncProducer, error) {

	kafkaConn = readFromENV("KAFKA_BROKER", "kafka:9092")

	fmt.Println("Kafka Broker - ", kafkaConn)

	// setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// producer config
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 10
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	// async producer
	//prd, err := sarama.NewAsyncProducer([]string{kafkaConn}, config)

	// sync producer
	prd, err := sarama.NewSyncProducer([]string{kafkaConn}, config)

	return prd, err
}

func readFromENV(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}

///////////////////////////////////////////////////////////

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// kafka connection
	// wait for kafka up
	fmt.Println("Wait for Kafka up 20sec")
	time.Sleep(20000 * time.Millisecond)
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
