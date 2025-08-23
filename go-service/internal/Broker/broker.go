package broker

import (
	"encoding/json"
	"fmt"

	model "github.com/SLANGERES/go-service/internal/Model"
	"github.com/rabbitmq/amqp091-go"
)

// Connect establishes a connection to RabbitMQ
func connect() (*amqp091.Connection, error) {
	conn, err := amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	return conn, nil
}

// SendDataToQueue publishes a Logging struct to RabbitMQ

func SendDataToQueue(logEntry model.Logging) error {
	// Step 1: Connect
	conn, err := connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	// Step 2: Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	// Step 3: Declare queue
	queueName := "logs"
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Step 4: Marshal struct into JSON
	body, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	// Step 5: Publish message
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
