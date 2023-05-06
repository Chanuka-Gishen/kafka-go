package main

import (
	db "backend/internal/config"
	"backend/internal/router"

	"net/http"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	// Initialize the Kafka producer
	var producer sarama.SyncProducer
	var err error
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	config.Producer.Return.Successes = true
	producer, err = sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	// Initialize the MySQL database connection
	db.SetupDB()
	defer db.CloseDB()

	// Start the HTTP server
	http.ListenAndServe(":8080", router.NewRouter())
}
