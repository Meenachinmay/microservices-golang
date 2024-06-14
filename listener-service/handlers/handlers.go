package handlers

// import (
//
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"github.com/confluentinc/confluent-kafka-go/kafka"
//	"log"
//	"net/http"
//
// )
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

//
//func ConsumeMessages(consumer *kafka.Consumer) {
//	for {
//		msg, err := consumer.ReadMessage(-1)
//		if err != nil {
//			fmt.Printf("consumer error %v (%v)\n", err, msg)
//			continue
//		}
//
//		// handle messages based on topic
//		switch *msg.TopicPartition.Topic {
//		case "new-log":
//			var payload LogPayload
//			err := json.Unmarshal(msg.Value, &payload)
//			if err != nil {
//				log.Printf("Error unmarshalling payload: %v\n", err)
//				continue
//			}
//			err = logEvent(payload)
//			if err != nil {
//				log.Printf("Error logging event: %v\n", err)
//				continue
//			} else {
//				log.Println("Event logged successfully:LISTENER<ConsumerMessage-Method>")
//			}
//		default:
//			log.Print("Unrecognized payload:\n")
//		}
//	}
//}
//
//func logEvent(payload LogPayload) error {
//	jsonData, _ := json.MarshalIndent(payload, "", "\t")
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
