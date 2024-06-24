package handlers

import (
	_ "broker/gRPC-client/logs"
	"broker/helpers"
	"broker/internal/config"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Meenachinmay/microservice-shared/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type RequestPayload struct {
	Action  string               `json:"action"`
	Auth    AuthPayload          `json:"auth,omitempty"`
	Log     LogJSONPayload       `json:"log,omitempty"`
	Mail    MailPayload          `json:"mail,omitempty"`
	Enquiry types.EnquiryPayload `json:"enquiry,omitempty"`
	Empty   EmptyPayload         `json:"empty,omitempty"`
}

type EmptyPayload struct{}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type EnquiryMailPayload struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogJSONPayload struct {
	ServiceName string `json:"service_name"`
	LogData     string `json:"log_data"`
}

type PaymentPayload struct {
	CardNumber     string  `json:"card_number"`
	CardHolderName string  `json:"card_holder_name"`
	CardCVV        string  `json:"card_cvv"`
	CardExpiry     string  `json:"card_expiry"`
	Amount         float64 `json:"amount"`
	Currency       string  `json:"currency"`
}

type LocalApiConfig struct {
	*config.Config
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
	case "fetch-log":
		//logItem(w, requestPayload.Log)
		lac.GetAllLogs(c)
	case "mail":
		//sendMail(c, requestPayload.Mail)
		lac.sendMailViaRabbit(c, requestPayload.Mail)
	case "enquiry_mail":
		lac.sendEnquiryMailViaRabbit(c, requestPayload.Enquiry)
	case "create_new_enquiry":
		//lac.routeEnquiryToEnquiryService(c, requestPayload.EnquiryPayload)
		lac.EnquiryViaGRPC(c, requestPayload.Enquiry)
	case "fetch-all-properties":
		lac.FetchAllProperties(c)

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
func logItem(c *gin.Context, log LogJSONPayload) {
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
