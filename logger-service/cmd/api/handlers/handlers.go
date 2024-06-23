package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log-service/helpers"
	"log-service/internal/database"
	"net/http"
)

// WriteLog handle to write log using http
func (apiConfig *LocalApiConfig) WriteLog(c *gin.Context) {
	// read json
	var requestPayload database.LogJSONPayload

	err := helpers.ReadJSON(c, &requestPayload)
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	// insert into the database
	newLog, err := apiConfig.DB.InsertLog(c, database.InsertLogParams{
		ServiceName: requestPayload.ServiceName,
		LogData:     requestPayload.LogData,
	})
	if err != nil {
		helpers.ErrorJSON(c, errors.New("error inserting log:[WriteLogHandler]"), http.StatusInternalServerError)
		return
	}

	response := helpers.JsonResponse{
		Error:   false,
		Message: "logged",
		Data:    newLog,
	}

	helpers.WriteJSON(c, http.StatusAccepted, response)
}
