package handlers

import (
	"broker/helpers"
	"broker/internal/config"
	"github.com/Meenachinmay/microservice-shared/routes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (lac *LocalApiConfig) FetchAllProperties(c *gin.Context) {
	respBody, err := helpers.MakeHTTPRequest(c, config.HttpGet, routes.EnquiryServiceFetchAllProperties, nil)
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
