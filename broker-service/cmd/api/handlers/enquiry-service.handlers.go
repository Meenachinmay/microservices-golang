package handlers

import (
	"broker/helpers"
	"broker/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (lac *LocalApiConfig) FetchAllProperties(c *gin.Context) {
	url := config.EnquiryServiceURL + "/fetch-properties"
	respBody, err := helpers.MakeHTTPRequest(c, config.HttpGet, url)
	if err != nil {
		helpers.ErrorJSON(c, err, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(c, http.StatusAccepted, respBody)
}

func (lac *LocalApiConfig) routeEnquiryToEnquiryService(c *gin.Context, enquiryPayload EnquiryPayload) {
	respBody, err := helpers.MakeHTTPRequest(c, config.HttpPost, config.EnquiryServiceURL, enquiryPayload)
	if err != nil {
		helpers.ErrorJSON(c, err, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(c, http.StatusAccepted, respBody)
}
