package handler

import (
	"fmt"
	"net/http"

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

const (
	//kafkaConn = "kafka:9092"
	topic = "test_topic"
)

func NewHandler(prd sarama.SyncProducer) *handler {
	return &handler{prd}
}

func (h handler) sendMessage(c echo.Context) error {
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
