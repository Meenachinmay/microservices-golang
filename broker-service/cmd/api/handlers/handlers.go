package handlers

import (
	"broker/actions"
	"broker/helpers"
	"broker/internal/config"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type LocalApiConfig struct {
	*config.Config
	//Producer *kafka.Producer
}

func (lac *LocalApiConfig) Broker(c *gin.Context) {
	payload := helpers.JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = helpers.WriteJSON(c, http.StatusOK, payload)
}

func (lac *LocalApiConfig) HandleSubmission(c *gin.Context) {
	var requestPayload RequestPayload

	err := helpers.ReadJSON(c, &requestPayload)
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	switch requestPayload.Action { // handling different actions from here
	case "auth":
		authenticate(c, requestPayload.Auth)
	case "log":
		//logItem(w, requestPayload.Log)
		lac.logEventViaRabbit(c, requestPayload.Log)
		//lac.logEventUsingKafka(w, requestPayload.Log, "new-log")
	case "mail":
		//sendMail(c, requestPayload.Mail)
		lac.sendMailViaRabbit(c, requestPayload.Mail)
	default:
		helpers.ErrorJSON(c, errors.New("invalid action"))
	}
}

// send mail
func sendMail(c *gin.Context, mail MailPayload) {
	jsonData, _ := json.MarshalIndent(mail, "", "\t")

	// call the mail service
	mailServiceURL := "http://mailer-service/send"

	// post to mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	// set header
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}
	defer response.Body.Close()

	// deal with response
	if response.StatusCode != http.StatusAccepted {
		helpers.ErrorJSON(c, errors.New("error calling mail service"))
		return
	}

	// send back json
	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = "message sent to " + mail.To

	helpers.WriteJSON(c, http.StatusAccepted, payload)
}

// log item in log service
func logItem(c *gin.Context, log LogPayload) {
	jsonData, _ := json.MarshalIndent(log, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	defer response.Body.Close()

	// response
	if response.StatusCode != http.StatusAccepted {
		helpers.ErrorJSON(c, err)
		return
	}

	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = "logged"

	helpers.WriteJSON(c, http.StatusAccepted, payload)

}

// authenticate service
func authenticate(c *gin.Context, a AuthPayload) {
	// create some json to send to the auth service
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back status code
	if response.StatusCode == http.StatusUnauthorized {
		helpers.ErrorJSON(c, errors.New("unauthorized"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		helpers.ErrorJSON(c, errors.New("error calling auth service"))
		return
	}

	// create a var we'll read response.body info
	var jsonFromService helpers.JsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	if jsonFromService.Error {
		helpers.ErrorJSON(c, err, http.StatusUnauthorized)
		return
	}

	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = "authentication successful"
	payload.Data = jsonFromService.Data

	helpers.WriteJSON(c, http.StatusAccepted, payload)
}

func (lac *LocalApiConfig) logEventViaRabbit(c *gin.Context, l LogPayload) {
	emitter, err := actions.NewEmitter(lac.Rabbit, "log_topics", "log.INFO")
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	j, _ := json.Marshal(&l)
	err = emitter.Emit(string(j))
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = "logged via rabbit"

	helpers.WriteJSON(c, http.StatusAccepted, payload)
}

func (lac *LocalApiConfig) sendMailViaRabbit(c *gin.Context, mail MailPayload) {
	emitter, err := actions.NewEmitter(lac.Rabbit, "mail_exchange", "mail_key")
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	j, _ := json.Marshal(&mail)
	err = emitter.Emit(string(j))
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = "sent mail via rabbit"

	helpers.WriteJSON(c, http.StatusAccepted, payload)
}

//func (lac *LocalApiConfig) logEventUsingKafka(w http.ResponseWriter, logPayload LogPayload, topic string) {
//	message, err := json.Marshal(logPayload)
//	if err != nil {
//		helpers.ErrorJSON(w, err)
//	}
//
//	err = lac.Producer.Produce(&kafka.Message{
//		TopicPartition: kafka.TopicPartition{
//			Topic: &topic, Partition: kafka.PartitionAny,
//		},
//		Value: message,
//	}, nil)
//	if err != nil {
//		log.Fatalln("failed to produce log message:", err)
//		helpers.ErrorJSON(w, err)
//	}
//
//	var payload helpers.JsonResponse
//	payload.Error = false
//	payload.Message = "logged via kafka"
//
//	helpers.WriteJSON(w, http.StatusAccepted, payload)
//}

// func (lac *LocalApiConfig) produceLogToKafka(producer *kafka.Producer, topic string, message []byte) {
//
// }

//func (lac *LocalApiConfig) pushToQueue(name, msg string) error {
//	emitter, err := actions.NewEventEmitter(lac.Rabbit)
//	if err != nil {
//		return err
//	}
//
//	payload := LogPayload{
//		Name: name,
//		Data: msg,
//	}
//
//	j, _ := json.Marshal(&payload)
//	err = emitter.Emit(string(j), "log.INFO")
//	if err != nil {
//		return err
//	}
//	return nil
//}
