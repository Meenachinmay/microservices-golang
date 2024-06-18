package handlers

import (
	"broker/actions"
	"broker/gRPC-client/logs"
	_ "broker/gRPC-client/logs"
	"broker/gRPC-client/payment"
	"broker/helpers"
	"broker/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"time"
)

type RequestPayload struct {
	Action         string         `json:"action"`
	Auth           AuthPayload    `json:"auth,omitempty"`
	Log            LogPayload     `json:"log,omitempty"`
	Mail           MailPayload    `json:"mail,omitempty"`
	EnquiryPayload EnquiryPayload `json:"enquiry,omitempty"`
}

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

type EnquiryPayload struct {
	UserID     string `json:"user_id"`
	PropertyID string `json:"property_id"`
	Name       string `json:"name"`
	Location   string `json:"location"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
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
	case "mail":
		//sendMail(c, requestPayload.Mail)
		lac.sendMailViaRabbit(c, requestPayload.Mail)
	case "enquiry_mail":
		lac.sendEnquiryMailViaRabbit(c, requestPayload.EnquiryPayload)

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

func (lac *LocalApiConfig) LogViaGRPC(c *gin.Context) {
	var requestPayload RequestPayload

	err := helpers.ReadJSON(c, &requestPayload)
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	conn, err := grpc.NewClient("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}
	defer conn.Close()

	cc := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = cc.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})

	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = "logged via grpc"

	helpers.WriteJSON(c, http.StatusAccepted, payload)
}

func (lac *LocalApiConfig) PaymentViaGRPC(c *gin.Context) {
	var paymentPayload PaymentPayload

	err := helpers.ReadJSON(c, &paymentPayload)
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	conn, err := grpc.NewClient("payment-service:50002", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}
	defer conn.Close()

	cc := payment.NewPaymentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	paymentResponse, err := cc.ProcessPayment(ctx, &payment.PaymentRequest{
		NewPayment: &payment.Payment{
			CardNumber:     paymentPayload.CardNumber,
			CardHolderName: paymentPayload.CardHolderName,
			CardCvv:        paymentPayload.CardCVV,
			CardExpiry:     paymentPayload.CardExpiry,
			Amount:         paymentPayload.Amount,
			Currency:       paymentPayload.Currency,
		},
	})

	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = "payment processed via grpc" + paymentResponse.String()

	helpers.WriteJSON(c, http.StatusAccepted, payload)

}

func (lac *LocalApiConfig) sendEnquiryMailViaRabbit(c *gin.Context, mail EnquiryPayload) {
	emitter, err := actions.NewEmitter(lac.Rabbit, "mail_exchange", "enquiry_mail")
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	// prepare enquiry email
	enquiryEmail := EnquiryMailPayload{
		From:      "chinmayanand896@icloud.com",
		To:        mail.UserID,
		Subject:   "Thank you for your enquiry.",
		Message:   fmt.Sprintf("Thank you for your enquiry about propery name %s and Id %s at location %s", mail.Name, mail.PropertyID, mail.Location),
		Timestamp: time.Now(),
	}

	j, _ := json.Marshal(&enquiryEmail)
	err = emitter.Emit(string(j))
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	var payload helpers.JsonResponse
	payload.Error = false
	payload.Message = "enquiry mail has been sent via rabbit"

	helpers.WriteJSON(c, http.StatusAccepted, payload)
}
