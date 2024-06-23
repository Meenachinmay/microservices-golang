package consumers

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type LogConsumer struct {
	conn *amqp.Connection
}

type LogPayload struct {
	ServiceName string `json:"service_name"`
	LogData     string `json:"log_data"`
}

func NewLogConsumer(conn *amqp.Connection) (*LogConsumer, error) {
	consumer := &LogConsumer{conn: conn}
	if err := consumer.setup(); err != nil {
		return nil, err
	}
	return consumer, nil
}

func (consumer *LogConsumer) setup() error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	return DeclareExchange(ch)
}

func (consumer *LogConsumer) ConsumeLogs(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	queue, err := DeclareRandomQueue(ch)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	for _, topic := range topics {
		if err := ch.QueueBind(queue.Name, topic, "log_topics", false, nil); err != nil {
			return fmt.Errorf("failed to bind a queue: %v", err)
		}
	}

	messages, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	go func() {
		for d := range messages {
			var payload LogPayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Printf("Failed to unmarshal message:CONSUMELOGS %v", err)
				continue
			}
			handleLogPayload(payload)
		}
	}()

	log.Printf("Waiting for log messages [Exchange, Queue] [log_topics, %s]\n", queue.Name)

	// handle termination signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	log.Printf("Received shutdown signal, exiting...\n")

	return nil
}

func handleLogPayload(payload LogPayload) {
	if err := logEvent(payload); err != nil {
		log.Printf("Failed to log event: %v", err)
	}
}

func logEvent(log LogPayload) error {
	jsonData, _ := json.MarshalIndent(log, "", "\t")
	logServiceURL := "http://logger-service/log"

	req, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request:LOGEVENT: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request:LOGEVENT: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected response status:LOGEVENT: %d", resp.StatusCode)
	}
	return nil
}
