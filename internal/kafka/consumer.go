package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func createLogFile() (*os.File, error) {
	dir := "/app/logs"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("не удалось создать директорию %s: %v", dir, err)
	}

	file, err := os.OpenFile("/app/logs/kafka_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error creating log file: %v", err)
		return nil, err
	}
	return file, nil
}

func StartConsumer(brokers []string, topic string) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
	defer consumer.Close()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatalf("Failed to get partitions for topic %s: %v", err)
	}

	logFile, err := createLogFile()
	if err != nil {
		log.Fatalf("Could not create log file: %v", err)
	}
	defer logFile.Close()

	log.Println("Kafka consumer started. Listening for messages...")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	for _, partition := range partitions {
		go func(partition int32) {
			partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				log.Printf("Failed to start partition consumer for partition %d: %v", partition, err)
				return
			}
			defer partitionConsumer.Close()

			// Прослушивание сообщений для этой партиции
			for {
				select {
				case msg := <-partitionConsumer.Messages():
					logEntry := fmt.Sprintf("[%s] %s\n", string(msg.Key), string(msg.Value))
					logFile.WriteString(logEntry)
					fmt.Print(logEntry)
				case err := <-partitionConsumer.Errors():
					log.Printf("Error consuming messages from partition %d: %v", partition, err)
				case <-signals:
					log.Println("Kafka consumer shutting down...")
					return
				}
			}
		}(partition)
	}

	// Ожидание сигнаов для завершения работы
	<-signals
}
