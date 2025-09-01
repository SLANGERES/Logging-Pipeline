package broker

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	model "github.com/SLANGERES/go-service/internal/Model"
	"github.com/rabbitmq/amqp091-go"
)

func rabbitURL() string {
	if v := os.Getenv("RABBITMQ_URL"); v != "" {
		return v
	}
	return "amqp://guest:guest@rabbitmq:5672/"
}

// Connect establishes a connection to RabbitMQ
func connect() (*amqp091.Connection, error) {
	conn, err := amqp091.Dial(rabbitURL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	return conn, nil
}

// Connection pool for RabbitMQ
type RabbitMQPool struct {
	conn    *amqp091.Connection
	channels chan *amqp091.Channel
	mu      sync.RWMutex
}

var (
	pool     *RabbitMQPool
	poolOnce sync.Once
)

// InitializePool creates a connection pool with multiple channels
func InitializePool() error {
	var err error
	poolOnce.Do(func() {
		pool = &RabbitMQPool{
			channels: make(chan *amqp091.Channel, 10), // Pool of 10 channels
		}

		// Create connection
		pool.conn, err = connect()
		if err != nil {
			return
		}

		// Pre-create channels
		for i := 0; i < 10; i++ {
			ch, chErr := pool.conn.Channel()
			if chErr != nil {
				err = fmt.Errorf("failed to create channel: %w", chErr)
				return
			}
			// Declare queue on each channel
			_, queueErr := ch.QueueDeclare(
				"logs", // name
				true,   // durable
				false,  // delete when unused
				false,  // exclusive
				false,  // no-wait
				nil,    // arguments
			)
			if queueErr != nil {
				err = fmt.Errorf("failed to declare queue: %w", queueErr)
				return
			}
			pool.channels <- ch
		}
	})
	return err
}

// GetChannel gets a channel from the pool
func (p *RabbitMQPool) getChannel() *amqp091.Channel {
	select {
	case ch := <-p.channels:
		return ch
	default:
		// If no channel available, create a new one
		ch, err := p.conn.Channel()
		if err != nil {
			return nil
		}
		// Declare queue
		_, err = ch.QueueDeclare(
			"logs", // name
			true,   // durable
			false,  // delete when unused
			false,  // exclusive
			false,  // no-wait
			nil,    // arguments
		)
		if err != nil {
			return nil
		}
		return ch
	}
}

// ReturnChannel returns a channel to the pool
func (p *RabbitMQPool) returnChannel(ch *amqp091.Channel) {
	select {
	case p.channels <- ch:
		// Channel returned to pool
	default:
		// Pool full, close channel
		ch.Close()
	}
}

// SendDataToQueue publishes a Logging struct to RabbitMQ using connection pool
func SendDataToQueue(logEntry model.Logging) error {
	if pool == nil {
		return fmt.Errorf("RabbitMQ pool not initialized")
	}

	ch := pool.getChannel()
	if ch == nil {
		return fmt.Errorf("failed to get channel from pool")
	}
	defer pool.returnChannel(ch)

	body, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	err = ch.Publish(
		"",     // exchange
		"logs", // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
			DeliveryMode: amqp091.Persistent, // Make messages persistent
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// SendBatchToQueue publishes multiple log entries efficiently
func SendBatchToQueue(logEntries []model.Logging) error {
	if pool == nil {
		return fmt.Errorf("RabbitMQ pool not initialized")
	}

	ch := pool.getChannel()
	if ch == nil {
		return fmt.Errorf("failed to get channel from pool")
	}
	defer pool.returnChannel(ch)

	// Use confirm mode for reliable publishing
	if err := ch.Confirm(false); err != nil {
		return fmt.Errorf("channel could not be put into confirm mode: %w", err)
	}

	confirms := ch.NotifyPublish(make(chan amqp091.Confirmation, len(logEntries)))

	// Publish all messages in batch
	for _, logEntry := range logEntries {
		body, err := json.Marshal(logEntry)
		if err != nil {
			return fmt.Errorf("failed to marshal log entry: %w", err)
		}

		err = ch.Publish(
			"",     // exchange
			"logs", // routing key
			false,  // mandatory
			false,  // immediate
			amqp091.Publishing{
				ContentType:  "application/json",
				Body:         body,
				DeliveryMode: amqp091.Persistent,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to publish message: %w", err)
		}
	}

	// Wait for all confirmations
	timeout := time.After(5 * time.Second)
	for i := 0; i < len(logEntries); i++ {
		select {
		case confirmation := <-confirms:
			if !confirmation.Ack {
				return fmt.Errorf("message %d not acknowledged", confirmation.DeliveryTag)
			}
		case <-timeout:
			return fmt.Errorf("timeout waiting for confirmations")
		}
	}

	return nil
}
