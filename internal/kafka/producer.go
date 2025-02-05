package kafka

import (
	"github.com/IBM/sarama"
	"log"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(address []string) *Producer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer(address, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	return &Producer{producer: p}
}

func (p *Producer) ProduceMessage(topic, key, message string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),     // тип запроса
		Value: sarama.StringEncoder(message), // текст о выполнении запроса
	}

	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func (p *Producer) Close() {
	p.producer.Close()
}
