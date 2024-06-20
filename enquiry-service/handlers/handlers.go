package handlers

import (
	"encoding/json"
	"enquiry-service/helpers"
	"enquiry-service/internal/database"
	"enquiry-service/mqactions"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type EnquiryPayload struct {
	UserID     int32  `json:"user_id"`
	PropertyID int32  `json:"property_id"`
	Name       string `json:"name"`
	Location   string `json:"location"`
}

type EnquiryMailPayload struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// HandleANewEnquiry To handle a new enquiry we need following data
/*
	start a counter to keep the track of time taken
	enquiryPayload, fetch user from database, fetch property from database, calculate priority

	NOTE: When a user makes an enquiry from the client side, we can send all the required details from the client itself
    and can save few database operations in this handler. But now for the practice purpose only we are fetching data here inside the handler.
*/
func (localApiConfig *LocalApiConfig) HandleANewEnquiry(c *gin.Context) {
	// starting counter
	var startTimer = time.Now()

	// read json
	var payload EnquiryPayload

	if err := helpers.ReadJSON(c, &payload); err != nil {
		helpers.ErrorJSON(c, err)
		return
	}
	// validate

	// insert into database
	newEnquiry, err := localApiConfig.DB.CreateEnquiry(c, database.CreateEnquiryParams{
		UserID:     payload.UserID,
		PropertyID: payload.PropertyID,
	})
	if err != nil {
		helpers.ErrorJSON(c, errors.New("couldn't save enquiry in DB"), http.StatusInternalServerError)
		return
	} else {
		log.Println("new enquiry created", newEnquiry)
	}

	// update the count of total enquiries made by current user
	updatedUser, err := localApiConfig.DB.AddNewEnquiryToUserById(c, payload.UserID)
	if err != nil {
		helpers.ErrorJSON(c, errors.New("couldn't update enquiry count in DB for user"), http.StatusInternalServerError)
		return
	} else {
		log.Println("enquiry count updated for user", updatedUser)
	}

	// decide user communication
	totalEnquiries, err := localApiConfig.getTotalEnquiriesLastWeek(c, updatedUser)
	if err != nil {
		helpers.ErrorJSON(c, errors.New("couldn't get total enquiries"), http.StatusInternalServerError)
		return
	} else {
		log.Printf("total enquiries count is %d for UserId %d", totalEnquiries, updatedUser.ID)
	}

	// fetch the property detail here - (only this database operation can be saved if we get data from client)
	foundProperty, err := localApiConfig.DB.GetAPropertyDetailsById(c, payload.PropertyID)
	if err != nil {
		helpers.ErrorJSON(c, errors.New("couldn't get property details"), http.StatusInternalServerError)
		return
	}

	// prepare notifyPayload (for now it's same as EnquiryEmailPayload)
	mailPayload := EnquiryMailPayload{
		From:      "chinmayanand896@gmail.com",
		To:        updatedUser.Email,
		Subject:   "Thank you for your enquiry.",
		Message:   fmt.Sprintf("Thank you for your enquiry about propery name %s at location %s", foundProperty.Name, foundProperty.Location),
		Timestamp: startTimer,
	}

	var responsePayload helpers.JsonResponse
	responsePayload.Error = false
	responsePayload.Message = "Thank you for your enquiry, We will reach you before you finish your coffee"
	responsePayload.Data = mailPayload

	helpers.WriteJSON(c, http.StatusAccepted, responsePayload)

	// execute communication
	go localApiConfig.notifyUserAboutEnquiry(c, updatedUser, totalEnquiries, mailPayload)
}

func (localApiConfig *LocalApiConfig) getTotalEnquiriesLastWeek(c *gin.Context, updatedUser database.User) (int, error) {
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	var count int64
	count, err := localApiConfig.DB.CountEnquiriesForUserInLastWeek(c, database.CountEnquiriesForUserInLastWeekParams{
		UserID:      updatedUser.ID,
		EnquiryDate: oneWeekAgo,
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (localApiConfig *LocalApiConfig) notifyUserAboutEnquiry(c *gin.Context, user database.User, totalEnquiries int, mailPayload EnquiryMailPayload) {
	if totalEnquiries >= 10 {
		log.Printf("Calling to the user %d...\n", user.ID)
		return
	} else if totalEnquiries >= 1 && totalEnquiries <= 3 {
		log.Printf("Sending sms to the user %d...\n", user.ID)
		return
	} else {
		log.Printf("Sending Email to the user %s...\n", user.Email)
		// communicate with mail-service to send an email using rabbitmq.
		localApiConfig.sendEmail(c, mailPayload)
		return
	}
	return
}

func (localApiConfig *LocalApiConfig) sendEmail(c *gin.Context, payload EnquiryMailPayload) {
	emitter, err := mqactions.NewEmitter(localApiConfig.Rabbit, "mail_exchange", "enquiry_mail")
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}

	j, _ := json.Marshal(&payload)
	err = emitter.Emit(string(j))
	if err != nil {
		helpers.ErrorJSON(c, err)
		return
	}
}
