package handlers

import (
	"broker/actions"
	"broker/helpers"
	"broker/internal/config"
	"bytes"
	"encoding/json"
	"errors"
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
}

func (lac *LocalApiConfig) Broker(w http.ResponseWriter, r *http.Request) {
	payload := helpers.JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = helpers.WriteJSON(w, http.StatusOK, payload)
}

func (lac *LocalApiConfig) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := helpers.ReadJSON(w, r, &requestPayload)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	switch requestPayload.Action { // handling different actions from here
	case "auth":
		authenticate(w, requestPayload.Auth)
	case "log":
		//logItem(w, requestPayload.Log)
		lac.logEventViaRabbit(w, requestPayload.Log)
	case "mail":
		sendMail(w, requestPayload.Mail)
	default:
		helpers.ErrorJSON(w, errors.New("invalid action"))
	}
}

// send mail
func sendMail(w http.ResponseWriter, mail MailPayload) {
	jsonData, _ := json.MarshalIndent(mail, "", "\t")

	// call the mail service
	mailServiceURL := "http://mailer-service/send"

	// post to mail service
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	// set header
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// deal with response
	if response.StatusCode != http.StatusAccepted {
		helpers.ErrorJSON(w, errors.New("error calling mail service"))
		return
	}

	// send back json
	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = "message sent to " + mail.To

	helpers.WriteJSON(w, http.StatusAccepted, payload)
}

// log item in log service
func logItem(w http.ResponseWriter, log LogPayload) {
	jsonData, _ := json.MarshalIndent(log, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	defer response.Body.Close()

	// response
	if response.StatusCode != http.StatusAccepted {
		helpers.ErrorJSON(w, err)
		return
	}

	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = "logged"

	helpers.WriteJSON(w, http.StatusAccepted, payload)

}

// authenticate service
func authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json to send to the auth service
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back status code
	if response.StatusCode == http.StatusUnauthorized {
		helpers.ErrorJSON(w, errors.New("unauthorized"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		helpers.ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	// create a var we'll read response.body info
	var jsonFromService helpers.JsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		helpers.ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		helpers.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = "authentication successful"
	payload.Data = jsonFromService.Data

	helpers.WriteJSON(w, http.StatusAccepted, payload)
}

func (lac *LocalApiConfig) logEventViaRabbit(w http.ResponseWriter, l LogPayload) {
	err := lac.pushToQueue(l.Name, l.Data)
	if err != nil {
		helpers.ErrorJSON(w, err)
	}

	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = "logged via rabbit"

	helpers.WriteJSON(w, http.StatusAccepted, payload)
}

func (lac *LocalApiConfig) pushToQueue(name, msg string) error {
	emitter, err := actions.NewEventEmitter(lac.Rabbit)
	if err != nil {
		return err
	}

	payload := LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.Marshal(&payload)
	err = emitter.Emit(string(j), "log.INFO")
	if err != nil {
		return err
	}
	return nil
}
