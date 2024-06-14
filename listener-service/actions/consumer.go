package actions

//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	amqp "github.com/rabbitmq/amqp091-go"
//	"log"
//	"net/http"
//)
//
//type Consumer struct {
//	conn *amqp.Connection
//}
//
//func NewConsumer(conn *amqp.Connection) (*Consumer, error) {
//	consumer := &Consumer{
//		conn: conn,
//	}
//
//	err := consumer.setup()
//	if err != nil {
//		return nil, err
//	}
//
//	return consumer, nil
//
//}
//
//func (consumer *Consumer) setup() error {
//	channel, err := consumer.conn.Channel()
//
//	if err != nil {
//		return err
//	}
//
//	return declareExchange(channel)
//}
//
//type LogPayload struct {
//	Name string `json:"name"'`
//	Data string `json:"data"'`
//}
//
//func (consumer *Consumer) ConsumeLogs(topics []string) error {
//	ch, err := consumer.conn.Channel()
//	if err != nil {
//		return err
//	}
//	defer ch.Close()
//
//	queue, err := declareRandomQueue(ch)
//	if err != nil {
//		return err
//	}
//
//	for _, s := range topics {
//		ch.QueueBind(queue.Name, s, "log_topics", false, nil)
//		if err != nil {
//			return err
//		}
//	}
//
//	messages, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
//	if err != nil {
//		return err
//	}
//
//	forever := make(chan bool)
//	go func() {
//		for d := range messages {
//			var payload LogPayload
//			_ = json.Unmarshal(d.Body, &payload)
//			go handlePayload(payload)
//		}
//	}()
//
//	fmt.Printf("Waiting for messages [Exchange, Queue] [log_topics, %s]\n", queue.Name)
//	<-forever
//
//	return nil
//}
//
//func handlePayload(payload LogPayload) {
//	switch payload.Name {
//	case "log":
//		err := logEvent(payload)
//		if err != nil {
//			log.Println(err)
//		}
//
//	case "auth":
//		// authenticate
//
//	default:
//		log.Printf("Unrecognized payload '%s'\n", payload.Name)
//	}
//}
//
//func logEvent(log LogPayload) error {
//	jsonData, _ := json.MarshalIndent(log, "", "\t")
//
//	logServiceURL := "http://logger-service/log"
//
//	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
//	if err != nil {
//		return err
//	}
//
//	request.Header.Set("Content-Type", "application/json")
//
//	client := &http.Client{}
//	response, err := client.Do(request)
//	if err != nil {
//		return err
//	}
//
//	defer response.Body.Close()
//
//	// response
//	if response.StatusCode != http.StatusAccepted {
//		return err
//	}
//
//	return nil
//}
