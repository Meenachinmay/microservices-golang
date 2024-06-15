package consumers

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

type MailConsumer struct {
	conn *amqp.Connection
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func NewMailConsumer(conn *amqp.Connection) (*MailConsumer, error) {
	consumer := &MailConsumer{conn: conn}
	if err := consumer.setup(); err != nil {
		return nil, err
	}
	return consumer, nil
}

func (consumer *MailConsumer) setup() error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	return DeclareMailExchange(ch)
}

func (consumer *MailConsumer) ConsumeMails() error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	queue, err := DeclareMailQueue(ch)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	if err := ch.QueueBind(queue.Name, "mail_key", "mail_exchange", false, nil); err != nil {
		return fmt.Errorf("failed to bind a queue: %v", err)
	}

	messages, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload MailPayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}
			sendMail(payload)
		}
	}()

	log.Printf("Waiting for mail messages [Exchange, Queue] [mail_exchange, %s]\n", queue.Name)
	<-forever
	return nil
}

func sendMail(payload MailPayload) error {
	jsonData, _ := json.MarshalIndent(payload, "", "\t")

	// call the mail service
	mailServiceURL := "http://mailer-service/send"

	// post to mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request:SENDMAIL: %v", err)
	}

	// set header
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request:SENDMAIL: %v", err)
	}
	defer response.Body.Close()

	// deal with response
	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected response status:SENDMAIL: %d", response.StatusCode)
	}
	log.Println("sent mail via rabbit:LISTENER_SERVICE-MailConsumer")
	return nil
}
