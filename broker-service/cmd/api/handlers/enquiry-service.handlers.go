package handlers

import (
	"broker/helpers"
	"broker/internal/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (lac *LocalApiConfig) FetchAllProperties(c *gin.Context) {
	url := config.EnquiryServiceURL + "/fetch-properties"
	respBody, err := helpers.MakeHTTPRequest(c, config.HttpGet, url, nil)
	log.Println("After the MakeHTTPRequest method completes")
	if err != nil {
		helpers.ErrorJSON(c, err, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(c, http.StatusAccepted, respBody)
}

func (lac *LocalApiConfig) routeEnquiryToEnquiryService(c *gin.Context, enquiryPayload EnquiryPayload) {
	url := config.EnquiryServiceURL + "/handle-enquiry"
	respBody, err := helpers.MakeHTTPRequest(c, config.HttpPost, url, enquiryPayload)
	if err != nil {
		helpers.ErrorJSON(c, err, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(c, http.StatusAccepted, respBody)
}
