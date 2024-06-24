package handlers

import (
	"broker/actions"
	"broker/helpers"
	"encoding/json"
	"fmt"
	"github.com/Meenachinmay/microservice-shared/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func (lac *LocalApiConfig) logEventViaRabbit(c *gin.Context, l LogJSONPayload) {
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

func (lac *LocalApiConfig) sendEnquiryMailViaRabbit(c *gin.Context, mail types.EnquiryPayload) {
	emitter, err := actions.NewEmitter(lac.Rabbit, "mail_exchange", "enquiry_mail")
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	// prepare enquiry email
	enquiryEmail := EnquiryMailPayload{
		From:      "chinmayanand896@icloud.com",
		To:        "TESTING@EMAIL.COM",
		Subject:   "Thank you for your enquiry.",
		Message:   fmt.Sprintf("Thank you for your enquiry about propery name %s and Id %d at location %s", mail.PropertyName, mail.PropertyID, mail.PropertyLocation),
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
