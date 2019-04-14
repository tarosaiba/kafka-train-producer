package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Shopify/sarama"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	message struct {
		ID   int    `json:"id"`
		Body string `json:"body"`
	}
)

const (
	kafkaConn = "kafka:9092"
	topic     = "test_topic"
)

var (
	seq = 1
)

//----------
// Handlers
//----------

func sendMessage(c echo.Context) error {
	m := &message{
		ID: seq,
	}
	if err := c.Bind(m); err != nil {
		return err
	}
	seq++

	// create producer
	producer, err := initProducer()
	if err != nil {
		fmt.Println("Error producer: ", err.Error())
		os.Exit(1)
	}
	publish(m.Body, producer)

	return c.JSON(http.StatusCreated, m)
}

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

func publish(message string, producer sarama.SyncProducer) {
	// publish sync
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	p, o, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Error publish: ", err.Error())
	}

	// publish async
	//producer.Input() <- &sarama.ProducerMessage{

	fmt.Println("Partition: ", p)
	fmt.Println("Offset: ", o)
}

///////////////////////////////////////////////////////////

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/kafka", sendMessage)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
