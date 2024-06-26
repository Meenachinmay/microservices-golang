package handlers

import (
	"enquiry-service/helpers"
	"enquiry-service/internal/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type EnquiryMailPayloadUsingSendgrid struct {
	To               string    `json:"to"`
	ToName           string    `json:"to_name"`
	Subject          string    `json:"subject"`
	PropertyName     string    `json:"name"`
	PropertyLocation string    `json:"location"`
	Timestamp        time.Time `json:"timestamp"`
}

// HandleFetchAllProperties
// Below is the method to fetch all the enquiries saved in the database.
// /*
func (localApiConfig *LocalApiConfig) HandleFetchAllProperties(c *gin.Context) {
	properties, err := localApiConfig.DB.FetchAllProperties(c)
	if err != nil {
		helpers.ErrorJSON(c, err, http.StatusInternalServerError)
		return
	}

	propertiesJSON := database.ConvertPropertiesToJSON(properties)
	responsePayload := helpers.JsonResponse{
		Error:   false,
		Message: "Fetched all properties",
		Data:    propertiesJSON,
	}

	helpers.WriteJSON(c, http.StatusAccepted, responsePayload)
}

func (localApiConfig *LocalApiConfig) HandleFetchScheduledTasks(c *gin.Context) {
	tasks, err := localApiConfig.DB.GetDueSchedules(c)
	if err != nil {
		helpers.ErrorJSON(c, err, http.StatusInternalServerError)
		return
	}
	tasksJSON := database.ConvertTasksToJSON(tasks)
	responsePayload := helpers.JsonResponse{
		Error:   false,
		Message: "Fetched scheduled tasks",
		Data:    tasksJSON,
	}
	helpers.WriteJSON(c, http.StatusAccepted, responsePayload)
}
