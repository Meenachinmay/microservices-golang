package handlers

import (
	"broker/helpers"
	"broker/internal/config"
	"github.com/Meenachinmay/microservice-shared/routes"
	"github.com/Meenachinmay/microservice-shared/types"
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

func (lac *LocalApiConfig) routeEnquiryToEnquiryService(c *gin.Context, enquiryPayload types.EnquiryPayload) {
	respBody, err := helpers.MakeHTTPRequest(c, config.HttpPost, routes.EnquiryServiceHandleEnquiryRoute, enquiryPayload)
	if err != nil {
		helpers.ErrorJSON(c, err, http.StatusInternalServerError)
		return
	}
	helpers.WriteJSON(c, http.StatusAccepted, respBody)
}

func (lac *LocalApiConfig) FetchAllScheduledTasks(c *gin.Context) {
	url := config.EnquiryServiceURL + "/handle-fetch-tasks"
	respBody, err := helpers.MakeHTTPRequest(c, config.HttpGet, url, nil)
	if err != nil {
		helpers.ErrorJSON(c, err, http.StatusInternalServerError)
		return
	}
	helpers.WriteJSON(c, http.StatusAccepted, respBody)
}
