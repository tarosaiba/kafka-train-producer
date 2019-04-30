package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Shopify/sarama"
	"github.com/labstack/echo"
)

type (
	handler struct {
		producer sarama.SyncProducer
	}

	message struct {
		ID   int    `json:"id"`
		Body string `json:"body"`
	}
)

var (
	seq = 1
)

//const (
//	topic = "test_topic"
//)
var topic string

func NewHandler(prd sarama.SyncProducer) *handler {
	return &handler{prd}
}

func (h handler) SendMessage(c echo.Context) error {
	m := &message{
		ID: seq,
	}
	if err := c.Bind(m); err != nil {
		return err
	}
	seq++

	// create producer
	publish(m.Body, h.producer)

	return c.JSON(http.StatusCreated, m)
}

func publish(message string, producer sarama.SyncProducer) {

	// Topic
	topic = readFromENV("KAFKA_TOPIC", "test_topic")

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

    fmt.Printf("Partition: %d Offset: %d ", p, o)
}

func readFromENV(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}
	return value
}
